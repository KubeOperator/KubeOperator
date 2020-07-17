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
- [x] GPU 支持;
- [x] 支持 Local Persistent Volumes;

Release Note: https://blog.fit2cloud.com/?p=1032

 ##  v2.4 （已发布）
 
- [x] 用户体系和权限；

Release Note:https://blog.fit2cloud.com/?p=1087

 ##  v2.5 （已发布）
 
- [x] LDAP/AD 对接; 
- [x] 消息中心；
- [x] 应用商店增加 [Argo CD](https://github.com/argoproj/argo-cd)，完整覆盖 CI 到 CD 的场景；
- [x] 集群健康评分;
- [x] 将部分内置应用移到应用商店；
- [x] 支持 ingress-nginx，用户可选 ingress-nginx 或者 Traefik
- [x] Excel批量导入主机；
- [x] 手动模式集群批量扩容；

Release Note:https://blog.fit2cloud.com/?p=1126

 ##  v2.6 （已发布）

- [x] 集成 [CIS 安全扫描](https://github.com/aquasecurity/kube-bench)；
- [x] 应用商店改成可选安装；
- [x] 支持 RHEL 7.4 以上操作系统；
- [x] vSphere 创建虚机时支持选用自有模板；
- [x] 支持更新 apiserver 证书。

 ##  V3.0 （开发中，8月18日发布）
 
- [ ] 开放 REST API;
- [ ] 支持 国际化 i18n；
- [ ] 支持 kubeadm 部署；
- [ ] 支持 arm64 平台架构；
- [ ] 使用 out-of-tree 的网络和存储插件；
- [ ] 支持在线和离线安装模式；
- [ ] 支持 Helm 3.x；
- [ ] 支持 Traefik 2.x；
- [ ] 集成 cert-manager；
- [ ] 支持 containerd；
- [ ] 支持 kubernetes 跨版本升级；
- [ ] 升级 kubeapps-plus 应用商店 (支持 Helm 3.x);
- [ ] 支持已有集群导入。
