package dto

import (
	v1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
)

type SourceSearch struct {
	Cluster   string `json:"cluster"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Limit     int64  `json:"limit"`
	Continue  string `json:"continue"`
}

type SourceList struct {
	Kind       string        `json:"kind"`
	ApiVersion string        `json:"apiVersion"`
	Metadata   string        `json:"metadata"`
	Items      []interface{} `json:"items"`
}

type SourceNsCreate struct {
	Cluster   string       `json:"cluster"`
	Kind      string       `json:"kind"`
	Namespace string       `json:"namespace"`
	Info      v1.Namespace `json:"info"`
}
type SourceScCreate struct {
	Cluster   string                 `json:"cluster"`
	Kind      string                 `json:"kind"`
	Namespace string                 `json:"namespace"`
	Info      storagev1.StorageClass `json:"info"`
}
type SourceSecretCreate struct {
	Cluster   string    `json:"cluster"`
	Kind      string    `json:"kind"`
	Namespace string    `json:"namespace"`
	Info      v1.Secret `json:"info"`
}
type SourcePvCreate struct {
	Cluster   string              `json:"cluster"`
	Kind      string              `json:"kind"`
	Namespace string              `json:"namespace"`
	Info      v1.PersistentVolume `json:"info"`
}

type SourceDelete struct {
	Cluster   string `json:"cluster"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
