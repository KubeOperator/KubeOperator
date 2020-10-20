package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

type F5Setting struct {
	common.BaseModel
	ID        string `json:"id" gorm:"type:varchar(64)"`
	ClusterID string `json:"cluster_id" gorm:"type:varchar(64)"`
	URL       string `json:"url" gorm:"type:varchar(64)"`
	USER      string `json:"user" gorm:"type:varchar(64)"`
	PASS      string `json:"pass" gorm:"type:varchar(64)""`
	PARTITION string `json:"partition" gorm:"type:varchar(64)" `
	PUBLICIP  string `json:"public_ip" gorm:"type:varchar(64)"`
}

func (s *F5Setting) BeforeCreate() (err error) {
	s.ID = uuid.NewV4().String()
	return err
}

func (s F5Setting) TableName() string {
	return "ko_f5_setting"
}
