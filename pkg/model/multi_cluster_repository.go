package model

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/model/common"
	"github.com/KubeOperator/KubeOperator/pkg/util/git"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
	"time"
)

type MultiClusterRepository struct {
	common.BaseModel
	ID           string    `json:"-"`
	Name         string    `json:"name"`
	Source       string    `json:"source"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	Branch       string    `json:"branch"`
	LastSyncHead string    `json:"lastSyncHead"`
	LastSyncTime time.Time `json:"lastSyncTime"`
	SyncInterval int64     `json:"syncInterval"`
	GitTimeout   int64     `json:"gitTimeout"`
	SyncEnable   bool      `json:"syncEnable"`
	SyncStatus   string    `json:"syncStatus"`
}

func (m *MultiClusterRepository) BeforeCreate() error {
	m.ID = uuid.NewV4().String()
	return nil
}

func (m *MultiClusterRepository) BeforeDelete() error {
	var mls []MultiClusterSyncLog
	if err := db.DB.Where(MultiClusterSyncLog{
		MultiClusterRepositoryID: m.ID,
	}).Find(&mls).Error; err != nil {
		return err
	}
	tx := db.DB.Begin()
	for m := range mls {
		if err := db.DB.Delete(&m).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *MultiClusterRepository) AfterDelete() error {
	_ = os.RemoveAll(path.Join(constant.DefaultDataDir, m.Name))
	return nil
}

func (m *MultiClusterRepository) Pull() error {
	if err := git.UpdateRepository(path.Join(constant.DefaultRepositoryDir, m.Name), m.Branch,
		&http.BasicAuth{Username: m.Username, Password: m.Password}); err != nil {
		return err
	}
	return nil
}
