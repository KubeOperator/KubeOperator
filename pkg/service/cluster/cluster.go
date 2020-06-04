package cluster

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
)

func Page(num, size int) (clusters []clusterModel.Cluster, total int, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).
		Count(&total).
		Offset((num - 1) * size).
		Limit(size).
		Preload("Status").
		Preload("Spec").
		Preload("Nodes").
		Find(&clusters).
		Error
	return
}

func List() (clusters []clusterModel.Cluster, err error) {
	err = db.DB.Model(clusterModel.Cluster{}).
		Preload("Spec").
		Preload("Status").
		Preload("Nodes").
		Find(&clusters).Error
	return
}

func Get(name string) (clusterModel.Cluster, error) {
	var result clusterModel.Cluster
	result.Name = name
	if err := db.DB.Where(result).First(&result).Error; err != nil {
		return result, err
	}
	if err := db.DB.First(&result).
		Related(&result.Spec).
		Related(&result.Status).Error; err != nil {
		return result, err
	}
	return result, nil
}

func Save(item *clusterModel.Cluster) error {
	if err := db.DB.Create(&item).Error; err != nil {
		return err
	}
	return InitCluster(*item)
}

func Delete(name string) error {
	var cluster clusterModel.Cluster
	if err := db.DB.Where(clusterModel.Cluster{Name: name}).
		First(&cluster).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&cluster).Error; err != nil {
		return err
	}
	return nil
}

func Batch(operation string, items []clusterModel.Cluster) ([]clusterModel.Cluster, error) {
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, c := range items {
			if err := Delete(c.Name); err != nil {
				tx.Rollback()
				return nil, err
			}

		}
		tx.Commit()
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return items, nil
}

func GetClusterStatus(clusterName string) (clusterModel.Status, error) {
	var cluster clusterModel.Cluster
	var status clusterModel.Status
	if err := db.DB.
		Where(&clusterModel.Cluster{Name: clusterName}).
		First(&cluster).Error; err != nil {
		return status, err
	}
	if err := db.DB.
		Where(clusterModel.Status{ID: cluster.StatusID}).
		First(&status).Error; err != nil {
		return status, err
	}
	if err := db.DB.
		First(&status).
		Order("last_probe_time asc").
		Related(&status.Conditions).Error; err != nil {
		return status, err
	}
	return status, nil
}

func SaveClusterStatus(status *clusterModel.Status) error {
	return db.DB.Save(status).Error
}

func GetClusterKubernetesApiEndpoint(clusterName string) (string, error) {
	c, err := Get(clusterName)
	if err != nil {
		return "", err
	}
	if c.Spec.LbKubeApiserverIp != "" {
		return KubernetesApiEndpointFromIp(c.Spec.LbKubeApiserverIp), nil
	}
	var node clusterModel.Node
	if err := db.DB.
		Where(clusterModel.Node{Role: constant.NodeRoleNameMaster}).
		First(node).
		Related(&(node.Host)).Error; err != nil {
		return "", err
	}
	return KubernetesApiEndpointFromIp(node.Host.Ip), nil
}

func KubernetesApiEndpointFromIp(ip string) string {
	return fmt.Sprintf("https://%s:%d", ip, facts.DefaultFacts[facts.LbKubeApiserverPortFactName])
}

func GetClusterSecret(name string) (secret clusterModel.Secret, err error) {
	var cluster clusterModel.Cluster
	if err := db.DB.Where(clusterModel.Cluster{Name: name}).
		First(&cluster).Error; err != nil {
		return
	}
	if err := db.DB.Where(clusterModel.Secret{ID: cluster.SecretID}).
		First(&secret).Error; err != nil {
		return
	}
	return
}
