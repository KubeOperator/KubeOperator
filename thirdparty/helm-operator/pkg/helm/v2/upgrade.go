package v2

import (
	"github.com/pkg/errors"

	"k8s.io/helm/pkg/chartutil"
	helmv2 "k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

type releaseResponse interface {
	GetRelease() *release.Release
}

func (h *HelmV2) UpgradeFromPath(chartPath string, releaseName string, values []byte,
	opts helm.UpgradeOptions) (*helm.Release, error) {
	// Load the chart from the given path
	chartRequested, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load chart from path [%s] for release [%s]", chartPath, releaseName)
	}

	var res releaseResponse
	if opts.Install {
		res, err = h.client.InstallReleaseFromChart(
			chartRequested,
			opts.Namespace,
			helmv2.ReleaseName(releaseName),
			helmv2.ValueOverrides(values),
			helmv2.InstallDisableHooks(opts.DisableHooks),
			helmv2.InstallDryRun(opts.DryRun),
			helmv2.InstallWait(opts.Wait || opts.Atomic),
			helmv2.InstallTimeout(int64(opts.Timeout.Seconds())),
		)
		if err != nil && opts.Atomic {
			h.logger.Log("warning", "installation failed with atomic flag set, uninstalling release")
			_, dErr := h.client.DeleteRelease(releaseName,
				helmv2.DeletePurge(true), helmv2.DeleteDisableHooks(opts.DisableHooks))
			if dErr != nil {
				return nil, errors.Wrapf(statusMessageErr(dErr), "failed to uninstall release, original installation error: %s", statusMessageErr(err))
			}
		}
	} else {
		res, err = h.client.UpdateReleaseFromChart(
			releaseName,
			chartRequested,
			helmv2.UpdateValueOverrides(values),
			helmv2.UpgradeDisableHooks(opts.DisableHooks),
			helmv2.UpgradeDryRun(opts.DryRun),
			helmv2.UpgradeForce(opts.Force),
			helmv2.UpgradeRecreate(opts.Recreate),
			helmv2.ReuseValues(opts.ReuseValues),
			helmv2.ResetValues(opts.ResetValues),
			helmv2.UpgradeTimeout(int64(opts.Timeout.Seconds())),
			helmv2.UpgradeWait(opts.Wait || opts.Atomic),
		)
		if err != nil && opts.Atomic {
			h.logger.Log("warning", "upgrade failed with atomic flag set, rolling back release")
			_, rErr := h.client.RollbackRelease(releaseName,
				helmv2.RollbackTimeout(int64(opts.Timeout.Seconds())),
				helmv2.RollbackWait(opts.Wait),
				helmv2.RollbackDisableHooks(opts.DisableHooks),
				helmv2.RollbackDryRun(opts.DryRun),
				helmv2.RollbackRecreate(opts.Recreate),
				helmv2.RollbackForce(opts.Force))
			return nil, errors.Wrapf(statusMessageErr(rErr), "failed to roll back release, original installation error: %s", statusMessageErr(err))
		}
	}
	if err != nil {
		return nil, statusMessageErr(err)
	}
	return releaseToGenericRelease(res.GetRelease()), nil
}
