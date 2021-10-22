<p align="center"><a href="https://kubeoperator.io"><img src="https://kubeoperator.oss-cn-beijing.aliyuncs.com/img/logo-red.png" alt="KubeOperator" width="300" /></a></p>
<h3 align="center">从这里开启您的 Kubernetes 之旅</h3>
<p align="center">
  <a href="http://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/github/license/kubeoperator/kubeoperator?color=%231890FF&style=flat-square" alt="License: Apache License v2"></a>
  <a href="https://app.codacy.com/gh/kubeoperator/kubeoperator?utm_source=github.com&utm_medium=referral&utm_content=kubeoperator/kubeoperator&utm_campaign=Badge_Grade_Dashboard"><img src="https://app.codacy.com/project/badge/Grade/da67574fd82b473992781d1386b937ef" alt="Codacy"></a>
  <a href="https://github.com/kubeoperator/kubeoperator/releases/latest"><img src="https://img.shields.io/github/v/release/kubeoperator/kubeoperator" alt="Latest release"></a>
  <a href="https://github.com/kubeoperator/kubeoperator"><img src="https://img.shields.io/github/stars/kubeoperator/kubeoperator?color=%231890FF&style=flat-square" alt="Stars"></a>
  <a href="https://github.com/kubeoperator/kubeoperator/releases/latest"><img src="https://img.shields.io/github/downloads/kubeoperator/kubeoperator/total" alt="Downloads"></a>
</p>
<hr />

KubeOperator 是一个开源的轻量级 Kubernetes 发行版，专注于帮助企业规划、部署和运营生产级别的 Kubernetes 集群。

KubeOperator 提供可视化的 Web UI，支持离线环境，支持物理机、VMware、OpenStack 和 FusionCompute 等 IaaS 平台，支持 x86_64 和 ARM64 架构，支持 GPU，内置应用商店，已通过 CNCF 的 Kubernetes 软件一致性认证。

KubeOperator 使用 Terraform 在 IaaS 平台上自动创建主机（用户也可以自行准备主机，比如物理机或者虚机），通过 Ansible 完成自动化部署和变更操作，支持 Kubernetes 集群 从 Day 0 规划，到 Day 1 部署，到 Day 2 运营的全生命周期管理。

### 组件项目

- **[KubePi](https://github.com/KubeOperator/KubePi)**: 简单易用的开源 Kubernetes 可视化管理面板
- **[Web Kubectl](https://github.com/KubeOperator/webkubectl)**: 在 Web 浏览器中运行 kubectl 命令

### KubeOperator 的优势

-   **简单易用**: 提供可视化的 Web UI，极大降低 K8s 部署和管理门槛，内置 [KubePi](https://github.com/KubeOperator/KubePi) 和 [Web Kubectl](https://github.com/KubeOperator/webkubectl)；
-   **按需创建**: 调用云平台 API，一键快速创建和部署 Kubernetes 集群；
-   **按需伸缩**: 快速伸缩 Kubernetes 集群，优化资源使用效率；
-   **按需修补**: 快速升级和修补 Kubernetes 集群，并与社区最新版本同步，保证安全性；
-   **离线部署**: 支持完全离线下的 K8s 集群部署；
-   **自我修复**: 通过重建故障节点确保集群可用性；
-   **全栈监控**: 提供从Pod、Node到集群的事件、监控、告警、和日志方案；
-   **Multi-AZ 支持**: 将 Master 节点分布在不同的故障域上确保集群高可用；
-   **应用商店**: 内置 [KubeApps](https://github.com/kubeapps/kubeapps) 应用商店；
-   **GPU 支持**: 支持 GPU 节点，助力运行深度学习等应用。

### UI 展示

![UI展示](https://kubeoperator.oss-cn-beijing.aliyuncs.com/img/demo.gif)

### 功能架构

![Architecture](https://kubeoperator.io/images/screenshot/ko-framework.svg)

### 在线体验

-   环境地址：<https://demo.kubeoperator.io/>
-   用户名：demo
-   密码：Password@123

| :warning: 注意                 |
| :--------------------------- |
| 该环境仅作体验目的使用，我们会定时清理、重置数据！    |
| 请勿修改体验环境用户的密码！               |
| 请勿在环境中添加业务生产环境地址、用户名密码等敏感信息！ |

### 快速开始

仅需两步快速安装 KubeOperator：

 1. 准备一台不小于 8 G内存的 64位 Linux 主机；
 2. 以 root 用户执行如下命令一键安装 KubeOperator。

```sh
curl -sSL https://github.com/KubeOperator/KubeOperator/releases/latest/download/quick_start.sh | sh
```

文档和演示视频：

- [完整文档](https://kubeoperator.io/docs/)
- [演示视频](https://www.bilibili.com/video/BV1jT4y1L7Ur/)
- [PPT 介绍](https://kubeoperator.io/download/KubeOperator_Intro.pdf)

### 微信交流群

![wechat-group](https://kubeoperator.io/docs/img/wechat-group.png)


### 致谢

- [Terraform](https://github.com/hashicorp/terraform): KubeOperator 采用 Terraform 来自动创建虚机；
- [Ansible](https://github.com/ansible/ansible): KubeOperator 采用 Ansible 作为自动化部署工具；
- [Kubeapps](https://github.com/kubeapps/kubeapps): KubeOperator 采用 Kubeapps 作为应用商店方案。

### License & Copyright

Copyright (c) 2014-2021 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubeOperator is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
