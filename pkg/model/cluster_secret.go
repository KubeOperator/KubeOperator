package model

import uuid "github.com/satori/go.uuid"

type ClusterSecret struct {
	ID              string
	KubeadmToken    string `gorm:"type:text(65535)"`
	KubernetesToken string `gorm:"type:text(65535)"`
}

func (n *ClusterSecret) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}

func (s ClusterSecret) TableName() string {
	return "ko_cluster_secret"
}
