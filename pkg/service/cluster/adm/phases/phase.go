package phases

import (
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
	Run(p kobe.Interface) (kobe.Result, error)
}

func RunPlaybookAndGetResult(b kobe.Interface, playbookName string) (result kobe.Result, err error) {
	taskId, err := b.RunPlaybook(playbookName)
	if err != nil {
		return
	}
	err = wait.Poll(PhaseInterval, PhaseTimeout, func() (done bool, err error) {
		res, err := b.GetResult(taskId)
		if err != nil {
			return true, err
		}
		if res.Finished {
			if res.Success {
				result, err = kobe.ParseResult(res.Content)
				if err != nil {
					return true, err
				}
			} else {
				if res.Content != "" {
					result, err = kobe.ParseResult(res.Content)
					if err != nil {
						return true, err
					}
					result.GatherFailedInfo()
				}
			}
			return true, nil
		}
		return false, nil
	})
	return
}
