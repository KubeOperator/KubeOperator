package helm

import "github.com/KubeOperator/KubeOperator/pkg/util/helm"

const phaseName = "init helm"

type InitHelmPhase struct {
	RepoName     string
	RepoUrl      string
	RepoUsername string
	RepoPassword string
}

func (i *InitHelmPhase) Init() error {

	rs, err := helm.ListRepo()
	if err != nil {
		return err
	}
	for _, r := range rs {
		if r.Name == i.RepoName {
			return nil
		}
	}
	err = helm.AddRepo(i.RepoName, i.RepoUrl, i.RepoUsername, i.RepoPassword)
	if err != nil {
		return err
	}
	return nil
}

func (i *InitHelmPhase) PhaseName() string {
	return phaseName
}
