package data

import (
	"os"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
)

var initDirs = []string{
	constant.DefaultDataDir,
}

const phaseName = "create data dir"

type InitDataPhase struct{}

func (i *InitDataPhase) Init() error {
	for _, d := range initDirs {
		if !file.Exists(d) {
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
