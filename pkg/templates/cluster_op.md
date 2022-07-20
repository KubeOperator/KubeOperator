### {{.title}}
> **项目**: {{ .projectName }} \
> **集群**: {{ .resourceName }} \
> **时间**: {{ .createdAt }} \
> **操作**: {{ .operator }} \
{{ if .errMsg }}
> **信息**: {{ .errMsg }} \
{{ end }}
<font color="info">本消息由KubeOperator自动发送</font>