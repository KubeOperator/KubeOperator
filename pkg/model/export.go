package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Interface interface {
	TableName() string
}

var Models = []Interface{
	Cluster{},
	ClusterStatusCondition{},
	ClusterSpec{},
	ClusterStatus{},
	ClusterNode{},
	ClusterSecret{},
	ClusterMonitor{},
	Credential{},
	Host{},
	Volume{},
	User{},
	Demo{},
	CloudProvider{},
	Region{},
}

var InitData = []Interface{
	CloudProvider{
		BaseModel: common.BaseModel{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		ID:   uuid.NewV4().String(),
		Name: "OpenStack",
		Vars: "{\"provider\": \"OpenStack\", \"imageVmdkPath\": \"/data/iso/openstack/kubeoperator_centos_7.6.1810-1.qcow2\", \"imageName\": \"kubeoperator_centos_7.6.1810\"}",
	},
	CloudProvider{
		BaseModel: common.BaseModel{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		ID:   uuid.NewV4().String(),
		Name: "vSphere",
		Vars: "{\"provider\": \"vSphere\", \"imageOvfPath\": \"/data/iso/vsphere/kubeoperator_centos_7.6.1810.ovf\", \"imageVmdkPath\": \"/data/iso/vsphere/kubeoperator_centos_7.6.1810-1.vmdk\", \"imageName\": \"kubeoperator_centos_7.6.1810\"}",
	},
}
