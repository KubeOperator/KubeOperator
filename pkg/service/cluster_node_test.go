package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"testing"
)

func TestClusterNodeService_Batch(t *testing.T) {
	Init()

	service := NewClusterNodeService()
	_, _ = service.Batch("zidong", dto.NodeBatch{
		Nodes:     []string{"worker-2"},
		Operation: "delete",
	})

}
