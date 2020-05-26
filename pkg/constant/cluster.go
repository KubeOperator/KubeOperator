package constant

const (
	ClusterRunning      = "Running"
	ClusterInitializing = "Initializing"
	ClusterFailed       = "Failed"
	ClusterTerminating  = "Terminating"
	ClusterWaiting      = "Waiting"

	ConditionTrue    = "True"
	ConditionFalse   = "False"
	ConditionUnknown = "Unknown"

	NodeRoleNameMaster = "master"
	NodeRoleNameWorker = "worker"

	ClusterProviderBareMetal = "bareMetal"
	ClusterProviderVSphere   = "vSphere"
)
