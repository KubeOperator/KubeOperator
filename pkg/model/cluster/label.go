package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Label struct {
	commonModel.BaseModel
	Value  string
	NodeID string
}

func (l *Label) BeforeCreate() error {
	l.ID = uuid.NewV4().String()
	l.CreatedDate = time.Now()
	l.UpdatedDate = time.Now()
	return nil
}

func (l *Label) BeforeUpdate() error {
	l.UpdatedDate = time.Now()
	return nil
}

func (l Label) TableName() string {
	return "ko_cluster_node_label"
}
