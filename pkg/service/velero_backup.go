package service

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/KubeOperator/KubeOperator/pkg/util/velero"
	"github.com/jinzhu/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VeleroBackupService interface {
	Create(operate string, backup dto.VeleroBackup) (string, error)
	GetBackups(cluster string) (*dto.VeleroBackupList, error)
	GetLogs(cluster, name, operate string) (string, error)
	GetDescribe(cluster, name, operate string) (string, error)
	Delete(cluster, name, operate string) (string, error)
	Install(cluster string, veleroInstall dto.VeleroInstall) (string, error)
	GetConfig(cluster string) (dto.VeleroInstall, error)
	UnInstall(cluster string) error
}

type veleroBackupService struct {
	ClusterService           ClusterService
	clusterRepo              repository.ClusterRepository
	BackupAccountService     BackupAccountService
	taskLogService           TaskLogService
	SystemRegistryRepository repository.SystemRegistryRepository
}

func NewVeleroBackupService() VeleroBackupService {
	return &veleroBackupService{
		ClusterService:           NewClusterService(),
		clusterRepo:              repository.NewClusterRepository(),
		BackupAccountService:     NewBackupAccountService(),
		taskLogService:           NewTaskLogService(),
		SystemRegistryRepository: repository.NewSystemRegistryRepository(),
	}
}

func (v veleroBackupService) Create(operate string, backup dto.VeleroBackup) (string, error) {
	var (
		result []byte
		err    error
	)
	cluster, err := v.clusterRepo.Get(backup.Cluster)
	if err != nil {
		return string(result), err
	}

	var clog model.TaskLog
	clog.ClusterID = cluster.ID
	if len(backup.BackupName) > 0 {
		clog.Type = constant.TaskLogTypeVeleroRestore
		_ = v.taskLogService.Start(&clog)
		cluster.CurrentTaskID = clog.ID
		if err := db.DB.Save(&cluster).Error; err != nil {
			logger.Log.Infof("save cluster failed, err: %v", err)
		}
		go func() {
			result, err := velero.Restore(backup.BackupName, v.handleArgs(backup))
			cluster.CurrentTaskID = ""
			if err := db.DB.Save(&cluster).Error; err != nil {
				logger.Log.Infof("save cluster failed, err: %v", err)
			}
			if err != nil {
				_ = v.taskLogService.End(&clog, false, string(result))
			}
			_ = v.taskLogService.End(&clog, true, string(result))
		}()
	} else {
		clog.Type = constant.TaskLogTypeVeleroBackup
		_ = v.taskLogService.Start(&clog)
		cluster.CurrentTaskID = clog.ID
		if err := db.DB.Save(&cluster).Error; err != nil {
			logger.Log.Infof("save cluster failed, err: %v", err)
		}
		go func() {
			result, err = velero.Create(backup.Name, operate, v.handleArgs(backup))
			cluster.CurrentTaskID = ""
			if err := db.DB.Save(&cluster).Error; err != nil {
				logger.Log.Infof("save cluster failed, err: %v", err)
			}
			if err != nil {
				_ = v.taskLogService.End(&clog, false, string(result))
			}
			_ = v.taskLogService.End(&clog, true, string(result))
		}()
	}
	return string(result), err
}

func (v veleroBackupService) GetBackups(cluster string) (*dto.VeleroBackupList, error) {
	if err := v.checkValid(cluster); err != nil {
		return nil, nil
	}

	var result dto.VeleroBackupList
	config, err := v.GetClusterConfig(cluster)
	if err != nil {
		return &result, err
	}
	args := []string{"--kubeconfig", config}

	schedules, err := velero.Get("schedule", args)
	if err != nil {
		return &result, err
	}

	backups, err := velero.Get("backup", args)
	if err != nil {
		return &result, err
	}
	result.Items = backups
	result.Items = append(result.Items, schedules...)

	return &result, err
}

func (v veleroBackupService) GetLogs(cluster, name, operate string) (string, error) {
	err := v.checkValid(cluster)
	if err != nil {
		return "", nil
	}
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
	err := v.checkValid(cluster)
	if err != nil {
		return "", nil
	}
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

func (v veleroBackupService) GetConfig(cluster string) (dto.VeleroInstall, error) {
	var clusterVelero model.ClusterVelero
	var result dto.VeleroInstall
	if err := db.DB.Where("cluster = ?", cluster).Find(&clusterVelero).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return result, err
	}
	result.ID = clusterVelero.ID
	result.Cluster = clusterVelero.Cluster
	result.BackupAccountName = clusterVelero.BackupAccountName
	result.Requests.Memory = clusterVelero.MemRequest
	result.Requests.Cpu = clusterVelero.CpuRequest
	result.Limits.Cpu = clusterVelero.CpuLimit
	result.Limits.Memory = clusterVelero.MemLimit
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

	clusterModel, err := v.ClusterService.Get(cluster)
	if err != nil {
		return result, err
	}
	arch := ""
	if clusterModel.Architectures == "amd64" {
		arch = "x86_64"
	} else {
		arch = "aarch64"
	}
	registry, err := v.SystemRegistryRepository.GetByArch(arch)
	if err != nil {
		return result, err
	}

	url := constant.LocalRepositoryDomainName + ":" + strconv.Itoa(registry.RegistryPort) + "/"

	if backupAccount.Type == "OSS" {
		args = append(args, "--provider", "alibabacloud")
		args = append(args, "--image", url+"velero/velero:v1.7.1")
		args = append(args, "--bucket", backupAccount.Bucket)
		args = append(args, "--plugins", url+"kubeoperator/velero-plugin-alibabacloud:v1.0.0-2d33b89")
		args = append(args, "--use-volume-snapshots", "false")

		endpoint := vars["endpoint"].(string)
		start := strings.Index(endpoint, "oss-") + 4
		end := strings.Index(endpoint, ".aliyuncs.com")
		region := endpoint[start:end]
		config := "region=" + region
		args = append(args, "--backup-location-config", config)
	}
	if backupAccount.Type == "MINIO" || backupAccount.Type == "S3" {
		args = append(args, "--provider", "aws")
		args = append(args, "--image", url+"velero/velero:v1.7.1")
		args = append(args, "--plugins", url+"velero/velero-plugin-for-aws:v1.2.1")
		args = append(args, "--bucket", backupAccount.Bucket)
		config := "region=minio,s3ForcePathStyle=true,s3Url=http://" + vars["endpoint"].(string)
		args = append(args, "--backup-location-config", config)
	}
	if veleroInstall.Requests.Cpu > 0 {
		args = append(args, "--velero-pod-cpu-request", strconv.Itoa(veleroInstall.Requests.Cpu)+"m")
	}
	if veleroInstall.Requests.Memory > 0 {
		args = append(args, "--velero-pod-mem-request", strconv.Itoa(veleroInstall.Requests.Memory)+"Mi")
	}
	if veleroInstall.Limits.Cpu > 0 {
		args = append(args, "--velero-pod-cpu-limit", strconv.Itoa(veleroInstall.Limits.Cpu)+"m")
	}
	if veleroInstall.Limits.Memory > 0 {
		args = append(args, "--velero-pod-mem-limit", strconv.Itoa(veleroInstall.Limits.Memory)+"Mi")
	}
	//args = append(args, "--wait")
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
		CpuLimit:          veleroInstall.Limits.Cpu,
		CpuRequest:        veleroInstall.Requests.Cpu,
		MemLimit:          veleroInstall.Limits.Memory,
		MemRequest:        veleroInstall.Requests.Memory,
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

func (v veleroBackupService) UnInstall(cluster string) error {
	secret, err := v.ClusterService.GetSecrets(cluster)
	if err != nil {
		return err
	}
	endpoints, err := v.ClusterService.GetApiServerEndpoints(cluster)
	if err != nil {
		return err
	}
	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		return err
	}
	err = kubeClient.CoreV1().Namespaces().Delete(context.Background(), "velero", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = kubeClient.RbacV1().ClusterRoleBindings().Delete(context.Background(), "velero", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	exClient, err := kubernetesUtil.NewKubernetesExtensionClient(&kubernetesUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		return err
	}
	listPvOptions := metav1.ListOptions{
		LabelSelector: "component=velero",
	}
	err = exClient.ApiextensionsV1().CustomResourceDefinitions().DeleteCollection(context.Background(), metav1.DeleteOptions{}, listPvOptions)
	if err != nil {
		return err
	}
	if err := db.DB.Where("cluster = ?", cluster).Delete(&model.ClusterVelero{}).Error; err != nil {
		return err
	}
	return nil
}

func (v veleroBackupService) checkValid(cluster string) error {
	secret, err := v.ClusterService.GetSecrets(cluster)
	if err != nil {
		return err
	}
	endpoints, err := v.ClusterService.GetApiServerEndpoints(cluster)
	if err != nil {
		return err
	}
	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		return err
	}
	velerons, err := kubeClient.CoreV1().Namespaces().Get(context.Background(), "velero", metav1.GetOptions{})
	if err != nil {
		return err
	}
	if velerons.Status.Phase == "Terminating" {
		return errors.New("velero is Terminating")
	}
	return nil
}

func CreateCredential(cluster string, backup dto.BackupAccount) (string, error) {
	var (
		filePath string
	)
	configPath := "/var/ko/velero/configs/" + cluster
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
		_, _ = file.WriteString("[default] \n")
		_, _ = file.WriteString("aws_access_key_id = " + vars["accessKey"].(string) + "\n")
		_, _ = file.WriteString("aws_secret_access_key = " + vars["secretKey"].(string) + "\n")
	}
	if backup.Type == "OSS" {
		vars := make(map[string]interface{})
		if err := json.Unmarshal([]byte(backup.Credential), &vars); err != nil {
			return filePath, err
		}
		_, _ = file.WriteString("ALIBABA_CLOUD_ACCESS_KEY_ID=" + vars["accessKey"].(string) + "\n")
		_, _ = file.WriteString("ALIBABA_CLOUD_ACCESS_KEY_SECRET=" + vars["secretKey"].(string) + "\n")
	}

	return filePath, err
}

func (v veleroBackupService) GetClusterConfig(cluster string) (string, error) {

	var (
		filePath string
	)
	configPath := "/var/ko/velero/configs/" + cluster
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
	_, _ = file.WriteString(config)
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
		schedule := "--schedule=" + backup.Schedule
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
		args = append(args, "--include-resources", handleArray(backup.IncludeResources))
	}
	if len(backup.ExcludeResources) > 0 {
		args = append(args, "--exclude-resources", handleArray(backup.ExcludeResources))
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
