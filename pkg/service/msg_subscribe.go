package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"reflect"
)

type MsgSubscribeService interface {
	Page(scope, resourceName string, num, size int, condition condition.Conditions) (page.Page, error)
	Update(updated dto.MsgSubscribeDTO) error
}

type msgSubScribeService struct {
	ClusterService ClusterService
}

func NewMsgSubscribeService() MsgSubscribeService {
	return &msgSubScribeService{
		ClusterService: NewClusterService(),
	}
}

func (m msgSubScribeService) Page(scope, resourceName string, num, size int, condition condition.Conditions) (page.Page, error) {
	var arr []dto.MsgSubscribeDTO
	var subscribes []model.MsgSubscribe
	var p page.Page
	d := db.DB.Model(model.MsgSubscribe{})
	if err := dbUtil.WithConditions(&d, model.MsgSubscribe{}, condition); err != nil {
		return p, err
	}
	var resourceID string
	if scope == constant.Cluster {
		cluster, err := m.ClusterService.Get(resourceName)
		if err != nil {
			return p, err
		}
		resourceID = cluster.ID
	}

	if err := d.Where("type = ? AND resource_id = ?", scope, resourceID).Order("CONVERT(name using gbk) asc").Count(&p.Total).Offset((num - 1) * size).Limit(size).Find(&subscribes).Error; err != nil {
		return p, err
	}
	for _, mo := range subscribes {
		msgDTO := dto.MsgSubscribeDTO{}
		msgDTO.CoverToDTO(mo)
		arr = append(arr, msgDTO)
	}
	p.Items = arr

	return p, nil
}

func (m msgSubScribeService) Update(updated dto.MsgSubscribeDTO) error {
	var old model.MsgSubscribe
	if err := db.DB.Model(model.MsgSubscribe{}).Where("name = ? AND type = ? AND resource_id = ?", updated.Name, updated.Type, updated.ResourceID).Find(&old).Error; err != nil {
		return err
	}
	if reflect.DeepEqual(updated.SubConfig, dto.MsgSubConfig{}) {
		return errors.New("can not set config null")
	}
	configB, err := json.Marshal(updated.SubConfig)
	if err != nil {
		return err
	}
	old.Config = string(configB)
	return db.DB.Save(&old).Error
}
