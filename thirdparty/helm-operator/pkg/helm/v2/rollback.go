package v2

import (
	helmv2 "k8s.io/helm/pkg/helm"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

func (h *HelmV2) Rollback(releaseName string, opts helm.RollbackOptions) (*helm.Release, error) {
	res, err := h.client.RollbackRelease(
		releaseName,
		helmv2.RollbackVersion(int32(opts.Version)),
		helmv2.RollbackTimeout(int64(opts.Timeout.Seconds())),
		helmv2.RollbackWait(opts.Wait),
		helmv2.RollbackDisableHooks(opts.DisableHooks),
		helmv2.RollbackDryRun(opts.DryRun),
		helmv2.RollbackRecreate(opts.Recreate),
		helmv2.RollbackForce(opts.Force),
	)
	if err != nil {
		return nil, statusMessageErr(err)
	}
	return releaseToGenericRelease(res.Release), nil
}
