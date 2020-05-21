package docker

import "github.com/KubeOperator/KubeOperator/pkg/util/kobe"

const (
	installPlaybookName = "docker.yml"
)

func Install(kobe kobe.Interface) (string, error) {
	return kobe.RunPlaybook(installPlaybookName)
}
