package model

import uuid "github.com/satori/go.uuid"

type ClusterSecret struct {
	ID              string
	KubeadmToken    string `gorm:"type:text(65535)" json:"kubeadmToken"`
	KubernetesToken string `gorm:"type:text(65535)" json:"kubernetesToken"`
	CertDataStr     string `json:"certDataStr" gorm:"type:text(65535)"`
	KeyDataStr      string `json:"keyDataStr" gorm:"type:text(65535)"`
	ConfigContent   string `json:"configContent" gorm:"type:text(65535)"`
}

func (n *ClusterSecret) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}
