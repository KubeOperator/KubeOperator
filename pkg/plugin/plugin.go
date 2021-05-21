package plugin

import (
	"io/ioutil"
	"path"
	"plugin"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/util/file"
)

const (
	phaseName = "plugin"
)

var plugins = make(map[string]*plugin.Plugin)

func GetPlugin(name string) *plugin.Plugin {
	return plugins[name]
}

const (
	releasePluginDir = "/usr/local/lib/ko/plugin"
	localPluginDir   = "./plugin"
)

var pluginDirs = []string{
	localPluginDir,
	releasePluginDir,
}

type InitPluginDBPhase struct {
}

func (i *InitPluginDBPhase) Init() error {
	var p string
	for _, pa := range pluginDirs {
		if file.Exists(pa) {
			p = pa
		}
	}
	if p == "" {
		logger.Log.Info("can not find plugin dir,skip")
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
				logger.Log.Errorf("can not load plugin: %s message: %s", pluginName, err.Error())
			} else {
				plugins[pluginName] = p
				logger.Log.Infof("load plugin: %s", pluginName)
			}
		}
	}
	return nil
}

func (i *InitPluginDBPhase) PhaseName() string {
	return phaseName
}
