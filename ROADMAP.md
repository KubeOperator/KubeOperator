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

 ##  v2.5 （开发中，20120.04.13 发布）
 
- [ ] LDAP/AD 对接; 
- [ ] 消息中心；
- [ ] 应用商店增加 [Argo CD](https://github.com/argoproj/argo-cd)，完整覆盖 CI 到 CD 的场景；
- [ ] 集群健康评分;
- [ ] 将部分内置应用移到应用商店；
- [ ] 支持 ingress-nginx，用户可选 ingress-nginx 或者 Traefik
- [ ] Excel批量导入主机；
- [ ] 手动模式集群批量扩容；

 ##  Backlog（计划中）
 
- [ ] 已有集群导入；
- [ ] 国际化支持；
- [ ] Helm 3.0 支持；
- [ ] Traefik 2.0 支持；
- [ ] 全容器化部署;
- [ ] Deprecate in-tree OpenStack and vSphere cloud controller；
- [ ] 自定义证书支持；
- [ ] 开放 REST API; 
- [ ] K8s 集群 API Server 的 高可用（ VIP ） 
- [ ] 外部存储支持 NetApp 存储（通过 [Trident](https://github.com/NetApp/trident)）； 
- [ ] 使用 [Sonobuoy](https://github.com/vmware-tanzu/sonobuoy) 进行 Kubernetes 集群合规检查并可视化展示结果；
- [ ] 支持 VMware NSX-T；
- [ ] 支持 containerd，用户可选 docker 或者 containerd
- [ ] 应用商店增加 [JumpServer 堡垒机](https://github.com/jumpserver/jumpserver)；
- [ ] 支持 ARM 64（鲲鹏） 
