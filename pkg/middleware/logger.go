package middleware

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/kataras/iris/v12/context"
	"strings"
)

const phaseName = "log middleware"

type InitLogMiddlewarePhase struct {
	systemLogService service.SystemLogService
	ApiMatchMap      map[string]string
}

var lM *InitLogMiddlewarePhase

func (l InitLogMiddlewarePhase) Init() error {
	apiMap := initApiMap()
	lM = &InitLogMiddlewarePhase{
		systemLogService: service.NewSystemLogService(),
		ApiMatchMap:      apiMap,
	}
	return nil
}

func (i *InitLogMiddlewarePhase) PhaseName() string {
	return phaseName
}

func LogMiddleware(ctx context.Context) {
	if ctx.GetStatusCode() == 200 {
		var u dto.SessionUser
		session := constant.Sess.Start(ctx)
		sessionUser := session.Get(constant.SessionUserKey)
		u = sessionUser.(*dto.Profile).User
		if ctx.Method() != "GET" {
			go saveSystemLogs(ctx, u.Name)
		}
		ctx.Next()
	}
}

func saveSystemLogs(ctx context.Context, name string) {
	apiStr := strings.ReplaceAll(ctx.Path(), "/api/v1/", "")
	operation := ""
	if strings.Index(apiStr, "/batch") != -1 {
		apiStr = strings.ReplaceAll(apiStr, "/batch", "")
		operation = "删除"
	} else {
		operation = matchMethod(ctx.Method())
	}
	operation_unit := matchApi(apiStr)
	if operation_unit == "" {
		operation_unit = apiStr
	}
	logInfo := dto.SystemLogCreate{
		Name:          name,
		OperationUnit: operation_unit,
		Operation:     operation,
		RequestPath:   ctx.Path(),
	}
	err := lM.systemLogService.Create(logInfo)
	if err != nil {
		fmt.Printf("save system logs err, err: %v\n", err)
	}
	return
}

func matchApi(api string) string {
	if strings.Index(api, "/") == -1 {
		backStr, ok := lM.ApiMatchMap[api]
		if ok {
			return backStr
		}
		return ""
	}
	backStr, ok := lM.ApiMatchMap[api]
	if ok {
		return backStr
	} else {
		lastIndex := strings.LastIndex(api, "/")
		api = api[0:lastIndex]
		return matchApi(api)
	}
}

func matchMethod(method string) string {
	switch method {
	case "POST":
		return "添加"
	case "DELETE":
		return "删除"
	case "PATCH":
		return "修改"
	}
	return method
}

func initApiMap() map[string]string {
	apiMap := make(map[string]string)

	apiMap["projects"] = "项目"
	apiMap["project/resources"] = "项目资源"
	apiMap["project/members"] = "项目成员"
	apiMap["clusters"] = "集群"
	apiMap["clusters/import"] = "集群导入"
	apiMap["clusters/backup"] = "备份"
	apiMap["clusters/provisioner"] = "存储供应商"
	apiMap["events/npd/create"] = "集群时间npd启用"
	apiMap["events/npd/delete"] = "集群时间npd启用"
	apiMap["clusters/tool/enable"] = "集群工具启用"
	apiMap["clusters/tool/disable"] = "集群工具停用"
	apiMap["clusters/cis"] = "集群CIS扫描"
	apiMap["clusters/backup/strategy"] = "集群备份策略"
	apiMap["clusters/backup/files"] = "集群文件"
	apiMap["clusters/backup/files/restore"] = "集群本地文件恢复"
	apiMap["clusters/backup/files/backup"] = "集群本地文件备份"

	apiMap["hosts"] = "主机"

	apiMap["regions"] = "区域"
	apiMap["regions/datacenter"] = "区域配置参数验证"
	apiMap["zones"] = "可用区"
	apiMap["zones/clusters"] = "可用区基本信息"
	apiMap["plans"] = "部署计划"
	apiMap["vm/configs"] = "虚拟机配置"

	apiMap["users"] = "用户"
	apiMap["users/change/password"] = "用户密码修改"

	apiMap["manifests"] = "版本管理"

	apiMap["settings"] = "系统设置"
	apiMap["credentials"] = "系统设置-凭证"
	apiMap["backupaccounts"] = "系统设置-备份"
	apiMap["backupaccounts/buckets"] = "系统设置-备份-获取桶"
	apiMap["settings/check/EMAIL"] = "系统设置-邮件验证"
	apiMap["license"] = "系统设置-许可"
	apiMap["message/setting"] = "消息配置"

	apiMap["logs"] = "日志"
	return apiMap
}
