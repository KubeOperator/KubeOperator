package constant

import (
	"path"
)

var (
	IngressChartPath = path.Join(ChartsDir, "nginx-ingress-0.5.0.tgz")
	EFKChartPath = path.Join(ChartsDir, "nginx-ingress-0.5.0.tgz")
)

const (
	IngressReleaseName = "nginx-ingress"
	EFKReleaseName     = "efk"
	KoNamespaceName    = "kube-operator"
)
