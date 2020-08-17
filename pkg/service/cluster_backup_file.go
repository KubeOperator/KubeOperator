package service

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/cloud_storage"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/backup"
	"github.com/jinzhu/gorm"
	"time"
)

type CLusterBackupFileService interface {
	Page(num, size int, clusterName string) (*page.Page, error)
	Create(creation dto.ClusterBackupFileCreate) (*dto.ClusterBackupFile, error)
	Batch(op dto.ClusterBackupFileOp) error
	Backup(creation dto.ClusterBackupFileCreate) error
	Restore(restore dto.ClusterBackupFileRestore) error
	Delete(name string) error
}

type cLusterBackupFileService struct {
	clusterBackupFileRepo           repository.ClusterBackupFileRepository
	clusterService                  ClusterService
	clusterLogService               ClusterLogService
	clusterBackupStrategyRepository repository.ClusterBackupStrategyRepository
	backupAccountRepository         repository.BackupAccountRepository
}

func NewClusterBackupFileService() CLusterBackupFileService {
	return &cLusterBackupFileService{
		clusterBackupFileRepo:           repository.NewClusterBackupFileRepository(),
		clusterService:                  NewClusterService(),
		clusterLogService:               NewClusterLogService(),
		clusterBackupStrategyRepository: repository.NewClusterBackupStrategyRepository(),
		backupAccountRepository:         repository.NewBackupAccountRepository(),
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

func (c cLusterBackupFileService) Delete(name string) error {
	backupFile, err := c.clusterBackupFileRepo.Get(name)
	if err != nil {
		return err
	}
	backupAccount, err := c.backupAccountRepository.Get(backupFile.ClusterBackupStrategy.BackupAccount.Name)
	if err != nil {
		return err
	}
	vars := make(map[string]interface{})
	json.Unmarshal([]byte(backupAccount.Credential), &vars)
	vars["type"] = backupAccount.Type
	vars["bucket"] = backupAccount.Bucket
	client, err := cloud_storage.NewCloudStorageClient(vars)
	if err != nil {
		return err
	}
	result, err := client.Exist(backupFile.Folder)
	if err != nil {
		return err
	}
	if result {
		_, err := client.Delete(backupFile.Folder)
		if err != nil {
			return err
		}
		return c.clusterBackupFileRepo.Delete(name)
	} else {
		return c.clusterBackupFileRepo.Delete(name)
	}
}

func (c cLusterBackupFileService) Backup(creation dto.ClusterBackupFileCreate) error {

	backupLog, err := c.clusterLogService.GetRunningLogWithClusterNameAndType(creation.ClusterName, constant.ClusterLogTypeBackup)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if backupLog.ID != "" {
		return errors.New("CLUSTER_IS_BACKUP")
	}

	restoreLog, err := c.clusterLogService.GetRunningLogWithClusterNameAndType(creation.ClusterName, constant.ClusterLogTypeRestore)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if restoreLog.ID != "" {
		return errors.New("CLUSTER_IS_RESTORE")
	}

	cluster, err := c.clusterService.Get(creation.ClusterName)
	if err != nil {
		return err
	}
	now := time.Now()
	day := now.Format("2006-01-02-15-04")
	fileName := cluster.Name + "-" + day + ".backup.tar.gz"
	creation.Name = fileName
	creation.Folder = cluster.Name + "/" + fileName
	go c.doBackup(cluster.Cluster, creation)
	return nil
}

func (c cLusterBackupFileService) doBackup(cluster model.Cluster, creation dto.ClusterBackupFileCreate) {

	var clog model.ClusterLog
	clog.Type = constant.ClusterLogTypeBackup
	clog.StartTime = time.Now()
	clog.EndTime = time.Now()
	err := c.clusterLogService.Save(cluster.Name, &clog)
	if err != nil {
		log.Error(err)
	}
	err = c.clusterLogService.Start(&clog)
	if err != nil {
		log.Error(err)
	}
	admCluster := adm.NewCluster(cluster)
	p := &backup.BackupClusterPhase{}
	err = p.Run(admCluster.Kobe)
	if err != nil {
		_ = c.clusterLogService.End(&clog, false, err.Error())
	} else {
		clusterBackupStrategy, err := c.clusterBackupStrategyRepository.Get(cluster.Name)
		if err != nil {
			_ = c.clusterLogService.End(&clog, false, err.Error())
			log.Error(err)
			return
		}
		backupAccount, err := c.backupAccountRepository.Get(clusterBackupStrategy.BackupAccount.Name)
		if err != nil {
			_ = c.clusterLogService.End(&clog, false, err.Error())
			log.Error(err)
			return
		}
		vars := make(map[string]interface{})
		json.Unmarshal([]byte(backupAccount.Credential), &vars)
		vars["type"] = backupAccount.Type
		vars["bucket"] = backupAccount.Bucket
		client, err := cloud_storage.NewCloudStorageClient(vars)
		if err != nil {
			_ = c.clusterLogService.End(&clog, false, err.Error())
			log.Error(err)
			return
		}
		srcFilePath := constant.BackupDir + "/" + cluster.Name + "/" + constant.BackupFileDefaultName
		_, err = client.Upload(srcFilePath, creation.Folder)
		if err != nil {
			_ = c.clusterLogService.End(&clog, false, err.Error())
			log.Error(err)
			return
		}
		_ = c.clusterLogService.End(&clog, true, "")
		_, err = c.Create(creation)
		if err != nil {
			_ = c.clusterLogService.End(&clog, false, err.Error())
			log.Error(err)
			return
		}
	}

}

func (c cLusterBackupFileService) Restore(restore dto.ClusterBackupFileRestore) error {

	file, err := c.clusterBackupFileRepo.Get(restore.Name)
	if err != nil {
		return err
	}
	restore.File = file
	clusterBackupStrategy, err := c.clusterBackupStrategyRepository.Get(restore.ClusterName)
	if err != nil {
		return err
	}
	backupAccount, err := c.backupAccountRepository.Get(clusterBackupStrategy.BackupAccount.Name)
	if err != nil {
		return err
	}
	restore.BackupAccount = *backupAccount
	go c.doRestore(restore)
	return nil
}

func (c cLusterBackupFileService) doRestore(restore dto.ClusterBackupFileRestore) {

	cluster, err := c.clusterService.Get(restore.ClusterName)
	if err != nil {
		return
	}

	var clog model.ClusterLog
	clog.Type = constant.ClusterLogTypeRestore
	clog.StartTime = time.Now()
	clog.EndTime = time.Now()
	err = c.clusterLogService.Save(cluster.Name, &clog)
	if err != nil {
		log.Error(err)
	}
	err = c.clusterLogService.Start(&clog)
	if err != nil {
		log.Error(err)
	}

	vars := make(map[string]interface{})
	json.Unmarshal([]byte(restore.BackupAccount.Credential), &vars)
	vars["type"] = restore.BackupAccount.Type
	vars["bucket"] = restore.BackupAccount.Bucket
	client, err := cloud_storage.NewCloudStorageClient(vars)
	if err != nil {
		_ = c.clusterLogService.End(&clog, false, err.Error())
		log.Error(err)
		return
	}

	srcFilePath := restore.File.Folder
	targetPath := constant.BackupDir + "/" + cluster.Name + "/" + constant.BackupFileDefaultName
	_, err = client.Download(srcFilePath, targetPath)
	if err != nil {
		log.Error(err)
		return
	}

	admCluster := adm.NewCluster(cluster.Cluster)
	p := &backup.RestoreClusterPhase{}
	err = p.Run(admCluster.Kobe)
	if err != nil {
		_ = c.clusterLogService.End(&clog, false, err.Error())
	} else {
		_ = c.clusterLogService.End(&clog, true, "")
	}
}
