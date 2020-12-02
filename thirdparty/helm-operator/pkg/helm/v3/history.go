package v3

import (
	"github.com/pkg/errors"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/releaseutil"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

func (h *HelmV3) History(releaseName string, opts helm.HistoryOptions) ([]*helm.Release, error) {
	cfg, err := newActionConfig(h.kubeConfig, h.infoLogFunc(opts.Namespace, releaseName), opts.Namespace, "")
	if err != nil {
		return nil, err
	}

	history := action.NewHistory(cfg)
	historyOptions(opts).configure(history)

	hist, err := history.Run(releaseName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve history for '%s'", releaseName)
	}

	releaseutil.Reverse(hist, releaseutil.SortByRevision)

	var rels []*helm.Release
	for i := 0; i < min(len(hist), history.Max); i++ {
		rels = append(rels, releaseToGenericRelease(hist[i]))
	}
	return rels, nil
}

type historyOptions helm.HistoryOptions

func (opts historyOptions) configure(action *action.History) {
	action.Max = opts.Max
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
