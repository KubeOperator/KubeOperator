package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type CredentialRepository interface {
	Get(name string) (model.Credential, error)
	List() ([]model.Credential, error)
	Page(num, size int) (int, []model.Credential, error)
	Save(credential *model.Credential) error
	Delete(name string) error
	GetById(id string) (model.Credential, error)
}

func NewCredentialRepository() CredentialRepository {
	return &credentialRepository{}
}

type credentialRepository struct {
}

func (c credentialRepository) Get(name string) (model.Credential, error) {
	var credential model.Credential
	credential.Name = name
	if err := db.DB.Where(credential).First(&credential).Error; err != nil {
		return credential, err
	}
	return credential, nil
}

func (c credentialRepository) List() ([]model.Credential, error) {
	var credentials []model.Credential
	err := db.DB.Model(model.Credential{}).Find(&credentials).Error
	return credentials, err
}

func (c credentialRepository) Page(num, size int) (int, []model.Credential, error) {
	var total int
	var credentials []model.Credential
	err := db.DB.Model(model.Credential{}).
		Count(&total).
		Find(&credentials).
		Offset((num - 1) * size).
		Limit(size).
		Error
	return total, credentials, err
}

func (c credentialRepository) Save(credential *model.Credential) error {
	if db.DB.NewRecord(credential) {
		return db.DB.Create(&credential).Error
	} else {
		return db.DB.Save(&credential).Error
	}
}

func (c credentialRepository) Delete(name string) error {
	var credential model.Credential
	credential.Name = name
	return db.DB.Delete(&credential).Error
}

func (c credentialRepository) GetById(id string) (model.Credential, error) {
	var credential model.Credential
	credential.ID = id
	if err := db.DB.Where(credential).First(&credential).Error; err != nil {
		return credential, err
	}
	return credential, nil
}
