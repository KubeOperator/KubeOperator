package velero

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func TestVeleroBackup(t *testing.T) {

	//args := []string{"--kubeconfig", "/Users/zk.wang/.kube/config"}
	//
	////result, err := Backup("cluster-zhengkun", args)
	////result, err :=GetBackupLogs("cluster-zhengkun-2022-02-14-15-01",args)
	//result, err := GetBackupDescribe("cluster-zhengkun-2022-02-14-15-01", args)

	//args := []string{"y|","/usr/local/bin/velero","delete","backup"}
	//result,err := ExecCommand("/bin","echo",args)
	//if err != nil {
	//	fmt.Println("failed")
	//	fmt.Println(err)
	//} else {
	//	fmt.Println("success")
	//	fmt.Println(string(result))
	//}

	args := []string{"-c", "echo y| /usr/local/bin/velero", "delete backup backup-1"}

	cmd := exec.Command("/bin/bash", args...)
	//cmd := exec.Command("/bin/bash","-c","echo y| /usr/local/bin/velero delete backup backup-1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("failed")
		fmt.Println(err)
	}
	cmd.Stderr = cmd.Stdout
	if err = cmd.Start(); err != nil {
		fmt.Println("failed")
		fmt.Println(err)
		return
	}

	var buffer bytes.Buffer
	for {
		out := make([]byte, 1024)
		length, err := stdout.Read(out)
		if err != nil {
			break
		}
		if length > 0 {
			buffer.Write(out[:length])
		}
	}

	if err = cmd.Wait(); err != nil {
		fmt.Println("failed")
		fmt.Println(buffer.String())
	}
	fmt.Println("success")
	fmt.Println(buffer.String())
}
