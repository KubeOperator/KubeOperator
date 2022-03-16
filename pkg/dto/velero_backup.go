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
	Cluster           string
	BackupAccountName string
	ID                string
}
