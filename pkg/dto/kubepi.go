package dto

type BindKubePI struct {
	SourceType   string `json:"sourceType"`
	Project      string `json:"project"`
	Cluster      string `json:"cluster"`
	BindUser     string `json:"bindUser"`
	BindPassword string `json:"bindPassword"`
}

type CheckConn struct {
	BindUser     string `json:"bindUser"`
	BindPassword string `json:"bindPassword"`
}

type SearchBind struct {
	SourceType string `json:"sourceType"`
	Project    string `json:"project"`
	Cluster    string `json:"cluster"`
}

type BindResponse struct {
	SourceType string `json:"sourceType"`
	Project    string `json:"project"`
	Cluster    string `json:"cluster"`
	BindUser   string `json:"bindUser"`
}
