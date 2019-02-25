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
if [ "$?" == "0" ];then    
    exit 1
else 
    colorMsg $green "[OK]"
fi
}




logPath="/opt/fit2openshift/logs/install/"
timestamp=$(date -d now +%F)
errorLogFile=${logPath}"error/install_error_"${timestamp}".log"
infoLogFile=${logPath}"info/install_info_"${timestamp}".log"
fullLogFile=${logPath}"install_"${timestamp}".log"

#install node,docker
#install node
 
printf "%-65s .......... " "Install Node:"
hasNode=`which node 2>&1`
if [[ "${hasNode}" =~ "no node" ]]; then
    wget https://nodejs.org/dist/v8.11.2/node-v8.11.2-linux-x64.tar.xz -O /tmp/node-v8.11.2-linux-x64.tar.xz 1>>$infoLogFile 2>>$errorLogFile  \
    && cd /tmp \
    && xz -d node-v8.11.2-linux-x64.tar.xz \
    && tar -xvf node-v8.11.2-linux-x64.tar -C  /usr/local/ 1>>$infoLogFile 2>>$errorLogFile >>$fullLogFile 1<&2 \
    && ln -s /usr/local/node-v8.11.2-linux-x64/bin/node /usr/bin/node \
    && ln -s /usr/local/node-v8.11.2-linux-x64/bin/npm /usr/bin/npm \
    && rm -fr /tmp/node-v8.11.2-linux-x64.tar.xz \
    && npm i npm@latest -g 1>>$infoLogFile 2>>$errorLogFile
    success
else 
   colorMsg $green "[OK]"
fi

printf "\n"
printf "%-65s .......... " "Install Docker::"

#install docker
hasDocker=`which docker 2>&1`
if [[ "${hasDocker}" =~ "no docker" ]]; then
    yum install docker 1>>$infoLogFile 2>>$errorLogFile   && service docker start
    success
else 
    colorMsg $green "[OK]"
fi
printf "\n"

printf "%-65s .......... " "Install Angular Cli:"
#install ng
hasNg=`which docker 2>&1`
if [[ "${hasNg}" =~ "no ng" ]]; then
     npm install -g @angular/cli  1>>$infoLogFile 2>>$errorLogFile 
     success
else 
   colorMsg $green "[OK]"
fi
