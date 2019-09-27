from django.db.models import Q

from cloud_provider import get_cloud_client
from cloud_provider.models import Plan, Zone
from kubeops_api.adhoc import drain_worker_node
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from kubeops_api.models.setting import Setting


def create_hosts(cluster):
    hosts = create_cluster_hosts(cluster)
    mix_vars = cluster.plan.mixed_vars
    mix_vars["hosts"] = hosts
    client = get_cloud_client(mix_vars)
    terraform_result = client.apply_terraform(cluster, mix_vars)
    if not terraform_result:
        raise RuntimeError("create host error!")
    for host in hosts:
        zone = Zone.objects.get(name=host["zone_name"])
        defaults = {
            "name": host['name'],
            "ip": host['ip'],
            "zone": zone,
            "username": 'root',
            "password": 'KubeOperator@2019'
        }
        h = Host.objects.update_or_create(defaults, name=host['name'])
        cluster.create_node(host['role'], h[0])


def is_worker(host):
    return host['role'] == 'worker'


def is_master(host):
    return host['role'] == 'master'


def scale_up(cluster, num):
    worker_hosts = cluster.get_current_worker_hosts()
    worker_size = len(worker_hosts)

    new_hosts = []
    worker_hosts_new = []
    master_hosts_new = []
    remove_list = []
    add_list = []
    if worker_size > num:
        hosts = create_cluster_hosts(cluster)
        worker_hosts_new = list(filter(is_worker, hosts))
        master_hosts_new = list(filter(is_master, hosts))
        for i in range(worker_size - num):
            rm_worker = worker_hosts_new.pop()
            remove_list.append(rm_worker)
        cluster.worker_size = num
        drain_workers(cluster, remove_list)
    elif worker_size < num:
        cluster.worker_size = num
        hosts = create_cluster_hosts(cluster)
        for h in hosts:
            if h.get('new', None):
                add_list.append(h)
    new_hosts.extend(worker_hosts_new)
    new_hosts.extend(master_hosts_new)
    mix_vars = cluster.plan.mixed_vars
    mix_vars["hosts"] = new_hosts
    client = get_cloud_client(mix_vars)
    terraform_result = client.apply_terraform(cluster, mix_vars)
    if not terraform_result:
        raise RuntimeError("create host error!")
    for host in add_list:
        zone = Zone.objects.get(name=host["zone_name"])
        defaults = {
            "name": host['name'],
            "ip": host['ip'],
            "zone": zone,
            "username": 'root',
            "password": 'KubeOperator@2019'
        }
        h = Host.objects.update_or_create(defaults, name=host['name'])
        node = cluster.create_node(host['role'], h[0])
        cluster.add_to_new_node(node)
    for host in remove_list:
        cluster.change_to()
        node = Node.objects.get(name=host['name'])
        node.host.delete()
    cluster.save()


def drain_workers(cluster, remove_list):
    master = cluster.get_first_master()
    for host in remove_list:
        drain_worker_node(master, host.name)


def create_cluster_hosts(cluster):
    roles = {
        "master": 1,
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
            name = role + "{}.".format(i) + "{}".format(domain)
            zone = get_zone(zones, i)
            host = {
                "role": role,
                "cpu": compute_model['cpu'],
                "memory": compute_model['memory'] * 1024,
                "name": name,
                "short_name": role + "{}".format(i),
                "domain": domain,
                "zone": zone,
                "zone_name": zone['zone_name'],
            }
            h = None
            try:
                h = Host.objects.get(name=name)
            except Exception as e:
                pass
            if h:
                host.update({
                    "ip": h.ip
                })
            else:
                host.update({
                    "ip": zone['ip_pool'].pop(),
                    "new": True
                })
            if not host['ip']:
                raise RuntimeError('zone: {}  ip address not enough!', zone['zone_name'])
            hosts.append(host)
    return hosts


def get_k8s_role_model(role, plan):
    k8s_model = None
    deploy_vars = plan.mixed_vars
    if role == 'master':
        k8s_model = deploy_vars['k8s_master_model']
    if role == 'worker':
        k8s_model = deploy_vars['k8s_worker_model']
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
            node.host.delete()
