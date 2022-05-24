package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/polaris"
)

type GradeService interface {
	GetGrade(clusterName string) (*dto.ClusterGrade, error)
}

type gradeService struct {
	clusterService ClusterService
}

func NewGradeService() GradeService {
	return &gradeService{
		clusterService: NewClusterService(),
	}
}

func (g gradeService) GetGrade(clusterName string) (*dto.ClusterGrade, error) {
	cluster, err := g.clusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}
	secret, err := g.clusterService.GetSecrets(cluster.Name)
	if err != nil {
		return nil, err
	}

	if cluster.Status == constant.ClusterRunning {
		result, err := polaris.RunGrade(&polaris.Config{
			Host:  cluster.Spec.LbKubeApiserverIp,
			Port:  cluster.Spec.KubeApiServerPort,
			Token: secret.KubernetesToken,
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return nil, errors.New("CLUSTER_IS_NOT_RUNNING")
	}
}
