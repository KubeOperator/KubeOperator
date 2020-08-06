package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type CLusterBackupStrategyService interface {
	Get(clusterName string) (*dto.ClusterBackupStrategy, error)
	Save(creation dto.ClusterBackupStrategyRequest) (*dto.ClusterBackupStrategy, error)
}

type cLusterBackupStrategyService struct {
	clusterBackupStrategyRepo repository.ClusterBackupStrategyRepository
	clusterService            ClusterService
	backupAccountService      BackupAccountService
}

func NewCLusterBackupStrategyService() CLusterBackupStrategyService {
	return &cLusterBackupStrategyService{
		clusterBackupStrategyRepo: repository.NewClusterBackupStrategyRepository(),
		clusterService:            NewClusterService(),
		backupAccountService:      NewBackupAccountService(),
	}
}

func (c cLusterBackupStrategyService) Get(clusterName string) (*dto.ClusterBackupStrategy, error) {
	var clusterBackupStrategyDTO dto.ClusterBackupStrategy
	mo, err := c.clusterBackupStrategyRepo.Get(clusterName)
	if err != nil {
		return nil, err
	}

	clusterBackupStrategyDTO = dto.ClusterBackupStrategy{
		ClusterBackupStrategy: *mo,
		BackupAccountName:     mo.BackupAccount.Name,
		ClusterName:           clusterName,
	}
	return &clusterBackupStrategyDTO, nil
}

func (c cLusterBackupStrategyService) Save(creation dto.ClusterBackupStrategyRequest) (*dto.ClusterBackupStrategy, error) {
	backupAccount, err := c.backupAccountService.Get(creation.BackupAccountName)
	if err != nil {
		return nil, err
	}
	cluster, err := c.clusterService.Get(creation.ClusterName)
	if err != nil {
		return nil, err
	}
	var id string
	old, err := c.Get(creation.ClusterName)
	if err != nil {
		return nil, err
	} else {
		id = old.ID
	}
	clusterBackupStrategy := model.ClusterBackupStrategy{
		ID:              id,
		ClusterID:       cluster.ID,
		Cron:            creation.Cron,
		Status:          creation.Status,
		BackupAccountID: backupAccount.ID,
		SaveNum:         creation.SaveNum,
	}

	err = c.clusterBackupStrategyRepo.Save(&clusterBackupStrategy)
	if err != nil {
		return nil, err
	}
	return &dto.ClusterBackupStrategy{ClusterBackupStrategy: clusterBackupStrategy}, err
}
