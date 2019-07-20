# KubeOperator 使用手册

KubeOperator 支持两种 kubernetes 集群部署模式，一种是一主多节点模式，另外一种是多主多节点模式。本手册仅描述一主多节点的部署和管理。

注：多主多节点模式适合在 MultiAZ（多故障域）下部署，实现双活环境下的高可用。KubeOperator 2.0 版本会支持 MultiAZ。





## 1 登录

访问 KubeOperator 控制台，执行登录操作。默认的登录用户名为 admin，默认密码为 kubeoperator@admin123。
> 为了保证系统的安全，请在完成登录后，点击控制台右上角的"修改密码"进行密码的重置。





## 2 设置

在使用 KubeOperator 之前，需要先对 KubeOperator 进行必要的参数设置。这些系统参数将影响到 kubernetes 集群的安装及相关服务的访问。

### 2.1 主机名

主机名主要用于 KubeOperator 所管理的集群用来访问 KubeOperator 使用，如果 KubeOperator 存在域名，请填写可以解析到本机的域名，否则可以使用本机的 IP 作为主机名。

### 2.2 域名后缀

域名后缀为集群节点访问地址的后缀，集群暴露出来的对外服务的 URL 都将以该域名后缀作为访问地址后缀。
例如: grafana.apps.cluster.f2c.com。

![KubeOperator设置](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/kubeops-setting.png?raw=true)





## 3 离线包

在离线包列表中可以查看 KubeOperator 当前所提供的 kubernetes 安装版本详细信息。在后续进行 kubernetes 集群部署时，可以从这些版本中选择其一进行部署。

![离线包-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/package-1.png?raw=true)





## 4 存储（当前仅支持 NFS）

在 kubernetes 集群中，服务的运行离不开数据的持久化存储，这就涉及到 k8s 的存储系统，常见的存储有 NFS、NetApp等。KubeOperator 当前版本支持的是 NFS。

### 4.1 添加存储

在左侧导航菜单中选择【存储】，进入【存储】页后点击【添加】，选择存储的类型：

![添加存储-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/storage-1.png?raw=true)

点击【下一步】，输入所选存储具体的配置信息，如 NFS 的服务器地址以及路径等。

![添加存储-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/storage-2.png?raw=true)





## 5 主机

### 5.1 节点准备

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

### 5.2 主机列表

在左侧导航菜单中选择【主机】，进入【主机】页后可以看到已添加主机的详细信息，包括 IP、CPU、内存、操作系统等。

![主机列表-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/hosts-1.png?raw=true)

### 5.3 添加主机

点击【添加】按钮添加新的主机。在输入完主机名称、IP、主机的 SSH 登录信息后，点击【提交】按钮即可完成一台主机的添加。

![添加主机-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/hosts-add-1.png?raw=true)





## 6 集群

### 6.1 集群列表

在左侧导航菜单中选择【集群】，进入【集群】页后可以看到已添加集群的详细信息，包括 集群部署的 kubernetes 版本、集群的部署模型、包含的节点数、运行状态等。

![集群列表-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-1.png?raw=true)

### 6.2 创建集群

#### 6.2.1 基本信息

点击【集群】页的【添加】按钮进行集群的创建。在【基本信息】里输入集群的名称，选择该集群所要部署的 kubernetes 版本。

![添加集群-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-1.png?raw=true)

#### 6.2.2 部署模型

选择 kubernetes 集群的部署模型。KubeOperator 当前版本仅支持一主多节点。选择部署模型后，将展示集群中各个角色的节点的详细配置要求。

![添加集群-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-2.png?raw=true)

#### 6.2.3 配置节点

【添加主机】环节中已经把集群所需的主机添加到了 KubeOperator 中，在【配置节点】环节，则可以根据不同的节点角色，选择主机列表中的各个主机。

![添加集群-3](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-3.png?raw=true)

#### 6.2.4 配置网络

【配置网络】环节，选择集群的网络插件，如下图的 flannel。如果集群是内网部署，网络模式请选择"host-gw"，若是公网部署，则选择"vxlan"。

![添加集群-4](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-4.png?raw=true)

#### 6.2.5 配置存储

选择在【添加存储】环节加入的存储系统，如 NFS 存储。kubernetes 集群将使用该存储做数据的持久化保存。

![添加集群-5](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-5.png?raw=true)

#### 6.2.6 配置检测

完成上述 5 个步骤后，KubeOperator 会对当前集群所选择的部署节点进行配置检测，包含CPU、内存和操作系统的检测。

![添加集群-6](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-6.png?raw=true)

#### 6.2.7 配置集群参数

完成检测后，可以对集群访问的一些域名参数进行配置，如无特殊要求，此处可保持 KubeOperator 设定的默认值。

![添加集群-7](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-7.png?raw=true)

#### 6.2.8 集群配置概览

创建集群的所有步骤完成后，会有一个集群配置概览页对之前步骤所设参数进行汇总，用户可在此页进行集群配置的最后检查。

![添加集群-8](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-create-8.png?raw=true)



### 6.3 部署集群

在集群列表中点击要进行部署的集群名称，默认展示的是该集群的【概览】信息。
【概览】页中展示了 kubernetes 集群的诸多详情，包括 kubernetes 版本、集群所用存储、网络等。
部署了 kubernetes 集群后，还可以在【概览】页中快捷的访问该 kubernetes 集群相关的监控采集、监控展示、镜像仓库、集群控制台等系统。

点击【概览】页最下方的【安装】按钮进行 kubernetes 集群部署。

![集群概览](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-overview.png?raw=true)

集群部署开始后，将会自动跳转到【任务】页。在【任务】页里可以看到集群部署当前所执行的具体任务信息。

![部署集群-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-deploy-1.png?raw=true)

如果是内网环境的话，一般一个典型的 5 节点集群的部署大概需要10分钟左右的时间，在出现类似下图的信息后，则表明集群已部署成功：

![部署集群-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-deploy-2.png?raw=true)

## 6.4 访问集群

回到集群的【概览】页，在该页提供了 Grafana、Prometheus、Registry-console、Dashboard 四个 kubernetes 集群相关的系统快捷访问方式。

> Grafana、Prometheus等系统的访问域名需要在 dns 服务器中添加相应的域名记录，如无条件，也可以通过修改 hosts 文件来达到相同的作用。

![集群概览](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster-overview.png?raw=true)

Prometheus 用来对整个 kubernetes 集群进行监控数据的采集。点击 Prometheus 下方的【转到】按钮即可访问 Prometheus 控制台。

![prometheus-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/prometheus-1.png?raw=true)

Grafana 对 Prometheus 采集到的监控数据进行了不同维度的图形化展示，更方便用户了解整个 kubernetes 集群的运行状况。点击 Grafana 下方的【转到】按钮访问 Grafana 控制台。

集群监控：

![grafana-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/grafana-1.png?raw=true)

节点监控：

![grafana-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/grafana-2.png?raw=true)

Registry-console 则用来存放 kubernetes 集群所使用到的 docker 镜像。

![regsitry-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/registry-1.png?raw=true)

Dashboard 对应的是 kubernetes 的控制台，从浏览器中访问 kubernetes 控制台需要用到【令牌】。点击【概览】页下方的【获取TOKEN】按钮获取令牌信息，将令牌信息复制到粘贴板。

![dashboard-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/dashboard-1.png?raw=true)

输入令牌信息后，点击【登录】，则可进入 kubernetes 控制台。

![dashboard-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/dashboard-2.png?raw=true)