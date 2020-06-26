package client

type OpenStackClient struct {
	Vars map[string]interface{}
}

func NewOpenStackClient(vars map[string]interface{}) *OpenStackClient {
	return &OpenStackClient{
		Vars: vars,
	}
}

func (v *OpenStackClient) listZones() string {
	return ""
}
