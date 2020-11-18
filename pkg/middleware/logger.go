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
		operation = "删除|delete"
	} else {
		operation = matchMethod(ctx.Method())
	}

	operation_unit := matchApi(apiStr)
	if operation_unit == "" {
		operation_unit = apiStr + "|" + apiStr
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
		return "添加|create"
	case "DELETE":
		return "删除|delete"
	case "PATCH":
		return "修改|update"
	}
	return method + "|" + method
}

func initApiMap() map[string]string {
	apiMap := make(map[string]string)
	apiMap["projects"] = "项目|project"
	apiMap["project/resources"] = "项目资源|project resources"
	apiMap["project/members"] = "项目成员|project members"
	apiMap["clusters"] = "集群|cluster"
	apiMap["clusters/import"] = "集群导入|cluster import"
	apiMap["clusters/backup"] = "备份|cluster backup"
	apiMap["clusters/provisioner"] = "存储供应商|cluster storage provisioner"
	apiMap["events/npd/create"] = "集群时间npd启用|cluster events npd create"
	apiMap["events/npd/delete"] = "集群时间npd启用|cluster events npd delete"
	apiMap["clusters/tool/enable"] = "集群工具启用|cluster tool enable"
	apiMap["clusters/tool/disable"] = "集群工具停用|cluster tool disable"
	apiMap["clusters/cis"] = "集群CIS扫描|cluster cis scan"
	apiMap["clusters/backup/strategy"] = "集群备份策略|cluster backup strategy"
	apiMap["clusters/backup/files"] = "集群文件|cluster backup files"
	apiMap["clusters/backup/files/restore"] = "集群本地文件恢复|cluster backup files restore"
	apiMap["clusters/backup/files/backup"] = "集群本地文件备份|cluster backup files backup"

	apiMap["hosts"] = "主机|host"

	apiMap["regions"] = "区域|region"
	apiMap["regions/datacenter"] = "区域配置参数验证|validation of region configuration parameters"
	apiMap["zones"] = "可用区|zone"
	apiMap["zones/clusters"] = "可用区基本信息|basic information of zone "
	apiMap["plans"] = "部署计划|deploy plan"
	apiMap["vm/configs"] = "虚拟机配置|virtual machine configuration"

	apiMap["users"] = "用户|user"
	apiMap["users/change/password"] = "用户密码修改|user password modification"

	apiMap["manifests"] = "版本管理|version"

	apiMap["settings"] = "系统设置|system settings"
	apiMap["credentials"] = "系统设置-凭证|system settings - credentials"
	apiMap["backupaccounts"] = "系统设置-备份|system settings - backups"
	apiMap["backupaccounts/buckets"] = "系统设置-备份-获取桶|system settings - backups - bucket"
	apiMap["settings/check/EMAIL"] = "系统设置-邮件验证|system settings - mail verification"
	apiMap["license"] = "系统设置-许可|system settings - license"
	apiMap["message/setting"] = "消息配置|message settings"

	apiMap["logs"] = "日志|cluster"
	return apiMap
}
