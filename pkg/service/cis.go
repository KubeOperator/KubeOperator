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
}

type cisService struct {
	clusterRepo    repository.ClusterRepository
	clusterService ClusterService
	systemRepo     repository.SystemSettingRepository
}

func NewCisService() CisService {
	return &cisService{
		clusterRepo:    repository.NewClusterRepository(),
		clusterService: NewClusterService(),
		systemRepo:     repository.NewSystemSettingRepository(),
	}
}

type CisSummary struct {
	Tests []CisTest `json:"tests"`
}

type CisTest struct {
	Results []CisResult `json:"results"`
}

type CisResult struct {
	TestNumber  string `json:"test_number"`
	TestDesc    string `json:"test_desc"`
	Remediation string `json:"remediation"`
	Status      string `json:"status"`
	Scored      bool   `json:"scored"`
}

func (*cisService) Page(num, size int, clusterName string) (*page.Page, error) {
	var cluster model.Cluster
	if err := db.DB.Where(&model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return nil, err
	}
	p := page.Page{}
	var tasks []model.CisTask
	if err := db.DB.Model(&model.CisTask{}).
		Count(&p.Total).
		Offset((num - 1) * size).
		Limit(size).
		Where(&model.CisTask{ClusterID: cluster.ID}).
		Preload("Results").
		Order("created_at desc").
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
	if err := db.DB.Where(&model.Cluster{Name: clusterName}).First(&cluster).Error; err != nil {
		return nil, err
	}
	var tasks []model.CisTask
	if err := db.DB.
		Where(&model.CisTask{ClusterID: cluster.ID}).
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

	var clusterTasks []model.CisTask
	db.DB.Where(&model.CisTask{Status: constant.ClusterRunning, ClusterID: cluster.ID}).Find(&clusterTasks)
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
	go Do(&cluster, client, &task)
	return &dto.CisTask{CisTask: task}, nil
}

func (c *cisService) Delete(clusterName, id string) error {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return err
	}
	if err := db.DB.Where(&model.CisTask{ID: id, ClusterID: cluster.ID}).Delete(&model.CisTask{}).Error; err != nil {
		return err
	}
	return nil
}

func Do(cluster *model.Cluster, client *kubernetes.Clientset, task *model.CisTask) {
	task.Status = CisTaskStatusRunning
	db.DB.Save(&task)
	systemRepo := repository.NewSystemSettingRepository()

	localIP, err := systemRepo.Get("ip")
	if err != nil || localIP.Value == "" {
		task.Message = "local ip is null"
		task.Status = CisTaskStatusFailed
		db.DB.Save(&task)
		return
	}
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
							Image:   fmt.Sprintf("%s:%d/kubeoperator/kube-bench:v0.0.1-%s", localIP.Value, constant.LocalDockerRepositoryPort, cluster.Spec.Architectures),
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
		task.Message = err.Error()
		task.Status = CisTaskStatusFailed
		db.DB.Save(&task)
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
					var summarys []CisSummary
					err = json.Unmarshal(bs, &summarys)
					if err != nil {
						return true, err
					}
					var results []model.CisTaskResult
					for _, summary := range summarys {
						for _, test := range summary.Tests {
							for _, res := range test.Results {
								results = append(results, model.CisTaskResult{
									ID:          uuid.NewV4().String(),
									ClusterID:   cluster.ID,
									CisTaskId:   task.ID,
									Number:      res.TestNumber,
									Desc:        res.TestDesc,
									Remediation: res.Remediation,
									Status:      res.Status,
									Scored:      res.Scored,
								})
							}
						}
					}
					task.Results = results
					task.Status = CisTaskStatusSuccess
					err = db.DB.Save(&task).Error
					if err != nil {
						log.Error(err)
					}
				}
			}
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		task.Message = err.Error()
		task.Status = CisTaskStatusFailed
		db.DB.Save(&task)
		return
	}
	err = client.BatchV1().Jobs(constant.DefaultNamespace).Delete(context.TODO(), resp.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Error(err.Error())
		return
	}
}
