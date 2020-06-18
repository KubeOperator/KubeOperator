package service

import (
	"testing"
)

func TestClusterInitService_Init(t *testing.T) {
	Init()
	service := NewClusterInitService()
	err := service.Init("zvc")
	if err != nil {
		t.Error(err)
	}
}

func TestGatherKubernetesToken(t *testing.T) {
	Init()
	cs := NewClusterService()
	cis := NewClusterInitService()
	c, err := cs.Get("test")
	if err != nil {
		t.Error(err)
	}
	err = cis.GatherKubernetesToken(c.Cluster)
	if err != nil {
		t.Error(err)
	}
}
