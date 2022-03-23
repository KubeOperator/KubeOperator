package service

import (
	"context"
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ClusterEventService interface {
	List(clusterName string) ([]dto.ClusterEventDTO, error)
	ListLimitOneDay(clusterName string) ([]dto.ClusterEventDTO, error)
	ExistEventUid(uid, clusterId string) (bool, error)
	Save(event model.ClusterEvent) error
	GetNpd(clusterName string) (bool, error)
	DeleteNpd(clusterName string) (bool, error)
	CreateNpd(clusterName string) (bool, error)
}

type clusterEventService struct {
	clusterEventRepo repository.ClusterEventRepository
	clusterRepo      repository.ClusterRepository
	clusterService   ClusterService
	systemRepo       repository.SystemSettingRepository
}

func NewClusterEventService() ClusterEventService {
	return &clusterEventService{
		clusterEventRepo: repository.NewClusterEventRepository(),
		clusterRepo:      repository.NewClusterRepository(),
		clusterService:   NewClusterService(),
		systemRepo:       repository.NewSystemSettingRepository(),
	}
}

func (c clusterEventService) List(clusterName string) ([]dto.ClusterEventDTO, error) {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}
	var eventDTOs []dto.ClusterEventDTO
	events, err := c.clusterEventRepo.List(cluster.ID)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		eventDTOs = append(eventDTOs, dto.ClusterEventDTO{
			ClusterEvent: event,
		})
	}
	return eventDTOs, nil
}

func (c clusterEventService) ListLimitOneDay(clusterName string) ([]dto.ClusterEventDTO, error) {
	cluster, err := c.clusterRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}
	var eventDTOs []dto.ClusterEventDTO
	events, err := c.clusterEventRepo.ListLimitOneDay(cluster.ID)
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		eventDTOs = append(eventDTOs, dto.ClusterEventDTO{
			ClusterEvent: event,
		})
	}
	return eventDTOs, nil
}

func (c clusterEventService) ExistEventUid(uid, clusterId string) (bool, error) {
	events, err := c.clusterEventRepo.ListByUidAndClusterId(uid, clusterId)
	if err != nil {
		return false, err
	}
	if len(events) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (c clusterEventService) Save(event model.ClusterEvent) error {
	return c.clusterEventRepo.Save(&event)
}

func (c clusterEventService) GetNpd(clusterName string) (bool, error) {
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return false, err
	}
	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return false, err
	}
	if cluster.Status == constant.ClusterRunning {
		client, err := kubernetes.NewKubernetesClient(&secret.KubeConf)
		if err != nil {
			return false, err
		}
		pod, err := client.CoreV1().Pods("kube-system").Get(context.Background(), "node-problem-detector", metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, nil
}

func (c clusterEventService) CreateNpd(clusterName string) (bool, error) {
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return false, err
	}
	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return false, err
	}
	if cluster.Status == constant.ClusterRunning {
		client, err := kubernetes.NewKubernetesClient(&secret.KubeConf)
		if err != nil {
			return false, err
		}
		cm := corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "node-problem-detector-config",
				Namespace: "kube-system",
			},
			Data: map[string]string{
				"abrt-adaptor.json":    constant.AbrtAdaptor,
				"docker-monitor.json":  constant.DockerMonitor,
				"kernel-monitor.json":  constant.KernelMonitor,
				"systemd-monitor.json": constant.SystemdMonitor,
			},
		}
		_, err = client.CoreV1().ConfigMaps("kube-system").Create(context.Background(), &cm, metav1.CreateOptions{})
		if err != nil {
			return false, err
		}
		var test = true
		_, err = client.AppsV1().DaemonSets("kube-system").Create(context.Background(), &v1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "node-problem-detector",
				Namespace: "kube-system",
				Labels:    map[string]string{"app": "node-problem-detector"},
			},
			Spec: v1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "node-problem-detector"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "node-problem-detector"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name: "node-problem-detector",
								Command: []string{
									"/node-problem-detector",
									"--logtostderr",
									"--config.system-log-monitor=/config/abrt-adaptor.json,/config/docker-monitor.json,/config/kernel-monitor.json,/config/systemd-monitor.json",
								},
								Image: fmt.Sprintf("%s:%d/kubeoperator/node-problem-detector:v0.8.1", constant.LocalRepositoryDomainName, constant.LocalDockerRepositoryPort),
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										"cpu":    resource.MustParse("10m"),
										"memory": resource.MustParse("80Mi"),
									},
									Requests: corev1.ResourceList{
										"cpu":    resource.MustParse("10m"),
										"memory": resource.MustParse("80Mi"),
									},
								},
								ImagePullPolicy: corev1.PullPolicy(corev1.PullIfNotPresent),
								SecurityContext: &corev1.SecurityContext{
									Privileged: &test,
								},
								Env: []corev1.EnvVar{
									{
										Name: "NODE_NAME",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "spec.nodeName",
											},
										},
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "log",
										MountPath: "/var/log",
										ReadOnly:  true,
									},
									{
										Name:      "kmsg",
										MountPath: "/dev/kmsg",
										ReadOnly:  true,
									},
									{
										Name:      "localtime",
										MountPath: "/etc/localtime",
										ReadOnly:  true,
									},
									{
										Name:      "config",
										MountPath: "/config",
										ReadOnly:  true,
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "log",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/var/log/",
									},
								},
							},
							{
								Name: "kmsg",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/dev/kmsg",
									},
								},
							},
							{
								Name: "localtime",
								VolumeSource: corev1.VolumeSource{
									HostPath: &corev1.HostPathVolumeSource{
										Path: "/etc/localtime",
									},
								},
							},
							{
								Name: "config",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "node-problem-detector-config",
										},
										Items: []corev1.KeyToPath{
											{
												Key:  "abrt-adaptor.json",
												Path: "abrt-adaptor.json",
											},
											{
												Key:  "docker-monitor.json",
												Path: "docker-monitor.json",
											},
											{
												Key:  "kernel-monitor.json",
												Path: "kernel-monitor.json",
											},
											{
												Key:  "systemd-monitor.json",
												Path: "systemd-monitor.json",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return false, err
		}
		err = wait.Poll(5*time.Second, 5*time.Minute, func() (done bool, err error) {
			ds, err := client.AppsV1().DaemonSets("kube-system").Get(context.Background(), "node-problem-detector", metav1.GetOptions{})
			if err != nil {
				err = client.CoreV1().ConfigMaps("kube-system").Delete(context.Background(), "node-problem-detector-config", metav1.DeleteOptions{})
				if err != nil {
					return true, err
				}
				return true, err
			}
			if ds.Status.DesiredNumberScheduled == ds.Status.NumberReady {
				return true, nil
			}
			return true, nil
		})
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, nil
}

func (c clusterEventService) DeleteNpd(clusterName string) (bool, error) {
	secret, err := c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return false, err
	}
	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return false, err
	}
	if cluster.Status == constant.ClusterRunning {
		client, err := kubernetes.NewKubernetesClient(&secret.KubeConf)
		if err != nil {
			return false, err
		}
		err = client.AppsV1().DaemonSets("kube-system").Delete(context.Background(), "node-problem-detector", metav1.DeleteOptions{})
		if err != nil {
			return false, err
		}
		err = client.CoreV1().ConfigMaps("kube-system").Delete(context.Background(), "node-problem-detector-config", metav1.DeleteOptions{})
		if err != nil {
			return false, err
		}
		return true, err
	}
	return false, err
}
