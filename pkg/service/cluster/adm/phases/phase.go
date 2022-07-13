package phases

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	PhaseInterval             = 5 * time.Second
	DefaultPhaseTimeoutMinute = 10
)

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
				logger.Log.Error(err)
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

func WaitForDeployRunning(namespace string, deploymentName string, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.Poll(5*time.Second, 2*time.Minute, func() (done bool, err error) {
		d, err := kubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.ReadyReplicas > 0 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func WaitForDaemonsetRunning(namespace string, daemonsetName string, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.Poll(5*time.Second, 2*time.Minute, func() (done bool, err error) {
		d, err := kubeClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonsetName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.NumberReady > 0 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func WaitForStatefulSetsRunning(namespace string, statefulSetsName string, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.Poll(5*time.Second, 2*time.Minute, func() (done bool, err error) {
		d, err := kubeClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetsName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.ReadyReplicas > 0 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
