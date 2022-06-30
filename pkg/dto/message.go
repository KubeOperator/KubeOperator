package dto

import "github.com/KubeOperator/KubeOperator/pkg/constant"

type MsgContent struct {
}

var MsgTitle = map[string]string{
	constant.ClusterInstall:      "集群安装",
	constant.ClusterDelete:       "集群删除",
	constant.ClusterUnInstall:    "集群卸载",
	constant.ClusterUpgrade:      "集群升级",
	constant.ClusterScale:        "集群伸缩",
	constant.ClusterAddWorker:    "集群扩容",
	constant.ClusterRemoveWorker: "集群缩容",
	constant.ClusterRestore:      "集群恢复",
	constant.ClusterBackup:       "集群备份",
	constant.ClusterEventWarning: "集群事件告警",
	constant.MsgTest:             "KubeOperator测试",
}

var Templates = map[string]map[string]string{
	constant.MsgTest: {
		constant.Email:      "pkg/templates/test.html",
		constant.DingTalk:   "pkg/templates/test.html",
		constant.WorkWeiXin: "pkg/templates/test.html",
	},

	constant.ClusterInstall: {
		constant.Email:      "",
		constant.DingTalk:   "",
		constant.WorkWeiXin: "",
	},
}
