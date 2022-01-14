package dto

type ClusterImport struct {
	Name        string `json:"name" validate:"clustername,required"`
	ApiServer   string `json:"apiServer" validate:"required"`
	Router      string `json:"router" validate:"required"`
	Token       string `json:"token" validate:"required"`
	ProjectName string `json:"projectName" validate:"koname,required,max=30"`
}
