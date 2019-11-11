# KubeOperator - 从这里开启您的 Kubernetes 之旅

[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/KubeOperatpr/KubeOperatpr/blob/master/LICENSE)
[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.6.5-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

KubeOperator 是一个开源项目，在离线网络环境下，通过可视化 Web UI 在 VMware、Openstack 或者物理机上规划、部署和管理生产级别的 Kubernetes 集群。KubeOperator 是 [Jumpserver](https://github.com/jumpserver/jumpserver) 明星开源团队在 Kubernetes 领域的的又一全新力作。

![overview](https://github.com/KubeOperator/docs/blob/master/website/static/img/overview.png?raw=true)

> 注： KubeOperator 2.1 已通过云原生基金会（CNCF）的 [Kubernetes 软件一致性认证](https://landscape.cncf.io/selected=kube-operator)。

## Web UI 展示

![overview](https://raw.githubusercontent.com/KubeOperator/website/master/images/kubeoperator-ui.jpg)

>更多功能截屏请查看：https://docs.kubeoperator.io/kubeoperator-v2.1/screenshot

## 整体架构

KubeOperator 使用 Terraform 在 IaaS 平台上自动创建主机（用户也可以自行准备主机，比如物理机或者虚机），通过 Ansible 完成自动化部署和变更操作，支持 Kubernetes 集群 从 Day 0 规划，到 Day 1 部署，到 Day 2 运维及变更的全生命周期管理。

![overview](https://github.com/KubeOperator/docs/blob/master/website/static/img/KubeOperator.jpeg?raw=true)

## 技术优势

-  按需创建：调用云平台 API，一键快速创建和部署 Kubernetes 集群 (即 Kubernetes as a Service)；
-  按需伸缩：快速伸缩 Kubernetes 集群，优化资源使用效率；
-  按需修补：快速升级和修补 Kubernetes 集群，并与社区最新版本同步，保证安全性；
-  自我修复：通过重建故障节点确保集群可用性；
-  离线部署：持续更新包括 Kubernetes 及常用组件的离线包；
-  Multi-AZ 支持：通过把 Kubernetes 集群 Master 节点分布在不同的故障域上确保的高可用；

 ## Demo 视频、使用文档

-  [:tv:8 分钟演示视频]( https://kubeoperator-1256577600.file.myqcloud.com/video/KubeOperator2.1.mp4)：详细演示 KubeOperator 的功能。
-  [:books:安装及使用文档](https://docs.kubeoperator.io/)：包括 KubeOperator 安装文档、使用文档、功能截屏、常见问题等。

 ## Kubernetes 离线安装包

KubeOperator 提供完整的离线 Kubernetes 安装包（包括 Kubernetes、Docker、etcd、Dashboard、Promethus、OS 补丁等），每个安装包会被构建成一个独立容器镜像供 KubeOperator 使用，具体信息请参考：[k8s-package](https://github.com/KubeOperator/k8s-package)。

## 版本规划

 v1.0 （已发布）

- [x] 提供原生 Kubernetes 的离线包仓库；
- [x] 支持一主多节点部署模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 内置 Kubernetes 常用系统应用的安装，包括 Registry、Promethus、Dashboard、Traefik Ingress、Helm 等；
- [x] 提供简易明了的 Kubernetes 集群运行状况面板；
- [x] 支持 NFS 作为持久化存储；
- [x] 支持 Flannel 网络插件；
- [x] 支持 Kubernetes 集群手动部署模式（自行准备主机和 NFS）；

 v2.0 （已发布）

- [x] 支持调用 VMware vCenter API 自动创建集群主机；
- [x] 支持 VMware vSAN 、VMFS/NFS 作为持久化存储；
- [x] 支持 Multi AZ，支持多主多节点部署模式；
- [x] 支持 Calico 网络插件；
- [x] 内置 Weave Scope；
- [x] 支持通过 F5 BIG-IP Controller 对外暴露服务（Nodeport mode, 七层和四层服务都支持）；
- [x] 支持 Kubernetes 1.15；

 v2.1 （已发布）
 
- [x] 支持 Openstack 云平台；
- [x] 支持 Openstack Cinder 作为持久化存储；
- [x] 支持 Kubernetes 集群升级 （Day 2）；
- [x] 支持 Kubernetes 集群扩缩容（Day 2）；
- [x] 支持 Kubernetes 集群备份与恢复（Day 2）；
- [x] 支持 Kubernetes 集群健康检查与诊断（Day 2）；
- [x] 支持 [webkubectl](https://github.com/webkubectl/webkubectl) ；

 v2.2 （进行中，2019.11.30 发布）

- [x] 集成 Loki 日志方案，实现监控、告警和日志技术栈的统一；
- [x] KubeOperator 自身的系统日志收集和管理；
- [x] 概览页面：展示关键信息，比如状态、容量、TOP 使用率、异常日志、异常容器等信息；
- [ ] 支持 Ceph RBD 存储 （通过 Rook）；
- [ ] 支持 Kubernetes 1.16；
- [ ] 支持全局的 NTP 设置；
- [ ] 支持操作系统版本扩大到：CentOS 7.4 / 7.5 / 7.6 / 7.7；

 v2.3 （计划中，2019.12.31 发布）

- [ ] 实现内置应用的统一认证；
- [ ] KubeApps 应用商店；
- [ ] 支持 NetApp 存储；

 v3.0 （计划中）
 
- [ ] 离线环境下使用 Sonobuoy 进行 Kubernetes 集群合规检查并可视化展示结果；
- [ ] 国际化支持；
- [ ] 支持 VMware NSX-T；

## 沟通交流
 
- 技术交流 QQ 群：825046920；
- 技术支持邮箱：support@fit2cloud.com；
- 微信群： 搜索微信号 wh_it0224，添加好友，备注（城市-github用户名）, 验证通过会加入群聊；

## 致谢

- [Terraform](https://github.com/hashicorp/terraform): KubeOperator 采用 Terraform 来自动创建虚机；
- [Clarity](https://github.com/vmware/clarity/): KubeOperator 采用 Clarity 作为前端 Web 框架；
- [Ansible](https://github.com/ansible/ansible): KubeOperator 采用 Ansible 作为自动化部署工具；
- [kubeasz](https://github.com/easzlab/kubeasz): 提供各种 Kubernetes Ansible 脚本；

## License

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubeOperator is licensed under the Apache License, Version 2.0.
