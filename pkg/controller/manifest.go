package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/manifest"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

type ManifestController struct {
	manifestService service.ManifestService
}

func NewManifestController() *ManifestController {
	return &ManifestController{
		manifestService: service.NewManifestService(),
	}
}

func (m *ManifestController) Get() []manifest.Manifest {
	return m.manifestService.List()
}
