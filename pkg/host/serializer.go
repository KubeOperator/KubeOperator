package host

type CreatHostRequest struct {
	Name           string `json:"name"`
	Ip             string `json:"ip"`
	Port           string `json:"port"`
	CredentialName string `json:"credential_name"`
}

type UpdateHostRequest struct {
	Id             string
	Name           string `json:"name"`
	Ip             string `json:"ip"`
	Port           string `json:"port"`
	CredentialName string `json:"credential_name"`
}

type PageHostResponse struct {
	Total int
	Items []Host
}
