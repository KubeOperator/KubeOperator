package client

type VSphereClient struct {
	Vars map[string]interface{}
}

func NewVSphereClient(vars map[string]interface{}) *VSphereClient {
	return &VSphereClient{
		Vars: vars,
	}
}

func (v *VSphereClient) listZones() string {
	return ""
}
