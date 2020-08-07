package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type ClusterBackupFile struct {
	model.ClusterBackupFile
}

type ClusterBackupFileCreate struct {
	ClusterName             string `json:"clusterName"`
	Name                    string `json:"name"`
	ClusterBackupStrategyID string `json:"clusterBackupStrategyId"`
	Folder                  string `json:"folder"`
}

type ClusterBackupFileOp struct {
	Operation string              `json:"operation" validate:"required"`
	Items     []ClusterBackupFile `json:"items" validate:"required"`
}
