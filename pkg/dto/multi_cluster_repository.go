package dto

import "github.com/KubeOperator/KubeOperator/pkg/model"

type MultiClusterRepository struct {
	model.MultiClusterRepository
}

type MultiClusterRepositoryCreateRequest struct {
	Name     string `json:"name"`
	Source   string `json:"source"`
	Branch   string `json:"branch"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MultiClusterRepositoryUpdateRequest struct {
	GitTimeout   int64 `json:"gitTimeout"`
	SyncInterval int64 `json:"syncInterval"`
	SyncEnable   bool  `json:"syncEnable"`
}

type UpdateRelationRequest struct {
	ClusterNames []string `json:"clusterNames"`
	Delete       bool     `json:"delete"`
}

type ClusterRelation struct {
	model.ClusterMultiClusterRepository
	ClusterName string `json:"clusterName"`
}

type MultiClusterSyncClusterLog struct {
	model.MultiClusterSyncClusterLog
	MultiClusterSyncClusterResourceLogs []model.MultiClusterSyncClusterResourceLog `json:"multiClusterSyncClusterResourceLogs"`
	ClusterName                         string                                     `json:"clusterName"`
}

type MultiClusterSyncLogDetail struct {
	model.MultiClusterSyncLog
	MultiClusterSyncClusterLogs []MultiClusterSyncClusterLog `json:"multiClusterSyncClusterLogs"`
}

type MultiClusterSyncLog struct {
	model.MultiClusterSyncLog
}
type MultiClusterRepositoryBatch struct {
	Operation string                   `json:"operation" validate:"required"`
	Items     []MultiClusterRepository `json:"items" validate:"required"`
}
