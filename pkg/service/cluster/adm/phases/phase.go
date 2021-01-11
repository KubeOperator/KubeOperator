package phases

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/spf13/viper"
	"io"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

const (
	PhaseInterval             = 5 * time.Second
	DefaultPhaseTimeoutMinute = 100
)

var log = logger.Default

type Interface interface {
	Name() string
	Run(p kobe.Interface, writer io.Writer) error
}

func RunPlaybookAndGetResult(b kobe.Interface, playbookName, tag string, writer io.Writer) error {
	taskId, err := b.RunPlaybook(playbookName, tag)
	var result kobe.Result
	if err != nil {
		return err
	}
	// 读取 ansible 执行日志
	if writer != nil {
		go func() {
			err = b.Watch(writer, taskId)
			if err != nil {
				log.Error(err)
			}
		}()
	}
	timeout := viper.GetInt("job.timeout")
	if timeout < DefaultPhaseTimeoutMinute {
		timeout = DefaultPhaseTimeoutMinute
	}
	err = wait.Poll(PhaseInterval, time.Duration(timeout)*time.Minute, func() (done bool, err error) {
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
