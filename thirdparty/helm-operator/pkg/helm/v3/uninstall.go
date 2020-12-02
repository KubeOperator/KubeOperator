package v3

import (
	"github.com/fluxcd/helm-operator/pkg/helm"
	"github.com/pkg/errors"

	"helm.sh/helm/v3/pkg/action"
)

func (h *HelmV3) Uninstall(releaseName string, opts helm.UninstallOptions) error {
	cfg, err := newActionConfig(h.kubeConfig, h.infoLogFunc(opts.Namespace, releaseName), opts.Namespace, "")
	if err != nil {
		return err
	}

	uninstall := action.NewUninstall(cfg)
	uninstallOptions(opts).configure(uninstall)

	if _, err := uninstall.Run(releaseName); err != nil {
		return errors.Wrapf(err, "failed to uninstall release '%s'", releaseName)
	}
	return nil
}

type uninstallOptions helm.UninstallOptions

func (opts uninstallOptions) configure(action *action.Uninstall) {
	action.DisableHooks = opts.DisableHooks
	action.DryRun = opts.DryRun
	action.KeepHistory = opts.KeepHistory
	action.Timeout = opts.Timeout
}
