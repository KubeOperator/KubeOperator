#! /bin/bash


red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}





echo -e " \n"





echo -e  " ███████╗██╗████████╗██████╗  ██████╗ ██████╗ ███████╗███╗   ██╗███████╗██╗  ██╗██╗███████╗████████╗\n██╔════╝██║╚══██╔══╝╚════██╗██╔═══██╗██╔══██╗██╔════╝████╗  ██║██╔════╝██║  ██║██║██╔════╝╚══██╔══╝\n█████╗  ██║   ██║    █████╔╝██║   ██║██████╔╝█████╗  ██╔██╗ ██║███████╗███████║██║█████╗     ██║   \n██╔══╝  ██║   ██║   ██╔═══╝ ██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║╚════██║██╔══██║██║██╔══╝     ██║   \n██║     ██║   ██║   ███████╗╚██████╔╝██║     ███████╗██║ ╚████║███████║██║  ██║██║██║        ██║   \n╚═╝     ╚═╝   ╚═╝   ╚══════╝ ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝╚═╝        ╚═╝"

echo -e " \n"

colorMsg $blue "FIT2OPENSHIFT"
colorMsg $blue "https://www.fit2openshift.org"
colorMsg $blue "Powered by FIT2CLOUD 飞致云"

echo -e " \n"

basepath=$(cd `dirname $0`; pwd)/..

#copy install files
cd $basepath/installer/scripts
./copy_files.sh

#build openshift
cd $basepath/fit2openshift
./build.sh


#build offline_packages

for file in $basepath/openshift-offline-resources//*; do
      $file/build.sh
done

#create service
cd $basepath/installer/scripts
./create_f2o_service.sh

colorMsg $green "Comppeted install openshift!"



