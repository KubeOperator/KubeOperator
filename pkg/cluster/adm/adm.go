package adm

import (
	"ko3-gin/pkg/cluster"
	"ko3-gin/pkg/cluster/adm/workflow"
	commonphases "ko3-gin/pkg/cluster/adm/workflow/phases/common"
	initphases "ko3-gin/pkg/cluster/adm/workflow/phases/init"
	joinphases "ko3-gin/pkg/cluster/adm/workflow/phases/join"
	"ko3-gin/pkg/host"
)

type Interface interface {
	Init(host host.Host) error
	JoinWorker(host host.Host) error
}

type ClusterAdm struct {
	Cluster cluster.Cluster
}

type InitData struct {
}

type JoinData struct {
}

func (ca *ClusterAdm) Init(host host.Host) error {
	runner := workflow.NewRunner()
	runner.AppendRunner(commonphases.NewKubeletInstallPhase())
	runner.AppendRunner(commonphases.NewKubeletStartPhase())
	runner.AppendRunner(initphases.NewCertsPhase())
	runner.AppendRunner(initphases.NewKubeConfigPhase())
	//runner.AppendRunner(initphases.NewControlPlanePhase())
	//runner.AppendRunner(initphases.NewEtcdPhase())
	//runner.AppendRunner(initphases.NewWaitControlPlanePhase())
	//runner.AppendRunner(initphases.NewMarkControlPlanePhase())
	//runner.AppendRunner(initphases.NewBootstrapTokenPhase())
	data := InitData{}
	if err := runner.Run(data, host); err != nil {
		return err
	}
	return nil
}

func (ca *ClusterAdm) JoinWorker(host host.Host) error {
	runner := workflow.NewRunner()
	runner.AppendRunner(commonphases.NewKubeletInstallPhase())
	runner.AppendRunner(joinphases.NewControlPlanePreparePhase())
	runner.AppendRunner(commonphases.NewKubeletStartPhase())
	runner.AppendRunner(joinphases.NewControlPlaneJoinPhase())
	data := JoinData{}
	if err := runner.Run(data, host); err != nil {
		return err
	}
	return nil
}
