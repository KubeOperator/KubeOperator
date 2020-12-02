package v2

import (
	"k8s.io/helm/pkg/downloader"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

func (h *HelmV2) DependencyUpdate(chartPath string) error {
	out := helm.NewLogWriter(h.logger)
	man := downloader.Manager{
		Out:       out,
		ChartPath: chartPath,
		HelmHome:  helmHome(),
		Getters:   getterProviders(),
	}
	return man.Update()
}
