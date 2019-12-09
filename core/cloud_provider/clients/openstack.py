import os
from cloud_provider.cloud_client import CloudClient
from keystoneclient.v3 import client as KeystoneClient
from openstack import connection
from urllib import request
import json

from cloud_provider.utils import download_plugins
from fit2ansible.settings import CLOUDS_RESOURCE_DIR
from kubeops_api.models.setting import Setting


class OpenStackCloudClient(CloudClient):
    cloud_config_path = os.path.join(CLOUDS_RESOURCE_DIR, 'openstack')
    working_path = None

    def list_region(self):
        keystone_client = get_keystone_client(self.vars)
        req = request.Request(self.vars.get('openstack_identity') + '/regions', headers={
            "X-Auth-Token": keystone_client.auth_token
        })
        response = request.urlopen(req).read()
        regions = json.loads(response)["regions"]
        data = []
        for region in regions:
            data.append(region["id"])
        return data

    def list_zone(self, region):
        openstack_client = get_openstack_client(self.vars, region)
        openstack_zones = openstack_client.list_availability_zone_names()
        openstack_networks = openstack_client.list_networks(get_filter(self.vars))
        openstack_floating_networks = openstack_client.list_networks({
            "router:external": True
        })
        openstack_security_groups = openstack_client.list_security_groups(get_filter(self.vars))
        openstack_datastores = openstack_client.list_volume_types()
        openstack_client.close()
        zones = []
        for zone_item in openstack_zones:
            zone = {"storages": [], "networkList": [], "floatingNetworkList": [], "securityGroups": [], "cluster": zone_item, "ipType": []}
            for floating_network in openstack_floating_networks:
                zone.get("floatingNetworkList").append({
                    "name": floating_network.name,
                    "id": floating_network.id,
                })
            for network in openstack_networks:
                subnetList = []
                openstack_subnets = openstack_client.list_subnets(get_filter(self.vars, network_id=network.id))
                for subnet in openstack_subnets:
                    subnetList.append({
                        "name": subnet.name,
                        "id": subnet.id,
                    })
                zone.get("networkList").append({
                    "name": network.name,
                    "id": network.id,
                    "subnetList": subnetList
                    })
            for sg in openstack_security_groups:
                zone.get("securityGroups").append(sg.name)
            for datastore in openstack_datastores:
                zone.get("storages").append({
                    "name": datastore.name,
                    "type": datastore.id
                })
            zone.get("ipType").append("private")
            zone.get("ipType").append("floating")
            zones.append(zone)
        return zones

    def get_flavors(self, region):
        openstack_client = get_openstack_client(self.vars, region)
        flavors = openstack_client.list_flavors()
        models = []
        for flavor in flavors:
            if flavor.disk == 40 and flavor.vcpus == 2 and flavor.ram / 1024 == 4:
                model = {
                    'name': flavor.name,
                    'meta': {
                        'id': flavor.id,
                        'cpu': 2,
                        'memory': 4,
                        'disk': 100
                    }
                }
                models.append(model)
                continue
            if flavor.disk == 100 and flavor.vcpus == 4 and flavor.ram / 1024 == 16:
                model = {
                    'name': flavor.name,
                    'meta': {
                        'id': flavor.id,
                        'cpu': 4,
                        'memory': 16,
                        'disk': 100
                    }
                }
                models.append(model)
                continue
            if flavor.disk == 100 and flavor.vcpus == 8 and flavor.ram / 1024 == 32:
                model = {
                    'name': flavor.name,
                    'meta': {
                        'id': flavor.id,
                        'cpu': 8,
                        'memory': 32,
                        'disk': 100
                    }
                }
                models.append(model)
                continue
            if flavor.disk == 100 and flavor.vcpus == 16 and flavor.ram / 1024 == 64:
                model = {'name': flavor.name,
                         'meta': {
                             'id': flavor.id,
                             'cpu': 16,
                             'memory': 64,
                             'disk': 100
                         }}
                models.append(model)
                continue
            if flavor.disk == 100 and flavor.vcpus == 22 and flavor.ram / 1024 == 128:
                model = {
                    'name': flavor.name,
                    'meta': {
                        'id': flavor.id,
                        'cpu': 32,
                        'memory': 128,
                        'disk': 100
                    }
                }
                models.append(model)
                continue
            if flavor.disk == 100 and flavor.vcpus == 64 and flavor.ram / 1024 == 256:
                model = {
                    'name': flavor.name,
                    'meta': {
                        'id': flavor.id,
                        'cpu': 64,
                        'memory': 256,
                        'disk': 100
                    }
                }
                models.append(model)
                continue
        return models

    def init_terraform(self, cluster):
        plugin_dir = os.path.join(self.working_path, '.terraform', 'plugins')
        if not os.path.exists(plugin_dir):
            os.makedirs(plugin_dir)
        hostname = Setting.objects.get(key='local_hostname').value
        port = cluster.package.repo_port
        url = "http://{}:{}/repository/raw/terraform/openstack.zip".format(hostname, port)
        download_plugins(url=url, target=plugin_dir)

    def apply_terraform(self, cluster, hosts_dict):
        return super().apply_terraform(cluster, hosts_dict)

    def create_image(self, zone):
        openstack_client = get_openstack_client(self.vars, zone.region.cloud_region)
        image = openstack_client.get_image(zone.region.image_name)
        if image is None:
            image = openstack_client.create_image(name=zone.region.image_name, filename=zone.region.image_vmdk_path,
                                          disk_format='qcow2')

def replace_params(vars):
    return {
        'uth_url': vars.get('openstack_identity', None),
        'username': vars.get('openstack_username', None),
        'password': vars.get('openstack_password', None),
        'project_id': vars.get('openstack_projectId', None),
        'user_domain_name': vars.get('openstack_domain_name', None),
    }


def get_keystone_client(vars):
    keystone = KeystoneClient.Client(auth_url=str.strip(vars.get('openstack_identity', '')),
                                     username=vars.get('openstack_username', None),
                                     password=vars.get('openstack_password', None),
                                     project_id=vars.get('openstack_projectId', None),
                                     user_domain_name=vars.get('openstack_domain_name', None))
    return keystone


def get_openstack_client(vars, region):
    client = connection.Connection(auth_url=str.strip(vars.get('openstack_identity', '')),
                                   username=vars.get('openstack_username', None),
                                   password=vars.get('openstack_password', None),
                                   project_id=vars.get('openstack_projectId', None),
                                   user_domain_name=vars.get('openstack_domain_name', None),
                                   region_name=region)
    return client


def get_filter(vars, **kwargs):
    filters = {
        'project_id': vars.get('openstack_projectId', None),
        'tenant_id': vars.get('openstack_projectId', None)
    }
    if kwargs.get('network_id'):
        filters['network_id'] = kwargs.get('network_id')
    return filters

