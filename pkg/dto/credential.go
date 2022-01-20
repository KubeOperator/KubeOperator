package dto

import "time"

type Credential struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"type:varchar(64)"`
	Type      string    `json:"type" gorm:"type:varchar(64)"`
	CreatedAt time.Time `json:"createdAt"`
}

type CredentialPage struct {
	Items []Credential `json:"items"`
	Total int          `json:"total"`
}

type CredentialCreate struct {
	Name       string `json:"name" validate:"koname,required,max=30"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"-"`
	PrivateKey string `json:"privateKey" validate:"-"`
	Type       string `json:"type" validate:"required"`
}

type CredentialUpdate struct {
	ID         string `json:"id" validate:"required"`
	Name       string `json:"name" validate:"koname,required,max=30"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"-"`
	PrivateKey string `json:"privateKey" validate:"-"`
	Type       string `json:"type" validate:"required"`
}

type CredentialBatchOp struct {
	Operation string       `json:"operation" validate:"required"`
	Items     []Credential `json:"items" validate:"required"`
}
