package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterBackupStrategy struct {
	model.ClusterBackupStrategy
	ClusterName       string `json:"clusterName"`
	BackupAccountName string `json:"backupAccountName"`
}

type ClusterBackupStrategyRequest struct {
	ID                string `json:"id"`
	Cron              int    `json:"cron"`
	SaveNum           int    `json:"saveNum"`
	BackupAccountName string `json:"backupAccountName"`
	ClusterName       string `json:"clusterName"`
	Status            string `json:"status"`
}
