#!/bin/bash
BASE_DIR=$(cd "$(dirname "$0")";pwd)
source ${BASE_DIR}/scripts/utils.sh
PROJECT_DIR=${BASE_DIR}
SCRIPT_DIR=${BASE_DIR}/scripts
action=$1
target=$2
args=$@


function usage() {
   echo "Fit2Openshift 部署安装脚本"
   echo
   echo "Usage: "
   echo "  f2octl [COMMAND] [ARGS...]"
   echo "  f2octl --help"
   echo
   echo "Commands: "
   echo "  install 部署安装 fit2openshift"
   echo "  start 启动 fit2openshift"
   echo "  restart [service] 重启, 并不会重建服务容器"
   echo "  reload [service] 重建容器如何需要并重启服务"
   echo "  status 查看Jumpserver状态"
   echo "  down [service] 删掉容器 不带参数删掉所有registry"
   echo "  python 进入api, 运行 python manage.py shell"
   echo "  db 连接数据库"
   echo "  ... 其他docker-compose执行的命令 如 logs 等等"
}


function service_to_docker_name() {
    service=$1
    if [[ "${service:0:3}" != "f2o" ]];then
        service=jms_${service}
    fi
    echo ${service}
}

function main() {
    EXE=docker-compose
    case "${action}" in
      install)
         bash ${SCRIPT_DIR}/5_install.sh
         ;;
      start)
         ${EXE} up -d
         ;;
      restart)
         ${EXE} restart ${target}
         ;;
      reload)
         ${EXE} up -d &> /dev/null
         ${EXE} restart ${target}
         ;;
      status)
         ${EXE} ps
         ;;
      down)
         if [[ -z "${target}" ]];then
             ${EXE} down
         else
             ${EXE} stop ${target} && ${EXE} rm ${target}
         fi
         ;;
      tail)
          if [[ -z "${target}" ]];then
              ${EXE} logs --tail 100 -f
          else
              docker_name=$(service_to_docker_name ${target})
              docker logs -f ${docker_name} --tail 100
          fi
          ;;
      python)
          docker exec -it f2o_api python manage.py shell
          ;;
      db)
          docker exec -it f2o_api python manage.py dbshell
          ;;
      exec)
          docker_name=$(service_to_docker_name ${target})
          docker exec -it ${docker_name} sh
          ;;
      help)
         usage
         ;;
      --help)
         usage
         ;;
      *)
         ${EXE} ${args}
    esac
}

main