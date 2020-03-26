import kubernetes
from django.core.cache import cache

from common.ssh import SSHClient
from kubeops_api.utils.health import parse_host_to_ssh_config


class KubernetesApi:
    def __init__(self, cluster):
        self.cluster = cluster

    def get_api_client(self):
        self.cluster.change_to()
        master = self.cluster.group_set.get(name='master').hosts.first()
        if master is not None and master.ip is not None:
            configuration = kubernetes.client.Configuration()
            configuration.api_key_prefix['authorization'] = 'Bearer'
            configuration.api_key['authorization'] = self.get_auth_token()
            configuration.debug = True
            configuration.host = 'https://' + master.ip + ":6443"
            configuration.verify_ssl = False
            return kubernetes.client.ApiClient(configuration)

    def get_auth_token(self):
        token = cache.get(self.cluster.name + "token", None)
        if not token:
            token = self.gather_auth_token()
        return token

    def gather_auth_token(self):
        cmd = "kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep tiller | awk '{print $1}') | grep token: | awk '{print $2}'"
        master = self.cluster.group_set.get(name='master').hosts.first()
        ssh_config = parse_host_to_ssh_config(master)
        ssh_client = SSHClient(ssh_config)
        connected = ssh_client.ping()
        if not connected:
            raise Exception("ssh connect error!")
        else:
            out, code = ssh_client.run_cmd(cmd)
        if code == 0:
            self.cache_auth_token(out)
            return out
        else:
            raise Exception('exec cmd error:' + out)

    def cache_auth_token(self, token):
        cache.set(self.cluster.name + "token", token)
