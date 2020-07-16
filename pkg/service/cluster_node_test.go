package service

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"testing"
)

func TestClusterNodeService_Batch(t *testing.T) {
	Init()

	service := NewClusterNodeService()
	_, err := service.Batch("zxv", dto.NodeBatch{
		Hosts:     []string{},
		Nodes:     []string{},
		Increase:  1,
		Operation: constant.BatchOperationCreate,
	})
	if err != nil {
		log.Fatal(err)
	}

}
