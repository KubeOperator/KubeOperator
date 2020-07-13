package kotf

import (
	"fmt"
	"testing"
)

func TestKotfIint(t *testing.T) {
	terraform := NewTerraform(&Config{
		Cluster: "clsuter3",
	})

	// vsphere example
	provider := `{
    "name": "vSphere",
    "username": "",
    "password": "",
    "host": ""
  }`
	cloudRegion := `{
    "datacenter": "Datacenter",
    "zones": [
      {
        "key": "x65623",
        "name": "Resources",
        "network": "VM Network",
        "datastore": "vsanDatastore",
        "imageName": "Centos7.6-template",
        "guestId": "centos7_64Guest"
      }
    ]
  }`
	hosts := `[
    {
      "shortName": "worker2",
      "name": "worker2.clsuter3.fit2cloud.com",
      "cpu": 1,
      "memory": 2,
      "domain": "clsuter3.fit2cloud.com",
      "ip": "172.16.10.223",

      "zone": {
        "key": "x65623",
        "name": "Resources",
        "network": "VM Network",
        "datastore": "vsanDatastore",
        "imageName": "Centos7.6-template",
        "guestId": "centos7_64Guest",
		  "netMask": 24,
		  "gateway": "172.16.10.254"
      }
    }
  ]`

	// openstack example

	//provider := `{
	//  "username": "admin",
	//  "password": "",
	//  "identity": "",
	//  "projectId": "",
	//  "domainName": ""
	//}`

	//cloudRegion := `{
	//  "datacenter": "RegionOne",
	//}`

	//hosts := `[
	//  {
	//    "shortName": "worker2",
	//    "name": "worker2.clsuter3.fit2cloud.com",
	//    "ip": "172.16.10.223",
	//    "model": "m1.small",
	//    "zone": {
	//      "name": "nova",
	//      "network": "fd989e3e-3d49-4658-b852-b65bbd552220",
	//      "imageName": "cirros",
	//      "securityGroup": "default",
	//      "subnet": "a6019a57-2a14-4551-8331-46822b3d3981",
	//      "ipType": "private"
	//    }
	//  }
	//]`

	test, err := terraform.Init("vSphere", provider, cloudRegion, hosts)
	if err != nil {
		fmt.Println(err)
	}
	if test.Success {
		test3, err := terraform.Apply()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(test3)
	} else {
		fmt.Println(test)
	}

}
