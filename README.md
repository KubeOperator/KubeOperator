<p align="center"><a href="https://kubeoperator.io"><img src="https://kubeoperator.oss-cn-beijing.aliyuncs.com/img/logo-red.png" alt="KubeOperator" width="300" /></a></p>
<h3 align="center">开源的轻量级 Kubernetes 发行版</h3>
<p align="center">
  <a href="http://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/github/license/kubeoperator/kubeoperator?color=%231890FF&style=flat-square" alt="License: Apache License v2"></a>
  <a href="https://github.com/kubeoperator/kubeoperator"><img src="https://img.shields.io/github/stars/kubeoperator/kubeoperator?color=%231890FF&style=flat-square" alt="Stars"></a>
</p>
<hr />

KubeOperator 是开源的轻量级 Kubernetes 发行版。KubeOperator 提供可视化的 Web UI，支持离线环境，通过 Ansible 完成自动化部署和变更操作，支持 Kubernetes 集群 从 Day 0 规划，到 Day 1 部署，到 Day 2 运营的全生命周期管理。

### KubeOperator 的优势

-   **简单易用**: 提供可视化的 Web UI，极大降低 K8s 部署和管理门槛；
-   **离线部署**: 支持完全离线下的 K8s 集群部署和管理；
-   **全栈监控**: 提供从Pod、Node到集群的事件、监控、告警、和日志方案；
-   **快速修补**: 快速升级 Kubernetes 集群，与社区最新版本同步，保证安全性；
-   **快速创建**: 通过 Terraform 调用 IaaS API，快速创建、部署和伸缩 Kubernetes 集群。

### 功能架构

![Architecture](https://kubeoperator.io/images/screenshot/ko.png)

### UI 展示

![UI展示](https://kubeoperator.oss-cn-beijing.aliyuncs.com/img/demo.gif)

### 致谢

- [Terraform](https://github.com/hashicorp/terraform): KubeOperator 采用 Terraform 来自动创建虚机；
- [Ansible](https://github.com/ansible/ansible): KubeOperator 采用 Ansible 作为自动化部署工具。

### License & Copyright

Copyright (c) 2014-2023 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubeOperator is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
