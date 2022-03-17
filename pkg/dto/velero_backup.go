package dto

type VeleroBackup struct {
	Name                    string
	Cluster                 string
	IncludeNamespaces       []string
	ExcludeNamespaces       []string
	IncludeResources        []string
	ExcludeResources        []string
	IncludeClusterResources bool
	Labels                  string
	Selector                string
	Ttl                     string
	Schedule                string
	BackupName              string
}

type VeleroInstall struct {
	Cluster           string        `json:"cluster"`
	BackupAccountName string        `json:"backupAccountName"`
	ID                string        `json:"id"`
	Limits            ResourceQuota `json:"limits"`
	Requests          ResourceQuota `json:"requests"`
}

type ResourceQuota struct {
	Cpu    int `json:"cpu"`
	Memory int `json:"memory"`
}
