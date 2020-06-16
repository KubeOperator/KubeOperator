package model

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
	Credential{},
	Host{},
	Volume{},
	User{},
	Demo{},
}
