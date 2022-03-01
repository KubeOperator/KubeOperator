package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/velero"
	"github.com/jinzhu/gorm"
	"os"
)

type VeleroBackupService interface {
	Create(operate string, backup dto.VeleroBackup) (string, error)
	Get(cluster, operate string) (map[string]interface{}, error)
	GetLogs(cluster, name, operate string) (string, error)
	GetDescribe(cluster, name, operate string) (string, error)
	Delete(cluster, name, operate string) (string, error)
	Install(cluster string, veleroInstall dto.VeleroInstall) (string, error)
	GetConfig(cluster string) (model.ClusterVelero, error)
}

type veleroBackupService struct {
	ClusterService       ClusterService
	BackupAccountService BackupAccountService
}

func NewVeleroBackupService() VeleroBackupService {
	return &veleroBackupService{
		ClusterService:       NewClusterService(),
		BackupAccountService: NewBackupAccountService(),
	}
}

func (v veleroBackupService) Create(operate string, backup dto.VeleroBackup) (string, error) {

	var (
		result []byte
		err    error
	)

	if len(backup.BackupName) > 0 {
		result, err = velero.Restore(backup.BackupName, v.handleArgs(backup))
		if err != nil {
			return string(result), err
		}
	} else {
		result, err = velero.Create(backup.Name, operate, v.handleArgs(backup))
		if err != nil {
			return string(result), err
		}
	}
	return string(result), err
}

func (v veleroBackupService) Get(cluster, operate string) (map[string]interface{}, error) {
	var result map[string]interface{}
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}
	args := []string{"--kubeconfig", config}
	res, err := velero.Get(operate, args)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (v veleroBackupService) GetLogs(cluster, name, operate string) (string, error) {
	var result string
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}
	res, err := velero.GetLogs(name, operate, args)
	if err != nil {
		return result, err
	}
	result = string(res)
	return result, err
}

func (v veleroBackupService) GetDescribe(cluster, name, operate string) (string, error) {
	var result string
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}
	res, err := velero.GetDescribe(name, operate, args)
	if err != nil {
		return result, err
	}
	result = string(res)
	return result, err
}

func (v veleroBackupService) Delete(cluster, name, operate string) (string, error) {
	var result string
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}
	res, err := velero.Delete(name, operate, args)
	if err != nil {
		return result, err
	}
	result = string(res)
	return result, err
}

func (v veleroBackupService) GetConfig(cluster string) (model.ClusterVelero, error) {
	var result model.ClusterVelero
	if err := db.DB.Where("cluster = ?", cluster).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return result, err
	}
	return result, nil
}

func (v veleroBackupService) Install(cluster string, veleroInstall dto.VeleroInstall) (string, error) {
	var result string
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}

	backupAccount, err := v.BackupAccountService.Get(veleroInstall.BackupAccountName)
	if err != nil {
		return result, err
	}
	vars := make(map[string]interface{})
	if err := json.Unmarshal([]byte(backupAccount.Credential), &vars); err != nil {
		return result, err
	}
	credentialPath, err := CreateCredential(cluster, *backupAccount)
	if err != nil {
		return result, err
	}
	args = append(args, "--secret-file", credentialPath)
	if backupAccount.Type == "OSS" {
		args = append(args, "--provider", "alibabacloud")
		args = append(args, "--image", "registry.cn-hangzhou.aliyuncs.com/acs/velero:latest")
		args = append(args, "--bucket", backupAccount.Bucket)
		args = append(args, "--plugins", "registry.cn-hangzhou.aliyuncs.com/acs/velero-plugin-alibabacloud:v1.0.0-2d33b89")
	}
	if backupAccount.Type == "MINIO" || backupAccount.Type == "S3" {
		args = append(args, "--provider", "aws")
		args = append(args, "--plugins", "velero/velero-plugin-for-aws:v1.2.1")
		args = append(args, "--bucket", backupAccount.Bucket)
		config := "srUrl=" + vars["endpoint"].(string)
		args = append(args, "--backup-location-config", config)
	}

	res, err := velero.Install(args)
	if err != nil {
		logger.Log.Errorf("install velero error: %s", err.Error())
		return result, err
	}

	clusterVelero := &model.ClusterVelero{
		Cluster:           cluster,
		BackupAccountName: veleroInstall.BackupAccountName,
		Bucket:            vars["bucket"].(string),
		Endpoint:          vars["endpoint"].(string),
	}
	if veleroInstall.ID != "" {
		clusterVelero.ID = veleroInstall.ID
		err = db.DB.Save(&clusterVelero).Error
	} else {
		err = db.DB.Create(&clusterVelero).Error
	}

	if err != nil {
		return result, err
	}
	result = string(res)
	return result, err
}

//func (v veleroBackupService)UnInstall(cluster string) error {
//	secret, err := v.ClusterService.GetSecrets(cluster)
//	if err != nil {
//		return  err
//	}
//	endpoints, err := v.ClusterService.GetApiServerEndpoints(cluster)
//	if err != nil {
//		return err
//	}
//	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
//		Hosts: endpoints,
//		Token: secret.KubernetesToken,
//	})
//
//	return nil
//}

func CreateCredential(cluster string, backup dto.BackupAccount) (string, error) {
	var (
		filePath string
	)
	configPath := "/Users/zk.wang/configs/" + cluster
	filePath = configPath + "/credentials-velero"
	_, err := os.Stat(filePath)
	if err == nil {
		err := os.Remove(filePath)
		if err != nil {
			return filePath, err
		}
	}
	_, err = os.Stat(configPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return filePath, err
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return filePath, err
	}
	defer file.Close()

	if backup.Type == "MINIO" || backup.Type == "S3" {
		vars := make(map[string]interface{})
		if err := json.Unmarshal([]byte(backup.Credential), &vars); err != nil {
			return filePath, err
		}
		file.WriteString("[default] \n")
		file.WriteString("aws_access_key_id = " + vars["accessKey"].(string) + "\n")
		file.WriteString("aws_secret_access_key = " + vars["secretKey"].(string) + "\n")
	}

	return filePath, err
}

func (v veleroBackupService) GetClusterConfig(cluster string) (string, error) {

	var (
		filePath string
	)
	configPath := "/Users/zk.wang/configs/" + cluster
	filePath = configPath + "/config"
	_, err := os.Stat(filePath)
	if err == nil {
		return filePath, nil
	}

	_, err = os.Stat(configPath)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(configPath, os.ModePerm)
		if err != nil {
			return filePath, err
		}
	}

	config, err := v.ClusterService.GetKubeconfig(cluster)
	if err != nil {
		return filePath, err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return filePath, err
	}
	file.WriteString(config)
	defer file.Close()

	return filePath, err
}

func (v veleroBackupService) handleArgs(backup dto.VeleroBackup) []string {
	args := []string{}

	config, err := v.GetClusterConfig(backup.Cluster)
	if err != nil {
		return args
	}
	configArg := []string{"--kubeconfig", config}
	args = append(configArg, args...)

	if len(backup.Schedule) > 0 {
		schedule := "--schedule=\"" + backup.Schedule + "\""
		args = append(args, schedule)
	}

	if !backup.IncludeClusterResources {
		args = append(args, "--include-cluster-resources=false")
	}
	if len(backup.Labels) > 0 {
		args = append(args, "--labels", backup.Labels)
	}
	if len(backup.IncludeNamespaces) > 0 {
		args = append(args, "--include-namespaces", handleArray(backup.IncludeNamespaces))
	}
	if len(backup.ExcludeNamespaces) > 0 {
		args = append(args, "--exclude-namespaces", handleArray(backup.ExcludeNamespaces))
	}
	if len(backup.IncludeResources) > 0 {
		args = append(args, "--include-resources", backup.IncludeResources)
	}
	if len(backup.ExcludeResources) > 0 {
		args = append(args, "--exclude-resources", backup.ExcludeResources)
	}
	if len(backup.Selector) > 0 {
		args = append(args, "--selector", backup.Selector)
	}
	if backup.Ttl != "" {
		args = append(args, "--ttl", backup.Ttl)
	}
	return args
}

func handleArray(arr []string) string {
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result = result + "," + arr[i]
	}
	return result
}
