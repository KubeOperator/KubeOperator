# KubeOperator - 容器集群管理平台

![Total visitor](https://visitor-count-badge.herokuapp.com/total.svg?repo_id=kubeoperator)
![Visitors in today](https://visitor-count-badge.herokuapp.com/today.svg?repo_id=kubeoperator)
[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.6.5-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

## 项目介绍

KubeOperator 是一个开源项目，帮助运维人员通过 Web-based UI，在完全离线和多云环境下，部署和管理生产级别的 Kubernetes 集群。KubeOperator 尤其适合在云平台（比如 VMware 及 Openstack）上部署和管理 Kubernetes 集群，实现 Kubernetes as a Service。
![overview](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/KubeOperator.jpeg?raw=true)

## 技术优势

-  按需创建：调用云平台 API，一键快速创建和部署 Kubernetes 集群 (即 Kubernetes as a Service)；
-  按需伸缩：快速伸缩 Kubernetes 集群，优化资源使用效率；
-  按需修补：快速升级和修补 Kubernetes 集群，并与社区最新版本同步，保证安全性；
-  自我修复：通过重建故障节点确保集群可用性；
-  离线部署：持续更新包括 Kubernetes 及常用组件的离线包；
-  Multi-AZ 支持：通过把 Kubernetes 集群 Master 节点分布在不同的故障域上确保的高可用；

## 版本规划

 v1.0.0 （已发布）

- [x] 提供原生 Kubernetes 的离线包仓库；
- [x] 支持一主多节点部署模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 集成 Kubernetes 常用插件的安装，包括 Registry、Promethus、Dashboard、Traefik Ingress、Helm 等；
- [x] 提供简易明了的 Kubernetes 集群运行状况面板；
- [x] 支持 NFS 作为持久化存储；
- [x] 支持 Flannel 作为网络方案；
- [x] 支持 Kubernetes 集群手动部署模式（自行准备主机和 NFS）；

 v2.0.0 （已发布）

- [x] 支持调用 VMware vCenter API 自动创建集群主机；
- [x] 支持 VMware vSAN 、VMFS/NFS 作为持久化存储；
- [x] 支持 Multi AZ，支持多主多节点部署模式；
- [x] 支持通过 F5 BIG-IP Controller 对外暴露服务（Nodeport mode, 七层和四层服务都支持）；
- [x] 集成 Weave Scope (支持 Web Shell)；
- [x] 支持 Calico 作为网络方案；

 v2.1.0 （开发中）
 
- [ ] 支持 Openstack 云平台；
- [ ] 支持 Ceph 作为持久化存储；
- [ ] 支持 Kubernetes 集群升级；
- [ ] 支持 Kubernetes 集群扩缩容；
- [ ] 支持 Kubernetes 集群备份与恢复；
- [ ] 支持 Kubernetes 集群健康检查与诊断；
- [ ] 集成 KubeApps（支持常用应用部署，如 Jenkins、GitLab、Harbor、Tekton、Sonarqube）；

 v3.0.0 （计划中）

- [ ] 支持 VMware NSX-T；
 
 ## 使用指南

-  [在线文档](https://docs.kubeoperator.io/)
-  [演示视频](http://kubeoperator.io/index.html#video)
-  [功能截屏](http://kubeoperator.io/index.html#screenshot)

 ## 沟通交流
 
- 技术交流 QQ 群：825046920
- 技术支持邮箱：support@fit2cloud.com

## 致谢

- [Terraform](https://github.com/hashicorp/terraform): KubeOperator 采用 Terraform 来自动创建虚机；
- [Clarity](https://github.com/vmware/clarity/): KubeOperator 采用 Clarity 作为前端 Web 框架；
- [Ansible](https://github.com/ansible/ansible): KubeOperator 采用 Ansible 作为自动化部署工具；
- [kubeasz](https://github.com/easzlab/kubeasz): 提供各种 Kubernetes Ansible 脚本；

## License & Copyright

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

KubeOperator is licensed under the Apache License, Version 2.0.
