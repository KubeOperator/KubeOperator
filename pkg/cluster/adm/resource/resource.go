package resource

import (
	"errors"
	"fmt"
	"github.com/thoas/go-funk"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/cluster/adm/spec"
	"ko3-gin/pkg/util/ssh"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	Docker = Package{
		Name:     "docker",
		Versions: spec.DockerVersions,
	}
	CNIPlugins = Package{
		Name:     "cni-plugins",
		Versions: spec.CNIPluginsVersions,
	}

	Kubeadm = Package{
		Name:     "kubeadm",
		Versions: spec.KubeadmVersions,
	}
	KubernetesNode = Package{
		Name:     "kubernetes-node",
		Versions: spec.K8sVersionsWithV,
	}
)

type Package struct {
	Name     string
	Versions []string
}

// CopyToNode copy package which use default version to node and return dst filename
func (p *Package) CopyToNodeWithDefault(s ssh.Interface) (string, error) {
	return p.CopyToNode(s, p.DefaultVersion())
}

// CopyToNode copy package which use default version to node and return dst filename
func (p *Package) CopyToNode(s ssh.Interface, version string) (string, error) {
	srcFile, err := p.ResourceForNode(s, version)
	if err != nil {
		return "", err
	}
	dstFile := path.Join(constants.TmpDir, filepath.Base(srcFile))
	err = s.CopyFile(srcFile, dstFile)
	if err != nil {
		return "", err
	}
	return dstFile, nil
}

func (p *Package) ResourceForNode(s ssh.Interface, version string) (string, error) {
	return p.Resource(Arch(s), version)
}

func (p *Package) Resource(arch, version string) (string, error) {
	version, err := p.NormalizeVersion(version)
	if err != nil {
		return "", err
	}
	basename := fmt.Sprintf("linux-%s/%s-linux-%s-%s.tar.gz", arch, p.Name, arch, version)
	srcFile := path.Join(constants.SrcDir, basename)
	if _, err := os.Stat(srcFile); err != nil {
		return "", err
	}
	return srcFile, nil
}

func (p *Package) DefaultVersion() string {
	return p.Versions[0]
}

func (p *Package) NormalizeVersion(version string) (string, error) {
	if p.Versions[0][0] == 'v' && version[0] != 'v' {
		version = "v" + version
	} else if p.Versions[0][0] != 'v' && version[0] == 'v' {
		version = version[1:]
	}

	if funk.ContainsString(p.Versions, version) {
		return version, nil
	}

	return "", errors.New("invalid version")
}

func Arch(s ssh.Interface) string {
	var arch string

	stdout, _, _, _ := s.Exec("arch")
	switch strings.TrimSpace(stdout) {
	case "x86_64":
		arch = "amd64"
	case "aarch64":
		arch = "arm64"
	}

	return arch
}