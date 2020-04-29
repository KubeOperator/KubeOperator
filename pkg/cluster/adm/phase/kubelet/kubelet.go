package kubelet

import (
	"fmt"
	"ko3-gin/pkg/cluster/adm/util"
	"ko3-gin/pkg/util/ssh"
)

func TryStartKubelet(ssh *ssh.SSH) {
	initSystem, err := util.GetInitSystem(ssh)
	if err != nil {
		fmt.Println("[kubelet-start] no supported init system detected, won't make sure the kubelet is running properly.")
		return
	}

	if !initSystem.ServiceExists("kubelet") {
		fmt.Println("[kubelet-start] couldn't detect a kubelet service, can't make sure the kubelet is running properly.")
	}

	if err := initSystem.ServiceRestart("kubelet"); err != nil {
		fmt.Printf("[kubelet-start] WARNING: unable to start the kubelet service: [%v]\n", err)
		fmt.Printf("[kubelet-start] Please ensure kubelet is reloaded and running manually.\n")
	}
}

func TryStopKubelet(ssh *ssh.SSH) {
	initSystem, err := util.GetInitSystem(ssh)
	if err != nil {
		fmt.Println("[kubelet-start] no supported init system detected, won't make sure the kubelet not running for a short period of time while setting up configuration for it.")
		return
	}

	if !initSystem.ServiceExists("kubelet") {
		fmt.Println("[kubelet-start] couldn't detect a kubelet service, can't make sure the kubelet not running for a short period of time while setting up configuration for it.")
	}

	if err := initSystem.ServiceStop("kubelet"); err != nil {
		fmt.Printf("[kubelet-start] WARNING: unable to stop the kubelet service momentarily: [%v]\n", err)
	}
}

func TryRestartKubelet(ssh *ssh.SSH) {
	initSystem, err := util.GetInitSystem(ssh)
	if err != nil {
		fmt.Println("[kubelet-start] no supported init system detected, won't make sure the kubelet not running for a short period of time while setting up configuration for it.")
		return
	}

	if !initSystem.ServiceExists("kubelet") {
		fmt.Println("[kubelet-start] couldn't detect a kubelet service, can't make sure the kubelet not running for a short period of time while setting up configuration for it.")
	}

	if err := initSystem.ServiceRestart("kubelet"); err != nil {
		fmt.Printf("[kubelet-start] WARNING: unable to restart the kubelet service momentarily: [%v]\n", err)
	}
}
