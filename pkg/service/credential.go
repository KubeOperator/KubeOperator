package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
)

type CredentialService interface {
	Get(name string) (dto.Credential, error)
	List() ([]dto.Credential, error)
	Page(num, size int) (int, []dto.Credential, error)
	Create(creation dto.Credential) error
	Delete(name string) error
}
type credentialService struct {
	credentialRepo repository.CredentialRepository
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

func (c credentialService) List() ([]dto.Credential, error) {
	var credentialDTOS []dto.Credential
	mos, err := c.credentialRepo.List()
	if err != nil {
		return credentialDTOS, err
	}
	for _, mo := range mos {
		credentialDTOS = append(credentialDTOS, dto.Credential{Credential: mo})
	}
	return credentialDTOS, err
}

func (c credentialService) Page(num, size int) (int, []dto.Credential, error) {
	var total int
	var credentialDTOS []dto.Credential
	total, mos, err := c.credentialRepo.Page(num, size)
	if err != nil {
		return total, credentialDTOS, err
	}
	for _, mo := range mos {
		credentialDTOS = append(credentialDTOS, dto.Credential{Credential: mo})
	}
	return total, credentialDTOS, err
}

func (c credentialService) Create(creation dto.Credential) error {
	credential := model.Credential{
		BaseModel:  common.BaseModel{},
		Name:       creation.Name,
		Password:   creation.Password,
		Username:   creation.Username,
		PrivateKey: creation.PrivateKey,
		Type:       creation.Type,
	}
	err := c.credentialRepo.Save(&credential)
	if err != nil {
		return err
	}
	return nil
}

func (c credentialService) Delete(name string) error {
	err := c.credentialRepo.Delete(name)
	if err != nil {
		return err
	}
	return nil
}
