import os

from django.db.models import Q

from cloud_provider import get_cloud_client
from cloud_provider.models import TerraformHost
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from kubeops_api.models.setting import Setting


def create_hosts(cluster):
    terraform_hosts = generate_host_model(cluster)
    cluster.set_terraform_hosts(terraform_hosts)
    cloud_provider = get_cloud_client(cluster.plan.mixed_vars)
    result = cloud_provider.apply_terraform(cluster=cluster)
    if result:
        for h in cluster.terraform_hosts.all():
            h.create_host()
        cluster.create_nodes_by_terraform()
    else:
        raise Exception('Create nodes error!')


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
    deploy_model = deploy_vars.get('k8s_deploy_model')
    domain = cluster.name + "." + Setting.objects.get(key="domain_suffix").value
    total_size = 0
    if deploy_model == 'multiple':
        roles["master"] = 3
    for role, size in roles.items():
        total_size = total_size + size
        role_model = get_k8s_role_model(role, deploy_vars)
        role_compute_model = find_compute_model(role_model, cluster.plan.compute_models)
        for i in range(0, size):
            host = TerraformHost(
                role=role,
                short_name=role + "{}".format(i + 1),
                host_name=role + "{}-{}".format(i + 1, cluster.name),
                name=role + "{}.".format(i + 1) + "{}".format(domain),
                domain=domain,
                folder=deploy_vars.get("vc_folder"),
                cpu=role_compute_model["cpu"],
                memory=role_compute_model["memory"] * 1024
            )
            hosts.append(host)
    available_ips = get_available_ips(ip_start, ip_end)
    if not total_size > len(available_ips):
        for no, host in enumerate(hosts):
            host.ip = available_ips[no]
    else:
        raise Exception("{} ip address not enough to create cluster".format(len(available_ips)))
    return TerraformHost.objects.bulk_create(hosts)


def find_compute_model(role, models):
    for model in models:
        if role == model["name"]:
            return model["meta"]


def get_k8s_role_model(role, deploy_vars):
    k8s_model = None
    if role == 'master':
        k8s_model = deploy_vars['k8s_master_model']
    if role == 'worker':
        k8s_model = deploy_vars['k8s_worker_model']
    if role == 'daemon':
        k8s_model = deploy_vars['k8s_daemon_model']
    return k8s_model


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


def delete_hosts(cluster):
    cloud_provider = get_cloud_client(cluster.plan.mixed_vars)
    result = cloud_provider.destroy_terraform(cluster.name)
    if not result:
        raise Exception('Destroy nodes error! ')
    else:
        cluster.change_to()
        nodes = Node.objects.filter(~Q(name__in=['::1', '127.0.0.1', 'localhost']))
        for node in nodes:
            node.delete()
        for host in cluster.terraform_hosts.all():
            host.host.delete()
