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
}

download ${openJava} dependencies/server-jre-8u192-linux-x64.tar.gz 
if [ "$?" == "0" ];then
   printf "download ${openJava} error!" 
else
    exit 1
fi
download ${chef} dependencies/chef-14.9.13-1.el7.x86_64.rpm
if [ "$?" == "0" ];then
   printf "download ${cher} error!"
else
    exit 1
fi
download ${nexus} nexus-okd-3.11.tar.gz
if [ "$?" == "0" ];then
   printf "download ${nexus} error!"
else
    exit 1
fi
tar -zvxf nexus-okd-3.11.tar.gz 
if [ "$?" == "0" ];then
else
    exit 1
fi
