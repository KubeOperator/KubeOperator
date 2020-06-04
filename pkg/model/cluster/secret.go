package cluster
import uuid "github.com/satori/go.uuid"

type Secret struct {
	ID              string
	KubeadmToken    string `gorm:"type:text(65535)"`
	KubernetesToken string `gorm:"type:text(65535)"`
}

func (n *Secret) BeforeCreate() (err error) {
	n.ID = uuid.NewV4().String()
	return nil
}

func (s Secret) TableName() string {
	return "ko_cluster_secret"
}
