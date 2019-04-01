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
printf "%-65s .......... " "Build fit2openshift webconsole ui:"
cd ui
npm install >/dev/null 2>&1
ng build --prod >>$fullLogFile 2>&1
if [ "$?" == "0" ];then
    colorMsg $green "[OK]"
else
    colorMsg $red "[DEFEATED]"
    printf "\n"
    printf "Build fit2openshift webconsole ui  defeated! An error log in :"${fullLogFile}
    printf "\n"
    exit 1
fi
printf "\n"
printf "%-65s .......... " "Build fit2openshift webconsole api: "
cd .. && docker build --rm=true --tag=registry.fit2cloud.com/fit2anything/fit2openshift/fit2openshift-app:latest . >>$fullLogFile 2>&1

if [ "$?" == "0" ];then
    colorMsg $green "[OK]"
else
    colorMsg $red "[DEFEATED]"
    printf "\n"
    printf "Build fit2openshift webconsole api  defeated! An error log in :"${fullLogFile}
    printf "\n"
    exit 1
fi

