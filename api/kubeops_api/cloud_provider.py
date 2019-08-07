import os

from ansible_api import ctx
from ansible_api.models import Project
from cloud_provider import get_cloud_client
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from kubeops_api.models.setting import Setting


def create_hosts(cluster):
    vars = cluster.plan.mixed_vars
    vars['host_username'] = 'root'
    vars['host_password'] = 'Calong@2015'
    vars["hosts"] = generate_host_model(cluster)
    cloud_provider = get_cloud_client(vars)

    # 创建hosts
    for _host in vars["hosts"]:
        name = _host['h_name']
        defaults = {
            "name": name,
            "ip": _host['ip'],
            "username": vars['host_username'],
            "password": vars['host_password']
        }
        Host.objects.update_or_create(defaults=defaults, name=name)
    # 创建远程机器
    result = cloud_provider.apply_terraform(cluster=cluster.name, vars=vars)
    if result:
        # 创建 node
        for _host in vars["hosts"]:
            host = Host.objects.get(name=_host['h_name'])
            host.gather_info()
            cluster.change_to()
            node = Node.objects.create(
                name=_host['name'],
                host=host
            )
            node.set_groups(group_names=[_host['role']])
            print('node:  {} created !'.format(node.name))
    else:
        raise RuntimeError('Create hosts error')


def generate_host_model(cluster):
    roles = {
        "master": 1,
        "daemon": 1,
        "worker": cluster.worker_size
    }
    hosts = []
    deploy_vars = cluster.plan.mixed_vars
    ip_start = deploy_vars.get('vc_ip_start')
    ip_end = deploy_vars.get('vc_ip_end')
    deploy_model = deploy_vars.get('k8s_master_model')
    domain = cluster.name + "." + Setting.objects.get(key="domain_suffix").value
    total_size = 0
    if deploy_model == 'multiple':
        roles["master"] = 3
    for role, size in roles.items():
        total_size = total_size + size
        role_compute_model = find_compute_model(deploy_model, cluster.plan.compute_models)
        for i in range(0, size):
            host = {
                "role": role,
                "s_name": role + "{}".format(i),
                "h_name": role + "{}-{}".format(i, cluster.name),
                "name": role + "{}.".format(i) + "{}".format(domain),
                "domain": domain,
                "folder": deploy_vars.get("vc_folder"),
                "cpu": role_compute_model["cpu"],
                "memory": role_compute_model["memory"] * 1024
            }
            hosts.append(host)
    available_ips = get_available_ips(ip_start, ip_end)
    if not total_size > len(available_ips):
        for no, host in enumerate(hosts):
            host["ip"] = available_ips[no]
    else:
        raise Exception("{} ip address not enough to create cluster".format(len(available_ips)))
    return hosts


def find_compute_model(role, models):
    for model in models:
        if role == model["name"]:
            return model["meta"]


def get_available_ips(start, end):
    sub_start = int(start.split('.')[3])
    sub_end = int(end.split('.')[3])
    ip_prefix = start[0:start.index(start.split('.')[3])]
    ip_list = []
    for i in range(sub_start, sub_end):
        ip_list.append(ip_prefix + str(i))
    hosts = Host.objects.filter(ip__in=ip_list)
    for host in hosts:
        ip_list.remove(host.ip)
    return ip_list
