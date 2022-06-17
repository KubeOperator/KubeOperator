package cron

import (
	"fmt"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/cron/job"
	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

const phaseName = "cron"

type InitCronPhase struct {
	Enable bool
}

func (c *InitCronPhase) Init() error {
	//Cron = cron.New()
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	Cron = cron.New(cron.WithLocation(nyc))
	if c.Enable {
		_, err := Cron.AddJob("@hourly", job.NewRefreshHostInfo())
		if err != nil {
			return fmt.Errorf("can not add corn job: %s", err.Error())
		}
		_, err = Cron.AddJob("@daily", job.NewClusterBackup())
		if err != nil {
			return fmt.Errorf("can not add backup corn job: %s", err.Error())
		}
		Cron.Start()
	}
	return nil
}

func (c *InitCronPhase) PhaseName() string {
	return phaseName
}
