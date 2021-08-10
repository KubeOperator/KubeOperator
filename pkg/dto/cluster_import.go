package dto

type ClusterImport struct {
	Name          string `json:"name"`
	ApiServer     string `json:"apiServer"`
	Router        string `json:"router"`
	Token         string `json:"token"`
	ProjectName   string `json:"projectName"`
	Architectures string `json:"architectures"`
}
