package client

type openStackClient struct {
	Vars map[string]interface{}
}

func NewOpenStackClient(vars map[string]interface{}) *openStackClient {
	return &openStackClient{
		Vars: vars,
	}
}

func (v *openStackClient) ListZones() string {
	return ""
}

func (v *openStackClient) ListDatacenter() ([]string, error) {
	return []string{}, nil
}
