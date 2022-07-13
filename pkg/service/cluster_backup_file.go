package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"

	"github.com/KubeOperator/KubeOperator/pkg/cloud_storage"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/page"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases/backup"
)

type CLusterBackupFileService interface {
	Page(num, size int, clusterName string) (*page.Page, error)
	Create(creation dto.ClusterBackupFileCreate) (*dto.ClusterBackupFile, error)
	Batch(op dto.ClusterBackupFileOp) error
	Backup(creation dto.ClusterBackupFileCreate) error
	Restore(restore dto.ClusterBackupFileRestore) error
	Delete(name string) error
	LocalRestore(clusterName string, file []byte) error
}

type cLusterBackupFileService struct {
	clusterBackupFileRepo           repository.ClusterBackupFileRepository
	clusterService                  ClusterService
	clusterRepo                     repository.ClusterRepository
	taskLogService                  TaskLogService
	clusterBackupStrategyRepository repository.ClusterBackupStrategyRepository
	backupAccountRepository         repository.BackupAccountRepository
	messageService                  MessageService
}

func NewClusterBackupFileService() CLusterBackupFileService {
	return &cLusterBackupFileService{
		clusterBackupFileRepo:           repository.NewClusterBackupFileRepository(),
		clusterService:                  NewClusterService(),
		clusterRepo:                     repository.NewClusterRepository(),
		taskLogService:                  NewTaskLogService(),
		clusterBackupStrategyRepository: repository.NewClusterBackupStrategyRepository(),
		backupAccountRepository:         repository.NewBackupAccountRepository(),
		messageService:                  NewMessageService(),
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
	if err := json.Unmarshal([]byte(backupAccount.Credential), &vars); err != nil {
		return err
	}
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
	isON := c.taskLogService.IsTaskOn(creation.ClusterName)
	if isON {
		return errors.New("TASK_IN_EXECUTION")
	}
	cluster, err := c.clusterRepo.GetWithPreload(creation.ClusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}

	now := time.Now()
	day := now.Format("2006-01-02-15-04")
	fileName := cluster.Name + "-" + day + ".backup.db"
	creation.Name = fileName
	creation.Folder = cluster.Name + "/" + fileName

	task, err := c.taskLogService.NewTerminalTask(cluster.ID, constant.TaskLogTypeBackup)
	if err != nil {
		return err
	}
	cluster.TaskLog = *task
	cluster.CurrentTaskID = task.ID
	_ = c.clusterRepo.Save(&cluster)

	go c.doBackup(cluster, creation, task)
	return nil
}

func (c cLusterBackupFileService) doBackup(cluster model.Cluster, creation dto.ClusterBackupFileCreate, task *model.TaskLog) {
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, task.ID)
	if err != nil {
		logger.Log.Errorf("create ansible log failed, error: %s", err.Error())
	}
	admCluster := adm.NewAnsibleHelper(cluster)
	p := &backup.BackupClusterPhase{}
	err = p.Run(admCluster.Kobe, writer)
	if err != nil {
		logger.Log.Errorf("run cluster log failed, error: %s", err.Error())
		_ = c.taskLogService.End(task, false, err.Error())
		cluster.CurrentTaskID = ""
		_ = c.clusterRepo.Save(&cluster)
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterBackup, false, err.Error()), cluster.Name, constant.ClusterBackup)
	} else {
		_ = c.taskLogService.End(task, true, "")
		cluster.CurrentTaskID = ""
		_ = c.clusterRepo.Save(&cluster)
		clusterBackupStrategy, err := c.clusterBackupStrategyRepository.Get(cluster.Name)
		if err != nil {
			logger.Log.Errorf("get backup strategy failed, error: %s", err.Error())
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterBackup, false, err.Error()), cluster.Name, constant.ClusterBackup)
			return
		}
		backupAccount, err := c.backupAccountRepository.Get(clusterBackupStrategy.BackupAccount.Name)
		if err != nil {
			logger.Log.Errorf("get backup account failed, error: %s", err.Error())
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterBackup, false, err.Error()), cluster.Name, constant.ClusterBackup)
			return
		}
		vars := make(map[string]interface{})
		if err := json.Unmarshal([]byte(backupAccount.Credential), &vars); err != nil {
			logger.Log.Errorf("backup account credential json.Unmarshal failed, error: %s", err.Error())
		}
		vars["type"] = backupAccount.Type
		vars["bucket"] = backupAccount.Bucket
		client, err := cloud_storage.NewCloudStorageClient(vars)
		if err != nil {
			logger.Log.Errorf("cloud storage new client failed, error: %s", err.Error())
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterBackup, false, err.Error()), cluster.Name, constant.ClusterBackup)
			return
		}
		srcFilePath := constant.BackupDir + "/" + cluster.Name + "/" + constant.BackupFileDefaultName
		_, err = client.Upload(srcFilePath, creation.Folder)
		if err != nil {
			logger.Log.Errorf("backup file upload failed, error: %s", err.Error())
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterBackup, false, err.Error()), cluster.Name, constant.ClusterBackup)
			return
		}
		creation.ClusterBackupStrategyID = clusterBackupStrategy.ID
		_, err = c.Create(creation)
		if err != nil {
			logger.Log.Errorf("backup file create failed, error: %s", err.Error())
			_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterBackup, false, err.Error()), cluster.Name, constant.ClusterBackup)
			return
		} else {
			go c.deleteBackupFile(cluster.Name)
			_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterBackup, true, ""), cluster.Name, constant.ClusterBackup)
		}
	}
}

func (c cLusterBackupFileService) Restore(restore dto.ClusterBackupFileRestore) error {
	isON := c.taskLogService.IsTaskOn(restore.ClusterName)
	if isON {
		return errors.New("TASK_IN_EXECUTION")
	}
	cluster, err := c.clusterRepo.GetWithPreload(restore.ClusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}

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

	task, err := c.taskLogService.NewTerminalTask(cluster.ID, constant.TaskLogTypeRestore)
	if err != nil {
		return err
	}
	cluster.TaskLog = *task
	cluster.CurrentTaskID = task.ID
	_ = c.clusterRepo.Save(&cluster)

	go c.doRestore(restore, cluster, task)
	return nil
}

func (c cLusterBackupFileService) doRestore(restore dto.ClusterBackupFileRestore, cluster model.Cluster, task *model.TaskLog) {
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, task.ID)
	if err != nil {
		logger.Log.Errorf("create ansible log failed, error: %s", err.Error())
	}

	vars := make(map[string]interface{})
	if err := json.Unmarshal([]byte(restore.BackupAccount.Credential), &vars); err != nil {
		logger.Log.Errorf("doRestore json.Unmarshal failed,  error: %s", err.Error())
	}
	vars["type"] = restore.BackupAccount.Type
	vars["bucket"] = restore.BackupAccount.Bucket
	client, err := cloud_storage.NewCloudStorageClient(vars)
	if err != nil {
		_ = c.taskLogService.End(task, false, err.Error())
		cluster.CurrentTaskID = ""
		_ = c.clusterRepo.Save(&cluster)
		logger.Log.Errorf("cloud storage new client failed, error: %s", err.Error())
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterRestore, false, err.Error()), cluster.Name, constant.ClusterRestore)
		return
	}

	srcFilePath := restore.File.Folder
	targetPath := constant.BackupDir + "/" + cluster.Name + "/" + constant.BackupFileDefaultName
	_, err = client.Download(srcFilePath, targetPath)
	if err != nil {
		_ = c.taskLogService.End(task, false, err.Error())
		cluster.CurrentTaskID = ""
		_ = c.clusterRepo.Save(&cluster)
		logger.Log.Errorf("cloud storage download failed, error: %s", err.Error())
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterRestore, false, err.Error()), cluster.Name, constant.ClusterRestore)
		return
	}

	admCluster := adm.NewAnsibleHelper(cluster)
	p := &backup.RestoreClusterPhase{}

	err = p.Run(admCluster.Kobe, writer)
	if err != nil {
		logger.Log.Errorf("restore cluster phase run failed, error: %s", err.Error())
		_ = c.taskLogService.End(task, false, err.Error())
		cluster.CurrentTaskID = ""
		_ = c.clusterRepo.Save(&cluster)
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterRestore, false, err.Error()), cluster.Name, constant.ClusterRestore)
	} else {
		_ = c.taskLogService.End(task, true, "")
		cluster.CurrentTaskID = ""
		_ = c.clusterRepo.Save(&cluster)
		_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterRestore, true, ""), cluster.Name, constant.ClusterRestore)
	}
}

func (c cLusterBackupFileService) LocalRestore(clusterName string, file []byte) error {
	isON := c.taskLogService.IsTaskOn(clusterName)
	if isON {
		return errors.New("TASK_IN_EXECUTION")
	}
	cluster, err := c.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "SpecNetwork", "SpecRuntime", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return err
	}

	clusterPath := constant.BackupDir + "/" + clusterName
	targetPath := clusterPath + "/" + constant.BackupFileDefaultName
	if _, err := os.Stat(targetPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(clusterPath, os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if _, err = os.Create(targetPath); err != nil {
		return err
	}
	if err = ioutil.WriteFile(targetPath, file, 0775); err != nil {
		return err
	}

	task, err := c.taskLogService.NewTerminalTask(cluster.ID, constant.TaskLogTypeRestore)
	if err != nil {
		return err
	}
	cluster.TaskLog = *task
	cluster.CurrentTaskID = task.ID
	_ = c.clusterRepo.Save(&cluster)

	go c.doLocalRestore(cluster, task)
	return nil
}

func (c cLusterBackupFileService) doLocalRestore(cluster model.Cluster, task *model.TaskLog) {
	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, task.ID)
	if err != nil {
		logger.Log.Errorf("create ansible log failed, error: %s", err.Error())
	}

	admCluster := adm.NewAnsibleHelper(cluster)
	p := &backup.RestoreClusterPhase{}
	err = p.Run(admCluster.Kobe, writer)
	if err != nil {
		logger.Log.Errorf("run cluster log failed, error: %s", err.Error())
		_ = c.taskLogService.End(task, false, err.Error())
		_ = c.messageService.SendMessage(constant.System, false, GetContent(constant.ClusterRestore, false, err.Error()), cluster.Name, constant.ClusterRestore)
	} else {
		_ = c.taskLogService.End(task, true, "")
		_ = c.messageService.SendMessage(constant.System, true, GetContent(constant.ClusterRestore, true, ""), cluster.Name, constant.ClusterRestore)
	}
}

func (c cLusterBackupFileService) deleteBackupFile(clusterName string) {
	clusterBackupStrategy, err := c.clusterBackupStrategyRepository.Get(clusterName)
	if err != nil {
		logger.Log.Errorf("get clusterBackupStrategy [%s]  error : %s", clusterName, err.Error())
		return
	}
	var backupFiles []model.ClusterBackupFile
	db.DB.Where("cluster_id = ?", clusterBackupStrategy.ClusterID).Order("created_at ASC").Find(&backupFiles)
	if len(backupFiles) > clusterBackupStrategy.SaveNum {
		var deleteFileNum = len(backupFiles) - clusterBackupStrategy.SaveNum
		for i := 0; i < deleteFileNum; i++ {
			logger.Log.Infof("delete backup file %s", backupFiles[i].Name)
			err := c.Delete(backupFiles[i].Name)
			if err != nil {
				logger.Log.Errorf("delete cluster [%s] backup file error : %s", clusterName, err.Error())
			}
		}
	}
}
