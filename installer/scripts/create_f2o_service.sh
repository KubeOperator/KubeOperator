#! /bin/bash


red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}


printf "%-65s .......... " "create f2o service"
cp ../../docker-compose.yml /opt/fit2openshift
cp -r  ../service/** /etc/init.d
colorMsg $green "[OK]"

