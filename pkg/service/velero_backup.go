package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/util/velero"
	"os"
)

type VeleroBackupService interface {
	CreateBackup(cluster string) (string, error)
	GetBackups(cluster string) (map[string]interface{}, error)
	GetBackupDescribe(cluster string, backupName string) (string, error)
	GetBackupLogs(cluster string, backupName string) (string, error)
}

type veleroBackupService struct {
	ClusterService ClusterService
}

func NewVeleroBackupService() VeleroBackupService {
	return &veleroBackupService{
		ClusterService: NewClusterService(),
	}
}

func (v veleroBackupService) CreateBackup(cluster string) (string, error) {
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return "", err
	}

	args := []string{"--kubeconfig", config}
	result, err := velero.Backup(cluster, args)
	if err != nil {
		return string(result), err
	}
	return string(result), err
}

func (v veleroBackupService) GetBackups(cluster string) (map[string]interface{}, error) {
	var result map[string]interface{}
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}
	res, err := velero.GetBackups(args)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		return result, err
	}
	return result, err
}

func (v veleroBackupService) GetBackupDescribe(cluster string, backupName string) (string, error) {
	var result string
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}
	res, err := velero.GetBackupDescribe(backupName, args)
	if err != nil {
		return result, err
	}
	result = string(res)
	return result, err
}

func (v veleroBackupService) GetBackupLogs(cluster string, backupName string) (string, error) {
	var result string
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}
	res, err := velero.GetBackupLogs(backupName, args)
	if err != nil {
		return result, err
	}
	result = string(res)
	return result, err
}

func (v veleroBackupService) GetRestores(cluster string) (map[string]interface{}, error) {
	var result map[string]interface{}
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return result, err
	}

	args := []string{"--kubeconfig", config}
	res, err := velero.GetRestores(args)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		return result, err
	}
	return result, err
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
	file.Close()

	return filePath, err
}
