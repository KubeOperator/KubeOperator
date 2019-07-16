# KubeOperator 使用手册

KubeOperator 支持两种 kubernetes 集群部署模式，一种是一主多节点模式，另外一种是多主多节点模式。本手册仅描述一主多节点的部署和管理。

注：多主多节点模式适合在MultiAZ（多故障域）下部署，实现双活环境下的高可用。KubeOperator 2.0 版本会 MultiAZ。

## 一、集群规划

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

## 二、KubeOperator 设置

1.设置主机名: 如 KubeOperator 存在域名,请填写可以解析到本机的域名,否则使用本机IP。

2.设置域名后缀: 此后缀为集群节点的域名后缀，例如: master-1.nmss.f2c.com。

![KubeOperator设置](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/kubeops_setting.png?raw=true)

## 三、准备存储（当前仅支持 NFS）

1.添加存储:

![添加存储-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/storage-1.png?raw=true)

![添加存储-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/storage-2.png?raw=true)

## 四、准备主机

1.添加主机:

![添加主机-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/hosts-1.png?raw=true)

## 五、创建集群

1.创建集群:

![添加集群-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/create_cluster-1.png?raw=true)

![添加集群-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/create_cluster-2.png?raw=true)

![添加集群-3](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/create_cluster-3.png?raw=true)

![添加集群-4](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/create_cluster-4.png?raw=true)

![添加集群-5](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/create_cluster-5.png?raw=true)

![添加集群-5](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/cluster_create-7.png?raw=true)

## 六、部署集群

1.开始部署:

![开始部署-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/begin-1.png?raw=true)

![开始部署-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/begin-2.png?raw=true)

![开始部署-3](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/begin-3.png?raw=true)

![开始部署-4](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/complete-2.png?raw=true)

## 七、访问集群

TBD

## 八、集群监控

1.集群监控: 

![监控-1](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/metrics.png?raw=true)

2.节点监控

![监控-2](https://github.com/KubeOperator/KubeOperator/blob/master/docs/images/metrics-nodes.png?raw=true)
