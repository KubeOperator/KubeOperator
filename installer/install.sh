#! /bin/bash


red=31
green=32
yellow=33
blue=34

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



echo -e " \n"





echo -e  " ███████╗██╗████████╗██████╗  ██████╗ ██████╗ ███████╗███╗   ██╗███████╗██╗  ██╗██╗███████╗████████╗\n██╔════╝██║╚══██╔══╝╚════██╗██╔═══██╗██╔══██╗██╔════╝████╗  ██║██╔════╝██║  ██║██║██╔════╝╚══██╔══╝\n█████╗  ██║   ██║    █████╔╝██║   ██║██████╔╝█████╗  ██╔██╗ ██║███████╗███████║██║█████╗     ██║   \n██╔══╝  ██║   ██║   ██╔═══╝ ██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║╚════██║██╔══██║██║██╔══╝     ██║   \n██║     ██║   ██║   ███████╗╚██████╔╝██║     ███████╗██║ ╚████║███████║██║  ██║██║██║        ██║   \n╚═╝     ╚═╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝        ╚═╝"

echo -e " \n"

colorMsg $blue "FIT2OPENSHIFT"
colorMsg $blue "https://www.fit2openshift.org"
colorMsg $blue "Powered by FIT2CLOUD 飞致云"

echo -e " \n"

basepath=$(cd `dirname $0`; pwd)/..

#install node and docker
cd $basepath/installer/dependencies && ./install_dependencies.sh
exit_error
#copy install files
cd $basepath/installer/scripts && ./copy_files.sh
exit_error
#build openshift
cd $basepath/fit2openshift && ./build.sh
exit_error
#build offline_packages
cd $basepath/openshift-offline-resources && ./bulid.sh
exit_error
#create service
cd $basepath/installer/scripts && ./create_f2o_service.sh
exit_error
colorMsg $green "Comppeted install openshift!"



