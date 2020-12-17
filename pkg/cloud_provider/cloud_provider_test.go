package cloud_provider

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	vars := make(map[string]interface{})
	vars["provider"] = "FusionCompute"
	vars["datacenter"] = "site"
	vars["server"] = "https://100.199.16.208:7443"
	vars["user"] = "kubeoperator"
	vars["password"] = "Calong@2015"

	client := NewCloudClient(vars)
	ips, err := client.GetIpInUsed("")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(ips)
	}

}
