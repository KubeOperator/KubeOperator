package cloud_provider

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	vars := make(map[string]interface{})
	vars["provider"] = "FusionCompute"
	vars["datacenter"] = "site"
	vars["server"] = ""
	vars["user"] = ""
	vars["password"] = ""

	client := NewCloudClient(vars)
	ips, err := client.GetIpInUsed("")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(ips)
	}

}
