package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
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
	Zone{},
	Plan{},
	PlanZones{},
}

var InitData = []Interface{
	CloudProvider{
		BaseModel: common.BaseModel{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		ID:   "065ca3f7-3208-4bce-bf0d-a5d30253f39c",
		Name: "OpenStack",
		Vars: "{\"provider\": \"OpenStack\", \"imageVmdkPath\": \"/data/iso/openstack/kubeoperator_centos_7.6.1810-1.qcow2\", \"imageName\": \"kubeoperator_centos_7.6.1810\"}",
	},
	CloudProvider{
		BaseModel: common.BaseModel{
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		},
		ID:   "214516d3-eccf-4275-bc65-e8b941738488",
		Name: "vSphere",
		Vars: "{\"provider\": \"vSphere\", \"imageOvfPath\": \"/data/iso/vsphere/kubeoperator_centos_7.6.1810.ovf\", \"imageVmdkPath\": \"/data/iso/vsphere/kubeoperator_centos_7.6.1810-1.vmdk\", \"imageName\": \"kubeoperator_centos_7.6.1810\"}",
	},
}
