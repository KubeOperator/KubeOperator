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
	MsgTest             = "MSG_TEST"
	LicenseExpires      = "LICENSE_EXPIRE"
	ClusterOperator     = "CLUSTER_OPERATOR"
)

//message level
const (
	MsgWarning = "Warning"
	MsgInfo    = "Info"
)

const (
	TestMessage = "KubeOperator消息测试"
)

const (
	StatusDisable = "DISABLE"
	StatusEnable  = "ENABLE"
)

var MsgTitle = map[string]string{
	ClusterInstall:      "集群安装",
	ClusterDelete:       "集群删除",
	ClusterUnInstall:    "集群卸载",
	ClusterUpgrade:      "集群升级",
	ClusterScale:        "集群伸缩",
	ClusterAddWorker:    "集群扩容",
	ClusterRemoveWorker: "集群缩容",
	ClusterRestore:      "集群恢复",
	ClusterBackup:       "集群备份",
	ClusterEventWarning: "集群事件告警",
	MsgTest:             "KubeOperator测试",
}

var Templates = map[string]map[string]string{
	MsgTest: {
		Email:      "pkg/templates/test.html",
		DingTalk:   "pkg/templates/test.html",
		WorkWeiXin: "pkg/templates/test.html",
	},
	ClusterInstall: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
}
