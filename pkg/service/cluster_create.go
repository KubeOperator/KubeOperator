package service

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/ansible"
	"github.com/sirupsen/logrus"
)

func (c clusterService) Create(creation dto.ClusterCreate) (*dto.Cluster, error) {
	loginfo, _ := json.Marshal(creation)
	logger.Log.WithFields(logrus.Fields{"cluster_creation": string(loginfo)}).Debugf("start to create the cluster %s", creation.Name)

	cluster := creation.ClusterCreateDto2Mo()
	tx := db.DB.Begin()
	var project model.Project
	if err := tx.Where("name = ?", creation.ProjectName).First(&project).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("select project failed, err: %v", err)
	}
	cluster.ProjectID = project.ID

	if err := tx.Create(&cluster.Secret).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	cluster.SecretID = cluster.Secret.ID
	if err := tx.Create(&cluster).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if cluster.Provider == constant.ClusterProviderPlan {
		if err := tx.Where("name = ?", creation.Plan).Preload("Zones").Preload("Region").First(&cluster.Plan).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("select plan %s failed, err: %s", creation.Plan, err.Error())
		}
	} else {
		if err := c.clusterIaasService.LoadMetalNodes(&creation, cluster, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	cluster.TaskLog.ClusterID = cluster.ID
	if err := c.tasklogService.Start(&cluster.TaskLog); err != nil {
		tx.Rollback()
		return nil, err
	}
	var manifest model.ClusterManifest
	if err := db.DB.Where("name = ?", cluster.Version).First(&manifest).Error; err != nil {
		return nil, err
	}
	var otherVars []dto.NameVersion
	if err := json.Unmarshal([]byte(manifest.OtherVars), &otherVars); err != nil {
		return nil, err
	}
	ingressType, ingressVersion := "ingress-nginx", ""
	if creation.IngressControllerType == "traefik" {
		ingressType = "traefik"
	}
	for _, otherVar := range otherVars {
		if otherVar.Name == ingressType {
			ingressVersion = otherVar.Version
			break
		}
	}
	cluster.SpecComponent = cluster.PrepareComponent(creation.IngressControllerType, ingressVersion, creation.EnableDnsCache, creation.SupportGpu)
	for _, component := range cluster.SpecComponent {
		if err := tx.Create(&component).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	cluster.SpecConf.ClusterID = cluster.ID
	if err := tx.Create(&cluster.SpecConf).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	cluster.SpecRuntime.ClusterID = cluster.ID
	if err := tx.Create(&cluster.SpecRuntime).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	cluster.SpecNetwork.ClusterID = cluster.ID
	if err := tx.Create(&cluster.SpecNetwork).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	projectResource := model.ProjectResource{
		ResourceID:   cluster.ID,
		ProjectID:    project.ID,
		ResourceType: constant.ResourceCluster,
	}
	if err := tx.Create(&projectResource).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("can not create project  %s resource reason %s", project.Name, err.Error())
	}

	subscribe := model.NewMsgSubscribe(constant.ClusterOperator, constant.Cluster, cluster.ID)
	if err := tx.Create(&subscribe).Error; err != nil {
		tx.Rollback()
	}

	writer, err := ansible.CreateAnsibleLogWriterWithId(cluster.Name, cluster.TaskLog.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"log_id": cluster.TaskLog.ID,
	}).Debugf("get ansible writer log of cluster %s successful, now start to init the cluster", cluster.Name)

	tx.Commit()

	logger.Log.Infof("init db data of cluster %s successful, now start to create cluster", cluster.Name)
	go c.clusterInitService.Init(*cluster, writer)

	return &dto.Cluster{Cluster: *cluster}, nil
}
