#! /bin/bash


red=31
green=32
yellow=33
blue=34


function echoMsg()
{
  echo -e "$1"
}

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

function exit_error()
{
if [ "$?" == "0" ];then
    printf "\n"
else 
    colorMsg $red "Can not  install openshift!"
    exit 1
fi
}

function printTitle()
{
  echoMsg "\n\n**********\t ${1} \t**********\n"
}

function printSubTitle()
{
  echoMsg "------\t ${1} \t------\n"
}




echo -e " \n"





echo -e  " ███████╗██╗████████╗██████╗  ██████╗ ██████╗ ███████╗███╗   ██╗███████╗██╗  ██╗██╗███████╗████████╗\n██╔════╝██║╚══██╔══╝╚════██╗██╔═══██╗██╔══██╗██╔════╝████╗  ██║██╔════╝██║  ██║██║██╔════╝╚══██╔══╝\n█████╗  ██║   ██║    █████╔╝██║   ██║██████╔╝█████╗  ██╔██╗ ██║███████╗███████║██║█████╗     ██║   \n██╔══╝  ██║   ██║   ██╔═══╝ ██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║╚════██║██╔══██║██║██╔══╝     ██║   \n██║     ██║   ██║   ███████╗╚██████╔╝██║     ███████╗██║ ╚████║███████║██║  ██║██║██║        ██║   \n╚═╝     ╚═╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝        ╚═╝"

echo -e " \n"

colorMsg $blue "FIT2OPENSHIFT"
colorMsg $blue "https://www.fit2openshift.org"
colorMsg $blue "Powered by FIT2CLOUD 飞致云"

echo -e " \n"
printTitle "开始安装Fit2Openshift"
basepath=$(cd `dirname $0`; pwd)/..
#copy install files
printSubTitle "拷贝文件..."
cd $basepath/installer/scripts && ./copy_files.sh
exit_error
printSubTitle "安装依赖..."
#install node and docker
cd $basepath/installer/dependencies && ./install_dependencies.sh
exit_error
printSubTitle "Build Fit2Openshift 镜像..."
#build openshift
cd $basepath/fit2openshift && ./build.sh
exit_error
printSubTitle "Build 离线包镜像..."
#build offline_packages
cd $basepath/openshift-offline-resources   && ./build.sh
exit_error
printSubTitle "创建Service f2o..."
#create service
cd $basepath/installer/scripts && ./create_f2o_service.sh
exit_error
colorMsg $green "Comppeted install openshift!"



