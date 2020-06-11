package tool

import "github.com/KubeOperator/KubeOperator/pkg/model/common"

type Tool struct {
	common.BaseModel
	ID        string
	Name      string
	Vars      string `gorm:"type:text(65535)"`
	ClusterID string
	StatusID  string
	Status    Status `gorm:"save_associations:false"`
}
