package v3

import (
	"sort"

	"github.com/ncabatoff/go-seq/seq"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

// releaseToGenericRelease transforms a v3 release structure
// into a generic `helm.Release`
func releaseToGenericRelease(r *release.Release) *helm.Release {
	return &helm.Release{
		Name:      r.Name,
		Namespace: r.Namespace,
		Chart:     chartToGenericChart(r.Chart),
		Info:      infoToGenericInfo(r.Info),
		Values:    configToGenericValues(r.Config),
		Manifest:  r.Manifest,
		Version:   r.Version,
	}
}

// chartToGenericChart transforms a v3 chart structure into
// a generic `helm.Chart`
func chartToGenericChart(c *chart.Chart) *helm.Chart {
	return &helm.Chart{
		Name:       c.Name(),
		Version:    formatVersion(c),
		AppVersion: c.AppVersion(),
		Values:     c.Values,
		Templates:  filesToGenericFiles(c.Templates),
	}
}

// filesToGenericFiles transforms a `chart.File` slice into
// an stable sorted slice with generic `helm.File`s
func filesToGenericFiles(f []*chart.File) []*helm.File {
	gf := make([]*helm.File, len(f))
	for i, ff := range f {
		gf[i] = &helm.File{Name: ff.Name, Data: ff.Data}
	}
	sort.SliceStable(gf, func(i, j int) bool {
		return seq.Compare(gf[i], gf[j]) > 0
	})
	return gf
}

// infoToGenericInfo transforms a v3 info structure into
// a generic `helm.Info`
func infoToGenericInfo(i *release.Info) *helm.Info {
	return &helm.Info{
		LastDeployed: i.LastDeployed.Time,
		Description:  i.Description,
		Status:       lookUpGenericStatus(i.Status),
	}
}

// configToGenericValues forces the `chartutil.Values` to be parsed
// as YAML so that the value types of the nested map always equal to
// what they would be when returned from storage, as a dry-run release
// seems to skip this step, resulting in differences for e.g. float
// values (as they will be string values when returned from storage).
// TODO(hidde): this may not be necessary for v3.0.0 (stable).
func configToGenericValues(c chartutil.Values) map[string]interface{} {
	s, err := c.YAML()
	if err != nil {
		return c.AsMap()
	}
	gv, err := chartutil.ReadValues([]byte(s))
	if err != nil {
		return c.AsMap()
	}
	return gv.AsMap()
}

// formatVersion formats the chart version, while taking
// into account that the metadata may actually be missing
// due to unknown reasons.
// https://github.com/kubernetes/helm/issues/1347
func formatVersion(c *chart.Chart) string {
	if c.Metadata == nil {
		return ""
	}
	return c.Metadata.Version
}

// lookUpGenericStatus looks up the generic status for the
// given v3 release status.
func lookUpGenericStatus(s release.Status) helm.Status {
	var statuses = map[release.Status]helm.Status{
		release.StatusUnknown:         helm.StatusUnknown,
		release.StatusDeployed:        helm.StatusDeployed,
		release.StatusUninstalled:     helm.StatusUninstalled,
		release.StatusSuperseded:      helm.StatusSuperseded,
		release.StatusFailed:          helm.StatusFailed,
		release.StatusUninstalling:    helm.StatusUninstalling,
		release.StatusPendingInstall:  helm.StatusPendingInstall,
		release.StatusPendingUpgrade:  helm.StatusPendingUpgrade,
		release.StatusPendingRollback: helm.StatusPendingRollback,
	}
	if status, ok := statuses[s]; ok {
		return status
	}
	return helm.StatusUnknown
}
