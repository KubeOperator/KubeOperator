package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
)

type ClusterSecretRepository interface {
	Get(id string) (model.ClusterSecret, error)
	Save(status *model.ClusterSecret) error
	Delete(id string) error
}

func NewClusterSecretRepository() ClusterSecretRepository {
	return &clusterSecretRepository{}
}

type clusterSecretRepository struct {
}

func (c clusterSecretRepository) Get(id string) (model.ClusterSecret, error) {
	status := model.ClusterSecret{
		ID: id,
	}
	if err := db.DB.First(&status).Error; err != nil {
		return status, err
	}

	if len(status.KubeadmToken) != 0 {
		admToken, err := encrypt.StringDecryptWithSalt(status.KubeadmToken)
		if err != nil {
			return status, err
		}
		status.KubeadmToken = admToken
	}
	if len(status.KubernetesToken) != 0 {
		token, err := encrypt.StringDecryptWithSalt(status.KubernetesToken)
		if err != nil {
			return status, err
		}
		status.KubernetesToken = token
	}
	return status, nil
}

func (c clusterSecretRepository) Save(status *model.ClusterSecret) error {
	if len(status.KubeadmToken) != 0 {
		admToken, err := encrypt.StringEncryptWithSalt(status.KubeadmToken)
		if err != nil {
			return err
		}
		status.KubeadmToken = admToken
	}
	if len(status.KubernetesToken) != 0 {
		token, err := encrypt.StringEncryptWithSalt(status.KubernetesToken)
		if err != nil {
			return err
		}
		status.KubernetesToken = token
	}

	if db.DB.NewRecord(status) {
		if err := db.DB.Create(&status).Error; err != nil {
			return err
		}
	} else {
		if err := db.DB.Save(&status).Error; err != nil {
			return err
		}
	}
	return nil
}

func (c clusterSecretRepository) Delete(id string) error {
	secret := model.ClusterSecret{ID: id}
	if err := db.DB.First(&secret).Error; err != nil {
		return err
	}
	if err := db.DB.Delete(&secret).Error; err != nil {
		return err
	}
	return nil
}
