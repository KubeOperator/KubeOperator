# KubeOperator - 从这里开启您的 Kubernetes 之旅

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/9cb491920f0d4058aa273500a38e3abf)](https://www.codacy.com/gh/KubeOperator/KubeOperator/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=KubeOperator/KubeOperator&amp;utm_campaign=Badge_Grade)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/kubeoperator/kubeoperator/blob/master/LICENSE)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/kubeoperator/kubeoperator)](https://github.com/kubeoperator/kubeoperator/releases/latest)
[![GitHub All Releases](https://img.shields.io/github/downloads/kubeoperator/kubeoperator/total)](https://github.com/kubeoperator/kubeoperator/releases)

> [English](README_EN.md) | 中文

KubeOperator 是一个开源的轻量级 Kubernetes 发行版，专注于帮助企业规划、部署和运营生产级别的 Kubernetes 集群。

KubeOperator 提供可视化的 Web UI，支持离线环境，支持物理机、VMware、OpenStack 和 FusionCompute 等 IaaS 平台，支持 x86 和 ARM64 架构，支持 GPU，内置应用商店，已通过 CNCF 的 Kubernetes 软件一致性认证。

KubeOperator 使用 Terraform 在 IaaS 平台上自动创建主机（用户也可以自行准备主机，比如物理机或者虚机），通过 Ansible 完成自动化部署和变更操作，支持 Kubernetes 集群 从 Day 0 规划，到 Day 1 部署，到 Day 2 运营的全生命周期管理。

## 整体架构

![Architecture](https://kubeoperator.io/images/screenshot/ko-framework.svg)

## Web UI 展示

![Web UI](https://kubeoperator.io/images/screenshot/02.jpg)

>更多功能截屏点击：[这里](https://kubeoperator.io/features.html)

## 快速开始

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

## KubeOperator 企业版

- [申请企业版试用](https://jinshuju.net/f/qc6g44/)

## 版本说明

KubeOperator 版本号命名规则为：v大版本.功能版本.Bug修复版本。比如：

```
v1.0.1 是 v1.0.0 之后的Bug修复版本；
v1.1.0 是 v1.0.0 之后的功能版本。
```
像其它优秀开源项目一样，KubeOperator 将每月发布一个功能版本。

## 技术优势

-  简单易用：提供可视化的 Web UI，极大降低 K8s 部署和管理门槛，内置 [Webkubectl](https://github.com/KubeOperator/webkubectl)；
-  按需创建：调用云平台 API，一键快速创建和部署 Kubernetes 集群；
-  按需伸缩：快速伸缩 Kubernetes 集群，优化资源使用效率；
-  按需修补：快速升级和修补 Kubernetes 集群，并与社区最新版本同步，保证安全性；
-  离线部署：支持完全离线下的 K8s 集群部署；
-  自我修复：通过重建故障节点确保集群可用性；
-  全栈监控：提供从Pod、Node到集群的事件、监控、告警、和日志方案；
-  Multi-AZ 支持：将 Master 节点分布在不同的故障域上确保集群高可用；
-  应用商店：内置 [KubeApps](https://github.com/kubeapps/kubeapps) 应用商店；
-  GPU 支持：支持 GPU 节点，助力运行深度学习等应用；

## 功能列表

<table class="subscription-level-table">
    <tr class="subscription-level-tr-border">
        <td class="features-first-td-background-style" rowspan="17">Day 0 规划</td>
        </td>
        <td class="features-third-td-background-style" rowspan="2">集群模式
        </td>
        <td class="features-third-td-background-style">1 个 Master 节点 n 个 Worker 节点模式：适合开发测试用途
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">3 个 Master 节点 n 个 Worker 节点模式：适合生产用途
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="4">计算方案
        </td>
        <td class="features-third-td-background-style">独立主机：支持自行准备的虚机、公有云主机和物理机
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">vSphere 平台：支持自动创建主机（使用 Terraform）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Openstack 平台：支持自动创建主机 （使用 Terraform）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">FusionCompute 平台：支持自动创建主机 （使用 Terraform）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="4">存储方案
        </td>
        <td class="features-third-td-background-style">独立主机：支持 NFS / Ceph RBD / Rook Ceph / Local Volume
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">vSphere 平台：支持 vSphere Datastore （vSAN 及 vSphere 兼容的集中存储）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">Openstack 平台：支持 Openstack Cinder （Ceph 及 Cinder 兼容的集中存储）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">FusionCompute 平台：支持 OceanStor
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="4">网络方案
        </td>
        <td class="features-third-td-background-style">支持 CoreDNS
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 Flannel / Calico / Cilium 网络插件
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 ingress-nginx / traefik
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持通过 F5 Big IP 对外暴露服务（X-PACK）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">GPU 方案
        </td>
        <td class="features-third-td-background-style">支持 NVIDIA GPU
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">操作系统
        </td>
        <td class="features-third-td-background-style">支持 RHEL / CentOS / EulerOS 操作系统
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">容器运行时
        </td>
        <td class="features-third-td-background-style">支持 Docker / Containerd
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-first-td-background-style" rowspan="6">Day 1 部署
        </td>
        <td class="features-third-td-background-style" rowspan="6">部署
        </td>
        <td class="features-third-td-background-style">支持在线和离线安装模式
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 Kubeadm 部署
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 x86_64 和 arm64 CPU 架构
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持可视化方式展示部署过程
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持一键自动化部署（使用 Ansible）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持已有集群导入
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-first-td-background-style" rowspan="22">Day 2 运营
        </td>
        <td class="features-third-td-background-style" rowspan="9">管理
        </td>
        <td class="features-third-td-background-style">支持以项目为核心的分级授权管理
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持系统管理员、项目管理员和集群管理员三种角色
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持多集群配置管理（X-PACK）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持对接 LDAP/AD（X-PACK）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持自定义 Logo 和 配色（X-PACK）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">对外开放 REST API
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持国际化 i18n
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">提供 Web Kubectl 界面
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">内置 Helm
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="4">可观察性
        </td>
        <td class="features-third-td-background-style">内置 Prometheus，支持对集群、节点、Pod、Container 的全方位监控和告警
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">内置 EFK、Loki 日志方案
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">内置 Grafana 作为监控和日志展示
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持消息中心，通过钉钉、微信通知各种集群异常事件（X-PACK）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">升级
        </td>
        <td class="features-third-td-background-style">支持集群升级
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">伸缩
        </td>
        <td class="features-third-td-background-style">支持增加或者减少 Worker 节点
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">备份
        </td>
        <td class="features-third-td-background-style">支持 etcd 定期备份和立即备份
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">恢复
        </td>
        <td class="features-third-td-background-style">支持 etcd 备份策略文件恢复和本地文件恢复
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style"  rowspan="2">安全合规
        </td>
        <td class="features-third-td-background-style">支持集群健康评分（X-PACK）
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 CSI 安全扫描
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="3">应用商店
        </td>
        <td class="features-third-td-background-style">提供 GitLab、Jenkins、Harbor、Argo CD、Sonarqube 等 CI/CD 工具
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">提供 Kuboard、Weave Scope、Redmine 等管理工具
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">提供深度学习AI 应用，比如 TensorFlow
        </td>
    </tr>
 </table>

具体版本路线图请参考：[Roadmap](https://github.com/KubeOperator/KubeOperator/blob/master/ROADMAP.md)

## 支持组件

- 核心
  - [kubernetes](https://github.com/kubernetes/kubernetes) v1.20.6
  - [etcd](https://github.com/coreos/etcd) v3.4.14
  - [docker](https://www.docker.com/) v19.03.15
  - [containerd](https://containerd.io/) v1.4.3
- 网络
  - [calico](https://github.com/projectcalico/calico) v3.16.5
  - [flanneld](https://github.com/coreos/flannel) v0.13.0
- 应用
  - [coredns](https://github.com/coredns/coredns) v1.7.0
  - [helm-v2](https://github.com/helm/helm) v2.17.0
  - [helm-v3](https://github.com/helm/helm) v3.6.0
  - [traefik](https://github.com/containous/traefik) v2.4.8
  - [ingress-nginx](https://github.com/kubernetes/ingress-nginx) v0.33.0
  - [metrics-server](https://github.com/kubernetes-sigs/metrics-server) v0.3.6
- 工具
  - [istio](https://github.com/istio/istio) 1.8.0
  - [dashboard](https://github.com/kubernetes/dashboard) v2.2.0
  - [kubeapps](https://github.com/kubeapps/kubeapps) v2.0.1
  - [prometheus](https://github.com/prometheus/prometheus) v2.20.1
  - [grafana](https://github.com/grafana/grafana) v7.3.3
  - [loki](https://github.com/grafana/loki) v2.1.0
  - [logging](https://github.com/elastic/elasticsearch) v7.6.2
  - [chartmuseum](https://github.com/helm/chartmuseum) v0.12.0
  - [docker-registry](https://github.com/docker/distribution) v2.7.1
- 应用商店
  - [argo-cd](https://github.com/argoproj/argo-cd) v2.0.3
  - [gitlab-ce](https://about.gitlab.com) v9.4.1
  - [harbor](https://github.com/goharbor/harbor) v1.10.2
  - [jenkins](https://github.com/jenkinsci/jenkins) v2.222.1
  - [kuboard](https://github.com/eip-work/kuboard-press) v2.0.5.1
  - [redmine](https://github.com/redmine/redmine) v4.1.1
  - [sonarqube](https://github.com/SonarSource/sonarqube) v7.9.2
  - [tensorflow-serving](https://github.com/tensorflow/serving) v1.14.0
  - [tensorflow-notebook](https://github.com/tensorflow/tensorflow) v1.6.0
  - [weave-scope](https://github.com/weaveworks/scope) v1.12.0

## 微信群

![wechat-group](https://kubeoperator.io/docs/img/wechat-group.png)

## 致谢

- [Terraform](https://github.com/hashicorp/terraform): KubeOperator 采用 Terraform 来自动创建虚机；
- [Ansible](https://github.com/ansible/ansible): KubeOperator 采用 Ansible 作为自动化部署工具；
- [Kubeapps](https://github.com/kubeapps/kubeapps): KubeOperator 采用 Kubeapps 作为应用商店方案。

## License

Copyright (c) 2014-2020 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubeOperator is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
