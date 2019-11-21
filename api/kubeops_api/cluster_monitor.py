import kubernetes.client
import redis
import json
import logging
import fit2ansible.settings
from kubernetes.client.rest import ApiException
from kubeops_api.cluster_data import ClusterData, Pod, NameSpace, Node, Container, Deployment, StorageClass, PVC
from kubeops_api.models.cluster import Cluster
from kubeops_api.prometheus_client import PrometheusClient
from kubeops_api.models.host import Host
from django.db.models import Q
from kubeops_api.cluster_health_data import ClusterHealthData
from django.utils import timezone

logger = logging.getLogger('kubeops')


class ClusterMonitor():

    def __init__(self, cluster):
        self.cluster = cluster
        self.retry_count = 0
        self.restart_pods = []
        self.warn_containers = []
        self.error_pods = []
        # init redis
        self.redis_cli = redis.StrictRedis(host=fit2ansible.settings.REDIS_HOST, port=fit2ansible.settings.REDIS_PORT)
        self.get_token()
        self.get_api_instance()

    def get_token(self):
        if self.redis_cli.exists(self.cluster.name):
            cluster_data = self.redis_cli.get(self.cluster.name)
            if cluster_data is not None:
                cluster_str = str(cluster_data, encoding='utf-8')
                cluster_d = json.loads(cluster_str)
                self.token = cluster_d['token']
            else:
                self.push_token_to_redis()
        else:
            self.push_token_to_redis()

    def push_token_to_redis(self):
        self.token = self.cluster.get_cluster_token()
        cluster_data = ClusterData(cluster=self.cluster, token=self.token, pods=[], nodes=[],
                                   namespaces=[], deployments=[], cpu_usage=0, cpu_total=0, mem_total=0,
                                   mem_usage=0, restart_pods=[], warn_containers=[], error_loki_containers=[],
                                   error_pods=[])
        self.redis_cli.set(self.cluster.name, json.dumps(cluster_data.__dict__))

    def get_api_instance(self):
        self.cluster.change_to()
        master = self.cluster.group_set.get(name='master').hosts.first()
        if master is not None and master.ip is not None:
            configuration = kubernetes.client.Configuration()
            configuration.api_key_prefix['authorization'] = 'Bearer'
            configuration.api_key['authorization'] = self.token
            configuration.debug = True
            configuration.host = 'https://' + master.ip + ":6443"
            configuration.verify_ssl = False
            self.api_instance = kubernetes.client.CoreV1Api(kubernetes.client.ApiClient(configuration))
            self.app_v1_api = kubernetes.client.AppsV1Api(kubernetes.client.ApiClient(configuration))
            self.storage_v1_Api = kubernetes.client.StorageV1Api(kubernetes.client.ApiClient(configuration))

    def check_authorization(self, retry_count):
        if retry_count > 2:
            raise Exception('init k8s client failed! retry_count=' + str(retry_count))
        self.retry_count = retry_count + 1
        try:
            self.api_instance.list_node()
        except ApiException as e:
            if e.status == 401:
                self.push_token_to_redis()
                self.get_api_instance()
            else:
                logger.error(msg='init k8s client failed ' + e.reason, exc_info=True)

    def list_pods(self):
        podList = []
        try:
            pods = self.api_instance.list_pod_for_all_namespaces()
            for p in pods.items:
                status = p.status
                containers = []
                restart_count = 0
                hostname = None
                host_ip = None
                if status.container_statuses is not None:
                    for c in status.container_statuses:
                        restart_count = restart_count + c.restart_count
                        container = Container(name=c.name, ready=c.ready, restart_count=c.restart_count,
                                              pod_name=p.metadata.name)
                        if container.ready == False:
                            self.warn_containers.append(container.__dict__)
                        containers.append(container.__dict__)
                if status.host_ip is not None:
                    host = Host.objects.get(ip=status.host_ip)
                    hostname = (host.name if host.name is not None else None)
                    host_ip = status.host_ip
                pod_ip = (status.pod_ip if status.pod_ip is not None else None)
                pod = Pod(name=p.metadata.name, cluster_name=self.cluster.name, restart_count=restart_count,
                          status=status.phase,
                          namespace=p.metadata.namespace,
                          host_ip=host_ip, pod_ip=pod_ip, host_name=hostname, containers=containers)
                if restart_count > 0:
                    self.restart_pods.append(pod.__dict__)
                if status.phase != 'Running' and status.phase != 'Succeeded':
                    self.error_pods.append(pod.__dict__)
                podList.append(pod.__dict__)
        except ApiException as e:
            logger.error(msg='list pod error ' + e.reason, exc_info=True)
        return podList

    def list_namespaces(self):
        namespace_list = []
        try:
            namespaces = self.api_instance.list_namespace()
            for n in namespaces.items:
                namespace = NameSpace(name=n.metadata.name, status=n.status.phase)
                namespace_list.append(namespace.__dict__)
        except ApiException as e:
            logger.error(msg='list namespace error ' + e.reason, exc_info=True)
        return namespace_list

    def list_nodes(self):
        node_list = []
        try:
            nodes = self.api_instance.list_node()
            for n in nodes.items:
                node = Node(name=n.metadata.name, status=n.status.phase, cpu=0, mem=0, cpu_usage=0, mem_usage=0)
                node = self.get_node_data(node)
                node_list.append(node.__dict__)
        except ApiException as e:
            logger.error(msg='list node error ' + e.reason, exc_info=True)
        return node_list

    def list_deployments(self):
        deployment_list = []
        try:
            deployments = self.app_v1_api.list_deployment_for_all_namespaces()
            for d in deployments.items:
                deployment = Deployment(name=d.metadata.name, ready_replicas=d.status.ready_replicas,
                                        replicas=d.status.replicas, namespace=d.metadata.namespace)
                deployment_list.append(deployment.__dict__)
        except ApiException as e:
            logger.error(msg='list namespace error ' + e.reason, exc_info=True)
        return deployment_list

    def set_cluster_data(self):
        self.check_authorization(self.retry_count)
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
                count = count - 1
        if count > 0:
            cpu_usage = cpu_usage / count
            mem_usage = mem_usage / count
        sort_restart_pod_list = quick_sort_pods(self.restart_pods)
        error_pods = quick_sort_pods(self.error_pods)

        cluster_data = ClusterData(cluster=self.cluster, token=self.token, pods=pods, nodes=nodes,
                                   namespaces=namespaces, deployments=deployments, cpu_usage=cpu_usage,
                                   cpu_total=cpu_total,
                                   mem_total=mem_total, mem_usage=mem_usage, restart_pods=sort_restart_pod_list,
                                   warn_containers=self.warn_containers, error_loki_containers=[],
                                   error_pods=error_pods)
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
            'host': host,
            'cluster': self.cluster
        }
        prometheus_client = PrometheusClient(config)
        return prometheus_client.get_node_resource(node)

    def get_loki_msg(self):
        host = "loki.apps." + self.cluster.name + "." + self.cluster.cluster_doamin_suffix
        config = {
            'host': host,
            'cluster': self.cluster
        }
        prometheus_client = PrometheusClient(config)
        return prometheus_client.get_msg_from_loki(self.cluster.name)

    def set_loki_data_to_cluster(self):
        cluster_data = self.redis_cli.get(self.cluster.name)
        if cluster_data is not None:
            cluster_str = str(cluster_data, encoding='utf-8')
            cluster_d = json.loads(cluster_str)
            cluster_d['error_loki_containers'] = quick_sort_error_loki_container(self.get_loki_msg())
            return self.redis_cli.set(self.cluster.name, json.dumps(cluster_d))
        else:
            return False

    def get_kubernetes_status(self):
        message = ''
        component_data, monitor_data, system_data = [], [], []
        try:
            components = self.api_instance.list_component_status()
            for c in components.items:
                status, msg = '', ''
                for condition in c.conditions:
                    if condition.type == 'Healthy':
                        msg = condition.message
                        if condition.status == 'True':
                            status = 'RUNNING'
                        elif condition.status == 'False':
                            status = 'ERROR'
                        else:
                            status = condition.status
                component = ClusterHealthData(namespace='component', name=c.metadata.name, status=status,
                                              ready='1/1', age=0, msg=msg, restart_count=0)
                component_data.append(component.__dict__)
            system_pods = self.api_instance.list_namespaced_pod('kube-system')
            system_data = self.get_pod_status(system_pods.items)
            monitor_pods = self.api_instance.list_namespaced_pod('kube-operator')
            monitor_data = self.get_pod_status(monitor_pods.items)
        except ApiException as e:
            message = e.reason
            logger.error(msg='list pod error ' + e.reason, exc_info=True)
        health_data = {
            'component': component_data,
            'kube-system': system_data,
            'monitoring': monitor_data,
            'message': message
        }
        return health_data

    def get_pod_status(self, items):
        pod_data = []
        for s in items:
            restart_count = 0
            if s.status.container_statuses is not None:
                count = len(s.status.container_statuses)
                ready = 0
                for c in s.status.container_statuses:
                    restart_count = restart_count + c.restart_count
                    if c.ready:
                        ready = ready + 1
                ready_status = str(ready) + '/' + str(count)
                # 计算存活时间
                now = timezone.now()
                age_time = now - s.status.start_time
                age = ''
                if age_time.days > 0:
                    age = str(age_time.days) + 'd'
                else:
                    seconds = age_time.seconds
                    hour = int(seconds / 60 / 60)
                    if hour >= 1:
                        age = str(hour) + 'h'
                        minute = int((seconds % 3600) / 60)
                    else:
                        minute = int(seconds / 60)
                    if minute >= 1:
                        age = age + str(minute) + 'm'
                system_pod = ClusterHealthData(namespace=s.metadata.namespace, name=s.metadata.name,
                                               status=s.status.phase,
                                               ready=ready_status, age=age, msg=s.status.message,
                                               restart_count=restart_count)
                pod_data.append(system_pod.__dict__)
        return pod_data

    def list_storage_class(self):
        sc_response = self.storage_v1_Api.list_storage_class()
        scs = []
        for item in sc_response.items:
            datastore = item.parameters.get('datastore', None)
            storage_class = StorageClass(name=item.metadata.name, provisioner=item.provisioner, datastore=datastore,
                                         create_time=str(item.metadata.creation_timestamp), pvcs=[])
            scs.append(storage_class.__dict__)
        pvc_response = self.api_instance.list_persistent_volume_claim_for_all_namespaces()
        for item in pvc_response.items:
            capacity = item.status.capacity.get('storage', None)
            pvc = PVC(name=item.metadata.name, namespace=item.metadata.namespace, status=item.status.phase,
                      capacity=capacity, storage_class=item.spec.storage_class_name, mount_by=item.metadata.namespace,
                      create_time=str(item.metadata.creation_timestamp))
            for sc in scs:
                if sc['name'] == item.spec.storage_class_name:
                    sc['pvcs'].append(pvc.__dict__)
        return scs


def delete_cluster_redis_data(cluster_name):
    redis_cli = redis.StrictRedis(host=fit2ansible.settings.REDIS_HOST, port=fit2ansible.settings.REDIS_PORT)
    return redis_cli.delete(cluster_name)


def quick_sort_pods(pod_list):
    if len(pod_list) < 2:
        return pod_list
    mid = pod_list[0]

    left, right = [], []
    pod_list.remove(mid)

    for item in pod_list:
        if item['restart_count'] <= mid['restart_count']:
            right.append(item)
        else:
            left.append(item)
    return quick_sort_pods(left) + [mid] + quick_sort_pods(right)


def quick_sort_error_loki_container(containers):
    if len(containers) < 2:
        return containers
    mid = containers[0]

    left, right = [], []
    containers.remove(mid)

    for item in containers:
        if item.get('error_count') <= mid.get('error_count'):
            right.append(item)
        else:
            left.append(item)
    return quick_sort_error_loki_container(left) + [mid] + quick_sort_error_loki_container(right)


def put_cluster_data_to_redis():
    clusters = Cluster.objects.filter(~Q(status=Cluster.CLUSTER_STATUS_READY),
                                      ~Q(status=Cluster.CLUSTER_STATUS_INSTALLING),
                                      ~Q(status=Cluster.CLUSTER_STATUS_DELETING))
    for cluster in clusters:
        cluster_monitor = ClusterMonitor(cluster)
        success = cluster_monitor.set_cluster_data()
        if success == False:
            logger.error(msg='put cluster data to redis error', exc_info=True)


def put_loki_data_to_redis():
    clusters = Cluster.objects.filter(~Q(status=Cluster.CLUSTER_STATUS_READY),
                                      ~Q(status=Cluster.CLUSTER_STATUS_INSTALLING),
                                      ~Q(status=Cluster.CLUSTER_STATUS_DELETING))
    for cluster in clusters:
        cluster_monitor = ClusterMonitor(cluster)
        success = cluster_monitor.set_loki_data_to_cluster()
        if success == False:
            logger.error(msg='put cluster loki data to redis error', exc_info=True)
