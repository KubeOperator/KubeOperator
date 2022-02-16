package velero

import (
	"fmt"
	"testing"
)

func TestVeleroBackup(t *testing.T) {

	args := []string{"--kubeconfig", "/Users/zk.wang/.kube/config"}

	//result, err := Backup("cluster-zhengkun", args)
	//result, err :=GetBackupLogs("cluster-zhengkun-2022-02-14-15-01",args)
	result, err := GetBackupDescribe("cluster-zhengkun-2022-02-14-15-01", args)
	if err != nil {
		fmt.Println("failed")
		fmt.Println(err)
	} else {
		fmt.Println("success")
		fmt.Println(result)
	}

}
