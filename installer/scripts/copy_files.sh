#! /bin/bash

red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

printf "%-65s .......... " "copy fit2openshift files"
cp -r  ../opt/** /opt
colorMsg $green "[OK]"
