package common

import "time"

type BaseModel struct {
	ID        string `gorm:"primary_key"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
