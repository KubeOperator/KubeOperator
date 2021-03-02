package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type BackupAccountRepository interface {
	Get(name string) (*model.BackupAccount, error)
	List(projectName string) ([]model.BackupAccount, error)
	Save(backupAccount *model.BackupAccount) error
	Page(num, size int) (int, []model.BackupAccount, error)
	Batch(operation string, items []model.BackupAccount) error
	Delete(name string) error
}

type backupAccountRepository struct {
	projectResourceRepository ProjectResourceRepository
}

func NewBackupAccountRepository() BackupAccountRepository {
	return &backupAccountRepository{
		projectResourceRepository: NewProjectResourceRepository(),
	}
}

func (b backupAccountRepository) Get(name string) (*model.BackupAccount, error) {
	var backupAccount model.BackupAccount
	if err := db.DB.Where("name = ?", name).First(&backupAccount).Error; err != nil {
		return nil, err
	}
	return &backupAccount, nil
}

func (b backupAccountRepository) List(projectName string) ([]model.BackupAccount, error) {
	var backupAccounts []model.BackupAccount
	if projectName == "" {
		err := db.DB.Find(&backupAccounts).Error
		if err != nil {
			return nil, err
		}
	} else {
		projectResources, err := b.projectResourceRepository.ListByProjectNameAndType(projectName, constant.ResourceBackupAccount)
		if err != nil {
			return nil, err
		}
		var resourceIds []string
		for _, pr := range projectResources {
			resourceIds = append(resourceIds, pr.ResourceID)
		}
		err = db.DB.Where("id in (?)", resourceIds).Find(&backupAccounts).Error
		if err != nil {
			return nil, err
		}
		return backupAccounts, nil
	}
	return nil, nil
}

func (b backupAccountRepository) Page(num, size int) (int, []model.BackupAccount, error) {
	var total int
	var backupAccounts []model.BackupAccount
	err := db.DB.Model(&model.BackupAccount{}).Count(&total).Offset((num - 1) * size).Limit(size).Find(&backupAccounts).Error
	return total, backupAccounts, err
}

func (b backupAccountRepository) Save(backupAccount *model.BackupAccount) error {
	if db.DB.NewRecord(backupAccount) {
		return db.DB.Create(backupAccount).Error
	} else {
		return db.DB.Save(&backupAccount).Error
	}
}

func (b backupAccountRepository) Batch(operation string, items []model.BackupAccount) error {

	tx := db.DB.Begin()
	switch operation {
	case constant.BatchOperationDelete:
		for i := range items {
			var backupAccount model.BackupAccount
			if err := db.DB.Where("name = ?", items[i].Name).First(&backupAccount).Error; err != nil {
				tx.Rollback()
				return err
			}
			if err := db.DB.Delete(&backupAccount).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	default:
		return constant.NotSupportedBatchOperation
	}
	tx.Commit()
	return nil
}

func (b backupAccountRepository) Delete(name string) error {
	backupAccount, err := b.Get(name)
	if err != nil {
		return err
	}
	return db.DB.Delete(&backupAccount).Error
}
