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
\cp -fr  ../opt/** /opt
rm -fr /opt/fit2openshift/data/mysql/*.md
chmod 777 -R /opt/fit2openshift/data/mysql/
colorMsg $green "[OK]"
