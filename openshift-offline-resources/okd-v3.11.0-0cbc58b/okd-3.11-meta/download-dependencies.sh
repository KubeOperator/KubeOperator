#!/bin/bash


openJava="http://fit2cloud-openshift.oss-cn-beijing.aliyuncs.com/okd-3.11/server-jre-8u192-linux-x64.tar.gz"
chef="http://fit2cloud-openshift.oss-cn-beijing.aliyuncs.com/okd-3.11/chef-14.9.13-1.el7.x86_64.rpm"
nexus="http://fit2cloud-openshift.oss-cn-beijing.aliyuncs.com/okd-3.11/nexus-okd-3.11.tar.gz"
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

if [ "$size_total" -eq "$size_current" ];then
      printf "download $1 success!"
    else
      exit 1
    fi
}

download ${chef} dependencies/chef-14.9.13-1.el7.x86_64.rpm
download ${openJava} dependencies/server-jre-8u192-linux-x64.tar.gz 
download ${nexus} nexus-okd-3.11.tar.gz
tar -zvxf nexus-okd-3.11.tar.gz 
if [ "$?" != "0" ];then
    exit 1
fi
