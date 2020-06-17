package prepare

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/facts"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareCertificates = "05-certificates.yml"
)

type CertificatesPhase struct {
	CertsExpired string
}

func (c CertificatesPhase) Name() string {
	return "GenerateCertificates"
}

func (c CertificatesPhase) Run(b kobe.Interface) error {
	if c.CertsExpired != "" {
		b.SetVar(facts.CertsExpiredFactName, c.CertsExpired)
	}
	return phases.RunPlaybookAndGetResult(b, prepareCertificates)
}
