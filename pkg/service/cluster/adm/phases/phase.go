package phases

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

const (
	PhaseInterval = 2 * time.Second
	PhaseTimeout  = 10 * time.Minute
)

type Interface interface {
	Name() string
	Run(p kobe.Interface) error
}

func RunPlaybookAndGetResult(b kobe.Interface, playbookName string) (res kobe.Result, err error) {
	taskId, err := b.RunPlaybook(playbookName)
	if err != nil {
		return
	}
	err = wait.Poll(PhaseInterval, PhaseTimeout, func() (done bool, err error) {
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
	return
}
