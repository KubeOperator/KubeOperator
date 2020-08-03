package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type BackupAccountService interface {
	Get(name string) (*dto.BackupAccount, error)
	List(projectName string) ([]dto.BackupAccount, error)
	Page(num, size int) (page.Page, error)
	Create(creation dto.BackupAccountCreate) (*dto.BackupAccount, error)
	Update(creation dto.BackupAccountUpdate) (*dto.BackupAccount, error)
	Batch(op dto.BackupAccountOp) error
}

type backupAccountService struct {
	backupAccountRepo repository.BackupAccountRepository
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
		json.Unmarshal([]byte(mo.Credential), &vars)
		backupDTO.CredentialVars = vars
		backupDTO.BackupAccount = mo

		backupAccountDTOs = append(backupAccountDTOs, *backupDTO)
	}
	page.Total = total
	page.Items = backupAccountDTOs
	return page, err
}

func (b backupAccountService) Create(creation dto.BackupAccountCreate) (*dto.BackupAccount, error) {

	credential, _ := json.Marshal(creation.CredentialVars)
	backupAccount := model.BackupAccount{
		Name:       creation.Name,
		Region:     creation.Region,
		Type:       creation.Type,
		Credential: string(credential),
		Status:     constant.Valid,
	}

	err := b.backupAccountRepo.Save(&backupAccount)
	if err != nil {
		return nil, err
	}

	return &dto.BackupAccount{BackupAccount: backupAccount}, err
}

func (b backupAccountService) Update(creation dto.BackupAccountUpdate) (*dto.BackupAccount, error) {

	credential, _ := json.Marshal(creation.CredentialVars)
	old, err := b.backupAccountRepo.Get(creation.Name)
	if err != nil {
		return nil, err
	}
	backupAccount := model.BackupAccount{
		ID:         old.ID,
		Name:       creation.Name,
		Region:     creation.Region,
		Type:       creation.Type,
		Credential: string(credential),
		Status:     constant.Valid,
	}

	err = b.backupAccountRepo.Save(&backupAccount)
	if err != nil {
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
