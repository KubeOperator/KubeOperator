#! /bin/bash


red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}


logPath="/opt/fit2openshift/logs/install/"
timestamp=$(date -d now +%F)
fullLogFile=${logPath}"install_"${timestamp}".log"

printf "%-65s .......... " "build image: okd_offline_package:v3.11.0-0cbc58b"

 docker build --rm=true --tag=registry.fit2cloud.com/fit2anything/fit2openshift/okd_offline_package:v3.11.0-0cbc58b . >>$fullLogFile 2>&1

if [ "$?" == "0" ];then
    colorMsg $green "[OK]"
else
    colorMsg $red "[DEFEATED]"
    printf "\n"
    printf "Build okd_ffline_package:v3.11.0-0cbc58b defeated! An error log in :"${fullLogFile}
    printf "\n"
    exit 1
fi



