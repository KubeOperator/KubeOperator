package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/polaris"
)

type GradeService interface {
	GetGrade(clusterName string) (*dto.ClusterGrade, error)
}

type gradeService struct {
	clusterRepo repository.ClusterRepository
}

func NewGradeService() GradeService {
	return &gradeService{
		clusterRepo: repository.NewClusterRepository(),
	}
}

func (g gradeService) GetGrade(clusterName string) (*dto.ClusterGrade, error) {
	cluster, err := g.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret"})
	if err != nil {
		return nil, err
	}

	if cluster.Status == constant.StatusRunning {
		result, err := polaris.RunGrade(&polaris.Config{
			Host:  cluster.SpecConf.LbKubeApiserverIp,
			Port:  cluster.SpecConf.KubeApiServerPort,
			Token: cluster.Secret.KubernetesToken,
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		return nil, errors.New("CLUSTER_IS_NOT_RUNNING")
	}
}
