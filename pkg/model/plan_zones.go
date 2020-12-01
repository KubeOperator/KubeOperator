package model

import uuid "github.com/satori/go.uuid"

type PlanZones struct {
	ID     string `json:"id" grom:"type:varchar(64)"`
	PlanID string `json:"planId" grom:"type:varchar(64)"`
	ZoneID string `json:"zoneId" grom:"type:varchar(64)"`
}

func (p *PlanZones) BeforeCreate() (err error) {
	p.ID = uuid.NewV4().String()
	return err
}
