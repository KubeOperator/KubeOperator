package model

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

var (
	DefaultProjectCanNotDelete = "DEFAULT_PROJECT_CAN_NOT_DELETE"
	ProjectHasClusterError     = "PROJECT_HAS_CLUSTER"
)

type Project struct {
	common.BaseModel
	ID          string    `json:"id" gorm:"type:varchar(64)"`
	Name        string    `json:"name" gorm:"type:varchar(64);not null;unique"`
	Description string    `json:"description" gorm:"type:varchar(128)"`
	Clusters    []Cluster `json:"-"`
}

func (p *Project) BeforeCreate() (err error) {
	p.ID = uuid.NewV4().String()
	return err
}

func (p *Project) BeforeDelete() (err error) {
	if p.Name == constant.DefaultResourceName {
		return errors.New(DefaultProjectCanNotDelete)
	}
	var projectResources []ProjectResource
	err = db.DB.Model(ProjectResource{}).Where(ProjectResource{ProjectID: p.ID, ResourceType: constant.ResourceCluster}).Find(&projectResources).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if len(projectResources) > 0 {
		return errors.New(ProjectHasClusterError)
	}
	var projectMembers []ProjectMember
	err = db.DB.Model(ProjectMember{}).Where(ProjectMember{ProjectID: p.ID}).Find(&projectMembers).Error
	if err != nil {
		return err
	}
	if len(projectMembers) > 0 {
		err := db.DB.Delete(&projectMembers).Error
		if err != nil {
			return err
		}
	}
	if err := db.DB.Where("project = ?", p.Name).Delete(&KubepiBind{}).Error; err != nil {
		return err
	}
	err = db.DB.Model(User{}).Where("current_project_id = ?", p.ID).Updates(&User{CurrentProjectID: ""}).Error
	if err != nil {
		return err
	}
	return err
}
