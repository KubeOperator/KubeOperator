package initial

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	initPost = "09-post.yml"
)

type PostPhase struct {
}

func (s PostPhase) Name() string {
	return "Post Init"
}

func (s PostPhase) Run(b kobe.Interface) error {
	return phases.RunPlaybookAndGetResult(b, initPost)
}
