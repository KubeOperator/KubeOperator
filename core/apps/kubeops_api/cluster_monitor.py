import kubernetes.client
import redis
import json
import logging
import kubeoperator.settings
import log.es
import datetime, time
import builtins

from kubernetes.client.rest import ApiException
from kubeops_api.cluster_data import ClusterData, Pod, NameSpace, Node, Container, Deployment, StorageClass, PVC, Event
from kubeops_api.models.cluster import Cluster
from kubeops_api.prometheus_client import PrometheusClient
from kubeops_api.models.host import Host
from django.db.models import Q
from kubeops_api.cluster_health_data import ClusterHealthData
from django.utils import timezone
from ansible_api.models.inventory import Host as C_Host
from common.ssh import SSHClient, SshConfig
from message_center.message_client import MessageClient
from kubeops_api.utils.date_encoder import DateEncoder

logger = logging.getLogger('kubeops')


class ClusterMonitor():

    def __init__(self, cluster):
        self.cluster = cluster
        self.retry_count = 0
        self.restart_pods = []
        self.warn_containers = []
        self.error_pods = []
        # init redis
        self.redis_cli = redis.StrictRedis(host=kubeoperator.settings.REDIS_HOST, port=kubeoperator.settings.REDIS_PORT)
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
        self.get_api_instance()

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
        self.check_authorization(3)
        nodes = self.list_nodes()
        pods = self.list_pods()
        namespaces = self.list_namespaces()
        deployments = self.list_deployments()

        cpu_usage = 0
        cpu_total = 0
        mem_total = 0
        mem_usage = 0
        count = len(nodes)
        warn_nodes = []
        for n in nodes:
            # 不计算异常node数据
            cpu_total = cpu_total + float(n['cpu'])
            cpu_usage = cpu_usage + float(n['cpu_usage'])
            mem_total = mem_total + float(n['mem'])
            mem_usage = mem_usage + float(n['mem_usage'])
            if float(n['cpu_usage']) == 0 and float(n['mem_usage']) == 0:
                count = count - 1
            elif float(n['cpu_usage']) > 0.8 or float(n['mem_usage']) > 0.8:
                warn_nodes.append(n)
        if count > 0:
            cpu_usage = cpu_usage / count
            mem_usage = mem_usage / count
        if len(warn_nodes) > 0:
            message_client = MessageClient()
            message = self.get_usage_message(warn_nodes)
            message_client.insert_message(message)

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
        try:
            res = prometheus_client.get_node_resource(node)
            return res
        except Exception as e:
            logger.error(msg='get node data error ', exc_info=True)
            return node

    def get_loki_msg(self):
        host = "loki.apps." + self.cluster.name + "." + self.cluster.cluster_doamin_suffix
        config = {
            'host': host,
            'cluster': self.cluster
        }
        prometheus_client = PrometheusClient(config)
        try:
            res = prometheus_client.get_msg_from_loki(self.cluster.name)
            return res
        except Exception as e:
            logger.error(msg='get loki meg error ', exc_info=True)
            return []

    def set_loki_data_to_cluster(self):
        cluster_data = self.redis_cli.get(self.cluster.name)
        if cluster_data is not None:
            cluster_str = str(cluster_data, encoding='utf-8')
            cluster_d = json.loads(cluster_str)
            cluster_d['error_loki_containers'] = quick_sort_error_loki_container(self.get_loki_msg())
            return self.redis_cli.set(self.cluster.name, json.dumps(cluster_d))
        else:
            return False

    def list_pod_status(self, namespace):
        message = ''
        pod_data = []
        try:
            if namespace == 'all':
                namespaces = self.api_instance.list_namespace()
                for ns in namespaces.items:
                    ns_name = ns.metadata.name
                    pods = self.api_instance.list_namespaced_pod(ns_name)
                    pod_d = self.get_pod_status(pods.items)
                    if len(pod_d) > 0:
                        pod_data = pod_data + pod_d
            else:
                pods = self.api_instance.list_namespaced_pod(namespace)
                pod_data = self.get_pod_status(pods.items)
        except ApiException as e:
            message = e.reason
            logger.error(msg='list pod error ' + e.reason, exc_info=True)
        health_data = {
            'pod_data': pod_data,
            'message': message
        }
        return health_data

    def list_namespace(self):
        namespaces = self.api_instance.list_namespace()
        ns_names = []
        for ns in namespaces.items:
            ns_name = ns.metadata.name
            ns_names.append(ns_name)
        return ns_names

    def get_component_status(self):
        component_data = []
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
        except ApiException as e:
            logger.error(msg='list component error ' + e.reason, exc_info=True)

        return component_data

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
            if item.parameters:
                datastore = item.parameters.get('datastore', None)
            else:
                datastore = ''
            storage_class = StorageClass(name=item.metadata.name, provisioner=item.provisioner, datastore=datastore,
                                         create_time=str(item.metadata.creation_timestamp), pvcs=[])
            scs.append(storage_class.__dict__)
        pvc_response = self.api_instance.list_persistent_volume_claim_for_all_namespaces()
        for item in pvc_response.items:
            if item.status.capacity:
                capacity = item.status.capacity.get('storage', None)
            else:
                capacity = ''
            pvc = PVC(name=item.metadata.name, namespace=item.metadata.namespace, status=item.status.phase,
                      capacity=capacity, storage_class=item.spec.storage_class_name, mount_by=item.metadata.namespace,
                      create_time=str(item.metadata.creation_timestamp))
            for sc in scs:
                if sc['name'] == item.spec.storage_class_name:
                    sc['pvcs'].append(pvc.__dict__)
        return scs

    def list_events(self):
        event_response = self.api_instance.list_event_for_all_namespaces()
        events = []
        actions = []
        year = datetime.datetime.now().year
        month = datetime.datetime.now().month
        index = (self.cluster.name + '-{}.{}').format(year, month)
        es_client = log.es.get_es_client()
        for item in event_response.items:
            # 过滤kubeapps-plus的同步事件
            if "apprepo-sync-chartmuseum" in item.metadata.name:
                continue
            component, host = '', ''
            if item.source is not None and item.source.component is not None:
                component = item.source.component
            elif item.reporting_component is not None:
                component = item.reporting_component
            if item.source is not None and item.source.host is not None:
                host = item.source.host
            elif item.reporting_instance is not None:
                host = item.reporting_instance

            if item.last_timestamp is not None:
                last_timestamp = item.last_timestamp
            elif item.event_time is not None:
                last_timestamp = item.event_time
            else:
                last_timestamp = item.metadata.creation_timestamp

            event = Event(uid=item.metadata.uid, name=item.metadata.name, type=item.type,
                          cluster_name=self.cluster.name, action=item.action,
                          reason=item.reason, count=item.count, host=host, component=component,
                          namespace=item.metadata.namespace,
                          message=item.message, last_timestamp=last_timestamp, first_timestamp=item.first_timestamp)
            events.append(event.__dict__)
            # 判断根据uid判断这个事件是否已经存入es
            if log.es.get_event_uid_exist(es_client, index, item.metadata.uid):
                action = {
                    '_op_type': 'index',
                    '_index': index,
                    '_type': 'event',
                    '_source': event.__dict__
                }
                actions.append(action)
                if event.type == 'Warning':
                    message_client = MessageClient()
                    message = self.get_event_message(event)
                    message_client.insert_message(message)
        return events, actions

    def get_event_message(self, event):
        message = {
            "item_id": self.cluster.item_id,
            "title": "集群事件告警",
            "content": self.get_event_content(event),
            "level": "WARNING",
            "type": "CLUSTER"
        }
        return message

    def get_event_content(self, event):
        content = {
            "item_name": self.cluster.item_name,
            "resource": "集群",
            "resource_name": self.cluster.name,
            "resource_type": 'CLUSTER_EVENT',
            "detail": json.dumps(event.__dict__, cls=DateEncoder),
            "status": self.cluster.status,
        }
        return content

    def get_usage_message(self, nodes):
        message = {
            "item_id": self.cluster.item_id,
            "title": "集群资源告警",
            "content": self.get_usage_content(nodes),
            "level": "WARNING",
            "type": "CLUSTER"
        }
        return message

    def get_usage_content(self, nodes):
        message = ''
        for n in nodes:
            cpu_usage = round(float(n['cpu_usage']),2) * 100
            mem_usage = round(float(n['mem_usage']),2) * 100
            m = '主机{0}的CPU使用率为:{1}%,内存使用率为{2}% \n\n'.format(n['name'], cpu_usage, mem_usage)
            if len(message)>0:
                message = message +'> '+m
            else:
                message = m

        content = {
            "item_name": self.cluster.item_name,
            "resource": "集群",
            "resource_name": self.cluster.name,
            "resource_type": 'CLUSTER_USAGE',
            "detail":json.dumps({'message': message}) ,
            "status": self.cluster.status,
        }
        return content


def delete_cluster_redis_data(cluster_name):
    redis_cli = redis.StrictRedis(host=kubeoperator.settings.REDIS_HOST, port=kubeoperator.settings.REDIS_PORT)
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
                                      ~Q(status=Cluster.CLUSTER_STATUS_DELETING),
                                      ~Q(status=Cluster.CLUSTER_STATUS_ERROR))
    for cluster in clusters:
        cluster_monitor = ClusterMonitor(cluster)
        success = cluster_monitor.set_cluster_data()
        if success == False:
            logger.error(msg='put cluster data to redis error', exc_info=True)


def put_loki_data_to_redis():
    clusters = Cluster.objects.filter(~Q(status=Cluster.CLUSTER_STATUS_READY),
                                      ~Q(status=Cluster.CLUSTER_STATUS_INSTALLING),
                                      ~Q(status=Cluster.CLUSTER_STATUS_DELETING),
                                      ~Q(status=Cluster.CLUSTER_STATUS_ERROR))
    for cluster in clusters:
        cluster_monitor = ClusterMonitor(cluster)
        success = cluster_monitor.set_loki_data_to_cluster()
        if success == False:
            logger.error(msg='put cluster loki data to redis error', exc_info=True)


def put_event_data_to_es():
    clusters = Cluster.objects.filter(~Q(status=Cluster.CLUSTER_STATUS_READY),
                                      ~Q(status=Cluster.CLUSTER_STATUS_INSTALLING),
                                      ~Q(status=Cluster.CLUSTER_STATUS_DELETING),
                                      ~Q(status=Cluster.CLUSTER_STATUS_ERROR))

    for cluster in clusters:
        cluster_monitor = ClusterMonitor(cluster)
        year = datetime.datetime.now().year
        month = datetime.datetime.now().month
        index = (cluster.name + '-{}.{}').format(year, month)
        es_client = log.es.get_es_client()
        index_exists = log.es.exists(es_client, index)
        if index_exists == False:
            index_exists = create_index(es_client, index)
        if index_exists:
            try:
                events, actions = cluster_monitor.list_events()
                if len(actions) > 0:
                    success, failed = log.es.batch_data(es_client, actions)
                    logger.info(
                        msg='put' + cluster.name + 'event to es success:' + str(success) + 'failed:' + str(failed),
                        exc_info=False)
            except ApiException as e:
                logger.error(msg='list event error' + e.reason, exc_info=True)
        else:
            logger.error(msg='create es index error', exc_info=True)


def create_index(client, index):
    index_mapping = {
        "properties": {
            "uid": {
                "type": "keyword"
            },
            "name": {
                "type": "text"
            },
            "type": {
                "type": "keyword"
            },
            "cluster_name": {
                "type": "text"
            },
            "reason": {
                "type": "text"
            },
            "action": {
                "type": "text"
            },
            "count": {
                "type": "integer"
            },
            "component": {
                "type": "text"
            },
            "namespace": {
                "type": "text"
            },
            "message": {
                "type": "text",
                "analyzer": "english"
            },
            "host": {
                "type": "text"
            },
            "last_timestamp": {
                "type": "date"
            },
            "first_timestamp": {
                "type": "date"
            }
        }
    }
    return log.es.create_index_and_mapping(client, index, 'event', index_mapping)


def delete_unused_node(cluster):
    cluster_monitor = ClusterMonitor(cluster)
    nodes = cluster_monitor.list_nodes()
    hosts = C_Host.objects.filter(
        Q(project_id=cluster.id) & ~Q(name='localhost') & ~Q(name='127.0.0.1') & ~Q(name='::1'))
    if len(nodes) > 0:
        for host in hosts:
            exist = False
            delete_name = host.name
            for node in nodes:
                if delete_name == node['name']:
                    exist = True
            if exist is False and delete_name != '':
                C_Host.objects.filter(name=delete_name).delete()
    return True


def sync_node_time(cluster):
    hosts = C_Host.objects.filter(
        Q(project_id=cluster.id) & ~Q(name='localhost') & ~Q(name='127.0.0.1') & ~Q(name='::1'))
    data = []
    times = []
    result = {
        'success': True,
        'data': []
    }
    for host in hosts:
        ssh_config = SshConfig(host=host.ip, port=host.port, username=host.username, password=host.password,
                               private_key=None)

        ssh_client = SSHClient(ssh_config)
        res = ssh_client.run_cmd('date')
        gmt_date = res[0]
        GMT_FORMAT = '%a %b %d %H:%M:%S CST %Y'
        date = time.strptime(gmt_date, GMT_FORMAT)
        timeStamp = int(time.mktime(date))
        times.append(timeStamp)
        show_time = time.strftime('%Y-%m-%d %H:%M:%S', date)
        time_data = {
            'hostname': host.name,
            'date': show_time,
        }
        data.append(time_data)
    result['data'] = data
    max = builtins.max(times)
    min = builtins.min(times)
    # 如果最大值减最小值超过5分钟 则判断有错
    if (max - min) > 300000:
        result['success'] = False
    return result
