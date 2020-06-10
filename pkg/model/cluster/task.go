package cluster

import (
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Task struct {
	commonModel.BaseModel
	ID        string
	Name      string
	ClusterID string
	Status    string
	Msg       string
	StartTime time.Time
	EndTime   time.Time
}

func (t *Task) BeforeCreate() (err error) {
	t.ID = uuid.NewV4().String()
	return nil
}

func (t Task) TableName() string {
	return "ko_cluster_task"
}
