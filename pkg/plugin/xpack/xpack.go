package xpack

import (
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/plugin"
	"github.com/pkg/errors"
)

var log = logger.Default

func LoadXpackPlugin() error {
	p := plugin.GetPlugin("xPack")
	if p != nil {
		f, err := p.Lookup("XpackRegister")
		if err != nil {
			log.Errorf("load xPack error: %s", err.Error())
		}
		fu, ok := f.(func() error)
		if !ok {
			log.Errorf("load xPack error: %v", ok)
		}
		if err := fu(); err != nil {
			return errors.Wrap(err, "register xpack err")
		}
	}
	return nil
}
