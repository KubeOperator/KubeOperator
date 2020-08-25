package service

import "github.com/KubeOperator/KubeOperator/pkg/manifest"

type ManifestService interface {
	List() []manifest.Manifest
}

type manifestService struct {
}

func NewManifestService() ManifestService {
	return &manifestService{}
}

func (*manifestService) List() []manifest.Manifest {
	return manifest.Manifests
}
