package plugin

import (
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
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

const (
	releasePluginDir = "/home/kops/config/plugin"
	localPluginDir   = "./plugin"
)

var pluginDirs = []string{
	localPluginDir,
	releasePluginDir,
}

type InitPluginDBPhase struct {
}

func (i *InitPluginDBPhase) Init() error {
	var log = logger.Default
	var p string
	for _, pa := range pluginDirs {
		if file.Exists(pa) {
			p = pa
		}
	}
	if p == "" {
		log.Info("can not find plugin dir,skip")
		return nil
	}
	fs, err := ioutil.ReadDir(p)
	if err != nil {
		return nil
	}
	for _, f := range fs {
		if !f.IsDir() && strings.Contains(f.Name(), ".so") {
			pluginName := strings.Replace(f.Name(), ".so", "", -1)
			p, err := plugin.Open(path.Join(p, f.Name()))
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
