#! /bin/bash


red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}


printf "%-65s .......... " "config firewall service"

firewall-cmd --zone=public --add-port=80/tcp --add-port=8080/tcp --add-port=5380/tcp --add-port=53/udp --add-port=53/tcp --add-port=8081/tcp --add-port=8082/tcp --permanent >/dev/null 2>&1

if [ "$?" != "0" ];then
    colorMsg $red "[Defeat]"
    exit 1
fi


systemctl restart firewalld >/dev/null 2>&1 

if [ "$?" != "0" ];then
    colorMsg $red "[Defeat]"
    exit 1
fi

colorMsg $green "[OK]"

