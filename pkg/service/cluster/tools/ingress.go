package tools

import (
	"context"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

const (
	ConditionNameDone = "EnsureDone"
)

func (cta *CLusterToolsAdm) IngressInstall(c *Tool) error {

}

func (cta *CLusterToolsAdm) EnsureIngressInstall(c *Tool) error {
	chart, err := helm.LoadCharts(constant.IngressChartPath)
	if err != nil {
		return nil
	}
	_, err = c.HelmClient.Install(constant.IngressReleaseName, chart, c.Values)
	if err != nil {
		return err
	}
	return nil
}

func (cta *CLusterToolsAdm) EnsureIngressRunning(c *Tool) error {
	timeout := true
	err := wait.Poll(5*time.Second, 10*time.Minute, func() (done bool, err error) {
		ds, err := c.KubernetesClient.
			AppsV1().Deployments(constant.KoNamespaceName).
			List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return true, err
		}
		for _, d := range ds.Items {
			if d.Status.AvailableReplicas > 0 {
				timeout = false
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	if timeout {
		return errors.New("wait nginx-ingress start timeout")
	}
	return nil
}
