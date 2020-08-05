package backup_acoount

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	vars := make(map[string]string)
	////vars["type"] = "OSS"
	////vars["endpoint"] = "http://oss-cn-hangzhou.aliyuncs.com"
	////vars["accessKey"] = "LTAI4Fd1E5YrrhsBjf7zrEBp"
	////vars["secretKey"] = "Z9iGtDe3qK4VUlBSHuAV9XvCcQ3r6c"
	////vars["bucket"]  = "kube-operator"
	//
	////vars["type"] = "AZURE"
	////vars["endpoint"] = "blob.core.chinacloudapi.cn"
	////vars["accountName"] = "zhengkun"
	////vars["accountKey"] = "UEcy3YKS46t/BXsKm8vAOY3tLkll1+EmJzr9GMekz8tniCJybmFnqVWAoC9WvlTzvPnz9rpfts8J2ma7nemQSg=="
	////vars["bucket"] = "test"
	//
	//vars["type"] = "S3"
	//vars["endpoint"] = "http://s3.cn-north-1.amazonaws.com.cn"
	//vars["accessKey"] = "AKIAWZUDY4HG4JJH4TXN"
	//vars["secretKey"] = "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"
	//vars["region"] = "cn-north-1"
	//vars["bucket"] = "kube-operator"
	client, err := NewBackupAccountClient(vars)
	if err != nil {
		fmt.Println(err)
		return
	}
	//res,err := client.ListBuckets()
	//res,err := client.Exist("test")
	//res,err :=client.Download("ceshi/test.replay","/opt/download/test2.replay")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//res,err :=client.Upload("/opt/download/test.replay","ceshi/test.replay")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//res,err := client.ListBuckets()
	//if err != nil {
	//	fmt.Println(err)
	//}

	res, err := client.Delete("ceshi/test.replay")
	if err != nil {
		fmt.Println(err)
	}
	//res ,err := client.Exist("kube-operator1")
	//if err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println(res)
}
