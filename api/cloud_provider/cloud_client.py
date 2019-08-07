import os
from abc import ABCMeta, abstractmethod

from python_terraform import Terraform, IsNotFlagged

from cloud_provider.utils import init_terraform, generate_terraform_file
from fit2ansible.settings import RESOURCE_DIR, ANSIBLE_PROJECTS_DIR


def get_cloud_client(vars):
    provider = vars.get('provider', {})
    from cloud_provider.clients.vsphere import VsphereCloudClient
    if provider == 'vsphere':
        return VsphereCloudClient(vars)
    else:
        return None


class CloudClient(metaclass=ABCMeta):
    cloud_config_path = RESOURCE_DIR

    def __init__(self, vars):
        self.vars = vars

    @abstractmethod
    def list_region(self):
        pass

    def apply_terraform(self, cluster, vars):
        target_path = os.path.join(ANSIBLE_PROJECTS_DIR, cluster, 'terraform')
        target = generate_terraform_file(target_path, self.cloud_config_path, vars)
        init_terraform(target_path, self.cloud_config_path)
        t = Terraform(working_dir=target)
        p = t.apply('./', refresh=True, skip_plan=True, no_color=IsNotFlagged, synchronous=False)
        p = p[0]
        for i in p.stdout:
            print(i.decode())
        p.communicate()
        code = p.returncode
        if not code == 0:
            for i in p.stderr:
                print(i.decode())
        return p.returncode == 0
