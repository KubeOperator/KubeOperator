package tools

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

type Interface interface {
	Install() error
}

func NewClusterTool(cluster dto.ClusterWithEndpoint, tool *model.ClusterTool) (Interface, error) {
	switch tool.Name {
	case "Prometheus":
		return NewPrometheus(cluster, tool)
	case "Logging":
		return NewLogging(cluster, tool)
	}
	return nil, nil
}
