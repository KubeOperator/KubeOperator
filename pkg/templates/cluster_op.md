### {{.title}}
> **项目**: {{ .projectName }} \n\n
> **集群**: {{ .resourceName }} \n\n
> **时间**: {{ .createdAt }} \n\n
> **操作**: {{ .operator }} \n\n
{{ if .errMsg }}
> **信息**: {{ .errMsg }} \n\n
{{ end }}
<font color="info">本消息由KubeOperator自动发送</font>