package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubepi"
	"github.com/jinzhu/gorm"
)

type KubepiService interface {
	GetKubePiUser() (*kubepi.ListUser, error)
	BindKubePi(req dto.BindKubePI) error
	GetKubePiBind(req dto.SearchBind) (*dto.BindResponse, error)
	CheckConn(req dto.CheckConn) error
}

func NewKubepiService() KubepiService {
	return &kubepiService{}
}

type kubepiService struct {
}

func (c kubepiService) GetKubePiUser() (*kubepi.ListUser, error) {
	var adminBind model.KubepiBind
	if err := db.DB.Where("source_type = ?", constant.SystemRoleAdmin).First(&adminBind).Error; err != nil {
		return nil, err
	}
	kubepiClient := kubepi.GetClient(kubepi.WithUsernameAndPassword(adminBind.BindUser, adminBind.BindPassword))
	users, err := kubepiClient.SearchUsers()
	if err != nil {
		logger.Log.Errorf("list kubepi users failed, err: %v", err)
		return users, err
	}

	return users, nil
}

func (s *kubepiService) BindKubePi(req dto.BindKubePI) error {
	var record model.KubepiBind
	_ = db.DB.Where("source_type = ? && source = ?", req.SourceType, req.Source).First(&record).Error
	if record.ID != "" && (req.BindUser != record.BindUser || req.BindPassword != record.BindPassword) {
		record.BindPassword = req.BindPassword
		record.BindUser = req.BindUser
		return db.DB.Save(&record).Error
	}
	bind := &model.KubepiBind{
		SourceType:   req.SourceType,
		Source:       req.Source,
		BindUser:     req.BindUser,
		BindPassword: req.BindPassword,
	}

	return db.DB.Create(bind).Error
}

func (s *kubepiService) GetKubePiBind(req dto.SearchBind) (*dto.BindResponse, error) {
	var record model.KubepiBind
	bind := &dto.BindResponse{
		SourceType: record.SourceType,
		Source:     record.Source,
	}
	if err := db.DB.Where("source_type = ? && source = ?", req.SourceType, req.Source).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bind, nil
		}
		return bind, err
	}
	bind.BindUser = record.BindUser
	return bind, nil
}

func (s *kubepiService) CheckConn(req dto.CheckConn) error {
	kubepiClient := kubepi.GetClient(kubepi.WithUsernameAndPassword(req.BindUser, req.BindPassword))
	return kubepiClient.CheckLogin()
}
