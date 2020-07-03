package phases

import (
	"encoding/json"
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

func RunPlaybookAndGetResult(b kobe.Interface, playbookName string) error {
	taskId, err := b.RunPlaybook(playbookName)
	var result kobe.Result
	if err != nil {
		return err
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
					if result.HostFailedInfo != nil && len(result.HostFailedInfo) > 0 {
						by, _ := json.Marshal(&result.HostFailedInfo)
						return true, errors.New(string(by))
					}
				}
			}
			return true, nil
		}
		return false, nil
	})
	return err
}
