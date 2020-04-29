package phases

import (
	"fmt"
	"ko3-gin/pkg/cluster/adm/workflow"
	"ko3-gin/pkg/host"
)

func NewKubeletStartPhase() workflow.Phase {
	return workflow.Phase{
		Name: "kubelet-start",
		Run:  runKubeletStart,
	}
}

func runKubeletStart(data workflow.RunData, host host.Host) error {
	fmt.Println("starting kubelet")
	return nil
}

func NewKubeletInstallPhase() workflow.Phase {
	return workflow.Phase{
		Name: "install-kubelet",
		Run:  runKubeletInstall,
	}
}

func runKubeletInstall(data workflow.RunData, host host.Host) error {
	fmt.Println("install kubelet...")
	return nil
}
