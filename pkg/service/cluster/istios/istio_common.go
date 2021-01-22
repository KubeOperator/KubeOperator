package istios

import (
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/client-go/kubernetes"
)

type IstioInterface interface {
	Install() error
	Uninstall() error
}

type IstioHelmInfo struct {
	Namespace     string
	Cluster       model.Cluster
	LocalhostName string
	HelmClient    helm.Interface
	KubeClient    *kubernetes.Clientset
}

func preInstallChart(h helm.Interface, istio *model.ClusterIstio) error {
	rs, err := h.List()
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Name == istio.Name {
			_, err := h.Uninstall(istio.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func installChart(h helm.Interface, istio *model.ClusterIstio, valueMap map[string]interface{}, chartName string) error {
	err := preInstallChart(h, istio)
	if err != nil {
		return err
	}
	m, err := MergeValueMap(valueMap)
	if err != nil {
		return err
	}
	_, err = h.Install(istio.Name, chartName, "", m)
	if err != nil {
		return err
	}
	return nil
}

func MergeValueMap(source map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	var valueStrings []string
	for k, v := range source {
		str := fmt.Sprintf("%s=%v", k, v)
		valueStrings = append(valueStrings, str)
	}
	for _, str := range valueStrings {
		err := strvals.ParseInto(str, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func uninstall(istio *model.ClusterIstio, h helm.Interface) error {
	rs, err := h.List()
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Name == istio.Name {
			_, _ = h.Uninstall(istio.Name)
		}
	}
	return nil
}
