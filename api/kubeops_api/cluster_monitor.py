import kubernetes.client
from kubernetes.client.rest import ApiException
import redis
from pprint import pprint
from kubeops_api.models.cluster import Cluster


class ClusterMonitor():

    def __init__(self, cluster):
        # init redis
        pool = redis.ConnectionPool(host='localhost', port=6379)
        self.redis_cli = redis.Redis(connection_pool=pool)
        self.cluster = cluster
        self.retry_count = 0
        self.get_authorization()
        self.get_api_instance()

    def list_pods(self):
        try:
            pods = self.api_instance.list_pod_for_all_namespaces()
            return pods
        except ApiException as e:
            raise Exception('list pod failed!' + e.reason)

    def get_authorization(self):
        try:
            if self.redis_cli.exists(self.cluster.name):
                self.token = str(self.redis_cli.get(self.cluster.name),encoding= 'utf-8')
            else:
                self.token = self.cluster.get_cluster_token()
                self.redis_cli.set(self.cluster.name,self.token)
        except ApiException as e:
            if e.status == 401:
                self.token = self.cluster.get_cluster_token()
            else:
                raise Exception('get authorization failed!' + e.reason)

    def get_api_instance(self):
        self.cluster.change_to()
        master = self.cluster.group_set.get(name='master').hosts.first()
        configuration = kubernetes.client.Configuration()
        configuration.api_key_prefix['authorization'] = 'Bearer'
        configuration.api_key['authorization'] = self.token
        configuration.debug = True
        configuration.host = 'https://' + master.ip + ":6443"
        configuration.verify_ssl = False
        self.api_instance = kubernetes.client.CoreV1Api(kubernetes.client.ApiClient(configuration))
        self.check_authorization(self.retry_count)

    def check_authorization(self,retry_count):
        if retry_count > 2:
            raise Exception('init k8s client failed! retry_count=' + retry_count)
        self.retry_count = retry_count + 1
        try:
            self.api_instance.list_node()
        except ApiException as e:
            if e.status == 401:
                self.get_authorization()
                self.get_api_instance()
            else:
                raise Exception('init k8s client failed!' + e.reason)


