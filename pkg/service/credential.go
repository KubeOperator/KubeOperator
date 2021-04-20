package service

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
)

var (
	hostIsNotNull       = "delete credential error, there are some hosts use this key"
	CredentialNameExist = "NAME_EXISTS"
)

type CredentialService interface {
	Get(name string) (dto.Credential, error)
	List(conditions condition.Conditions) ([]dto.Credential, error)
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
	Create(creation dto.CredentialCreate) (*dto.Credential, error)
	Delete(name string) error
	Batch(op dto.CredentialBatchOp) error
	GetById(id string) (dto.Credential, error)
	Update(name string, update dto.CredentialUpdate) (*dto.Credential, error)
}

type credentialService struct {
	credentialRepo repository.CredentialRepository
	hostRepo       repository.HostRepository
}

func NewCredentialService() CredentialService {
	return &credentialService{
		credentialRepo: repository.NewCredentialRepository(),
		hostRepo:       repository.NewHostRepository(),
	}
}

func (c credentialService) Get(name string) (dto.Credential, error) {
	var credentialDTO dto.Credential
	mo, err := c.credentialRepo.Get(name)
	if err != nil {
		return credentialDTO, err
	}
	credentialDTO.Credential = mo
	return credentialDTO, err
}

func (c credentialService) GetById(id string) (dto.Credential, error) {
	var credentialDTO dto.Credential
	mo, err := c.credentialRepo.GetById(id)
	if err != nil {
		return credentialDTO, err
	}
	credentialDTO.Credential = mo
	return credentialDTO, err
}

func (c credentialService) List(conditions condition.Conditions) ([]dto.Credential, error) {
	var (
		credentialDTOS []dto.Credential
		mos            []model.Credential
	)
	d := db.DB.Model(model.Credential{})
	if err := dbUtil.WithConditions(&d, model.Credential{}, conditions); err != nil {
		return nil, err
	}
	if err := d.Order("name").
		Find(&mos).Error; err != nil {
		return nil, err
	}
	for _, mo := range mos {
		var credentialDTO dto.Credential
		credentialDTO.Credential = mo
		credentialDTOS = append(credentialDTOS, credentialDTO)
	}
	return credentialDTOS, nil
}

func (c credentialService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {
	var (
		p              page.Page
		credentialDTOS []dto.Credential
		mos            []model.Credential
	)

	d := db.DB.Model(model.Credential{})
	if err := dbUtil.WithConditions(&d, model.Credential{}, conditions); err != nil {
		return nil, err
	}
	if err := d.
		Count(&p.Total).
		Order("name").
		Offset((num - 1) * size).
		Limit(size).
		Find(&mos).Error; err != nil {
		return nil, err
	}
	for _, mo := range mos {
		var credentailDTO dto.Credential
		credentailDTO.Credential = mo
		credentialDTOS = append(credentialDTOS, credentailDTO)
	}
	p.Items = credentialDTOS
	return &p, nil
}

func (c credentialService) Create(creation dto.CredentialCreate) (*dto.Credential, error) {

	old, _ := c.Get(creation.Name)
	if old.ID != "" {
		return nil, errors.New(CredentialNameExist)
	}
	password, err := encrypt.StringEncrypt(creation.Password)
	if err != nil {
		return nil, err
	}

	credential := model.Credential{
		BaseModel:  common.BaseModel{},
		Name:       creation.Name,
		Password:   password,
		Username:   creation.Username,
		PrivateKey: creation.PrivateKey,
		Type:       creation.Type,
	}
	err = c.credentialRepo.Save(&credential)
	if err != nil {
		return nil, err
	}
	return &dto.Credential{Credential: credential}, nil
}

func (c credentialService) Update(name string, update dto.CredentialUpdate) (*dto.Credential, error) {

	credential, err := c.Get(name)
	if err != nil {
		return nil, err
	}
	if update.Username != "" {
		credential.Username = update.Username
	}
	credential.Type = update.Type
	if update.Type == constant.Password {
		if update.Password == "" {
			return nil, errors.New("PASSWORD_CAN_NOT_NULL")
		} else {
			credential.Password, err = encrypt.StringEncrypt(update.Password)
			credential.PrivateKey = ""
			if err != nil {
				return nil, err
			}
		}
	}
	if update.Type == constant.PrivateKey {
		if update.PrivateKey == "" {
			return nil, errors.New("PRIVATE_KEY_CAN_NOT_NULL")
		} else {
			credential.PrivateKey = update.PrivateKey
			credential.Password = ""
		}
	}
	if err := db.DB.Save(&credential).Error; err != nil {
		return nil, err
	}
	return &credential, nil
}

func (c credentialService) Delete(name string) error {

	credential, err := c.Get(name)
	if err != nil {
		return err
	}
	hosts, err := c.hostRepo.ListByCredentialID(credential.ID)
	if err != nil {
		return err
	}
	if len(hosts) > 0 {
		return errors.New(hostIsNotNull)
	}
	err = c.credentialRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}

func (c credentialService) Batch(op dto.CredentialBatchOp) error {
	var deleteItems []model.Credential
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.Credential{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	err := c.credentialRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}
