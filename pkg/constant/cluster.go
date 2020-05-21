package constant

const (
	ClusterRunning      = "Running"
	ClusterInitializing = "Initializing"
	ClusterFailed       = "Failed"
	ClusterTerminating  = "Terminating"

	ConditionTrue    = "True"
	ConditionFalse   = "False"
	ConditionUnknown = "Unknown"

	NodeRoleLabelKey   = "kubernetes.io/role"
	NodeRoleNameMaster = "master"
	NodeRoleNameWorker = "worker"
)
