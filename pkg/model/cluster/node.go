package cluster

import (
	uuid "github.com/satori/go.uuid"
	"ko3-gin/pkg/model/common"
	"time"
)

type Node struct {
	common.BaseModel
	Cluster   Cluster
	ClusterID string
}

func (n *Node) BeforeCreate() error {
	n.ID = uuid.NewV4().String()
	n.CreatedDate = time.Now()
	n.UpdatedDate = time.Now()
	return nil
}

func (n *Node) BeforeUpdate() error {
	n.UpdatedDate = time.Now()
	return nil
}

func (n Node) TableName() string {
	return "ko_node"
}
