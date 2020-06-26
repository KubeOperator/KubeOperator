package client

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
)

type CloudClient interface {
	listZones() string
}

func NewCloudClient(vars string) CloudClient {
	var cloudVars map[string]interface{}
	if err := json.Unmarshal([]byte(vars), &cloudVars); err == nil {
		if cloudVars["provider"] == constant.OpenStack {
			return NewOpenStackClient(cloudVars)
		}
		if cloudVars["provider"] == constant.VSphere {
			return NewVSphereClient(cloudVars)
		}
	}
	return nil
}
