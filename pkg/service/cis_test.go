package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	client, _ := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Host:  "https://172.16.10.101",
		Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6InZVM2pGWFQzQW02eUtKNlBXUk1iSGlCQUhpbG5rbEF2MzctVlRsT1NFejAifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi1xNHRqaiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImRhOWI1NmIxLTJiNjYtNGFjMS04ZjhiLTI4MWU0YjI1YjdhMyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.jdTAgrLYv31C0Ru4iAIAFoODF7S9FVY5jQ-Pi3iAvbv7O8zpkuFFV1hI4iQr6Jnl7UGjrNeBbMsDE14Ttaft-Fj0UFfntL8BxjYtlTEdZ1yJ6zzfi2Yv5bgQCkyVaQ07IMupxGhZezbtQrRuLNPOx7O0Hz6xlcPDy9diUAXoGFErK963DpLxlmRkSzOnkO21H_EKuOoa13pk97YQhfAORMcvxuN5oV7_AlvvU0TcV9pGXifbkKUH7q0tBMQhBdZ588Kva8tCI6CeicC4u60UU-YpmO8O3M1MgSffLHWl-q_fUSixq7pB5X02NWF4hJLWGBijrQfbN4lSXI05xcODhg",
		Port:  8443,
	})
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
							Image:   "aquasec/kube-bench:latest",
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
		fmt.Println(err.Error())
		return
	}

	err = wait.Poll(5*time.Second, 5*time.Minute, func() (done bool, err error) {
		job, err := client.BatchV1().Jobs(constant.DefaultNamespace).Get(context.TODO(), resp.Name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if job.Status.Succeeded > 0 {
			pds, err := client.CoreV1().Pods(constant.DefaultNamespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: "job-name=kube-bench",
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
					fmt.Println(summarys)
				}
				break
			}
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = client.BatchV1().Jobs(constant.DefaultNamespace).Delete(context.TODO(), "kube-bench", metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
