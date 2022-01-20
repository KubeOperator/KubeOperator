package service

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_storage"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
)

const (
	checkFailed            = "CHECK_FAILED"
	backupAccountNameExist = "NAME_EXISTS"
	encryptBackupKeys      = "secretKey, accountKey, password"
)

type BackupAccountService interface {
	GetAfterDecrypt(name string) (*dto.BackupAccount, error)
	List(projectName string) ([]dto.BackupAccount, error)
	Page(num, size int) (page.Page, error)
	Create(creation dto.BackupAccountRequest) (*dto.BackupAccount, error)
	Update(creation dto.BackupAccountRequest) (*dto.BackupAccount, error)
	Batch(op dto.BackupAccountOp) error
	GetBuckets(request dto.CloudStorageRequest) ([]interface{}, error)
	Delete(name string) error
}

type backupAccountService struct {
	backupAccountRepo repository.BackupAccountRepository
}

func NewBackupAccountService() BackupAccountService {
	return &backupAccountService{
		backupAccountRepo: repository.NewBackupAccountRepository(),
	}
}

func (b backupAccountService) GetAfterDecrypt(name string) (*dto.BackupAccount, error) {
	var backupAccountDTO dto.BackupAccount
	mo, err := b.backupAccountRepo.Get(name)
	if err != nil {
		return nil, err
	}
	vars := make(map[string]interface{})
	if err := json.Unmarshal([]byte(mo.Credential), &vars); err != nil {
		return nil, err
	}

	encrypt.VarsDecrypt("ahead", encryptBackupKeys, vars)

	varsAfterHandle, _ := json.Marshal(vars)
	mo.Credential = string(varsAfterHandle)

	backupAccountDTO = dto.BackupAccount{
		BackupAccount: *mo,
	}
	return &backupAccountDTO, nil
}

func (b backupAccountService) List(projectName string) ([]dto.BackupAccount, error) {
	var backupAccountDTOs []dto.BackupAccount

	mos, err := b.backupAccountRepo.List(projectName)
	if err != nil {
		return nil, err
	}

	for _, mo := range mos {
		vars := make(map[string]interface{})
		if err := json.Unmarshal([]byte(mo.Credential), &vars); err != nil {
			return nil, err
		}
		encrypt.DeleteVarsDecrypt("ahead", encryptBackupKeys, vars)

		varsAfterHandle, _ := json.Marshal(vars)
		mo.Credential = string(varsAfterHandle)

		backupAccountDTOs = append(backupAccountDTOs, dto.BackupAccount{BackupAccount: mo})
	}
	return backupAccountDTOs, nil
}

func (b backupAccountService) Page(num, size int) (page.Page, error) {
	var page page.Page
	var backupAccountDTOs []dto.BackupAccount
	total, mos, err := b.backupAccountRepo.Page(num, size)
	if err != nil {
		return page, err
	}
	for _, mo := range mos {
		backupDTO := new(dto.BackupAccount)
		vars := make(map[string]interface{})
		if err := json.Unmarshal([]byte(mo.Credential), &vars); err != nil {
			return page, err
		}

		backupDTO.CredentialVars = encrypt.DeleteVarsDecrypt("ahead", encryptBackupKeys, vars)
		varsAfterHandle, _ := json.Marshal(vars)
		mo.Credential = string(varsAfterHandle)
		backupDTO.BackupAccount = mo

		backupAccountDTOs = append(backupAccountDTOs, *backupDTO)
	}
	page.Total = total
	page.Items = backupAccountDTOs
	return page, err
}

func (b backupAccountService) Create(creation dto.BackupAccountRequest) (*dto.BackupAccount, error) {
	old, _ := b.GetAfterDecrypt(creation.Name)
	if old != nil && old.ID != "" {
		return nil, errors.New(backupAccountNameExist)
	}

	if err := b.CheckValid(creation); err != nil {
		return nil, err
	}

	encrypt.VarsEncrypt("ahead", encryptBackupKeys, creation.CredentialVars)

	credential, _ := json.Marshal(creation.CredentialVars)
	backupAccount := model.BackupAccount{
		Name:       creation.Name,
		Bucket:     creation.Bucket,
		Type:       creation.Type,
		Credential: string(credential),
		Status:     constant.Valid,
	}

	if err := b.backupAccountRepo.Save(&backupAccount); err != nil {
		return nil, err
	}

	return &dto.BackupAccount{BackupAccount: backupAccount}, nil
}

func (b backupAccountService) Update(creation dto.BackupAccountRequest) (*dto.BackupAccount, error) {
	if err := b.CheckValid(creation); err != nil {
		return nil, err
	}

	old, err := b.backupAccountRepo.Get(creation.Name)
	if err != nil {
		return nil, err
	}

	credential, _ := json.Marshal(creation.CredentialVars)
	backupAccount := model.BackupAccount{
		ID:         old.ID,
		Name:       creation.Name,
		Bucket:     creation.Bucket,
		Type:       creation.Type,
		Credential: string(credential),
		Status:     constant.Valid,
	}

	if err = b.backupAccountRepo.Save(&backupAccount); err != nil {
		return nil, err
	}

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
	vars := request.CredentialVars
	encrypt.VarsDecrypt("ahead", encryptBackupKeys, vars)

	vars["type"] = request.Type
	client, err := cloud_storage.NewCloudStorageClient(vars)
	if err != nil {
		return nil, err
	}
	return client.ListBuckets()
}

func (b backupAccountService) CheckValid(create dto.BackupAccountRequest) error {
	vars := create.CredentialVars
	encrypt.VarsDecrypt("ahead", encryptBackupKeys, vars)

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
		return errors.New(checkFailed)
	} else {
		deleteSuccess, err := client.Delete(constant.DefaultFireName)
		if err != nil {
			return err
		}
		if !deleteSuccess {
			return errors.New(checkFailed)
		}
	}
	return nil
}

func (b backupAccountService) Delete(name string) error {
	return b.backupAccountRepo.Delete(name)
}
