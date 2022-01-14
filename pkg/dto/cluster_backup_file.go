package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterBackupFile struct {
	model.ClusterBackupFile
}

type ClusterBackupFileCreate struct {
	ClusterName             string `json:"clusterName" validate:"clustername,required"`
	Name                    string `json:"name" validate:"required"`
	ClusterBackupStrategyID string `json:"clusterBackupStrategyId" validate:"required"`
	Folder                  string `json:"folder" validate:"required"`
}

type ClusterBackupFileOp struct {
	Operation string              `json:"operation" validate:"required"`
	Items     []ClusterBackupFile `json:"items" validate:"required"`
}

type ClusterBackupFileRestore struct {
	ClusterName   string                  `json:"clusterName" validate:"clustername,required"`
	Name          string                  `json:"name" validate:"required"`
	File          model.ClusterBackupFile `json:"file" validate:"-"`
	BackupAccount model.BackupAccount     `json:"backupAccount" validate:"-"`
}
