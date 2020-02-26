#!/usr/bin/env bash
BASE_DIR=$(dirname "$0")
source ${BASE_DIR}/utils.sh
KUBEOPS_DIR="${INSTALL_DIR}/kubeoperator/"
validationPassed=1

echo -ne "root 用户检测 \t\t........................ "
isRoot=`id -u -n | grep root | wc -l`
if [ "x$isRoot" == "x1" ];then
  echo_green "[OK]"
else
  echo_red "[ERROR] 请用 root 用户执行安装脚本"
  validationPassed=0
fi


#操作系统检测
echo -ne  "操作系统检测 \t\t........................ "
if [ -f /etc/redhat-release ];then
  osVersion=`cat /etc/redhat-release | grep -oE '[0-9]+\.[0-9]+'`
  majorVersion=`echo $osVersion | awk -F. '{print $1}'`
  minorVersion=`echo $osVersion | awk -F. '{print $2}'`
  if [ "x$majorVersion" == "x" ];then
    echo_red "[ERROR] 操作系统类型版本不符合要求，请使用 CentOS 7.4 / 7.5 / 7.6 / 7.7 64 位版本"
    validationPassed=0
  else
    if [[ $majorVersion == 7 ]] && [[ $minorVersion > 3 ]];then
      is64bitArch=`uname -m`
      if [ "x$is64bitArch" == "xx86_64" ];then
         echo_green "[OK]"
      else
         echo_red "[ERROR] 操作系统必须是 64 位的，32 位的不支持"
         validationPassed=0
      fi
    else
      echo_red "[ERROR] 操作系统类型版本不符合要求，请使用 CentOS 7.4 / 7.5 / 7.6 / 7.7 版本"
      validationPassed=0
    fi
  fi
else
    echo_red "[ERROR] 操作系统类型版本不符合要求，请使用 CentOS 7.4 / 7.5 / 7.6 / 7.7 版本"
    validationPassed=0
fi


#CPU检测
echo -ne "CPU检测 \t\t........................ "
processor=`cat /proc/cpuinfo| grep "processor"| wc -l`
if [ $processor -lt 2 ];then
  echo_red "[ERROR] CPU 小于 2核，KubeOperator 所在机器的 CPU 需要至少 2核"
  validationPassed=0
else
  echo_green "[OK]"
fi


#内存检测
echo -ne "内存检测 \t\t........................ "
memTotal=`cat /proc/meminfo | grep MemTotal | awk '{print $2}'`
if [ $memTotal -lt 7500000 ];then
  echo_red "[ERROR] 内存小于 8G，KubeOperator 所在机器的内存需要至少 8G"
  validationPassed=0
else
  echo_green "[OK]"
fi


#磁盘剩余空间检测
echo -ne "磁盘剩余空间检测 \t........................ "
tmp_path=${INSTALL_DIR#*/}
path="/${tmp_path%%/*}"

IFSOld=$IFS
IFS=$'\n'
lines=$(df)
for line in ${lines};do
  linePath=`echo ${line} | awk -F' ' '{print $6}'`
  lineAvail=`echo ${line} | awk -F' ' '{print $4}'`
  if [ "${linePath:0:1}" != "/" ]; then
    continue
  fi
  
  if [ "${linePath}" == "/" ]; then
    rootAvail=${lineAvail}
    continue
  fi
  
  pathLength=${#path}
  if [ "${linePath:0:${pathLength}}" == "${path}" ]; then
    pathAvail=${lineAvail}
    break
  fi
done
IFS=$IFSOld

if test -z "${pathAvail}"
then
  pathAvail=${rootAvail}
fi

if [ $pathAvail -lt 50000000 ]; then
  echo_red "[ERROR] 安装目录剩余空间小于 50G，KubeOperator 所在机器的安装目录可用空间需要至少 50G"
  validationPassed=0
else
  echo_green "[OK]"
fi

if [ $validationPassed -eq 0 ]; then
  echo_red "安装环境检测未通过，请查阅上述环境检测结果"
  exit 1
fi
