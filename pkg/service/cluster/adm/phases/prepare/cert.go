package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareCertificates = "05-certificates.yml"
)

type CertificatesPhase struct {
}

func (c CertificatesPhase) Name() string {
	return "GenerateCertificates"
}

func (c CertificatesPhase) Run(b kobe.Interface, fileName string) error {
	return phases.RunPlaybookAndGetResult(b, prepareCertificates, "", fileName)
}
