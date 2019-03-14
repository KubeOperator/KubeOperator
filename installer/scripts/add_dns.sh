#! /bin/bash

red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

function isDNSExists()
{
  cat /etc/resolv.conf | grep $1
  if [[ "$?" != "0" ]]; then
    echo "nameserver $1" >> /etc/resolv.conf
  fi
}



printf "%-65s .......... " "Set dns"
isDNSExists "localhost"
isDNSExists "114.114.114.114"
colorMsg $green "[OK]"

