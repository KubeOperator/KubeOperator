package ops

import (
	"ko3-gin/pkg/cluster"
	"ko3-gin/pkg/cluster/adm"
	"ko3-gin/pkg/cluster/util"
	"time"
)

type ClusterConfig struct {
	Name string
}
type ClusterCreateCreation struct {
	Config ClusterConfig
	Nodes  []cluster.Node
}

func (cc *ClusterCreateCreation) Start() string {
	c := cluster.Cluster{
		Name: cc.Config.Name,
	}
	go func(c *cluster.Cluster) {
		ad := adm.ClusterAdm{Cluster: *c}
		n := util.NewNodeSelector(cc.Nodes)
		deployNode := n.SelectFistNode("role", "master")
		if err := ad.Init(deployNode.Host); err != nil {
			// 处理初始化错误
			return
		}
		workerNodes := n.SelectNodes("role", "worker")
		cjc := ClusterJoinCreation{
			Cluster: c,
			Nodes:   workerNodes,
		}
		cjc.Start()
	}(&c)
	time.Sleep(3 * time.Minute)
	return c.Name
}
