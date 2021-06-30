package hook

import "github.com/KubeOperator/KubeOperator/pkg/util/repo"

func init() {
	BeforeApplicationStart.AddFunc(loadRegistery)
}

func loadRegistery() error {
	repo.LoadRegistery()
	return nil
}
