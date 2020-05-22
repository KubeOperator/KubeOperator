package serializer

import (
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
)

type Host struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Memory       string `json:"memory"`
	CpuCore      int    `json:"cpu_core"`
	Os           string `json:"os"`
	OsVersion    string `json:"os_version"`
	CpuNum       int    `json:"cpu_num"`
	GpuInfo      string `json:"gpu_info"`
	Ip           string `json:"ip"`
	Port         int    `json:"port"`
	CredentialId string `json:"credential_id"`
	ClusterId    string `json:"cluster_id"`
	NodeId       string `json:"node_id"`
	Status       string `json:"status"`
}

func FromModel(h hostModel.Host) Host {
	return Host{
		Name:   h.Name,
		ID:     h.ID,
		Memory: h.Memory,
	}
}

func ToModel(h Host) hostModel.Host {
	return hostModel.Host{
		Name: h.Name,
		ID:   h.ID,
	}
}

type ListHostResponse struct {
	Items []Host `json:"items"`
	Total int    `json:"total"`
}

type GetHostResponse struct {
	Item Host `json:"item"`
}

type CreateHostRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateHostResponse struct {
	Item Host `json:"item"`
}

type DeleteHostRequest struct {
	Name string `json:"name"`
}

type DeleteHostResponse struct {
}

type UpdateHostRequest struct {
	Item Host `json:"item" binding:"required"`
}

type UpdateHostResponse struct {
	Item Host `json:"item"`
}

type BatchHostRequest struct {
	Operation string `json:"operation" binding:"required"`
	Items     []Host `json:"items"`
}

type BatchHostResponse struct {
	Items []Host `json:"items"`
}
