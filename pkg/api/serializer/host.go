package serializer

import "ko3-gin/pkg/model"

type CreatHostRequest struct {
	Name           string `json:"name"`
	Ip             string `json:"ip"`
	Port           int    `json:"port"`
	CredentialName string `json:"credential_name"`
}

type UpdateHostRequest struct {
	Name           string `json:"name"`
	Ip             string `json:"ip"`
	Port           string `json:"port"`
	CredentialName string `json:"credential_name"`
}

type PageHostResponse struct {
	Total int
	Items []model.Host
}
