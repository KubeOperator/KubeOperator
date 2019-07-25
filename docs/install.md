# KubeOperator 安装手册

## 环境要求

+ 推荐硬件配置: 2个CPU核心,4G 内存,50G 硬盘
+ 操作系统要求: `CentOS 7 Minimal`
+ 配置基础网络、更新源等

## 安装准备



## 开始安装

### 1.源代码安装


安装前请自行至百度云盘下载数据文件:  
https://pan.baidu.com/XXXX/nexus-data.tar.gz 提取码 0304 
  

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
```

### 2.离线包安装

安装前请自行至百度云盘下载完整离线包:  
https://pan.baidu.com/XXXX/nexus-data.tar.gz 提取码 0304 

``` bash
# 解压离线包
$ unzip kubeOperator-release-xx.zip
# 进入项目目录
$ cd kubeOperator-release
# 运行安装脚本
$ ./kubeopsctl install
# 启动 KubeOperator
$  ./kubeopsctl start
# 查看 KubeOperator 状态
$ ./kubeopsctl status
```
