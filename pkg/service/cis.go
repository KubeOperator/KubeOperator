package service

import (
	"context"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CisService interface {
}

type cisService struct {
}

func NewCisService() CisService {
	return &cisService{}
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

func Do(client kubernetes.Clientset) {
	j := v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kube-bench",
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
							Image:   "",
							Command: []string{"kube-bench", "--json", "--benchmark cis-1.5"},
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
		log.Error(err.Error())
		return
	}
	fmt.Println(resp)

}
