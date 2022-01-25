package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterBackupStrategy struct {
	model.ClusterBackupStrategy
	ClusterName       string `json:"clusterName"`
	BackupAccountName string `json:"backupAccountName"`
}

type ClusterBackupStrategyRequest struct {
	ID                string `json:"id" validate:"-"`
	Cron              int    `json:"cron" validate:"gte=1,lte=300" en:"Backup Interval" zh:"备份间隔"`
	SaveNum           int    `json:"saveNum"  validate:"gte=1,lte=100" en:"Keep Copies" zh:"保留份数"`
	BackupAccountName string `json:"backupAccountName" validate:"required"`
	ClusterName       string `json:"clusterName" validate:"required"`
	Status            string `json:"status" validate:"oneof=ENABLE DISABLE"`
}
