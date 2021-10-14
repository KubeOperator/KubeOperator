package constant

const (
	Email          = "EMAIL"
	DingTalk       = "DING_TALK"
	WorkWeiXin     = "WORK_WEIXIN"
	LocalMail      = "LOCAL"
	Enable         = "ENABLE"
	Disable        = "DISABLE"
	Cluster        = "CLUSTER"
	System         = "SYSTEM"
	Read           = "READ"
	UnRead         = "UNREAD"
	SendSuccess    = "SUCCESS"
	SendFailed     = "FAILED"
	OperateSuccess = "SUCCESS"
	OperateFailed  = "FAILED"
)

//message type
const (
	ClusterInstall      = "CLUSTER_INSTALL"
	ClusterImport       = "CLUSTER_IMPORT"
	ClusterUnInstall    = "CLUSTER_UN_INSTALL"
	ClusterUpgrade      = "CLUSTER_UPGRADE"
	ClusterDelete       = "CLUSTER_DELETE"
	ClusterScale        = "CLUSTER_SCALE"
	ClusterAddWorker    = "CLUSTER_ADD_WORKER"
	ClusterRemoveWorker = "CLUSTER_REMOVE_WORKER"
	ClusterRestore      = "CLUSTER_RESTORE"
	ClusterBackup       = "CLUSTER_BACKUP"
	ClusterEventWarning = "CLUSTER_EVENT_WARNING"
)

//message level
const (
	MsgWarning = "Warning"
	MsgInfo    = "Info"
)
