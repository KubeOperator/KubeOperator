package storage

import (
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
)

const oceanStor = "10-plugin-cluster-storage-oceanstor.yml"

type OceanStorPhase struct {
	OceanStorType           string
	OceanstorProduct        string
	OceanstorURLs           string
	OceanstorUser           string
	OceanstorPassword       string
	OceanstorPools          string
	OceanstorPortal         string
	OceanstorControllerType string
	OceanstorIsMultipath    string
}

func (o OceanStorPhase) Name() string {
	return "CrateOceanStorStorage"
}

func (o OceanStorPhase) Run(b kobe.Interface, writer io.Writer) error {
	if o.OceanStorType != "" {
		b.SetVar("oceanstor_type", o.OceanStorType)
	}
	if o.OceanstorProduct != "" {
		b.SetVar("oceanstor_product", o.OceanstorProduct)
	}
	if o.OceanstorURLs != "" {
		b.SetVar("oceanstor_urls", o.OceanstorURLs)
	}
	if o.OceanstorUser != "" {
		b.SetVar("oceanstor_user", o.OceanstorUser)
	}
	if o.OceanstorPassword != "" {
		b.SetVar("oceanstor_password", o.OceanstorPassword)
	}
	if o.OceanstorPools != "" {
		b.SetVar("oceanstor_pools", o.OceanstorPools)
	}
	if o.OceanstorPortal != "" {
		b.SetVar("oceanstor_portal", o.OceanstorPortal)
	}
	if o.OceanstorControllerType != "" {
		b.SetVar("oceanstor_controller_type", o.OceanstorControllerType)
	}
	if o.OceanstorIsMultipath != "" {
		b.SetVar("oceanstor_is_multipath", o.OceanstorIsMultipath)
	}
	return phases.RunPlaybookAndGetResult(b, oceanStor, "", writer)
}
