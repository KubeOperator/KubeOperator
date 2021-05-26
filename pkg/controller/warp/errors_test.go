package warp

import (
	"fmt"
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func TestError(t *testing.T) {
	config.Init()
	logger.Init()
	dbi := db.InitDBPhase{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Name:     viper.GetString("db.name"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
	}
	_ = dbi.Init()

	id := "1213123"
	_, err := getCluster(id)
	if errors.Is(errors.Cause(err), gorm.ErrRecordNotFound) {
		fmt.Print("这两个错误是同一种，可以直接一起输出\n")
	}
	fmt.Printf("stack trace: \n%+v\n", err)
	logger.Log.Info(fmt.Sprintf("%+v", err))
}

func getCluster(id string) (model.Cluster, error) {
	var cluster model.Cluster
	err := db.DB.Where("id = ?", id).Find(&cluster).Error
	if err != nil {
		err = errors.Wrapf(err, "get cluster by %s err", id)
		return cluster, err
	}
	return cluster, nil
}
