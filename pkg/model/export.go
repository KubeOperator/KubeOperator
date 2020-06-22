package model

import ()

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
}
