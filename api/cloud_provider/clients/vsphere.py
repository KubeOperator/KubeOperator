import os
from threading import Thread
from time import sleep
from cloud_provider.cloud_client import CloudClient
from pyVim import connect
from pyVmomi import vim

from cloud_provider.utils import download_plugins
from fit2ansible.settings import CLOUDS_RESOURCE_DIR
from kubeops_api.models.setting import Setting


class VsphereCloudClient(CloudClient):
    cloud_config_path = os.path.join(CLOUDS_RESOURCE_DIR, 'vsphere')
    working_path = None

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

    def init_terraform(self):
        plugin_dir = os.path.join(self.working_path, '.terraform', 'plugins')
        if not os.path.exists(plugin_dir):
            os.makedirs(plugin_dir)
        hostname = Setting.objects.get(key='local_hostname').value
        port = 8082
        url = "http://{}:{}/repository/raw/terraform/vsphere.zip".format(hostname, port)
        download_plugins(url=url, target=plugin_dir)

    def apply_terraform(self, cluster):
        vars = cluster.plan.mixed_vars
        st = connect.SmartConnectNoSSL(host=vars['vc_host'], user=vars['vc_username'],
                                       pwd=vars['vc_password'], port=int(443))
        content = st.RetrieveContent()
        container = content.rootFolder
        dc = get_obj(content, [vim.Datacenter], container, vars['region'])
        folder = get_obj(content, [vim.Folder], container, 'kubeoperator')
        if not folder:
            dc.vmFolder.CreateFolder('kubeoperator')
        return super().apply_terraform(cluster)

    def create_image(self, zone):
        params = replace_params(self.vars)
        st = get_service_instance(params)
        content = st.RetrieveContent()
        container = content.rootFolder
        viewType = [vim.Folder]
        folder = get_obj(content, viewType, container, 'kubeoperator')
        viewType = [vim.VirtualMachine]
        vm = get_obj(content, viewType, folder, zone.region.image_name)
        ds = get_obj(content, [vim.Datastore], container, zone.vars['vc_storage'])
        dc = get_obj(content, [vim.Datacenter], container, zone.region.cloud_region)
        cluster = get_obj(content, [vim.ClusterComputeResource], container, zone.cloud_zone)
        if not vm:
            manager = st.content.ovfManager
            spec_params = vim.OvfManager.CreateImportSpecParams()
            ovf_path = zone.region.image_ovf_path
            vmdk_path = zone.region.image_vmdk_path
            ovfd = get_ovf_descriptor(ovf_path)
            resource_pool = cluster.resourcePool
            import_spec = manager.CreateImportSpec(ovfd,
                                                   resource_pool,
                                                   ds,
                                                   spec_params)
            lease = resource_pool.ImportVApp(import_spec.importSpec,
                                             dc.vmFolder)
            while True:
                if lease.state == vim.HttpNfcLease.State.ready:
                    url = lease.info.deviceUrl[0].url.replace('*', self.vars['vc_host'])
                    keepalive_thread = Thread(target=keep_lease_alive, args=(lease,))
                    keepalive_thread.start()
                    curl_cmd = (
                            "curl -Ss -X POST --insecure -T %s -H 'Content-Type: \
                            application/x-vnd.vmware-streamVmdk' %s" %
                            (vmdk_path, url))
                    os.system(curl_cmd)
                    lease.HttpNfcLeaseComplete()
                    keepalive_thread.join()
                    vm = get_obj(content, [vim.VirtualMachine], container, zone.region.image_name)
                    vm.MarkAsTemplate()
                    break
                elif lease.state == vim.HttpNfcLease.State.error:
                    print("Lease error: " + lease.state.error)
                    break


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


def get_ovf_descriptor(ovf_path):
    if os.path.exists(ovf_path):
        with open(ovf_path, 'r') as f:
            try:
                ovfd = f.read()
                f.close()
                return ovfd
            except:
                print("Could not read file: {}".format(ovf_path))


def keep_lease_alive(lease):
    while (True):
        sleep(5)
        try:
            print('模版上传中...')
            lease.HttpNfcLeaseProgress(50)
            if lease.state == vim.HttpNfcLease.State.done:
                return
        except:
            return


def replace_params(vars):
    return {
        "host": vars.get('vc_host', None),
        "username": vars.get('vc_username', None),
        "password": vars.get('vc_password', None),
    }
