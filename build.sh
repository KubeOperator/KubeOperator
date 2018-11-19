#!/bin/bash
#

args=$1

image=registry.fit2cloud.com/fit2anything/fit2openshift

function build(){
   docker build -t $image .
}

function push(){
  docker push $image
}

build

if [ -n "$args" ];then
    push
fi
