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
errorLogFile=${logPath}"error/install_error_"${timestamp}".log"
infoLogFile=${logPath}"info/install_info_"${timestamp}".log"
fullLogFile=${logPath}"install_"${timestamp}".log"

printf "%-65s .......... " "build image: okd_ffline_package:v3.11.0-0cbc58b"

cd okd-3.11-meta \
&& ./download-dependencies.sh 1>>$infoLogFile 2>>$errorLogFile >>$fullLogFile 1>&2 \
&&  docker build --rm=true --tag=registry.fit2cloud.com/fit2anything/fit2openshift/okd_ffline_package:v3.11.0-0cbc58b . 1>>$infoLogFile 2>>$errorLogFile >>$fullLogFile 1>&2 \

if [ "$?" == "0" ];then
    colorMsg $green "[OK]"
else
    colorMsg $red "[DEFEATED]"
    printf "\n"
    printf "Build okd_ffline_package:v3.11.0-0cbc58b defeated! An error log in :"${errorLogFile}
    printf "\n"
    exit 1
fi

colorMsg $green "[OK]" \


