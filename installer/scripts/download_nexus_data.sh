#!/bin/bash
red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

printf "%-65s .......... " "copy fit2openshift files"




nexus="http://fit2openshift.oss-cn-beijing.aliyuncs.com/okd-3.11//tmp/nexus-data.tar.gz"
if [ ! -d dependencies  ];then
  mkdir dependencies
else
  echo "dependencies dir exist skip it ..."
fi


function download()
{
    wget $1 -O $2
    size_total=`curl -Is $1 | grep Content-Length  | awk -F': ' '{print $2}' | tr -d '\r'`
    size_current=`du -b  $2 | awk -F' ' '{print $1}'`

if [ "x$size_total" == "x$size_current" ];then
      printf "download $1 success!"
    else
      exit 1
    fi
}

download ${nexus} /tmp/nexus-data.tar.gz \
&& tar -zvxf /tmp/nexus-data.tar.gz -C /opt/fit2openshift/data/nexus \
&& chmox -R 777 /opt/fit2openshift/data/nexus 
if [ "$?" != "0" ];then
    colorMsg $red "[Defeat]"
    exit 1
fi
colorMsg $green "[OK]"
