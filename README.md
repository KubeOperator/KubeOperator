# KubeOperator

[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.6.5-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

KubeOperator 是一个开源项目，在离线网络环境下，通过可视化 Web UI 在 VMware、Openstack 或者物理机上部署和管理生产级别的 Kubernetes 集群。

![overview](https://github.com/KubeOperator/docs/blob/master/website/static/img/kubeoperator-ui.jpg?raw=true)

> Note: 可以点击查看大图。

## 整体架构

KubeOperator 使用 Terraform 在 IaaS 平台上自动创建主机，通过 Ansible 完成自动化部署和变更操作，支持 Kubernetes 集群 从 Day 0 规划，到 Day 1 部署，到 Day 2 变更的全生命周期管理。
![overview](https://github.com/KubeOperator/docs/blob/master/website/static/img/KubeOperator.jpeg?raw=true)

## 技术优势

-  按需创建：调用云平台 API，一键快速创建和部署 Kubernetes 集群 (即 Kubernetes as a Service)；
-  按需伸缩：快速伸缩 Kubernetes 集群，优化资源使用效率；
-  按需修补：快速升级和修补 Kubernetes 集群，并与社区最新版本同步，保证安全性；
-  自我修复：通过重建故障节点确保集群可用性；
-  离线部署：持续更新包括 Kubernetes 及常用组件的离线包；
-  Multi-AZ 支持：通过把 Kubernetes 集群 Master 节点分布在不同的故障域上确保的高可用；

## 版本规划

 v1.0 （已发布）

- [x] 提供原生 Kubernetes 的离线包仓库；
- [x] 支持一主多节点部署模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 内置 Kubernetes 常用系统应用的安装，包括 Registry、Promethus、Dashboard、Traefik Ingress、Helm 等；
- [x] 提供简易明了的 Kubernetes 集群运行状况面板；
- [x] 支持 NFS 作为持久化存储；
- [x] 支持 Flannel 作为网络方案；
- [x] 支持 Kubernetes 集群手动部署模式（自行准备主机和 NFS）；

 v2.0 （已发布）

- [x] 支持调用 VMware vCenter API 自动创建集群主机；
- [x] 支持 VMware vSAN 、VMFS/NFS 作为持久化存储；
- [x] 支持 Multi AZ，支持多主多节点部署模式；
- [x] 支持通过 F5 BIG-IP Controller 对外暴露服务（Nodeport mode, 七层和四层服务都支持）；
- [x] 内置 Weave Scope (支持 Web Shell)；
- [x] 支持 Calico 作为网络方案；

 v2.1 （开发中，预计 2019.10.31 发布）
 
- [ ] 支持 Openstack 云平台；
- [ ] 支持 Openstack Cinder 作为持久化存储；
- [ ] 支持 Kubernetes 集群升级 （Day 2）；
- [ ] 支持 Kubernetes 集群扩缩容（Day 2）；
- [ ] 支持 Kubernetes 集群备份与恢复（Day 2）；
- [ ] 支持 Kubernetes 集群健康检查与诊断（Day 2）；

 v2.2 （计划中，预计 2019.12.31 发布）

- [ ] 集成 KubeApps（支持常用应用部署，如 Jenkins、GitLab、Harbor、Tekton、Sonarqube）；
- [ ] 支持 VMware NSX-T；
 
 ## 使用指南

-  [演示视频](https://kubeoperator-1256577600.file.myqcloud.com/video/KubeOperator_2.0.mp4)
-  [在线文档](https://docs.kubeoperator.io/)

 ## 沟通交流
 
- 技术交流 QQ 群：825046920
- 技术支持邮箱：support@fit2cloud.com

## 致谢

- [Terraform](https://github.com/hashicorp/terraform): KubeOperator 采用 Terraform 来自动创建虚机；
- [Clarity](https://github.com/vmware/clarity/): KubeOperator 采用 Clarity 作为前端 Web 框架；
- [Ansible](https://github.com/ansible/ansible): KubeOperator 采用 Ansible 作为自动化部署工具；
- [kubeasz](https://github.com/easzlab/kubeasz): 提供各种 Kubernetes Ansible 脚本；

## License

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

KubeOperator is licensed under the Apache License, Version 2.0.
