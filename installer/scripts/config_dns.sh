#! /bin/bash

red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}


printf "%-65s .......... " "Config DNS"



\cp -fr  ./configs/dnsmasq.conf  /etc/ 
\cp -fr  ./configs/fit2openshift.dns.conf /etc/dnsmasq.d/

service dnsmasq restart >>/dev/null 2>&1
if [[ "$?" != "0" ]]; then
  exit 0;
fi


colorMsg $green "[OK]"

