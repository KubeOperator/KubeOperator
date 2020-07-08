# KubeOperator - 从这里开启您的 Kubernetes 之旅

[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/KubeOperatpr/KubeOperatpr/blob/master/LICENSE)

KubeOperator 是一个开源的轻量级 Kubernetes 发行版，专注于帮助企业规划、部署和运营生产级别的 Kubernetes 集群。

KubeOperator 提供可视化的 Web UI，支持离线环境，支持物理机、VMware 和 OpenStack 等 IaaS 平台，支持 x86 和 arm64 架构，支持 GPU，内置应用商店，已通过 CNCF 的 Kubernetes 软件一致性认证。

KubeOperator 使用 Terraform 在 IaaS 平台上自动创建主机（用户也可以自行准备主机，比如物理机或者虚机），通过 Ansible 完成自动化部署和变更操作，支持 Kubernetes 集群 从 Day 0 规划，到 Day 1 部署，到 Day 2 运营的全生命周期管理。

KubeOperator 的整体架构：

![Architecture](https://kubeoperator.io/images/screenshot/ko-framework.svg)

> 注： KubeOperator 已通过云原生基金会（CNCF）的 [Kubernetes 软件一致性认证](https://landscape.cncf.io/selected=kube-operator)。

## Web UI 展示

![Web UI](https://kubeoperator.io/images/screenshot/6.png)

>更多功能截屏点击：[这里](https://docs.kubeoperator.io/kubeoperator-v2.2/screenshot)

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
        <td class="features-first-td-background-style" rowspan="15">Day 0 规划</td>
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
        <td class="features-third-td-background-style" rowspan="3">计算方案
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
        <td class="features-third-td-background-style" rowspan="3">存储方案
        </td>
        <td class="features-third-td-background-style">独立主机：支持 NFS / Ceph RBD / Local Volume
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
        <td class="features-third-td-background-style" rowspan="4">网络方案
        </td>
        <td class="features-third-td-background-style">支持 Flannel / Calico 网络插件
        </td>
    </tr>
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持通过 F5 Big IP 对外暴露服务
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 Traefik / ingress-nginx
        </td>
    </tr>    
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 CoreDNS
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
        <td class="features-third-td-background-style">支持 RHEL/CentOS 7.4 以上版本
        </td>
    </tr>  
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">容器运行时
        </td>
        <td class="features-third-td-background-style">支持 Docker / containerd
        </td>
    </tr>     
    <tr class="subscription-level-tr-border">
        <td class="features-first-td-background-style" rowspan="3">Day 1 部署
        </td>
        <td class="features-third-td-background-style" rowspan="3">部署
        </td>  
        <td class="features-third-td-background-style">提供离线环境下的完整安装包
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
        <td class="features-first-td-background-style" rowspan="21">Day 2 运营
        </td>
        <td class="features-third-td-background-style" rowspan="9">管理
        </td>  
        <td class="features-third-td-background-style">支持以项目为核心的分级授权管理
        </td>         
    </tr>
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">支持系统管理员、项目管理员和只读用户等三种角色
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">支持对接 LDAP/AD
        </td>
    </tr>    
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">对外开放 REST API
        </td>
    </tr>    
    <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">可通过应用商店安装 K8s Dashboard 管理应用
        </td>
    </tr>     
     <tr class="subscription-level-tr-border">
         <td class="features-third-td-background-style">可通过应用商店安装 Weave Scope 管理应用
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
         <td class="features-third-td-background-style">支持更新证书
        </td>
    </tr>     
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="4">可观察性
        </td>
         <td class="features-third-td-background-style">内置 Promethus，支持对集群、节点、Pod、Container的全方位监控和告警
        </td>
    </tr>
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">内置 Loki 日志方案
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">内置 Grafana 作为监控和日志展示
        </td>
    </tr> 
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持消息中心，通过钉钉、微信通知各种集群异常事件；
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
         <td class="features-third-td-background-style">支持 etcd 定期备份
        </td>
    </tr>  
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style"  rowspan="2">安全合规
        </td>
         <td class="features-third-td-background-style">支持集群健康评分
        </td>
    </tr>   
    <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">支持 CSI 安全扫描
        </td>
    </tr>    
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style" rowspan="2">应用商店
        </td>
         <td class="features-third-td-background-style">提供 GitLab、Jenkins、Harbor、Argo CD 等 CI/CD 工具
        </td>
    </tr> 
     <tr class="subscription-level-tr-border">
        <td class="features-third-td-background-style">提供深度学习AI 应用，比如 TensorFlow
        </td>
    </tr>    
 </table>

具体版本路线图请参考：[Roadmap](https://github.com/KubeOperator/KubeOperator/blob/master/ROADMAP.md)

 ## PPT、Demo 视频、使用文档

-  [PPT 介绍文档]( ./KubeOperator_Intro.pdf)：KubeOperator 的介绍 PPT。
-  [8 分钟的演示视频]( https://kubeoperator-1256577600.file.myqcloud.com/video/KubeOperator2.1.mp4)：详细演示 KubeOperator 的功能。
-  [安装及使用文档](https://docs.kubeoperator.io/)：包括 KubeOperator 安装文档、使用文档、功能截屏、常见问题等。

## 沟通交流
 
- 技术交流 QQ 群：825046920
- 微信群： 搜索微信号 wzk727912172，添加好友，备注（城市-github用户名）, 验证通过会加入群聊
- 技术咨询：https://jinshuju.net/f/QrZJpt

## 致谢

- [Terraform](https://github.com/hashicorp/terraform): KubeOperator 采用 Terraform 来自动创建虚机；
- [Clarity](https://github.com/vmware/clarity/): KubeOperator 采用 Clarity 作为前端 Web 框架；
- [Ansible](https://github.com/ansible/ansible): KubeOperator 采用 Ansible 作为自动化部署工具；
- [Kubeapps](https://github.com/kubeapps/kubeapps): KubeOperator 采用 Kubeapps 作为应用商店方案。

## License

Copyright (c) 2014-2019 FIT2CLOUD 飞致云

[https://www.fit2cloud.com](https://www.fit2cloud.com)<br>

KubeOperator is licensed under the Apache License, Version 2.0.
