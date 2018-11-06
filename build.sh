#!/bin/bash
#

args=$1

image=registry.fit2cloud.com/jumpserver/fit2ansible

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
