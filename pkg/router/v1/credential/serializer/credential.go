package serializer

import (
	credentialModel "github.com/KubeOperator/KubeOperator/pkg/model/credential"
)

type Credential struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Username string `json:"username"`
}

func FromModel(model credentialModel.Credential) Credential {
	return Credential{
		Name: model.Name,
	}
}

func ToModel(c Credential) credentialModel.Credential {
	return credentialModel.Credential{
		Name: c.Name,
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
	Name string `json:"name" binding:"required"`
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
