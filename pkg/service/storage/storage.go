package storage

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"k8s.io/client-go/kubernetes"
)

type Cluster struct {
	model.Cluster
	Helm       helm.Interface
	Kubernetes *kubernetes.Clientset
}

func NewCluster(cluster dto.ClusterWithEndpoint) (*Cluster, error) {
	c := &Cluster{Cluster: cluster.Cluster}
	h, err := helm.NewClient(helm.Config{
		ApiServer:   fmt.Sprintf("https://%s:%d", cluster.Endpoint.Address, cluster.Endpoint.Port),
		BearerToken: cluster.Cluster.Secret.KubernetesToken,
	})
	if err != nil {
		return nil, err
	}
	k, err := kubeUtil.NewKubernetesClient(&kubeUtil.Config{
		Host:  cluster.Endpoint.Address,
		Token: cluster.Cluster.Secret.KubernetesToken,
		Port:  cluster.Endpoint.Port,
	})
	if err != nil {
		return nil, err
	}
	c.Helm = h
	c.Kubernetes = k
	return c, nil
}

type ClassCreation interface {
	PreCreate()
	CreateProvisioner()
	CreateStorageClass()
}

func NewStorageClassCreation(cluster dto.ClusterWithEndpoint, class dto.StorageClass) (ClassCreation, error) {
	c, err := NewCluster(cluster)
	if err != nil {
		return nil, err
	}
	switch class.Provisioner {
	case "nfs":
		return NewNfsStorageClassCreation(c, class), nil
	default:
		return nil, errors.New("not supported")
	}
}
