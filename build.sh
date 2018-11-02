#!/bin/bash
#

args=$1

image=registry.fit2cloud.com/jumpserver/ansible_ui

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
