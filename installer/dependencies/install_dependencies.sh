red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

colorMsg $blue "INSTALL BUILD TOOLS"

printf "%-65s .......... " "Install Node.js:"

#install node,docker

#install node
node -v
if [ "$?" == "0" ];then
    echo "Skip install node..."
else
    wget https://nodejs.org/dist/v8.11.2/node-v8.11.2-linux-x64.tar.xz -O /tmp/node-v8.11.2-linux-x64.tar.xz \
    && cd /tmp \
    && xz -d node-v8.11.2-linux-x64.tar.xz \
    && tar -xvf node-v8.11.2-linux-x64.tar -C  /usr/local/ \
    && ln -s /usr/local/node-v8.11.2-linux-x64/bin/node /usr/bin/node \
    && ln -s /usr/local/node-v8.11.2-linux-x64/bin/npm /usr/bin/npm \
    && rm -fr /tmp/node-v8.11.2-linux-x64.tar.xz	
fi
printf "%-65s .......... " "Install Docker::"

#install docker
docker -v
if [ "$?" == "0" ];then
    echo "Skip install docker..."
else
    yum install docker && service docker start
fi

colorMsg $green "[OK]"
