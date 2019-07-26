# KubeOperator 使用手册

KubeOperator 支持两种 Kubernetes 集群部署模式，一种是一主多节点模式，另外一种是多主多节点模式。本手册仅描述一主多节点的部署和管理。

> 多主多节点模式适合在 MultiAZ（多故障域）下部署，实现双活环境下的高可用。KubeOperator 2.1 版本会支持 MultiAZ。

## 1 登录

KubeOperator 完全启动后，访问 KubeOperator 控制台，进行登录。默认的登录用户名为 admin，默认密码为 kubeoperator@admin123。

> 为了保证系统的安全，请在完成登录后，点击控制台右上角的"修改密码"进行密码的重置。

## 2 设置

在使用 KubeOperator 之前，需要先对 KubeOperator 进行必要的参数设置。这些系统参数将影响到 Kubernetes 集群的安装及相关服务的访问。

### 2.1 主机 IP

主机 IP 指 KubeOperator 机器自身的 IP。KubeOperator 所管理的集群将使用该 IP 来访问 KubeOperator。

### 2.2 集群域名后缀

集群域名后缀为集群节点访问地址的后缀，集群暴露出来的对外服务的 URL 都将以该域名后缀作为访问地址后缀。例如: grafana.apps.cluster.f2c.com。

![setting-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/setting-1.png?raw=true)

## 3 凭据

### 3.1 凭据列表

凭据为 KubeOperator 连接主机资产的凭证，可以使用 password 或者 private key 。

![credential-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/credential-1.png?raw=true)

### 3.1 创建凭据

点击【添加】按钮添加新的主机。

![add_credential-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/add_credential-1?raw=true)



## 3 离线包

在离线包列表中可以查看 KubeOperator 当前所提供的 Kubernetes 安装版本详细信息。在后续进行 Kubernetes 集群部署时，可以从这些版本中选择其一进行部署（当前仅支持1.15.0,后续会跟随 Kubernetes 社区发布离线包）。


## 4 主机

### 4.1 节点准备

KubeOperator 当前版本仅支持一主多节点的部署和管理，对于集群中各节点的要求如下：
<table>
    <tr>
        <td>角色</td>
        <td>数量</td>
        <td>操作系统</td>
        <td>最低配置</td>
        <td>推荐配置</td>
        <td>描述</td>
    </tr>
    <tr>
        <td>daemon</td>
        <td>1</td>
        <td>CentOS 7.6</td>
        <td>1C 2G</td>
        <td>2C 4G</td>
        <td>集群内 NTP 服务和 DNS 服务。</td>
    </tr>
    <tr>
        <td>master</td>
        <td>1</td>
        <td>CentOS 7.6</td>
        <td>2C 4G</td>
        <td>4C 16G</td>
        <td>运行 etcd、kube-apiserver、kube-scheduler、kube-apiserver。</td>
    </tr>
    <tr>
        <td>worker</td>
        <td>3+</td>
        <td>CentOS 7.6</td>
        <td>2C 8G</td>
        <td>8C 32G</td>
        <td>运行 kubelet、应用工作负载。</td>
    </tr>
</table>

### 4.2 主机列表

在左侧导航菜单中选择【主机】，进入【主机】页后可以看到已添加主机的详细信息，包括 IP、CPU、内存、操作系统等。


### 4.3 添加主机

点击【添加】按钮添加新的主机。在输入完主机名称、IP、主机的 SSH 登录信息后，点击【提交】按钮即可完成一台主机的添加。


## 5 集群

### 5.1 集群列表

在左侧导航菜单中选择【集群】，进入【集群】页后可以看到已添加集群的详细信息，包括 集群部署的 Kubernetes 版本、部署模式、节点数及运行状态等。


### 5.2 创建集群

#### 5.2.1 基本信息

点击【集群】页的【添加】按钮进行集群的创建。在【基本信息】里输入集群的名称，选择该集群所要部署的 Kubernetes 版本。


#### 5.2.2 部署模型

选择 Kubernetes 集群的部署模型。KubeOperator 当前版本仅支持一主多节点。选择部署模型后，KubeOperator 将展示集群中各个角色节点的详细配置要求。


#### 5.2.3 配置节点

【添加主机】环节，把集群所需的主机添加到了 KubeOperator 中。在【配置节点】环节，则可以根据不同的节点角色，选择主机列表中的各个主机。

#### 5.2.4 配置网络

【配置网络】环节，选择集群的网络插件，当前版本仅支持 Flannel。

> 如果集群节点全部都在同一个二层网络下，请选择"host-gw"。如果不是，则选择"vxlan"。"host-gw" 性能优于 "vxlan"。


#### 5.2.5 配置存储

【添加存储】环节，选择外部持久化存储。


#### 5.2.6 配置检测

完成上述 5 个步骤后，KubeOperator 会对当前集群所选择的部署节点进行配置检测，包含 CPU、内存和操作系统的检测。


#### 5.2.7 配置集群参数

完成检测后，可以对集群的域名参数进行配置，如无特殊要求，推荐使用默认值。


#### 5.2.8 集群配置概览

所有步骤完成后，会有一个集群配置概览页对之前步骤所设参数进行汇总，用户可在此页进行集群配置的最后检查。


### 5.3 部署集群

在集群列表中点击要进行部署的集群名称，默认展示的是该集群的【概览】信息。【概览】页中展示了 Kubernetes 集群的诸多详情，包括 Kubernetes 版本、集群所用存储、网络模式等。点击【概览】页最下方的【安装】按钮进行 Kubernetes 集群的部署。


集群部署开始后，将会自动跳转到【任务】页。在【任务】页里可以看到集群部署当前所执行的具体任务信息。


如果是内网环境的话，一个典型的 5 节点集群的部署大概需要10分钟左右的时间。在出现类似下图的信息后，表明集群已部署成功：


## 5.4 访问 Kubernetes 集群

回到集群的【概览】页，该页提供了 Grafana、Prometheus、Registry-console、Dashboard 等四个管理系统快捷访问方式。

> 这四个系统的访问域名需要在 DNS 服务器中添加相应的域名记录。如无条件，也可以通过修改本地 hosts 文件来达到相同的作用。


#### 5.4.1 访问 Dashboard

Dashboard 对应的是 Kubernetes 的控制台，从浏览器中访问 Kubernetes 控制台需要用到【令牌】。点击【概览】页下方的【获取TOKEN】按钮获取令牌信息，将令牌信息复制到粘贴板。

![dashboard-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/dashboard-1.png?raw=true)

输入令牌信息后，点击【登录】，则可进入 Kubernetes 控制台。

![dashboard-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/dashboard-2.png?raw=true)

#### 5.4.2 访问 Grafana

Grafana 对 Prometheus 采集到的监控数据进行了不同维度的图形化展示，更方便用户了解整个 Kubernetes 集群的运行状况。点击 Grafana 下方的【转到】按钮访问 Grafana 控制台。

集群级别的监控面板：

![grafana-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/grafana-1.png?raw=true)

节点级别的监控面板：

![grafana-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/grafana-2.png?raw=true)

#### 5.4.3 访问 Registry

Registry 则用来存放 Kubernetes 集群所使用到的 Docker 镜像。

![regsitry-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/registry-1.png?raw=true)

#### 5.4.4 访问 Prometheus

Prometheus 用来对整个 kubernetes 集群进行监控数据的采集。点击 Prometheus 下方的【转到】按钮即可访问 Prometheus 控制台。

![prometheus-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/prometheus-1.png?raw=true)
