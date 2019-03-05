red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}

function success()
{
if [ "$?" != "0" ];then    
    exit 1
else 
    colorMsg $green "[OK]"
fi
}

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


logPath="/opt/fit2openshift/logs/install/"
timestamp=$(date -d now +%F)

fullLogFile=${logPath}"install_"${timestamp}".log"

#install node,docker
#install node
 
printf "%-65s .......... " "Install Node:"
hasNode=`which node 2>&1`
if [[ "${hasNode}" =~ "no node" ]]; then
    download https://nodejs.org/dist/v8.11.2/node-v8.11.2-linux-x64.tar.xz /tmp/node-v8.11.2-linux-x64.tar.xz >>$fullLogFile 2>&1
    cd /tmp \
    && xz -d node-v8.11.2-linux-x64.tar.xz \
    && tar -xvf node-v8.11.2-linux-x64.tar -C  /usr/local/ >>$fullLogFile 2>&1  \
    && ln -s /usr/local/node-v8.11.2-linux-x64/bin/node /usr/bin/node \
    && ln -s /usr/local/node-v8.11.2-linux-x64/bin/npm /usr/bin/npm \
    && rm -fr /tmp/node-v8.11.2-linux-x64.tar.xz \
    && npm i npm@latest -g >>$fullLogFile 2>&1
    success
else 
   colorMsg $green "[OK]"
fi

printf "\n"
printf "%-65s .......... " "Install Docker:"

#install docker
hasDocker=`which docker 2>&1`
if [[ "${hasDocker}" =~ "no docker" ]]; then
    yum install -y docker >>$fullLogFile 2>&1   && /bin/systemctl start docker.service
    success
else 
    colorMsg $green "[OK]"
fi
printf "\n"

printf "%-65s .......... " "Install Docker-compose:"
#install docker-compose
hasDockerCompose=`which docker-compose 2>&1`
if [[ "${hasDockerCompose}" =~ "no docker-compose" ]]; then
    yum install -y docker-compose >>$fullLogFile 2>&1 
    success
else
    colorMsg $green "[OK]"
fi
printf "\n"



printf "%-65s .......... " "Install Angular Cli:"
#install ng

hasNg=`which ng 2>&1`
if [[ "${hasNg}" =~ "no ng" ]]; then
     npm install -g @angular/cli  >>$fullLogFile 2>&1 && ln -s /usr/local/node-v8.11.2-linux-x64/lib/node_modules/@angular/cli/bin/ng /usr/bin/ng 
     success
else 
   colorMsg $green "[OK]"
fi
