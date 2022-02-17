package service

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/velero"
	"os"
)

type VeleroBackupService interface {
	CreateBackup(backup dto.VeleroBackup) (string, error)
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

func (v veleroBackupService) CreateBackup(backup dto.VeleroBackup) (string, error) {
	result, err := velero.Backup(backup.Name, v.handleArgs(backup))
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

func (v veleroBackupService) handleArgs(backup dto.VeleroBackup) []string {
	args := []string{}

	config, err := v.GetClusterConfig(backup.Cluster)
	if err != nil {
		return args
	}
	configArg := []string{"--kubeconfig", config}
	args = append(configArg, args...)

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
