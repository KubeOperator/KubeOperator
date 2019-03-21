#!/bin/bash
red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

printf "%-65s .......... " "generate ssh key"
ls /root/.ssh | grep id_rsa >/dev/null || /bin/expect ./ssh-keygen.exp > /dev/null 1>&2 

if [ "$?" != "0" ];then
    colorMsg $red "[Defeat]"
    exit 1
fi
colorMsg $green "[OK]"
