package phase

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	PlaybookNameBase = "01-base.yml"
)

type SystemConfigPhase struct {
}

func (s SystemConfigPhase) Name() string {
	return "ConfigSystem"
}

func (s SystemConfigPhase) Run(b kobe.Interface) error {
	taskId, err := b.RunPlaybook(PlaybookNameBase)
	if err != nil {
		return err
	}
	var res kobe.Result
	return wait.Poll(phases.PhaseInterval, phases.PhaseTimeout, func() (done bool, err error) {
		result, err := b.GetResult(taskId)
		if err != nil {
			return true, err
		}
		if result.Finished {
			if result.Success {
				res, err = kobe.ParseResult(result.Content)
				if err != nil {
					return true, err
				}
			}
			return true, errors.New(result.Message)
		}
		return false, nil
	})
}
