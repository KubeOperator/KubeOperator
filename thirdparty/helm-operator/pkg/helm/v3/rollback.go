package v3

import (
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

func (h *HelmV3) Rollback(releaseName string, opts helm.RollbackOptions) (*helm.Release, error) {
	cfg, err := newActionConfig(h.kubeConfig, h.infoLogFunc(opts.Namespace, releaseName), opts.Namespace, "")
	if err != nil {
		return nil, err
	}

	rollback := action.NewRollback(cfg)
	rollbackOptions(opts).configure(rollback)

	if err := rollback.Run(releaseName); err != nil {
		return nil, errors.Wrapf(err, "failed to perform rollback for release '%s'", releaseName)
	}

	// As rolling back does no longer return information about
	// the release in v3, we need to make an additional call to
	// get the release we rolled back to.
	return h.Get(releaseName, helm.GetOptions{Namespace: opts.Namespace})
}

type rollbackOptions helm.RollbackOptions

func (opts rollbackOptions) configure(action *action.Rollback) {
	action.Timeout = opts.Timeout
	action.Version = opts.Version
	action.Wait = opts.Wait
	action.DisableHooks = opts.DisableHooks
	action.DryRun = opts.DryRun
	action.Recreate = opts.Recreate
}
