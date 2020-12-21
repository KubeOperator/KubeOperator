package hook

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
)

func init() {
	BeforeApplicationStart.AddFunc(markClusterDirtyData)
	BeforeApplicationStart.AddFunc(markClusterNodeDirtyData)
}

var clusterService = service.NewClusterService()

// cluster
func markClusterDirtyData() error {
	clusters, err := clusterService.List()
	if err != nil {
		return err
	}
	var clusterIds []string
	for _, cluster := range clusters {
		if cluster.Status == constant.StatusCreating || cluster.Status == constant.StatusTerminating || cluster.Status == constant.StatusInitializing || cluster.Status == constant.StatusWaiting {
			clusterIds = append(clusterIds, cluster.ID)
		}
	}
	if err := db.DB.Model(model.Cluster{}).Where("id in (?)", clusterIds).Updates(map[string]interface{}{"dirty": 1}).Error; err != nil {
		return err
	}
	return nil
}

// cluster node
func markClusterNodeDirtyData() error {
	var status = []string{constant.StatusTerminating, constant.StatusInitializing, constant.StatusCreating, constant.StatusWaiting}
	if err := db.DB.Model(model.ClusterNode{}).Where("status in (?)", status).Updates(map[string]interface{}{"dirty": 1}).Error; err != nil {
		return err
	}
	return nil
}
