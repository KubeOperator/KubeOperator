package cloud_storage

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	vars := make(map[string]interface{})
	vars["type"] = "OSS"
	vars["endpoint"] = "http://oss-cn-hangzhou.aliyuncs.com"
	vars["accessKey"] = "LTAI4GHUhAMVs5niWf59aFst"
	vars["secretKey"] = "ZPoxtXdV9HhzZcBHpLaURaL2Mxl8dN"
	vars["bucket"] = "kube-operator"

	client, err := NewCloudStorageClient(vars)
	if err != nil {
		fmt.Println(err)
		return
	}
	//res,err := client.ListBuckets()
	//res,err := client.Exist("test")
	res, err := client.Download("ceshi/test.replay", "/opt/download/test2.replay")
	if err != nil {
		fmt.Println(err)
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
	fmt.Println(res)
}
