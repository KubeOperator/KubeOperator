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
	UpdateSubscribeUser(msgSubscribeUserDTO dto.MsgSubscribeUserDTO) error
	AddSubscribeUser(msgSubscribeUserDTO dto.MsgSubscribeUserDTO) error
	DeleteSubscribeUser(msgSubscribeUserDTO dto.MsgSubscribeUserDTO) error
	List(scope, resourceName string, condition condition.Conditions) ([]dto.MsgSubscribeDTO, error)
	GetSubscribeUser(cluster, search string, user dto.SessionUser) (dto.AddSubscribeResponse, error)
}

type msgSubScribeService struct {
	ClusterService       ClusterService
	ClusterMemberService ClusterMemberService
}

func NewMsgSubscribeService() MsgSubscribeService {
	return &msgSubScribeService{
		ClusterService:       NewClusterService(),
		ClusterMemberService: NewClusterMemberService(),
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
		var users []model.User
		db.DB.Raw("select * from ko_user where id in (select user_id from ko_msg_subscribe_user where subscribe_id = ?)", mo.ID).Scan(&users)
		msgDTO := dto.NewMsgSubscribeDTO(mo)
		msgDTO.Users = users
		arr = append(arr, msgDTO)
	}
	p.Items = arr

	return p, nil
}

func (m msgSubScribeService) List(scope, resourceName string, condition condition.Conditions) ([]dto.MsgSubscribeDTO, error) {
	var (
		arr        []dto.MsgSubscribeDTO
		subscribes []model.MsgSubscribe
		resourceID string
	)
	d := db.DB.Model(model.MsgSubscribe{})
	if err := dbUtil.WithConditions(&d, model.MsgSubscribe{}, condition); err != nil {
		return arr, err
	}
	if scope == constant.Cluster {
		cluster, err := m.ClusterService.Get(resourceName)
		if err != nil {
			return arr, err
		}
		resourceID = cluster.ID
	}
	if err := d.Where("type = ? AND resource_id = ?", scope, resourceID).Order("CONVERT(name using gbk) asc").Find(&subscribes).Error; err != nil {
		return arr, err
	}
	for _, mo := range subscribes {
		var users []model.User
		db.DB.Raw("select * from ko_user where id in (select user_id from ko_msg_subscribe_user where subscribe_id = ?)", mo.ID).Scan(&users)
		msgDTO := dto.NewMsgSubscribeDTO(mo)
		msgDTO.Users = users
		arr = append(arr, msgDTO)
	}
	return arr, nil
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

func (m msgSubScribeService) UpdateSubscribeUser(msgSubscribeUserDTO dto.MsgSubscribeUserDTO) error {
	tx := db.DB.Begin()
	if err := tx.Where("subscribe_id = ?", msgSubscribeUserDTO.MsgSubscribeID).Delete(&model.MsgSubscribeUser{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, userId := range msgSubscribeUserDTO.Users {
		mo := model.MsgSubscribeUser{
			UserID:      userId,
			SubscribeID: msgSubscribeUserDTO.MsgSubscribeID,
		}
		if err := tx.Create(&mo).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (m msgSubScribeService) AddSubscribeUser(msgSubscribeUserDTO dto.MsgSubscribeUserDTO) error {
	tx := db.DB.Begin()
	for _, userId := range msgSubscribeUserDTO.Users {
		mo := model.MsgSubscribeUser{
			UserID:      userId,
			SubscribeID: msgSubscribeUserDTO.MsgSubscribeID,
		}
		if err := tx.Create(&mo).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m msgSubScribeService) DeleteSubscribeUser(msgSubscribeUserDTO dto.MsgSubscribeUserDTO) error {
	tx := db.DB.Begin()
	for _, userId := range msgSubscribeUserDTO.Users {
		db.DB.Where("subscribe_id = ? AND user_id = ?", msgSubscribeUserDTO.MsgSubscribeID, userId).Delete(model.MsgSubscribeUser{})
	}
	tx.Commit()
	return nil
}

func (m msgSubScribeService) GetSubscribeUser(cluster, search string, user dto.SessionUser) (dto.AddSubscribeResponse, error) {
	var (
		clusterMembers []model.ClusterMember
		admins         []model.User
		userIds        []string
		items          []model.User
		res            dto.AddSubscribeResponse
	)
	clu, err := m.ClusterService.Get(cluster)
	if err != nil {
		return dto.AddSubscribeResponse{}, err
	}
	db.DB.Model(model.ClusterMember{}).Where("cluster_id = ?", clu.ID).Find(&clusterMembers)

	for _, m := range clusterMembers {
		userIds = append(userIds, m.UserID)
	}

	if user.IsAdmin {
		db.DB.Model(model.User{}).Where("is_admin = 1").Find(&admins)
		for _, u := range admins {
			userIds = append(userIds, u.ID)
		}
	}
	db.DB.Where("name LIKE ? AND id in (?)", "%"+search+"%", userIds).Find(&items)
	res.Items = items
	return res, nil
}
