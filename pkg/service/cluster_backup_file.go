package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
)

type CLusterBackupFileService interface {
	Page(num, size int, clusterName string) (*page.Page, error)
	Create(creation dto.ClusterBackupFileCreate) (*dto.ClusterBackupFile, error)
	Batch(op dto.ClusterBackupFileOp) error
}

type cLusterBackupFileService struct {
	clusterBackupFileRepo repository.ClusterBackupFileRepository
	clusterService        ClusterService
}

func NewClusterBackupFileService() CLusterBackupFileService {
	return &cLusterBackupFileService{
		clusterBackupFileRepo: repository.NewClusterBackupFileRepository(),
		clusterService:        NewClusterService(),
	}
}

func (c cLusterBackupFileService) Page(num, size int, clusterName string) (*page.Page, error) {

	cluster, err := c.clusterService.Get(clusterName)
	if err != nil {
		return nil, err
	}

	var page page.Page
	var fileDTOs []dto.ClusterBackupFile
	total, mos, err := c.clusterBackupFileRepo.Page(num, size, cluster.ID)
	if err != nil {
		return nil, err
	}
	for _, mo := range mos {
		fileDTO := new(dto.ClusterBackupFile)
		fileDTO.ClusterBackupFile = mo
		fileDTOs = append(fileDTOs, *fileDTO)
	}
	page.Total = total
	page.Items = fileDTOs
	return &page, err
}

func (c cLusterBackupFileService) Create(creation dto.ClusterBackupFileCreate) (*dto.ClusterBackupFile, error) {

	var cluster dto.Cluster
	cluster, err := c.clusterService.Get(creation.ClusterName)
	if err != nil {
		return nil, err
	}

	file := model.ClusterBackupFile{
		Name:                    creation.Name,
		ClusterBackupStrategyID: creation.ClusterBackupStrategyID,
		Folder:                  creation.Folder,
		ClusterID:               cluster.ID,
	}

	err = c.clusterBackupFileRepo.Save(&file)
	if err != nil {
		return nil, err
	}

	return &dto.ClusterBackupFile{ClusterBackupFile: file}, err
}

func (c cLusterBackupFileService) Batch(op dto.ClusterBackupFileOp) error {
	var deleteItems []model.ClusterBackupFile
	for _, item := range op.Items {
		deleteItems = append(deleteItems, model.ClusterBackupFile{
			BaseModel: common.BaseModel{},
			ID:        item.ID,
			Name:      item.Name,
		})
	}
	err := c.clusterBackupFileRepo.Batch(op.Operation, deleteItems)
	if err != nil {
		return err
	}
	return nil
}
