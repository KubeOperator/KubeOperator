package ops

import (
	"ko3-gin/pkg/cluster"
	"ko3-gin/pkg/cluster/adm"
)

type ClusterJoinCreation struct {
	Cluster *cluster.Cluster
	Nodes   []cluster.Node
}

func (cjc *ClusterJoinCreation) Start() {
	for _, node := range cjc.Nodes {
		go func() {
			a := adm.ClusterAdm{Cluster: *cjc.Cluster}
			if err := a.JoinWorker(node.Host); err != nil {
				// 处理 JOIN 失败
			}
		}()
	}
}
