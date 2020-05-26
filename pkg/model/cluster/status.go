package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type Status struct {
	commonModel.BaseModel
	ID        string
	ClusterID string
	Version   string
	Message   string
	Phase     string
}

func (s *Status) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return nil
}

func (s Status) TableName() string {
	return "ko_cluster_status"
}
