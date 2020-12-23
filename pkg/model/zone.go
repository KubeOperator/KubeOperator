package model

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
)

var (
	DeleteZoneError = "DELETE_ZONE_FAILED_RESOURCE"
)

type Zone struct {
	common.BaseModel
	ID           string `json:"id" gorm:"type:varchar(64)"`
	Name         string `json:"name" gorm:"type:varchar(256);not null;unique"`
	Vars         string `json:"vars" gorm:"type: text(65535)"`
	Status       string `json:"status" gorm:"type:varchar(64)"`
	RegionID     string `json:"regionID" gorm:"type:varchar(64)"`
	CredentialID string `json:"credentialId" gorm:"type:varchar(64)"`
	IpPoolID     string `json:"ipPoolId"`
	Region       Region `json:"-"`
	IpPool       IpPool `json:"_"`
}

func (z *Zone) BeforeCreate() (err error) {
	z.ID = uuid.NewV4().String()
	return err
}

func (z *Zone) BeforeDelete() (err error) {
	var planZones []PlanZones
	err = db.DB.Where(PlanZones{ZoneID: z.ID}).Find(&planZones).Error
	if err != nil {
		return err
	}
	if len(planZones) > 0 {
		return errors.New(DeleteZoneError)
	}
	return err
}
