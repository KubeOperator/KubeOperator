package service

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_storage"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/condition"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	dbUtil "github.com/KubeOperator/KubeOperator/pkg/util/db"
)

var (
	CheckFailed            = "CHECK_FAILED"
	BackupAccountNameExist = "NAME_EXISTS"
)

type BackupAccountService interface {
	Get(name string) (*dto.BackupAccount, error)
	List(projectName string, conditions condition.Conditions) ([]dto.BackupAccount, error)
	Page(num, size int, conditions condition.Conditions) (*page.Page, error)
	Create(creation dto.BackupAccountRequest) (*dto.BackupAccount, error)
	Update(name string, creation dto.BackupAccountUpdate) (*dto.BackupAccount, error)
	Batch(op dto.BackupAccountOp) error
	GetBuckets(request dto.CloudStorageRequest) ([]interface{}, error)
	Delete(name string) error
	ListByClusterName(clusterName string) ([]dto.BackupAccount, error)
}

type backupAccountService struct {
	projectResourceRepository repository.ProjectResourceRepository
	backupAccountRepo         repository.BackupAccountRepository
}

func NewBackupAccountService() BackupAccountService {
	return &backupAccountService{
		backupAccountRepo: repository.NewBackupAccountRepository(),
	}
}

func (b backupAccountService) Get(name string) (*dto.BackupAccount, error) {
	var backupAccountDTO dto.BackupAccount
	mo, err := b.backupAccountRepo.Get(name)
	if err != nil {
		return nil, err
	}
	vars := make(map[string]interface{})
	if err := json.Unmarshal([]byte(mo.Credential), &vars); err != nil {
		return nil, err
	}
	backupAccountDTO = dto.BackupAccount{
		CredentialVars: vars,
		BackupAccount:  *mo,
	}
	return &backupAccountDTO, nil
}

func (b backupAccountService) List(projectName string, conditions condition.Conditions) ([]dto.BackupAccount, error) {
	var (
		backupAccountDTO []dto.BackupAccount
		mos              []model.BackupAccount
	)
	d := db.DB.Model(model.BackupAccount{})
	if err := dbUtil.WithConditions(&d, model.BackupAccount{}, conditions); err != nil {
		return nil, err
	}
	if projectName == "" {
		if err := d.Order("name").
			Find(&mos).Error; err != nil {
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
		if err := d.Order("name").
			Where("id in (?)", resourceIds).
			Find(mos).Error; err != nil {
			return nil, err
		}
	}
	for _, mo := range mos {
		backupAccountDTO = append(backupAccountDTO, dto.BackupAccount{
			BackupAccount: mo,
		})
	}
	return backupAccountDTO, nil
}

func (b backupAccountService) ListByClusterName(clusterName string) ([]dto.BackupAccount, error) {
	var (
		backupAccountDTOs []dto.BackupAccount
		backupAccounts    []model.BackupAccount
	)
	err := db.DB.Raw("SELECT * FROM ko_backup_account WHERE id IN (SELECT resource_id FROM ko_cluster_resource WHERE  resource_type = 'BACKUP_ACCOUNT' AND cluster_id = (SELECT DISTINCT id FROM ko_cluster WHERE `name` = ?) )", clusterName).Scan(&backupAccounts).Error
	if err != nil {
		return nil, err
	}
	for _, mo := range backupAccounts {
		backupAccountDTOs = append(backupAccountDTOs, dto.BackupAccount{BackupAccount: mo})
	}
	return backupAccountDTOs, nil
}

func (b backupAccountService) Page(num, size int, conditions condition.Conditions) (*page.Page, error) {
	var (
		p                 page.Page
		backupAccountDTOs []dto.BackupAccount
		mos               []model.BackupAccount
		projectResources  []model.ProjectResource
		clusterResources  []model.ClusterResource
	)

	d := db.DB.Model(model.BackupAccount{})
	if err := dbUtil.WithConditions(&d, model.BackupAccount{}, conditions); err != nil {
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
		vars := make(map[string]interface{})
		if err := json.Unmarshal([]byte(mo.Credential), &vars); err != nil {
			return &p, err
		}
		if err := db.DB.Where("resource_id = ?", mo.ID).Preload("Project").Find(&projectResources).Error; err != nil {
			return nil, err
		}
		if err := db.DB.Where("resource_id = ?", mo.ID).Preload("Cluster").Find(&clusterResources).Error; err != nil {
			return nil, err
		}

		var projects string
		for _, pr := range projectResources {
			projects += (pr.Project.Name + ",")
		}
		var clusters string
		for _, cr := range clusterResources {
			clusters += (cr.Cluster.Name + ",")
		}
		backupDTO := dto.BackupAccount{
			CredentialVars: vars,
			BackupAccount:  mo,
			Projects:       projects,
			Clusters:       clusters,
		}
		backupAccountDTOs = append(backupAccountDTOs, backupDTO)
	}
	p.Items = backupAccountDTOs
	return &p, nil
}

func (b backupAccountService) Create(creation dto.BackupAccountRequest) (*dto.BackupAccount, error) {
	old, _ := b.Get(creation.Name)
	if old != nil && old.ID != "" {
		return nil, errors.New(BackupAccountNameExist)
	}

	err := b.CheckValid(creation)
	if err != nil {
		return nil, err
	}
	tx := db.DB.Begin()
	credential, _ := json.Marshal(creation.CredentialVars)
	backupAccount := model.BackupAccount{
		Name:       creation.Name,
		Bucket:     creation.Bucket,
		Type:       creation.Type,
		Credential: string(credential),
		Status:     constant.Valid,
	}

	err = tx.Create(&backupAccount).Error
	if err != nil {
		return nil, err
	}

	var backAccount []model.BackupAccount
	err = tx.Where("name in (?)", creation.Name).Find(&backAccount).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var projects []model.Project
	err = tx.Where("name in (?)", creation.Projects).Find(&projects).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, project := range projects {
		err = tx.Create(&model.ProjectResource{
			ResourceType: constant.ResourceBackupAccount,
			ResourceID:   backupAccount.ID,
			ProjectID:    project.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(creation.Clusters) != 0 {
		var clusters []model.Cluster
		err = tx.Where("name in (?)", creation.Clusters).Find(&clusters).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, cluster := range clusters {
			err = tx.Create(&model.ClusterResource{
				ResourceType: constant.ResourceBackupAccount,
				ResourceID:   backupAccount.ID,
				ClusterID:    cluster.ID,
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	tx.Commit()
	return &dto.BackupAccount{BackupAccount: backupAccount}, err
}

func (b backupAccountService) Update(name string, update dto.BackupAccountUpdate) (*dto.BackupAccount, error) {
	err := b.CheckValid(dto.BackupAccountRequest{
		Name:           update.Name,
		Type:           update.Type,
		CredentialVars: update.CredentialVars,
		Bucket:         update.Bucket,
	})
	if err != nil {
		return nil, err
	}
	var projects []model.Project
	tx := db.DB.Begin()
	err = tx.Where("name in (?)", update.Projects).Find(&projects).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Model(model.ProjectResource{}).Where("resource_id = ?", update.ID).Delete(&model.ProjectResource{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, project := range projects {
		err = tx.Create(&model.ProjectResource{
			ResourceType: constant.ResourceBackupAccount,
			ResourceID:   update.ID,
			ProjectID:    project.ID,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	credential, _ := json.Marshal(update.CredentialVars)
	backupAccount := model.BackupAccount{
		ID:         update.ID,
		Name:       name,
		Bucket:     update.Bucket,
		Type:       update.Type,
		Credential: string(credential),
		Status:     constant.Valid,
	}
	if err := tx.Save(&backupAccount).Error; err != nil {
		return nil, err
	}
	tx.Commit()
	return &dto.BackupAccount{BackupAccount: backupAccount}, err
}

func (b backupAccountService) Batch(op dto.BackupAccountOp) error {
	var deleteItems []model.BackupAccount
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.BackupAccount{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	err := b.backupAccountRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}

func (b backupAccountService) GetBuckets(request dto.CloudStorageRequest) ([]interface{}, error) {
	vars := request.CredentialVars.(map[string]interface{})
	vars["type"] = request.Type
	client, err := cloud_storage.NewCloudStorageClient(vars)
	if err != nil {
		return nil, err
	}
	return client.ListBuckets()
}

func (b backupAccountService) CheckValid(create dto.BackupAccountRequest) error {
	vars := create.CredentialVars.(map[string]interface{})
	vars["type"] = create.Type
	vars["bucket"] = create.Bucket
	client, err := cloud_storage.NewCloudStorageClient(vars)
	if err != nil {
		return err
	}
	file, err := os.Create(constant.DefaultFireName)
	if err != nil {
		return err
	}
	defer file.Close()
	success, err := client.Upload(constant.DefaultFireName, constant.DefaultFireName)
	if err != nil {
		return err
	}
	if !success {
		return errors.New(CheckFailed)
	} else {
		deleteSuccess, err := client.Delete(constant.DefaultFireName)
		if err != nil {
			return err
		}
		if !deleteSuccess {
			return errors.New(CheckFailed)
		}
	}
	return nil
}

func (b backupAccountService) Delete(name string) error {
	return b.backupAccountRepo.Delete(name)
}
