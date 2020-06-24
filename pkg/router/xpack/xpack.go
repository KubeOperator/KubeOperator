package xpack

import (
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/plugin"
	"github.com/kataras/iris/v12"
)

var log = logger.Default

func XPack(parent iris.Party) {
	p := plugin.GetPlugin("xPack")
	if p != nil {
		f, err := p.Lookup("RouterRegister")
		if err != nil {
			log.Errorf("load xPack error: %s", err.Error())
		}
		fu, ok := f.(func(parent iris.Party))
		if !ok {
			log.Errorf("load xPack error: %s", ok)
		}
		fu(parent)
	}
}
