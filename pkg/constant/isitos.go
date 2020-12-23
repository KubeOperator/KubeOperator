package constant

const (
	IstioNamespace         = "istio-system"
	IstioOperatorNamespace = "istio-operator-help"

	BaseChartName     = "nexus/base"
	CniChartName      = "nexus/istio-cni"
	PilotChartName    = "nexus/istio-discovery"
	IngressChartName  = "nexus/istio-ingress"
	EgressChartName   = "nexus/istio-egress"
	OperatorChartName = "nexus/istio-operator"
	RemoteChartName   = "nexus/istiod-remote"
	CorednsChartName  = "nexus/istiocoredns"
)
