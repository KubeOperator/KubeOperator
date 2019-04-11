# Welcome to FIT2OPENSHIFT (aka Fit to OpenShift)

[![Python3](https://img.shields.io/badge/python-3.6-green.svg?style=plastic)](https://www.python.org/)
[![Django](https://img.shields.io/badge/django-2.1-brightgreen.svg?style=plastic)](https://www.djangoproject.com/)
[![Ansible](https://img.shields.io/badge/ansible-2.4.2.0-blue.svg?style=plastic)](https://www.ansible.com/)
[![Angular](https://img.shields.io/badge/angular-7.0.4-red.svg?style=plastic)](https://www.angular.cn/)

**说明：该项目当前处于 V1.0 Beta 阶段，仅支持 okd 3.11**

## 什么是 FIT2OPENSHIFT？

FIT2OPENSHIFT 是一个开源项目，帮助运维人员通过 Web 控制台，在完全离线环境下实现 [OpenShift](https://www.openshift.com/) 社区版（[okd](https://www.okd.io/)）集群的可视化部署及运维管理。

## 为什么需要 FIT2OPENSHIFT? 

K8S 是未来的 Linux，OpenShift 是首选的 K8S 发行版。OpenShift 社区版（okd）集群部署及后续持续运营的门槛较高（尤其是离线环境下高可用多节点集群的部署）。

## FIT2OPENSHIFT 有什么功能？

- [x] 提供 OpenShift 社区版各个版本离线包仓库；
- [x] 提供两种 OpenShift 部署模式：单节点模式，高可用模式；
- [x] 支持离线环境下的一键自动化部署，可视化展示集群部署进展和结果；
- [x] 支持集群进行扩容；
- [x] 提供简易明了的集群运行状况面板；
- [x] 默认支持 GlusterFS 做为持久化存储；
- [ ] 支持其他外部存储做为持久化存储；
- [ ] 支持集群的升级；
- [ ] 支持集群备份及恢复；

## FIT2OPENSHIFT 支持哪些 okd 版本？

3.11 及以上。

## FIT2OPENSHIFT 背后是谁在支持？

[FIT2CLOUD](https://www.fit2cloud.com) 是 FIT2OPENSHIFT 项目的发起者及核心贡献者。该项目由第一开源堡垒机 [Jumpserver](http://www.jumpserver.org/) 的原班团队打造。

## FIT2OPENSHIFT 的架构

![架构图](https://raw.githubusercontent.com/fit2anything/fit2openshift/master/docs/images/overview.png)

## 安装 FIT2OPENSHIFT

 [安装手册](https://github.com/fit2anything/fit2openshift/blob/master/docs/install.md)

## 使用 FIT2OPENSHIFT

 [使用手册](https://github.com/fit2anything/fit2openshift/blob/master/docs/user-guide.md)
 
## License & Copyright

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

Licensed under The GNU General Public License version 2 (GPLv2)  (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

https://www.gnu.org/licenses/gpl-2.0.html

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

