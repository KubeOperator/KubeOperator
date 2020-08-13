package cloud_storage

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	vars := make(map[string]interface{})
	vars["type"] = "OSS"

	client, err := NewCloudStorageClient(vars)
	if err != nil {
		fmt.Println(err)
		return
	}
	//res,err := client.ListBuckets()
	//res,err := client.Exist("test")
	//res, err := client.Download("ceshi/test.replay", "/opt/download/test2.replay")
	//if err != nil {
	//	fmt.Println(err)
	//}
	result, err := client.Delete("")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
	//res,err :=client.Upload("/opt/download/test.replay","ceshi/test.replay")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//res,err := client.ListBuckets()
	//if err != nil {
	//	fmt.Println(err)
	//}

	//res, err := client.Delete("ceshi/test.replay")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//res ,err := client.Exist("kube-operator1")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)
}
