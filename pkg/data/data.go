package data

import (
	"os"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"k8s.io/kubernetes/third_party/forked/etcd237/pkg/fileutil"
)

var initDirs = []string{
	constant.DefaultDataDir,
}

const phaseName = "create data dir"

type InitDataPhase struct{}

func (i *InitDataPhase) Init() error {
	for _, d := range initDirs {
		if !fileutil.Exist(d) {
			err := os.MkdirAll(d, 0750)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func (i *InitDataPhase) PhaseName() string {
	return phaseName
}
