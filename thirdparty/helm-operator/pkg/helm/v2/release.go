package v2

import (
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/ncabatoff/go-seq/seq"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

// releaseToGenericRelease transforms a v2 release structure
// into a generic `helm.Release`
func releaseToGenericRelease(r *release.Release) *helm.Release {
	return &helm.Release{
		Name:      r.Name,
		Namespace: r.Namespace,
		Chart:     chartToGenericChart(r.Chart),
		Info:      infoToGenericInfo(r.Info),
		Values:    valuesToGenericValues(r.Config),
		Manifest:  r.Manifest,
		Version:   int(r.Version),
	}
}

// chartToGenericChart transforms a v3 chart structure into
// a generic `helm.Chart`, while taking into account that the
// metadata may actually be missing due to unknown reasons.
// https://github.com/kubernetes/helm/issues/1347
func chartToGenericChart(c *chart.Chart) *helm.Chart {
	if c == nil || c.Metadata == nil {
		return nil
	}

	return &helm.Chart{
		Name:       c.Metadata.Name,
		Version:    c.Metadata.Version,
		AppVersion: c.Metadata.AppVersion,
		Values:     valuesToGenericValues(c.Values),
		Templates:  templatesToGenericFiles(c.Templates),
	}
}

// filesToGenericFiles transforms an `any.Any` slice into a
// stable sorted slice with generic `helm.File`s
func filesToGenericFiles(f []*any.Any) []*helm.File {
	gf := make([]*helm.File, len(f))
	for i, ff := range f {
		gf[i] = &helm.File{Name: ff.TypeUrl, Data: ff.Value}
	}
	sort.SliceStable(gf, func(i, j int) bool {
		return seq.Compare(gf[i], gf[j]) > 0
	})
	return gf
}

// filesToGenericFiles transforms a `chart.Template` slice into
// a stable sorted slice with generic `helm.File`s
func templatesToGenericFiles(t []*chart.Template) []*helm.File {
	gf := make([]*helm.File, len(t))
	for i, tf := range t {
		gf[i] = &helm.File{Name: tf.Name, Data: tf.Data}
	}
	sort.SliceStable(gf, func(i, j int) bool {
		return seq.Compare(gf[i], gf[j]) > 0
	})
	return gf
}

// infoToGenericInfo transforms a v2 info structure into
// a generic `helm.Info`
func infoToGenericInfo(i *release.Info) *helm.Info {
	if i == nil {
		return nil
	}
	return &helm.Info{
		LastDeployed: time.Unix(i.LastDeployed.Seconds, int64(i.LastDeployed.Nanos)),
		Description:  i.Description,
		Status:       lookUpGenericStatus(i.Status.Code),
	}
}

// valuesToGenericValues transforms a v2 values structure into
// a generic `map[string]interface{}`
func valuesToGenericValues(c *chart.Config) map[string]interface{} {
	vals, _ := chartutil.ReadValues([]byte(c.GetRaw()))
	return vals.AsMap()
}

// lookUpGenericStatus looks up the generic status for the
// given v2 status code
func lookUpGenericStatus(s release.Status_Code) helm.Status {
	var statuses = map[int32]helm.Status{
		0: helm.StatusUnknown,
		1: helm.StatusDeployed,
		2: helm.StatusUninstalled,
		3: helm.StatusSuperseded,
		4: helm.StatusFailed,
		5: helm.StatusUninstalling,
		6: helm.StatusPendingInstall,
		7: helm.StatusPendingUpgrade,
		8: helm.StatusPendingRollback,
	}
	if status, ok := statuses[int32(s)]; ok {
		return status
	}
	return helm.StatusUnknown
}
