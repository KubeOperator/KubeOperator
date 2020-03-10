from django.db.models import Q

from cloud_provider import get_cloud_client
from cloud_provider.compute_model import get_compute_model_meta
from cloud_provider.models import Plan, Zone
from kubeops_api.adhoc import drain_worker_node
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from time import sleep
from kubeops_api.models.item_resource import ItemResource

def create_compute_resource(cluster):
    hosts_dict = create_cluster_hosts_dict(cluster)
    create_nodes(cluster, hosts_dict)


def scale_compute_resource(cluster, num):
    worker_size = cluster.worker_size
    change_list = []
    hosts_dict = []
    delete = False
    add = False
    if worker_size == num:
        return
    if worker_size > num:
        _hosts_dict = create_cluster_hosts_dict(cluster)
        result = create_cluster_scale_down_hosts_dict(_hosts_dict, worker_size - num)
        change_list = result['change_list']
        hosts_dict = result['hosts_dict']
        drain_workers(cluster, change_list)
        cluster.worker_size = num
        delete = True
    elif worker_size < num:
        cluster.worker_size = num
        _hosts_dict = create_cluster_hosts_dict(cluster)
        result = create_cluster_scale_up_hosts_dict(_hosts_dict)
        change_list = result['change_list']
        hosts_dict = result['hosts_dict']
        add = True
    create_nodes(cluster, hosts_dict)
    cluster.save()
    for host_dict in change_list:
        host = Host.objects.get(name=host_dict['name'])
        if delete:
            ItemResource.objects.filter(resource_id=host.id.hex).delete()
            host.delete()
        if add:
            cluster.add_to_new_node(host.node)


def create_cluster_scale_down_hosts_dict(hosts_dict, num):
    change_list = []
    worker_hosts_dict = list(filter(is_worker, hosts_dict))
    master_hosts_dict = list(filter(is_master, hosts_dict))
    for i in range(num):
        rm_worker = worker_hosts_dict.pop()
        change_list.append(rm_worker)
    hosts_dict = []
    hosts_dict.extend(master_hosts_dict)
    hosts_dict.extend(worker_hosts_dict)
    return {
        "hosts_dict": hosts_dict,
        "change_list": change_list
    }


def create_cluster_scale_up_hosts_dict(hosts_dict):
    change_list = []
    worker_hosts_dict = list(filter(is_worker, hosts_dict))
    master_hosts_dict = list(filter(is_master, hosts_dict))
    for host_dict in hosts_dict:
        if host_dict.get('new', None):
            change_list.append(host_dict)
    hosts_dict = []
    hosts_dict.extend(master_hosts_dict)
    hosts_dict.extend(worker_hosts_dict)
    return {
        "hosts_dict": hosts_dict,
        "change_list": change_list
    }


def create_nodes(cluster, hosts_dict):
    hosts = []
    new_nodes = []
    for host_dict in hosts_dict:
        zone = Zone.objects.get(name=host_dict["zone_name"])
        defaults = {
            "name": host_dict['name'],
            "ip": host_dict['ip'],
            "zone": zone,
            "status": Host.HOST_STATUS_CREATING,
            "auto_gather_info": False
        }
        if host_dict.get('new', False):
            result = Host.objects.update_or_create(defaults, name=host_dict['name'])
            host = result[0]
            node = cluster.create_node(host_dict['role'], host)
            item_resource = ItemResource.objects.get(resource_id=cluster.id)
            item_r = ItemResource(item_id=item_resource.item_id,resource_id=host.id,resource_type=ItemResource.RESOURCE_TYPE_HOST)
            item_r.save()
            new_nodes.append(node)
            hosts.append(host)
    client = get_cloud_client(cluster.plan.mixed_vars)
    terraform_result = client.apply_terraform(cluster, hosts_dict)
    if not terraform_result:
        for node in new_nodes:
            node.host.delete()
            raise RuntimeError("create host error!")
    if cluster.plan.mixed_vars.get('provider') == 'openstack':
        print("sleep 20s,等待sshd服务可用")
        sleep(20)
    for host in hosts:
        host.gather_info(retry=5)


def is_worker(host):
    return host['role'] == 'worker'


def is_master(host):
    return host['role'] == 'master'


def create_cluster_hosts_dict(cluster):
    roles = {
        "master": 1,
        "worker": cluster.worker_size
    }
    hosts = []
    deploy_template = cluster.plan.deploy_template
    domain = cluster.name + '.' + cluster.cluster_doamin_suffix
    if deploy_template == Plan.DEPLOY_TEMPLATE_MULTIPLE:
        roles['master'] = 3
    for role, size in roles.items():
        compute_model_name = cluster.plan.compute_models[role]
        compute_model = {}
        if cluster.plan.region.template.name == 'openstack':
            for model in cluster.plan.vars['compute_models']:
                if model['name'] == compute_model_name:
                    compute_model = model['meta']
        else:
            compute_model = get_compute_model_meta(compute_model_name)
        for i in range(1, size + 1):
            name = role + "{}.".format(i) + "{}".format(domain)
            zone = get_zone(cluster.plan.get_zones(), i)
            ## 选择到了zone 后更新 zone参数
            if zone:
                cluster.configs.update(zone.vars)
                cluster.save()
            if not zone:
                raise RuntimeError('Can not find  available ip address!')
            host = {
                "role": role,
                "cpu": compute_model['cpu'],
                "memory": compute_model['memory'] * 1024,
                "name": name,
                "short_name": role + "{}".format(i),
                "domain": domain,
                "zone": zone.to_dict(),
                "zone_name": zone.name,
            }
            host_set = Host.objects.filter(name=name)
            if host_set:
                host.update({
                    "ip": host_set.first().ip
                })
            else:
                host.update({
                    "ip": zone.allocate_ip(),
                    "new": True
                })
            hosts.append(host)
    return hosts


def drain_workers(cluster, remove_list):
    master = cluster.get_first_master()
    for host in remove_list:
        drain_worker_node(master, host['name'])


def get_zone(zones, index):
    select_zone = None
    zone_size = len(zones)
    if zone_size == 1:
        select_zone = zones[0]
    elif zone_size > 1:
        hash = index % len(zones)
        if zones[hash].ip_available_size() > 0:
            select_zone = zones[hash]
        else:
            zones.pop(zones[hash])
            for zone in zones:
                if zone.ip_available_size() > 0:
                    select_zone = zone
    return select_zone


def delete_hosts(cluster):
    cloud_provider = get_cloud_client(cluster.plan.mixed_vars)
    result = cloud_provider.destroy_terraform(cluster.name)
    if not result:
        raise Exception('Destroy nodes error! ')
    else:
        cluster.change_to()
        nodes = Node.objects.filter(~Q(name__in=['::1', '127.0.0.1', 'localhost']))
        for node in nodes:
            ItemResource.objects.filter(resource_id=node.host.id.hex).delete()
            node.host.delete()
