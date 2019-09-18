from django.db.models import Q

from cloud_provider import get_cloud_client
from cloud_provider.models import Plan, Zone
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from kubeops_api.models.setting import Setting


# def scale_hosts(cluster):
#     worker_size = cluster.worker_size
#     current_worker_size = len(cluster.current_workers)
#     if worker_size < current_worker_size:
#         remove_host(cluster, current_worker_size - worker_size)
#     elif worker_size > current_worker_size:
#         add_new_host(cluster, worker_size - current_worker_size)
#
#
# def add_new_host(cluster, num):
#     hosts = []
#     role = "worker"
#     compute_model = get_k8s_role_model(role, cluster.plan)
#     domain = cluster.name + "." + Setting.objects.get(key="domain_suffix").value
#     for worker in cluster.current_workers:
#         hosts.append({
#             "role": role,
#             "cpu": compute_model['cpu'],
#             "memory": compute_model['memory'] * 1024,
#             "name": worker.name,
#             "domain": domain,
#             "host_name": role + "{}-{}".format(i, cluster.name),
#             "zone_vars": zone_name,
#             "ip": worker.ip
#         })
#
# def remove_host(cluster, num):


def create_hosts(cluster):
    hosts = create_cluster_hosts(cluster)
    mix_vars = cluster.plan.mixed_vars
    mix_vars["hosts"] = hosts
    client = get_cloud_client(mix_vars)
    terraform_result = client.apply_terraform(cluster, mix_vars)
    if not terraform_result:
        raise RuntimeError("create host error!")
    for host in hosts:
        defaults = {
            "name": host['name'],
            "ip": host['ip'],
            "username": 'root',
            "password": 'KubeOperator@2019'
        }
        h = Host.objects.update_or_create(defaults, name=host['name'])
        cluster.create_node(host['role'], h)


def create_cluster_hosts(cluster):
    roles = {
        "master": 1,
        "daemon": 1,
        "worker": cluster.worker_size
    }
    hosts = []
    deploy_vars = cluster.plan.mixed_vars
    deploy_template = cluster.plan.deploy_template
    domain = cluster.name + "." + Setting.objects.get(key="domain_suffix").value
    zones = deploy_vars['zones']
    gen_zone_ip_pool(zones)
    if deploy_template == Plan.DEPLOY_TEMPLATE_MULTIPLE:
        roles['master'] = 3
    for role, size in roles.items():
        compute_model = get_k8s_role_model(role, cluster.plan)
        for i in range(1, size + 1):
            zone = get_zone(zones, i)
            ip = zone['ip_pool'].pop()
            zone_name = zone['name']
            if not ip:
                raise RuntimeError('zone: {}  ip address not enough!', zone_name)
            host = {
                "role": role,
                "cpu": compute_model['cpu'],
                "memory": compute_model['memory'] * 1024,
                "name": role + "{}.".format(i) + "{}".format(domain),
                "domain": domain,
                "ip": ip,
                "zone": zone
            }
            hosts.append(host)
    return hosts


def get_k8s_role_model(role, plan):
    k8s_model = None
    deploy_vars = plan.mixed_vars
    if role == 'master':
        k8s_model = deploy_vars['k8s_master_model']
    if role == 'worker':
        k8s_model = deploy_vars['k8s_worker_model']
    if role == 'daemon':
        k8s_model = deploy_vars['k8s_daemon_model']
    return find_compute_model(k8s_model, plan.compute_models)


def get_zone(zones, index):
    hash = index % len(zones)
    return zones[hash]


def find_compute_model(role, models):
    for model in models:
        if role == model["name"]:
            return model["meta"]


def gen_zone_ip_pool(zones):
    for zone_dic in zones:
        zone = Zone.objects.get(name=zone_dic['zone_name'])
        print(zone.ip_pools())
        zone_dic['ip_pool'] = zone.ip_pools()


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
