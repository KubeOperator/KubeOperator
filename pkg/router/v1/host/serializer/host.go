package serializer

import (
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	"time"
)

type Host struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Memory       int       `json:"memory"`
	CpuCore      int       `json:"cpuCore"`
	Os           string    `json:"os"`
	OsVersion    string    `json:"osVersion"`
	GpuNum       int       `json:"gpuNum"`
	GpuInfo      string    `json:"gpuInfo"`
	Ip           string    `json:"ip"`
	Port         int       `json:"port"`
	CredentialID string    `json:"credentialId"`
	ClusterID    string    `json:"clusterId"`
	NodeID       string    `json:"nodeId"`
	Status       string    `json:"status"`
	CreateAt     time.Time `json:"createAt"`
	UpdateAt     time.Time `json:"updateAt"`
}

func FromModel(h hostModel.Host) Host {
	return Host{
		Name:         h.Name,
		ID:           h.ID,
		Memory:       h.Memory,
		CpuCore:      h.CpuCore,
		Os:           h.Os,
		OsVersion:    h.OsVersion,
		Ip:           h.Ip,
		Port:         h.Port,
		CredentialID: h.CredentialID,
		NodeID:       h.NodeID,
		Status:       h.Status,
		GpuNum:       h.GpuNum,
		GpuInfo:      h.GpuInfo,
		CreateAt:     h.CreatedAt,
		UpdateAt:     h.UpdatedAt,
	}
}

func ToModel(h Host) hostModel.Host {
	return hostModel.Host{
		Name:         h.Name,
		ID:           h.ID,
		Memory:       h.Memory,
		CpuCore:      h.CpuCore,
		Os:           h.Os,
		OsVersion:    h.OsVersion,
		Ip:           h.Ip,
		Port:         h.Port,
		CredentialID: h.CredentialID,
		NodeID:       h.NodeID,
		Status:       h.Status,
		GpuNum:       h.GpuNum,
		GpuInfo:      h.GpuInfo,
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
	Name         string `json:"name" binding:"required"`
	Ip           string `json:"ip"  binding:"required"`
	Port         int    `json:"port"  binding:"required"`
	CredentialID string `json:"credentialId" binding:"required"`
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
