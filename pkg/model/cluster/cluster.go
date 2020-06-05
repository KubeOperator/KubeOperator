package cluster

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/kobe/api"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Cluster struct {
	common.BaseModel
	ID       string
	Name     string
	SpecID   string
	SecretID string
	StatusID string
	Status   Status `gorm:"save_associations:false"`
	Spec     Spec   `gorm:"save_associations:false"`
	Secret   Secret `gorm:"save_associations:false"`
	Nodes    []Node `gorm:"save_associations:false"`
}

func (c Cluster) TableName() string {
	return "ko_cluster"
}

func (c *Cluster) BeforeCreate(scope *gorm.Scope) error {
	c.ID = uuid.NewV4().String()
	c.Spec.ID = uuid.NewV4().String()
	c.Status.ID = uuid.NewV4().String()
	c.Status = Status{
		Message: "",
		Phase:   constant.ClusterWaiting,
	}
	if err := scope.DB().
		Create(&(c.Spec)).Error; err != nil {
		return err
	}
	if err := scope.DB().
		Create(&(c.Status)).Error; err != nil {
		return err
	}
	c.SpecID = c.Spec.ID
	c.StatusID = c.Status.ID
	return nil
}

func (c *Cluster) AfterCreate(scope *gorm.Scope) error {
	workerNo := 1
	masterNo := 1
	for _, node := range c.Nodes {
		node.ClusterID = c.ID
		switch node.Role {
		case constant.NodeRoleNameMaster:
			node.Name = fmt.Sprintf("%s-%d", constant.NodeRoleNameMaster, masterNo)
			masterNo++
		case constant.NodeRoleNameWorker:
			node.Name = fmt.Sprintf("%s-%d", constant.NodeRoleNameWorker, workerNo)
			workerNo++
		}
		if err := db.DB.
			Where(hostModel.Host{Name: node.Host.Name}).
			First(&node.Host).Error; err != nil {
			return err
		}
		if err := db.DB.Create(&node).Error; err != nil {
			return err
		}
		node.Host.NodeID = node.ID
		if err := db.DB.Save(&node.Host).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c Cluster) BeforeDelete(scope *gorm.Scope) error {
	return nil
}

func (c Cluster) AfterDelete(scope *gorm.Scope) error {
	err := scope.DB().
		Delete(Spec{ID: c.SpecID}).
		Delete(Status{ID: c.StatusID}).Error
	if err != nil {
		return err
	}
	if err := scope.DB().
		Where(Node{ClusterID: c.ID}).
		Delete(Node{}).Error; err != nil {
		return err
	}
	return nil
}

func (c Cluster) ParseInventory() api.Inventory {
	var masters []string
	var workers []string
	var chrony []string
	var hosts []*api.Host
	for _, node := range c.Nodes {
		hosts = append(hosts, node.ToKobeHost())
		switch node.Role {
		case constant.NodeRoleNameMaster:
			masters = append(masters, node.Name)
		case constant.NodeRoleNameWorker:
			workers = append(workers, node.Name)
		}
	}
	if len(masters) > 0 {
		chrony = append(chrony, masters[0])
	}

	return api.Inventory{
		Hosts: hosts,
		Groups: []*api.Group{
			{
				Name:     "kube-master",
				Hosts:    masters,
				Children: []string{},
				Vars:     map[string]string{},
			},
			{
				Name:     "kube-worker",
				Hosts:    workers,
				Children: []string{},
				Vars:     map[string]string{},
			},
			{
				Name:     "new-worker",
				Hosts:    []string{},
				Children: []string{},
				Vars:     map[string]string{},
			}, {

				Name:     "lb",
				Hosts:    []string{},
				Children: []string{},
				Vars:     map[string]string{},
			},
			{
				Name:     "etcd",
				Hosts:    masters,
				Children: []string{"master"},
				Vars:     map[string]string{},
			}, {
				Name:     "chrony",
				Hosts:    chrony,
				Children: []string{},
				Vars:     map[string]string{},
			},
		},
	}
}

func (c *Cluster) SetSecret(secret Secret) error {
	c.Secret = secret
	if db.DB.NewRecord(secret) {
		if err := db.DB.
			Create(&(c.Secret)).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.
			Save(&(c.Secret)).Error; err != nil {
			return err
		}
	}
	c.SecretID = c.Secret.ID
	if err := db.DB.Save(&c).Error; err != nil {
		return err
	}
	return nil
}

func (c Cluster) FistMaster() Node {
	var master Node
	for i, _ := range c.Nodes {
		if c.Nodes[i].Role == constant.NodeRoleNameMaster {
			master = c.Nodes[i]
			break
		}
	}
	return master
}
