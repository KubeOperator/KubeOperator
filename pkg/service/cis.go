package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	uuid "github.com/satori/go.uuid"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

type CisService interface {
	Page(num, size int, clusterName string) (*page.Page, error)
	List(clusterName string) ([]dto.CisTask, error)
	Create(clusterName string) (*dto.CisTask, error)
	Delete(clusterName, id string) error
	Get(clusterName, id string) (*dto.CisTaskDetail, error)
}

type cisService struct {
	clusterRepo    repository.ClusterRepository
	clusterService ClusterService
}

func NewCisService() CisService {
	return &cisService{
		clusterRepo:    repository.NewClusterRepository(),
		clusterService: NewClusterService(),
	}
}

//type CisResultList []model.CisTaskResult
//
//func (c CisResultList) Len() int {
//	return len(c)
//}
//
//func (c CisResultList) Less(i, j int) bool {
//	c1 := c[i].Number
//	c2 := c[j].Number
//
//	c1s := strings.Split(c1, ".")
//	c2s := strings.Split(c2, ".")
//
//	var maxLen int
//	if len(c1s) > len(c2s) {
//		maxLen = len(c1s)
//	} else {
//		maxLen = len(c2s)
//	}
//	for i := 0; i < maxLen; i++ {
//		a, _ := strconv.Atoi(c1s[i])
//		b, _ := strconv.Atoi(c2s[i])
//		if a == b {
//			continue
//		}
//		return a < b
//	}
//	return false
//}
//
//func (c CisResultList) Swap(i, j int) {
//	c[i], c[j] = c[j], c[i]
//
//}

func (c *cisService) Get(clusterName, id string) (*dto.CisTaskDetail, error) {
	var cisTask model.CisTaskWithResult
	if err := db.DB.First(&cisTask, &model.CisTaskWithResult{CisTask: model.CisTask{ID: id}}).Error; err != nil {
		return nil, err
	}
	var nodeList dto.CisNodeList
	if err := json.Unmarshal([]byte(cisTask.Result), &nodeList); err != nil {
		return nil, err
	}

	return &dto.CisTaskDetail{CisTaskWithResult: cisTask, NodeList: nodeList}, nil
}

func (*cisService) Page(num, size int, clusterName string) (*page.Page, error) {
	var cluster model.Cluster
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}
	p := page.Page{}
	var tasks []model.CisTask
	if err := db.DB.Model(&model.CisTask{}).
		Where("cluster_id = ?", cluster.ID).
		Count(&p.Total).
		Order("created_at desc").
		Offset((num - 1) * size).
		Limit(size).
		//Preload("Results").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	var dtos []dto.CisTask
	for _, task := range tasks {
		dtos = append(dtos, dto.CisTask{CisTask: task})
	}
	p.Items = dtos
	return &p, nil
}

const (
	CisTaskStatusCreating = "Creating"
	CisTaskStatusRunning  = "Running"
	CisTaskStatusSuccess  = "Success"
	CisTaskStatusFailed   = "Failed"
)

func (c *cisService) List(clusterName string) ([]dto.CisTask, error) {
	var cluster model.Cluster
	if err := db.DB.Where("name = ?", clusterName).First(&cluster).Error; err != nil {
		return nil, err
	}
	var tasks []model.CisTask
	if err := db.DB.
		Where("cluster_id = ?", cluster.ID).
		Preload("Results").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	var dtos []dto.CisTask
	for _, task := range tasks {
		dtos = append(dtos, dto.CisTask{CisTask: task})
	}
	return dtos, nil
}

func (c *cisService) Create(clusterName string) (*dto.CisTask, error) {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}
	var registery model.SystemRegistry
	if cluster.Spec.Architectures == constant.ArchAMD64 {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfAMD64).First(&registery).Error; err != nil {
			return nil, errors.New("load image pull port of arm failed")
		}
	} else {
		if err := db.DB.Where("architecture = ?", constant.ArchitectureOfARM64).First(&registery).Error; err != nil {
			return nil, errors.New("load image pull port of arm failed")
		}
	}
	localRepoPort := registery.RegistryPort

	var clusterTasks []model.CisTask
	db.DB.Where("status = ? AND cluster_id = ?", constant.ClusterRunning, cluster.ID).Find(&clusterTasks)
	if len(clusterTasks) > 0 {
		return nil, errors.New("CIS_TASK_ALREADY_RUNNING")
	}
	tx := db.DB.Begin()
	task := model.CisTask{
		ClusterID: cluster.ID,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Status:    CisTaskStatusCreating,
	}
	err = tx.Create(&task).Error
	if err != nil {
		return nil, err
	}

	endpoints, err := c.clusterService.GetApiServerEndpoints(cluster.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	secret, err := c.clusterService.GetSecrets(cluster.Name)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	client, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	go Do(&cluster, client, &task, localRepoPort)
	return &dto.CisTask{CisTask: task}, nil
}

func (c *cisService) Delete(clusterName, id string) error {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return err
	}
	if err := db.DB.Where("id = ? AND cluster_id = ?", id, cluster.ID).Delete(&model.CisTask{}).Error; err != nil {
		return err
	}
	return nil
}

func Do(cluster *model.Cluster, client *kubernetes.Clientset, task *model.CisTask, port int) {
	taskWithResult := &model.CisTaskWithResult{CisTask: *task}

	taskWithResult.Status = CisTaskStatusRunning
	db.DB.Save(&taskWithResult)

	jobId := fmt.Sprintf("kube-bench-%s", uuid.NewV4().String())
	j := v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobId,
		},
		Spec: v1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "kube-bench"}},
				Spec: corev1.PodSpec{
					HostPID:       true,
					RestartPolicy: "Never",
					Containers: []corev1.Container{
						{
							Name:    "kube-bench",
							Image:   fmt.Sprintf("%s:%d/kubeoperator/kube-bench:v0.0.1", constant.LocalRepositoryDomainName, port),
							Command: []string{"kube-bench", "--json"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "var-lib-etcd",
									MountPath: "/var/lib/etcd",
									ReadOnly:  true,
								},
								{
									Name:      "var-lib-kubelet",
									MountPath: "/var/lib/kubelet",
									ReadOnly:  true,
								},
								{
									Name:      "etc-systemd",
									MountPath: "/etc/systemd",
									ReadOnly:  true,
								},
								{
									Name:      "etc-kubernetes",
									MountPath: "/etc/kubernetes",
									ReadOnly:  true,
								},
								{
									Name:      "usr-bin",
									MountPath: "/usr/local/mount-from-host/bin",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "var-lib-etcd",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/etcd",
								},
							},
						},
						{
							Name: "var-lib-kubelet",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/kubelet",
								},
							},
						},
						{
							Name: "etc-systemd",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/systemd",
								},
							},
						},
						{
							Name: "etc-kubernetes",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/kubernetes",
								},
							},
						},
						{
							Name: "usr-bin",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/usr/bin",
								},
							},
						},
					},
				},
			},
		},
	}

	resp, err := client.BatchV1().Jobs(constant.DefaultNamespace).Create(context.TODO(), &j, metav1.CreateOptions{})
	if err != nil {
		taskWithResult.Message = err.Error()
		taskWithResult.Status = CisTaskStatusFailed
		db.DB.Save(&taskWithResult)
		return
	}

	err = wait.Poll(5*time.Second, 5*time.Minute, func() (done bool, err error) {
		job, err := client.BatchV1().Jobs(constant.DefaultNamespace).Get(context.TODO(), resp.Name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if job.Status.Succeeded > 0 {
			pds, err := client.CoreV1().Pods(constant.DefaultNamespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: fmt.Sprintf("job-name=%s", resp.Name),
			})
			if err != nil {
				return true, err
			}
			for _, p := range pds.Items {
				if p.Status.Phase == corev1.PodSucceeded {
					r := client.CoreV1().Pods(constant.DefaultNamespace).GetLogs(p.Name, &corev1.PodLogOptions{})
					bs, err := r.DoRaw(context.TODO())
					if err != nil {
						return true, err
					}
					taskWithResult.Result = string(bs)
					var nodeList dto.CisNodeList
					if err := json.Unmarshal(bs, &nodeList); err != nil {
						return true, err
					}
					for i := range nodeList {
						taskWithResult.TotalPass += nodeList[i].TotalPass
						taskWithResult.TotalFail += nodeList[i].TotalFail
						taskWithResult.TotalWarn += nodeList[i].TotalWarn
						taskWithResult.TotalInfo += nodeList[i].TotalInfo
					}
					taskWithResult.Status = CisTaskStatusSuccess
					db.DB.Save(&taskWithResult)
				}
			}
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		taskWithResult.Message = err.Error()
		taskWithResult.Status = CisTaskStatusFailed
		db.DB.Save(&taskWithResult)
		return
	}
	err = client.BatchV1().Jobs(constant.DefaultNamespace).Delete(context.TODO(), resp.Name, metav1.DeleteOptions{})
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
}
