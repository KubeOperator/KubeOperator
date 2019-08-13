import os
from abc import ABCMeta, abstractmethod

from python_terraform import Terraform, IsNotFlagged

from cloud_provider.utils import generate_terraform_file, create_terrafrom_working_dir
from fit2ansible.settings import  CLOUDS_RESOURCE_DIR


def get_cloud_client(vars):
    provider = vars.get('provider', {})
    from cloud_provider.clients.vsphere import VsphereCloudClient
    if provider == 'vsphere':
        return VsphereCloudClient(vars)
    else:
        return None


class CloudClient(metaclass=ABCMeta):
    cloud_config_path = CLOUDS_RESOURCE_DIR
    working_path = None

    def __init__(self, vars):
        self.vars = vars

    @abstractmethod
    def list_region(self):
        pass

    @abstractmethod
    def init_terraform(self):
        pass

    def apply_terraform(self, cluster, vars):
        if self.working_path:
            self.working_path = create_terrafrom_working_dir(cluster_name=cluster.name)
        generate_terraform_file(self.working_path, self.cloud_config_path, vars)
        self.init_terraform()
        t = Terraform(working_dir=self.working_path)
        p = t.apply('./', refresh=True, skip_plan=True, no_color=IsNotFlagged, synchronous=False)
        p = p[0]
        for i in p.stdout:
            print(i.decode())
        p.communicate()
        return p.returncode == 0
