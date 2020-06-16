package repository

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type UserRepository interface {
	Page(num, size int) (int, []model.User, error)
	List() ([]model.User, error)
	Get(name string) (model.User, error)
	Save(item *model.User) error
	Delete(name string) error
	Batch(operation string, items []model.User) ([]model.User, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (u userRepository) Page(num, size int) (int, []model.User, error) {
	var total int
	var users []model.User
	err := db.DB.Model(model.User{}).Count(&total).Find(&users).Offset((num - 1) * size).Limit(size).Error
	return total, users, err
}

func (u userRepository) List() ([]model.User, error) {
	var users []model.User
	err := db.DB.Model(model.User{}).Find(&users).Error
	return users, err
}

func (u userRepository) Get(name string) (model.User, error) {
	var user model.User
	user.Name = name
	if err := db.DB.Where(user).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u userRepository) Save(item *model.User) error {
	if db.DB.NewRecord(item) {
		return db.DB.Create(&item).Error
	} else {
		return db.DB.Save(&item).Error
	}
}

func (u userRepository) Delete(name string) error {
	var user model.User
	user.Name = name
	return db.DB.Delete(&user).Error
}

func (u userRepository) Batch(operation string, items []model.User) ([]model.User, error) {
	var deleteItems []model.User
	switch operation {
	case constant.BatchOperationDelete:
		tx := db.DB.Begin()
		for _, item := range items {
			err := db.DB.Model(model.User{}).First(&item).Delete(&item).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			deleteItems = append(deleteItems, item)
			tx.Commit()
		}
	default:
		return nil, constant.NotSupportedBatchOperation
	}
	return deleteItems, nil
}
