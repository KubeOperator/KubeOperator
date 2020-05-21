package credential

import (
	"ko3-gin/pkg/model/common"
	credentialModel "ko3-gin/pkg/model/credential"
)

type Credential struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Username string `json:"username"`
}

func FromModel(model credentialModel.Credential) Credential {
	return Credential{
		Name:   model.Name,
		Status: model.Status,
	}
}

func ToModel(c Credential) credentialModel.Credential {
	return credentialModel.Credential{
		BaseModel: common.BaseModel{
			Name: c.Name,
		},
	}
}

type ListResponse struct {
	Items []Credential `json:"items"`
	Total int          `json:"total"`
}

type GetResponse struct {
	Item Credential `json:"item"`
}

type CreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateResponse struct {
	Item Credential `json:"item"`
}

type DeleteRequest struct {
	Name string `json:"name"`
}

type DeleteResponse struct {
}

type UpdateRequest struct {
	Item Credential `json:"item" binding:"required"`
}

type UpdateResponse struct {
	Item Credential `json:"item"`
}

type BatchRequest struct {
	Operation string       `json:"operation" binding:"required"`
	Items     []Credential `json:"items"`
}

type BatchResponse struct {
	Items []Credential `json:"items"`
}
