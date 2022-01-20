package service

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
)

var (
	hostIsNotNull       = "delete credential error, there are some hosts use this key"
	CredentialNameExist = "NAME_EXISTS"
)

type CredentialService interface {
	Get(name string) (dto.Credential, error)
	List() ([]dto.Credential, error)
	Page(num, size int) (page.Page, error)
	Create(creation dto.CredentialCreate) (*dto.Credential, error)
	Delete(name string) error
	Batch(op dto.CredentialBatchOp) error
	GetById(id string) (dto.Credential, error)
	Update(update dto.CredentialUpdate) (dto.Credential, error)
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
	credentialDTO = dto.Credential{
		ID:        mo.ID,
		Name:      mo.Name,
		Type:      mo.Type,
		Username:  mo.Username,
		CreatedAt: mo.CreatedAt,
	}
	return credentialDTO, err
}

func (c credentialService) GetById(id string) (dto.Credential, error) {
	var credentialDTO dto.Credential
	mo, err := c.credentialRepo.GetById(id)
	if err != nil {
		return credentialDTO, err
	}
	credentialDTO = dto.Credential{
		ID:        mo.ID,
		Name:      mo.Name,
		Type:      mo.Type,
		Username:  mo.Username,
		CreatedAt: mo.CreatedAt,
	}
	return credentialDTO, err
}

func (c credentialService) List() ([]dto.Credential, error) {
	var credentialDTOS []dto.Credential
	mos, err := c.credentialRepo.List()
	if err != nil {
		return credentialDTOS, err
	}
	for _, mo := range mos {
		var credentialDTO = dto.Credential{
			ID:        mo.ID,
			Name:      mo.Name,
			Type:      mo.Type,
			Username:  mo.Username,
			CreatedAt: mo.CreatedAt,
		}
		credentialDTOS = append(credentialDTOS, credentialDTO)
	}
	return credentialDTOS, err
}

func (c credentialService) Page(num, size int) (page.Page, error) {

	var page page.Page
	var total int
	var credentialDTOS []dto.Credential
	total, mos, err := c.credentialRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		var credentailDTO = dto.Credential{
			ID:        mo.ID,
			Name:      mo.Name,
			Type:      mo.Type,
			Username:  mo.Username,
			CreatedAt: mo.CreatedAt,
		}
		credentialDTOS = append(credentialDTOS, credentailDTO)
	}
	page.Total = total
	page.Items = credentialDTOS
	return page, err
}

func (c credentialService) Create(creation dto.CredentialCreate) (*dto.Credential, error) {

	old, _ := c.Get(creation.Name)
	if old.Name != "" {
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
	return &dto.Credential{Name: credential.Name, Username: credential.Username, Type: credential.Type, CreatedAt: credential.CreatedAt}, nil
}

func (c credentialService) Update(update dto.CredentialUpdate) (dto.Credential, error) {

	password, err := encrypt.StringEncrypt(update.Password)
	if err != nil {
		return dto.Credential{}, err
	}

	credential := model.Credential{
		ID:         update.ID,
		Name:       update.Name,
		Password:   password,
		Username:   update.Username,
		PrivateKey: update.PrivateKey,
		Type:       update.Type,
	}
	err = c.credentialRepo.Save(&credential)
	if err != nil {
		return dto.Credential{}, err
	}
	return dto.Credential{Name: credential.Name, Username: credential.Username, Type: credential.Type, CreatedAt: credential.CreatedAt}, nil
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
