package v3

import (
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage/driver"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

func (h *HelmV3) Get(releaseName string, opts helm.GetOptions) (*helm.Release, error) {
	cfg, err := newActionConfig(h.kubeConfig, h.infoLogFunc(opts.Namespace, releaseName), opts.Namespace, "")
	if err != nil {
		return nil, err
	}

	get := action.NewGet(cfg)
	getOptions(opts).configure(get)

	res, err := get.Run(releaseName)
	switch err {
	case nil:
		return releaseToGenericRelease(res), nil
	case driver.ErrReleaseNotFound:
		return nil, nil
	default:
		return nil, errors.Wrapf(err, "failed to retrieve release '%s'", releaseName)
	}
}

type getOptions helm.GetOptions

func (opts getOptions) configure(action *action.Get) {
	action.Version = opts.Version
}
