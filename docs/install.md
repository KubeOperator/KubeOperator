# KubeOperator 安装手册

## 端口说明

<table border="0">
    <tr>
        <td>协议</td>
        <td>服务名称</td>
        <td>端口</td>
    </tr>
    <tr>
        <td>tcp</td>
        <td>Nginx</td>
        <td>80</td>
    </tr>
    <tr>
        <td>tcp</td>
        <td>Nexus</td>
        <td>8082,8092</td>
    </tr>
    <tr>
        <td>tcp</td>
        <td>Redis</td>
        <td>6379</td>
    </tr>
</table>



## 环境要求

+ 推荐硬件配置: 2个CPU核心,4G 内存,50G 硬盘
+ 操作系统要求: `CentOS 7 Minimal`
+ 配置基础网络、更新源等

## 安装准备

安装前请自行至百度云盘下载数据文件:
+ nexus-data.tar.gz (离线安装所需的rpm包和镜像): https:// pan.baidu.com/XXXX/nexus-data.tar.gz 提取码 0304 


## 开始安装

``` bash
# 文档中脚本默认均以root用户执行
$ yum update -y 
# 安装wget,git
$ yum install -y wget,git
# 下载KubeOperator
$  cd /opt/
$  git clone https://github.com/fit2anything/KubeOperator.git
# 请自行到百度网盘下载 nexus 数据文件: nexus-data.tar.gz
$ cp nexus-data.tar.gz /opt/KubeOperator/docker/nexus/
# 解压 nexus 数据文件
$ tar -zvxf nexus-data.tar.gz
# 运行安装脚本
$  ./kubeopsctl install
# 启动 KubeOperator
$  ./kubeopsctl start
# 查看 KubeOperator 状态
$ ./kubeopsctl status
