import logging

import docker

from fit2ansible.settings import PACKAGE_IMAGE_NAME, PACKAGE_PATH_PREFIX

Logger = logging.getLogger(__name__)


def get_docker_client():
    client = docker.DockerClient(base_url='unix://var/run/docker.sock')
    return client


def is_package_container_exists(package_name):
    result = False
    client = get_docker_client()
    containers = client.containers.list(all=True)
    for container in containers:
        if container.name == package_name:
            result = True
    return result


def is_package_container_start(package_name):
    client = get_docker_client()
    container = client.containers.get(package_name)
    return str(container.status) == 'running'


def create_package_container(package):
    client = get_docker_client()
    image_name = PACKAGE_IMAGE_NAME
    package_path = "{}{}/data".format(PACKAGE_PATH_PREFIX, package.name)
    volumes = {
        package_path: {"bind": "/nexus-data", "mode": "rw"}
    }
    ports = {
        '8081/tcp': package.repo_port,
        '8092/tcp': package.registry_port
    }
    Logger.info("离线包 {} 容器不存在,开始创建容器".format(package.name))
    container = client.containers.run(detach=True, volumes=volumes, name=package.name, ports=ports,
                                      image=image_name)
    return container


def start_package_container(package):
    client = get_docker_client()
    container = client.containers.get(package.name)
    container.start()
    Logger.info("离线包 {} 容器已启动".format(package.name))
    return container


def stop_package_container(package):
    client = get_docker_client()
    container = client.containers.get(package.name)
    container.remove()
    Logger.info("离线包 {} 容易已停止".format(package.name))


def wait_for_container_running(container):
    time = 120
    while time > 0:
        time = time - 1
        container.reload()
        if str(container.status) == 'running':
            break
