 ## v1.0 （已发布）

- [x] 提供原生 Kubernetes 的离线包仓库；
- [x] 支持一主多节点部署模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 内置 Kubernetes 常用系统应用的安装，包括 Registry、Promethus、Dashboard、Traefik Ingress、Helm 等；
- [x] 提供简易明了的 Kubernetes 集群运行状况面板；
- [x] 支持 NFS 作为持久化存储；
- [x] 支持 Flannel 网络插件；
- [x] 支持 Kubernetes 集群手动部署模式（自行准备主机和 NFS）；

 ## v2.0 （已发布）

- [x] 支持调用 VMware vCenter API 自动创建集群主机；
- [x] 支持 VMware vSAN 、VMFS/NFS 作为持久化存储；
- [x] 支持 Multi AZ，支持多主多节点部署模式；
- [x] 支持 Calico 网络插件；
- [x] 内置 Weave Scope；
- [x] 支持通过 F5 BIG-IP Controller 对外暴露服务（Nodeport mode, 七层和四层服务都支持）；
- [x] 支持 Kubernetes 1.15；

 ## v2.1 （已发布）
 
- [x] 支持 Openstack 云平台；
- [x] 支持 Openstack Cinder 作为持久化存储；
- [x] 支持 Kubernetes 集群升级 （Day 2）；
- [x] 支持 Kubernetes 集群扩缩容（Day 2）；
- [x] 支持 Kubernetes 集群备份与恢复（Day 2）；
- [x] 支持 Kubernetes 集群健康检查与诊断（Day 2）；
- [x] 支持 [webkubectl](https://github.com/webkubectl/webkubectl) ；

 ## v2.2 （已发布）

- [x] 集成 [Loki](https://github.com/grafana/loki) 日志方案，实现监控、告警和日志技术栈的统一；
- [x] KubeOperator 自身的系统日志收集和管理；
- [x] 概览页面：展示关键信息，比如状态、容量、TOP 使用率、异常日志、异常容器等信息；
- [x] 支持 Ceph RBD 存储 （通过 [Rook](https://github.com/rook/rook)）；
- [x] 支持 Kubernetes 1.16；
- [x] 支持全局的和 NTP 设置；
- [x] 支持操作系统版本扩大到：CentOS 7.4 / 7.5 / 7.6 / 7.7；
- [x] 支持在 Web UI 上面查看集群事件；

Release Note: https://blog.fit2cloud.com/?p=980

 ## v2.3 （已发布）

- [x] KubeApps Plus 应用商店；
- [x] GPU 支持；
- [x] 支持 Local Persistent Volumes；

Release Note: https://blog.fit2cloud.com/?p=1032

 ##  v2.4 （已发布）
 
- [x] 用户体系和权限；

Release Note:https://blog.fit2cloud.com/?p=1087

 ##  v2.5 （已发布）
 
- [x] LDAP/AD 对接； 
- [x] 消息中心；
- [x] 应用商店增加 [Argo CD](https://github.com/argoproj/argo-cd)，完整覆盖 CI 到 CD 的场景；
- [x] 集群健康评分；
- [x] 将部分内置应用移到应用商店；
- [x] 支持 ingress-nginx，用户可选 ingress-nginx 或者 Traefik
- [x] Excel批量导入主机；
- [x] 手动模式集群批量扩容；

Release Note: https://blog.fit2cloud.com/?p=1126

 ##  v2.6 （已发布）

- [x] 集成 [CIS 安全扫描](https://github.com/aquasecurity/kube-bench)；
- [x] 应用商店改成可选安装；
- [x] 支持 RHEL 7.4 以上操作系统；
- [x] vSphere 创建虚机时支持选用自有模板；
- [x] 支持更新 apiserver 证书。

Release Note: https://blog.fit2cloud.com/?p=1219

 ##  v3.0 （已发布）
 
- [x] 开放 REST API；
- [x] 支持 国际化 i18n；
- [x] 支持 kubeadm 部署；
- [x] 支持 arm64 平台架构；
- [x] 支持在线和离线安装模式；
- [x] 组件升级，包括 Helm 3.x、Traefik 2.x、Kubeapps等；
- [x] 集成 cert-manager；
- [x] 支持 containerd；
- [x] 支持已有集群导入。

Release Note: https://blog.fit2cloud.com/?p=1416

 ##  v3.1 （已发布）
 
 - [x] K8s 版本历史及集群版本升级功能优化；
 - [x] 集成 CIS 安全扫描；
 - [x] 集群事件；
 - [x] 集群 ETCD 定时备份和自定义恢复；
 - [x] 自定义 Logo 和 配色 （X-Pack）；
 - [x] LDAP 对接（X-Pack）；
 
 Release Note: https://blog.fit2cloud.com/?p=1480
 
 ##  v3.2 （已发布）
 
 - [x] 增加消息中心（X-Pack）；
 - [x] 支持邮箱、钉钉、企业微信告警（X-Pack）；
 - [x] 支持实时查看任务执行返回日志；
 - [x] 应用商店增加 Redmine；
 
 Release Note: https://blog.fit2cloud.com/?p=1516
 
 ## v3.3 (已发布)

 - [x] FusionCompute 支持自动部署模式；
 - [x] 持久化存储支持 OceanStor；
 - [x] 集群日志，支持 EFK；
 - [x] 集群健康评估（X-Pack）；
 - [x] F5 对接（X-Pack）；
 - [x] 支持动态管理 Kubernetes 及组件版本；
 - [x] 自动模式支持自定义 cpu、内存规格；
 - [x] 集群创建支持指定网卡信息、helm版本；
 - [x] 支持登录验证码；
 - [x] 支持添加、删除 namespace；
 - [x] 集群事件支持启用、禁用 npd；
 - [x] restapi 开启 rbac 认证；
 - [x] 支持 session 和 jwt 两种认证方式；
 
 Release Note: https://blog.fit2cloud.com/?p=1612

 ## v3.4 (已发布)
 
 - [x] 多集群配置管理（X-Pack）；
 - [x] GPU 支持；
 - [x] 系统操作日志；
 - [x] 集群日志，支持 Loki；
 - [x] 备份支持到本地盘和 SFTP；
 - [x] 批量导入主机；
 - [x] 忘记密码；
 - [x] 应用商店增加 Kuboard、TensorFlow；
 
 Release Note: https://blog.fit2cloud.com/?p=1669

 ## v3.5 (开发中)
 
 - [ ] 集群异常状态诊断及修复；
 - [ ] 支持 Istio；
 - [ ] 自动模式创建主机支持 IP 池；
 - [ ] 支持自定义 Ansible、Terraform 并发参数；
 - [ ] 版本管理支持上传 K8s 离线包；
