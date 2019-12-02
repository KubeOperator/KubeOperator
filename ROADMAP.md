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
- [x] 支持全局的 DNS 和 NTP 设置；
- [x] 支持操作系统版本扩大到：CentOS 7.4 / 7.5 / 7.6 / 7.7；
- [ ] 集成 node-problem-detector，支持在 Web UI 上面查看集群事件；

 ## v2.3 （进行中，2019.12.31 发布）

- [ ] KubeApps 应用商店； 

 ##  v3.0 （计划中）
 
- [ ] 外部存储支持 NetApp 存储（通过 [Trident](https://github.com/NetApp/trident)）； 
- [ ] 支持用户和权限管理；
- [ ] 支持消息中心；
- [ ] 离线环境下使用 [Sonobuoy](https://github.com/vmware-tanzu/sonobuoy) 进行 Kubernetes 集群合规检查并可视化展示结果；
- [ ] 国际化支持；
- [ ] 支持 VMware NSX-T；
- [ ] 实现内置应用的统一认证；
- [ ] Helm 3.0 支持；
- [ ] Traefik 2.0 支持；
