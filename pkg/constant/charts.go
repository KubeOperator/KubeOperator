package constant

import (
	"path"
)

var (
	IngressChartPath = path.Join(ChartsDir, "nginx-ingress-0.5.0.tgz")
)

const (
	IngressReleaseName = "nginx-ingress"
	KoNamespaceName    = "kube-operator"
)
