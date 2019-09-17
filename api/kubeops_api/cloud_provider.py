from django.db.models import Q

from cloud_provider import get_cloud_client
from cloud_provider.models import TerraformHost, Plan, Zone
from kubeops_api.models.node import Node
from kubeops_api.models.setting import Setting


def create_hosts(cluster):
    terraform_hosts = create_terraform_hosts(cluster)
    cluster.set_terraform_hosts(terraform_hosts)
    cloud_provider = get_cloud_client(cluster.plan.mixed_vars)
    result = cloud_provider.apply_terraform(cluster=cluster)
    if result:
        for h in cluster.terraform_hosts.all():
            h.create_host()
        cluster.create_nodes_by_terraform()
    else:
        for host in terraform_hosts:
            host.delete()
        raise Exception('Create nodes error!')


def add_hosts(cluster, num):
    terraform_hosts = add_terraform_host(cluster, num)
    cluster.set_terraform_hosts(terraform_hosts)
    cloud_provider = get_cloud_client(cluster.plan.mixed_vars)
    result = cloud_provider.apply_terraform(cluster=cluster)
    if result:
        for h in cluster.terraform_hosts.all():
            h.create_host()
        cluster.add_nodes_by_terraform()
    else:
        for host in terraform_hosts:
            host.delete()
        raise Exception('Add nodes error!')


def create_terraform_hosts(cluster):
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
            host = TerraformHost(
                role=role,
                cpu=compute_model['cpu'],
                memory=compute_model["memory"] * 1024,
                name=role + "{}.".format(i) + "{}".format(domain),
                domain=domain,
                short_name=role + "{}".format(i),
                host_name=role + "{}-{}".format(i, cluster.name),
                zone_vars=zone_name,
                ip=ip
            )
            hosts.append(host)
    return TerraformHost.objects.bulk_create(hosts)


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


def add_terraform_host(cluster, num):
    deploy_vars = cluster.plan.mixed_vars
    zones = deploy_vars['zones']
    gen_zone_ip_pool(zones)
    compute_model = get_k8s_role_model('worker', cluster.plan)
    domain = cluster.name + "." + Setting.objects.get(key="domain_suffix").value
    role = 'worker'
    hosts = []
    for i in range(1, num + 1):
        no = cluster.node_size + i
        zone = get_zone(zones, no)
        zone_name = zone['zone_name']
        ip = zone['ip_pool'].pop()
        if not ip:
            raise RuntimeError('zone: {}  ip address not enough!', zone_name)
        host = TerraformHost(
            role=role,
            cpu=compute_model['cpu'],
            memory=compute_model["memory"] * 1024,
            name=role + "{}.".format(no) + "{}".format(domain),
            domain=domain,
            short_name=role + "{}".format(no),
            host_name=role + "{}-{}".format(no, cluster.name),
            zone_vars=zone,
            ip=ip
        )
        hosts.append(host)
    return TerraformHost.objects.bulk_create(hosts)
