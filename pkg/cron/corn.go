package cron

import (
	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

const phaseName = "cron"

type InitCronPhase struct {
}

func (c *InitCronPhase) Init() error {
	Cron := cron.New()
	Cron.Start()
	return nil
}

func (c *InitCronPhase) PhaseName() string {
	return phaseName
}
