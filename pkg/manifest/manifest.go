package manifest

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/ghodss/yaml"
	"io/ioutil"
)

var log = logger.Default

const (
	localPath   = "./manifest/manifest.yml"
	releasePath = "/usr/local/lib/ko/manifest.yml"
)

var supportedPath = []string{localPath, releasePath}

type Manifest struct {
	Name     string     `json:"name"`
	Category []Category `json:"category"`
}

type Category struct {
	Name  string `json:"name"`
	Items []Item `json:"items"`
}

type Item struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var Manifests []Manifest

type InitManifestPhase struct {
}

func (*InitManifestPhase) PhaseName() string {
	return "init manifest"
}

func (*InitManifestPhase) Init() error {
	var manifestPath = ""
	for _, p := range supportedPath {
		if fileutil.Exist(p) {
			manifestPath = p
			log.Debugf("using manifest: %s", localPath)
			break
		}
	}
	if manifestPath == "" {
		return errors.New("can not find manifest file")
	}
	bs, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bs, &Manifests)
	if err != nil {
		return err
	}
	return nil
}
