package v2

import (
	"github.com/pkg/errors"

	helmv2 "k8s.io/helm/pkg/helm"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

func (h *HelmV2) History(releaseName string, opts helm.HistoryOptions) ([]*helm.Release, error) {
	max := helmv2.WithMaxHistory(256)
	if opts.Max != 0 {
		max = helmv2.WithMaxHistory(int32(opts.Max))
	}
	res, err := h.client.ReleaseHistory(releaseName, max)
	if err != nil {
		return nil, errors.Wrapf(statusMessageErr(err), "failed to retrieve history for [%s]", releaseName)
	}
	var rels []*helm.Release
	for _, r := range res.Releases {
		rels = append(rels, releaseToGenericRelease(r))
	}
	return rels, nil
}
