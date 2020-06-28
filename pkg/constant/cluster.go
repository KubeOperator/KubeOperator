package constant

const (
	ClusterRunning      = "Running"
	ClusterInitializing = "Initializing"
	ClusterNotConnected = "NotConnected"
	ClusterFailed       = "Failed"
	ClusterTerminating  = "Terminating"
	ClusterTerminated   = "Terminated"
	ClusterWaiting      = "Waiting"

	ClusterSourceLocal    = "local"
	ClusterSourceExternal = "external"

	ConditionTrue    = "True"
	ConditionFalse   = "False"
	ConditionUnknown = "Unknown"

	NodeRoleNameMaster = "master"
	NodeRoleNameWorker = "worker"

	ClusterProviderBareMetal = "bareMetal"
	ClusterProviderVSphere   = "vSphere"

	ToolRunning      = "Running"
	ToolInitializing = "Initializing"
	ToolFailed       = "Failed"

	DefaultNamespace     = "kube-operator"
	DefaultApiServerPort = 8443
)
