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
	ClusterInstall            = "CLUSTER_INSTALL"
	ClusterImport             = "CLUSTER_IMPORT"
	ClusterUpgrade            = "CLUSTER_UPGRADE"
	ClusterDelete             = "CLUSTER_DELETE"
	ClusterScale              = "CLUSTER_SCALE"
	ClusterAddWorker          = "CLUSTER_ADD_WORKER"
	ClusterRemoveWorker       = "CLUSTER_REMOVE_WORKER"
	ClusterRestore            = "CLUSTER_RESTORE"
	ClusterBackup             = "CLUSTER_BACKUP"
	ClusterEnableProvisioner  = "CLUSTER_ENABLE_PROVISIONER"
	ClusterDisableProvisioner = "CLUSTER_DISABLE_PROVISIONER"
	ClusterEnableComponent    = "CLUSTER_ENABLE_COMPONENT"
	ClusterDisableComponent   = "CLUSTER_DISABLE_COMPONENT"
	ClusterEventWarning       = "CLUSTER_EVENT_WARNING"
	MsgTest                   = "MSG_TEST"
	LicenseExpires            = "LICENSE_EXPIRE"
	ClusterOperator           = "CLUSTER_OPERATOR"
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
	ClusterInstall:            "集群安装",
	ClusterImport:             "集群导入",
	ClusterDelete:             "集群删除",
	ClusterUpgrade:            "集群升级",
	ClusterScale:              "集群伸缩",
	ClusterAddWorker:          "集群扩容",
	ClusterRemoveWorker:       "集群缩容",
	ClusterRestore:            "集群恢复",
	ClusterBackup:             "集群备份",
	ClusterEnableProvisioner:  "启用存储提供商",
	ClusterDisableProvisioner: "禁用存储提供商",
	ClusterEnableComponent:    "启用集群组件",
	ClusterDisableComponent:   "禁用集群组件",
	ClusterEventWarning:       "集群事件告警",
	MsgTest:                   "KubeOperator测试",
	LicenseExpires:            "License到期提醒",
}

var Templates = map[string]map[string]string{
	MsgTest: {
		Email:      "pkg/templates/test.html",
		DingTalk:   "pkg/templates/test.html",
		WorkWeiXin: "pkg/templates/test.html",
	},
	LicenseExpires: {
		Email:      "pkg/templates/license_expire.html",
		DingTalk:   "pkg/templates/license_expire.md",
		WorkWeiXin: "pkg/templates/license_expire.md",
	},
	ClusterInstall: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterDelete: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterUpgrade: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterScale: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterAddWorker: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterRestore: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterBackup: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterEnableProvisioner: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterDisableProvisioner: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterEnableComponent: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterDisableComponent: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
	ClusterEventWarning: {
		Email:      "pkg/templates/cluster_op.html",
		DingTalk:   "pkg/templates/cluster_op.md",
		WorkWeiXin: "pkg/templates/cluster_op.md",
	},
}
