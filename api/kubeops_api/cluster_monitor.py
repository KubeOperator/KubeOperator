import kubernetes.client
from kubernetes.client.rest import ApiException
from pprint import pprint
from kubeops_api.models.cluster import Cluster

class ClusterMonitor():

    def __init__(self,cluster):
        self.cluster = cluster
        self.token = self.cluster.get_cluster_token()
        self.cluster.change_to()
        master = self.cluster.group_set.get(name='master').hosts.first()
        configuration = kubernetes.client.Configuration()
        configuration.api_key_prefix['authorization'] = 'Bearer'
        configuration.api_key['authorization'] = self.token
        print('---token----')
        print(self.token)
        configuration.debug = True
        configuration.host = 'https://'+master.ip+":6443"
        configuration.verify_ssl = False
        print('https://'+master.ip+":6443")
        self.api_instance = kubernetes.client.CoreV1Api(kubernetes.client.ApiClient(configuration))

    def list_pods(self):
        pods = self.api_instance.list_pod_for_all_namespaces()
        return pods
