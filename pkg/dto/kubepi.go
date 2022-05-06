package dto

type BindKubePI struct {
	SourceType   string `json:"sourceType"`
	Source       string `json:"source"`
	BindUser     string `json:"bindUser"`
	BindPassword string `json:"bindPassword"`
}

type CheckConn struct {
	BindUser     string `json:"bindUser"`
	BindPassword string `json:"bindPassword"`
}

type SearchBind struct {
	SourceType string `json:"sourceType"`
	Source     string `json:"source"`
}

type BindResponse struct {
	SourceType string `json:"sourceType"`
	Source     string `json:"source"`
	BindUser   string `json:"bindUser"`
}
