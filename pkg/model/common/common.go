package common

import "time"

type BaseModel struct {
	ID          string `gorm:"primary_key"`
	Name        string
	CreatedDate time.Time
	UpdatedDate time.Time
	Status      string
}
