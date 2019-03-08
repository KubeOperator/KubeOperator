#! /bin/bash

red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

printf "%-65s .......... " "set dns"
echo "nameserver 114.114.114.114" >> /etc/resolv.conf
colorMsg $green "[OK]"