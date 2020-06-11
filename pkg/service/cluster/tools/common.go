package tools

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (cta *CLusterToolsAdm) EnsureClusterReachable(c *Tool) error {
	_, err := c.KubernetesClient.CoreV1().
		Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	return nil
}
