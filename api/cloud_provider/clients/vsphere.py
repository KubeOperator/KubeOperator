import json
import os

from pyVmomi.VmomiSupport import VmomiJSONEncoder

from cloud_provider.cloud_client import CloudClient
from pyVim import connect
from pyVmomi import vim

from fit2ansible.settings import RESOURCE_DIR


class VsphereCloudClient(CloudClient):
    cloud_config_path = os.path.join(RESOURCE_DIR, 'vsphere')

    def list_region(self):
        params = replace_params(self.vars)
        st = get_service_instance(params)
        content = st.RetrieveContent()
        container = content.rootFolder
        viewType = [vim.Datacenter]
        regions = get_obj_list(content, viewType, container)
        data = []
        for region in regions:
            data.append(region.name)
        return data

    def list_zone(self, region):
        params = replace_params(self.vars)
        st = get_service_instance(params)
        content = st.RetrieveContent()
        container = content.rootFolder
        viewType = [vim.Datacenter]
        region = get_obj(content, viewType, container, region)
        zones = []
        for entity in region.hostFolder.childEntity:
            zone = {
                "storages": [],
                "networks": [],
                "images": []
            }
            if isinstance(entity, vim.ClusterComputeResource):
                zone["name"] = entity.name
                for network in entity.network:
                    zone.get("networks").append(network.name)
                for datastore in entity.datastore:
                    zone.get("storages").append(datastore.name)
                for host in entity.host:
                    for vm in host.vm:
                        if vm.summary.config.template:
                            zone.get('images').append(vm.name)
                zones.append(zone)
        return zones


def get_obj(content, vimtype, folder, name):
    obj = None
    container = content.viewManager.CreateContainerView(folder, vimtype, True)
    for c in container.view:
        if c.name == name:
            obj = c
    return obj


def get_obj_list(content, vimtype, folder):
    objs = []
    container = content.viewManager.CreateContainerView(folder, vimtype, True)
    for c in container.view:
        objs.append(c)
    return objs


def get_service_instance(kwargs):
    host = kwargs.get('host')
    username = kwargs.get('username')
    password = kwargs.get('password')
    service_instance = connect.SmartConnectNoSSL(host=host, user=username, pwd=password, port=int(443))
    if not service_instance:
        raise Exception('Could not connect to the specified host using specified username and password')
    return service_instance


def replace_params(vars):
    return {
        "host": vars.get('vc_host', None),
        "username": vars.get('vc_username', None),
        "password": vars.get('vc_password', None),
    }
