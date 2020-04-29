package ops

import (
	"ko3-gin/pkg/cluster"
	"ko3-gin/pkg/host"
	"testing"
)

func TestClusterCreateCreation_Start(t *testing.T) {
	cc := ClusterCreateCreation{
		Config: ClusterConfig{Name: "test"},
		Nodes: []cluster.Node{
			{
				Host: host.Host{},
				Labels: map[string]string{
					"role": "worker",
				},
			},
			{
				Host: host.Host{},
				Labels: map[string]string{
					"role": "master",
				},
			},
		},
	}
	cc.Start()
}
