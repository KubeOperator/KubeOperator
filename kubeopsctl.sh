BASE_DIR=$(cd "$(dirname "$0")";pwd)
PROJECT_DIR=${BASE_DIR}
SCRIPT_DIR=${BASE_DIR}/scripts
action=$1
target=$2
args=$@


function usage() {
   echo_green "kubeOperator 部署安装脚本"
   echo
   echo_green "Usage: "
   echo_green "  kubeopsctl [COMMAND] [ARGS...]"
   echo_green "  kubeopsctl --help"
   echo_green
   echo_green "Commands: "
   echo_green "  install [安装路径] 部署安装 kubeOperator"
   echo_red   "  如果自定义安装路径请务必指定安装路径，默认/opt"
   echo_green "  start 启动 kubeOperator"
   echo_green "  uninstall 卸载 kubeOperator"
   echo_green "  upgrade [安装路径] 升级 kubeOperator"
   echo_red   "  如果自定义安装路径请务必指定安装路径，默认/opt"
   echo_green "  restart [service] 重启, 并不会重建服务容器"
   echo_green "  reload [service] 重建容器如何需要并重启服务"
   echo_green "  status 查看 kubeOperator 状态"
   echo_green "  down [service] 删掉容器 不带参数删掉所有"
   echo_green "  python 进入 api, 运行 python manage.py shell"
   echo_green "  db 连接数据库"
   echo_green "  ... 其他 docker-compose 执行的命令 如 logs 等等"
}


function service_to_docker_name() {
    service=$1
    if [[ "${service:0:3}" != "f2o" ]];then
        service=jms_${service}
    fi
    echo ${service}
}

function main() {
    if [[ "${action}" == "install" || "${action}" == "upgrade" ]];then
      if [[ -n $target ]];then
         sed_pattern="s|^export INSTALL_DIR=.*$|export INSTALL_DIR=\"${target}\"|g" 
         utils_file="${BASE_DIR}/scripts/utils.sh"
         sed -i -r "${sed_pattern}" ${utils_file}
      fi
    fi
    source ${BASE_DIR}/scripts/utils.sh
    EXE="docker-compose -f ${INSTALL_DIR}/kubeoperator/docker-compose.yml"
    case "${action}" in
      install)
         bash ${SCRIPT_DIR}/5_install.sh
         ;;
      uninstall)
         bash ${SCRIPT_DIR}/6_uninstall.sh
         ;;
      upgrade)
         bash ${SCRIPT_DIR}/7_upgrade.sh
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
          docker exec -it kubeops_api python manage.py shell
          ;;
      db)
          docker exec -it kubeops_api python manage.py dbshell
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

