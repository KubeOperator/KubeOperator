package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubepi"
	"github.com/jinzhu/gorm"
)

type KubepiService interface {
	GetKubePiUser() (*kubepi.ListUser, error)
	BindKubePi(req dto.BindKubePI) error
	GetKubePiBind(req dto.SearchBind) (*dto.BindResponse, error)
	CheckConn(req dto.CheckConn) error
	LoadInfo(name string) (*ConnInfo, error)
}

func NewKubepiService() KubepiService {
	return &kubepiService{}
}

type kubepiService struct {
}

type ConnInfo struct {
	Name     string `josn:"name"`
	Password string `josn:"password"`
}

func (c kubepiService) GetKubePiUser() (*kubepi.ListUser, error) {
	var adminBind model.KubepiBind
	if err := db.DB.Where("source_type = ?", constant.SystemRoleAdmin).First(&adminBind).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("NO_KUBEPI_ADMIN")
		}
		return nil, err
	}
	password, err := encrypt.StringDecrypt(adminBind.BindPassword)
	if err != nil {
		return nil, err
	}
	kubepiClient := kubepi.GetClient(kubepi.WithUsernameAndPassword(adminBind.BindUser, password))
	users, err := kubepiClient.SearchUsers()
	if err != nil {
		logger.Log.Errorf("list kubepi users failed, err: %v", err)
		return users, err
	}

	return users, nil
}

func (s *kubepiService) BindKubePi(req dto.BindKubePI) error {
	var record model.KubepiBind
	password, err := encrypt.StringEncrypt(req.BindPassword)
	if err != nil {
		return err
	}
	_ = db.DB.Where("source_type = ? AND project = ? AND cluster = ?", req.SourceType, req.Project, req.Cluster).First(&record).Error
	if record.ID != "" {
		if req.BindUser != record.BindUser || password != record.BindPassword {
			record.BindPassword = password
			record.BindUser = req.BindUser
			return db.DB.Save(&record).Error
		}
		return nil
	}

	bind := &model.KubepiBind{
		SourceType:   req.SourceType,
		Project:      req.Project,
		Cluster:      req.Cluster,
		BindUser:     req.BindUser,
		BindPassword: password,
	}

	return db.DB.Create(bind).Error
}

func (s *kubepiService) GetKubePiBind(req dto.SearchBind) (*dto.BindResponse, error) {
	var record model.KubepiBind
	bind := &dto.BindResponse{
		SourceType: record.SourceType,
		Project:    record.Project,
		Cluster:    record.Cluster,
	}
	if err := db.DB.Where("source_type = ? AND project = ? AND cluster = ?", req.SourceType, req.Project, req.Cluster).First(&record).Error; err != nil {
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

func (s *kubepiService) LoadInfo(name string) (*ConnInfo, error) {
	var bind model.KubepiBind
	if err := db.DB.Where("cluster = ? AND source_type = ?", name, constant.ResourceCluster).First(&bind).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if bind.ID != "" {
		return &ConnInfo{Name: bind.BindUser, Password: bind.BindPassword}, nil
	}
	if err := db.DB.Where("cluster = ? AND source_type = ?", name, constant.ResourceProject).First(&bind).Error; err != nil {
		return nil, err
	}
	return &ConnInfo{Name: bind.BindUser, Password: bind.BindPassword}, nil
}
