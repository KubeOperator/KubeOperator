package plugin

import (
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"io/ioutil"
	"path"
	"plugin"
	"strings"
)

const (
	phaseName = "plugin"
)

var plugins = make(map[string]*plugin.Plugin)

func GetPlugin(name string) *plugin.Plugin {
	return plugins[name]
}

type InitPluginDBPhase struct {
}

func (i *InitPluginDBPhase) Init() error {
	var log = logger.Default
	fs, err := ioutil.ReadDir("pkg/plugin")
	if err != nil {
		return err
	}
	for _, f := range fs {
		if !f.IsDir() && strings.Contains(f.Name(), ".so") {
			pluginName := strings.Replace(f.Name(), ".so", "", -1)
			p, err := plugin.Open(path.Join("plugin", f.Name()))
			if err != nil {
				log.Errorf("can not load plugin: %s message: %s", pluginName, err.Error())
			} else {
				plugins[pluginName] = p
				log.Infof("load plugin: %s", pluginName)
			}
		}
	}
	return nil
}

func (i *InitPluginDBPhase) PhaseName() string {
	return phaseName
}
