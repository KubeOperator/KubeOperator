package model

import uuid "github.com/satori/go.uuid"

type ClusterSecret struct {
	ID              string
	KubeadmToken    string `gorm:"type:text(65535)" json:"kubeadmToken"`
	KubernetesToken string `gorm:"type:text(65535)" json:"kubernetesToken"`
}

func (n *ClusterSecret) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}
