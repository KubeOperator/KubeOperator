package storage

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
)

type NfsStorageClassCreation struct {
	cluster      *Cluster
	storageClass dto.StorageClass
}

func NewNfsStorageClassCreation(cluster *Cluster, storageClass dto.StorageClass) *NfsStorageClassCreation {
	return &NfsStorageClassCreation{
		cluster:      cluster,
		storageClass: storageClass,
	}
}

func (NfsStorageClassCreation) PreCreate() {

}
func (NfsStorageClassCreation) CreateProvisioner() {

}
func (NfsStorageClassCreation) CreateStorageClass() {

}
