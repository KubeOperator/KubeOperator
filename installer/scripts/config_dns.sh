#! /bin/bash

red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

function isConfigExists()
{
  cat /etc/named.rfc1912.zones | grep $1 >>/dev/null 2>&1
  if [[ "$?" != "0" ]]; then
    echo -e  "$2" >> /etc/named.rfc1912.zones
  fi
}

function isNamedZoneExists()
{
  ls /var/named/ | grep $1
  if [[ "$?" != "0" ]]; then
    cp ./configs/fit2openshift.io.zone /var/named/
  fi
}



printf "%-65s .......... " "Config DNS"

isConfigExists "fit2openshift.io"  "zone \"fit2openshift.io\" IN {
        type master;
        file \"fit2openshift.io.zone\";
};"

isNamedZoneExists "fit2openshift.io.zone"

\cp -fr  ./configs/named.conf  /etc/

service named restart >>/dev/null 2>&1
if [[ "$?" != "0" ]]; then
  exit 0;
fi


colorMsg $green "[OK]"

