package serializer

import (
	credentialModel "github.com/KubeOperator/KubeOperator/pkg/model/credential"
	"time"
)

type Credential struct {
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	PrivateKey string    `json:"privateKey"`
	Type       string    `json:"type"`
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
}

func FromModel(model credentialModel.Credential) Credential {
	return Credential{
		Name:       model.Name,
		Username:   model.Username,
		Password:   model.Password,
		PrivateKey: model.PrivateKey,
		Type:       model.Type,
		CreateAt:   model.CreatedAt,
		UpdateAt:   model.UpdatedAt,
	}
}

func ToModel(c Credential) credentialModel.Credential {
	return credentialModel.Credential{
		Name:       c.Name,
		Username:   c.Username,
		Password:   c.Password,
		PrivateKey: c.PrivateKey,
		Type:       c.Type,
	}
}

type ListCredentialResponse struct {
	Items []Credential `json:"items"`
	Total int          `json:"total"`
}

type GetCredentialResponse struct {
	Item Credential `json:"item"`
}

type CreateCredentialRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	Username   string `json:"username"`
}

type DeleteCredentialRequest struct {
	Name string `json:"name"`
}

type UpdateCredentialRequest struct {
	Item Credential `json:"item" binding:"required"`
}

type BatchCredentialRequest struct {
	Operation string       `json:"operation" binding:"required"`
	Items     []Credential `json:"items"`
}

type BatchCredentialResponse struct {
	Items []Credential `json:"items"`
}
