package model

import uuid "github.com/satori/go.uuid"

type Demo struct {
	ID   string
	Name string
}

func (d *Demo) BeforeCreate() {
	d.ID = uuid.NewV4().String()
}
