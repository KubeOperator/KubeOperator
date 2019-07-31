from abc import ABCMeta, abstractmethod


def get_cloud_client(vars):
    provider = vars.get('provider', {})
    from cloud_provider.clients.vsphere import VsphereCloudClient
    if provider == 'vsphere':
        return VsphereCloudClient(vars)
    else:
        return None


class CloudClient(metaclass=ABCMeta):

    def __init__(self, vars):
        self.vars = vars

    @abstractmethod
    def list_region(self):
        pass
