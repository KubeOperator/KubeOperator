import kubernetes.client
import redis
import json
from kubernetes.client.rest import ApiException
from kubeops_api.cluster_data import ClusterData, Pod, NameSpace, Node, Container, Deployment
from kubeops_api.models.cluster import Cluster
from kubeops_api.prometheus_client import PrometheusClient
from django.db.models import Q


class ClusterMonitor():

    def __init__(self, cluster):
        # init redis
        self.redis_cli = redis.StrictRedis(host='localhost', port=6379)
        self.cluster = cluster
        self.retry_count = 0
        self.get_authorization()
        self.get_api_instance()

    def get_authorization(self):
        try:
            if self.redis_cli.exists(self.cluster.name):
                cluster_data = self.redis_cli.get(self.cluster.name)
                if cluster_data is not None:
                    cluster_str = str(cluster_data, encoding='utf-8')
                    cluster_d = json.loads(cluster_str)
                    self.token = cluster_d['token']
                else:
                    self.token = self.cluster.get_cluster_token()
            else:
                self.token = self.cluster.get_cluster_token()
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
        self.app_v1_api = kubernetes.client.AppsV1Api(kubernetes.client.ApiClient(configuration))
        self.check_authorization(self.retry_count)

    def check_authorization(self, retry_count):
        if retry_count > 2:
            raise Exception('init k8s client failed! retry_count=' + str(retry_count))
        self.retry_count = retry_count + 1
        try:
            self.api_instance.list_node()
        except ApiException as e:
            if e.status == 401:
                self.get_authorization()
                self.get_api_instance()
            else:
                raise Exception('init k8s client failed!' + e.reason)

    def list_pods(self):
        try:
            pods = self.api_instance.list_pod_for_all_namespaces()
            podList = []
            for p in pods.items:
                status = p.status
                containers = []
                restart_count = 0
                for c in status.container_statuses:
                    restart_count = restart_count + c.restart_count
                    container = Container(name=c.name, ready=c.ready, restart_count=c.restart_count,
                                          pod_name=p.metadata.name)
                    containers.append(container.__dict__)
                pod = Pod(name=p.metadata.name, cluster_name=self.cluster.name, restart_count=restart_count,
                          status=status.phase,
                          namespace=p.metadata.namespace,
                          host_ip=status.host_ip, pod_ip=status.pod_ip, host_name=None, containers=containers)
                podList.append(pod.__dict__)
            return podList
        except ApiException as e:
            raise Exception('list pod failed!' + e.reason)

    def list_namespaces(self):
        namespaces = self.api_instance.list_namespace()
        namespace_list = []
        for n in namespaces.items:
            namespace = NameSpace(name=n.metadata.name, status=n.status.phase)
            namespace_list.append(namespace.__dict__)
        return namespace_list

    def list_nodes(self):
        nodes = self.api_instance.list_node()
        node_list = []
        for n in nodes.items:
            node = Node(name=n.metadata.name, status=n.status.phase, cpu=0, mem=0, cpu_usage=0, mem_usage=0)
            node = self.get_node_data(node)
            node_list.append(node.__dict__)
        return node_list

    def list_deployments(self):
        deployments = self.app_v1_api.list_deployment_for_all_namespaces()
        deployment_list = []
        for d in deployments.items:
            deployment = Deployment(name=d.metadata.name, ready_replicas=d.status.ready_replicas,
                                    replicas=d.status.replicas, namespace=d.metadata.namespace)
            deployment_list.append(deployment.__dict__)
        return deployment_list

    def set_cluster_data(self):
        nodes = self.list_nodes()
        pods = self.list_pods()
        namespaces = self.list_namespaces()
        deployments = self.list_deployments()

        cpu_usage = 0
        cpu_total = 0
        mem_total = 0
        mem_usage = 0
        count = len(nodes)
        for n in nodes:
            # 不计算异常node数据
            cpu_total = cpu_total + float(n['cpu'])
            cpu_usage = cpu_usage + float(n['cpu_usage'])
            mem_total = mem_total + float(n['mem'])
            mem_usage = mem_usage + float(n['mem_usage'])
            if n['cpu_usage'] == 0 and n['mem_usage'] == 0:
                count = count -1
        cpu_usage = cpu_usage / count
        mem_usage = mem_usage / count

        cluster_data = ClusterData(cluster=self.cluster, token=self.token, pods=pods, nodes=nodes,
                                   namespaces=namespaces, deployments=deployments, cpu_usage=cpu_usage,
                                   cpu_total=cpu_total,
                                   mem_total=mem_total, mem_usage=mem_usage)
        return self.redis_cli.set(self.cluster.name, json.dumps(cluster_data.__dict__))

    def list_cluster_data(self):
        cluster_data = self.redis_cli.get(self.cluster.name)
        result = {}
        if cluster_data is not None:
            cluster_str = str(cluster_data, encoding='utf-8')
            result = json.loads(cluster_str)
        return result

    def get_node_data(self, node):
        host = "prometheus.apps." + self.cluster.name + "." + self.cluster.cluster_doamin_suffix
        config = {
            'host': host
        }
        prometheus_client = PrometheusClient(config)
        return prometheus_client.get_node_resource(node)
