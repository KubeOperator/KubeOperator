import requests
import time

from kubeops_api.apps_client import AppsClient
from kubeops_api.models.host import Host
from kubeops_api.cluster_data import LokiContainer


class PrometheusClient():

    def __init__(self, config, cluster):
        self.host = config.get("host", None)
        self.table_name = config.get("table_name", None)
        self.param = config.get("param", None)
        self.start = config.get("start", None)
        self.end = config.get("end", None)
        self.cluster = cluster

    def query(self):
        url = "http://{host}/api/v1/query?query={table_name}{param}&start={start}&end={end}"
        query_url = url.format(host=self.host, table_name=self.table_name, param=self.param, start=self.start,
                               end=self.end)
        app_client = AppsClient(cluster=self.cluster)
        req = app_client.get('prometheus', url)
        return req.json()

    def targets(self):
        url = "http://{host}/api/v1/targets"
        query_url = url.format(host=self.host)
        app_client = AppsClient(cluster=self.cluster)
        req = app_client.get('prometheus', query_url)
        return req.json()

    def handle_targets_message(self, json):
        result = {
            'success': True,
            'data': [],
            'rate': 0
        }
        if json.get('status') and json.get('status') == 'success':
            keys = ['kubernetes-control-manager', 'etcd', 'kubernetes-nodes', 'kubernetes-schedule',
                    'kubernetes-apiservers']
            for key in keys:
                result['data'].append({
                    'job': key,
                    'data': [],
                    'rate': 0
                })
            active_targets = json.get('data').get('activeTargets')
            for target in active_targets:
                key = target.get('labels').get('job')
                if key in keys:
                    index = keys.index(target.get('labels').get('job'))
                    instance_address = target.get('discoveredLabels').get('__address__').split(':')[0]
                    hostName = Host.objects.get(ip=instance_address).name
                    health = target.get('health')
                    status = 'NotReady'
                    if health == 'up':
                        status = 'Ready'
                    if key != 'kubernetes-nodes':
                        result['data'][index]['data'].append({
                            'key': hostName,
                            'value': status
                        })
                    if key == 'kubernetes-nodes' and 'master' not in hostName:
                        result['data'][index]['data'].append({
                            'key': hostName,
                            'value': status
                        })
        else:
            result['success'] = False
        self.calculate_available_rate(result)
        return result

    def calculate_available_rate(self, result):
        service_up = 0
        for res in result['data']:
            job_up = 0
            for job in res['data']:
                if job['value'] == 'Ready':
                    job_up = job_up + 1
            res['rate'] = job_up / len(res['data']) * 100 if len(res['data']) > 0 else 0
            if res['rate'] == 100:
                service_up = service_up + 1
        result['rate'] = service_up / len(result['data']) * 100 if len(result['data']) > 0 else 0

    def get_node_resource(self, node):
        app_client = AppsClient(cluster=self.cluster)
        # cpu usage
        cpu_usage_url = 'http://{host}/api/v1/query?query=sum(rate(container_cpu_usage_seconds_total{{id=\"/\",kubernetes_io_hostname="{hostname}"}}[5m]))/sum(machine_cpu_cores{{kubernetes_io_hostname="{hostname}"}})'
        cpu_usage_query_url = cpu_usage_url.format(host=self.host, hostname=node.name)
        req = requests.get(cpu_usage_query_url)
        cpu_usage_json = req.json()
        if cpu_usage_json['status'] == 'success' and len(cpu_usage_json['data']['result']) > 0:
            node.cpu_usage = cpu_usage_json['data']['result'][0]['value'][1]
        # cpu total
        cpu_total_url = 'http://{host}/api/v1/query?query=sum(machine_cpu_cores{{kubernetes_io_hostname="{hostname}"}})'
        cpu_total_query_url = cpu_total_url.format(host=self.host, hostname=node.name)
        cpu_total_req = requests.get(cpu_total_query_url)
        cpu_total_json = cpu_total_req.json()
        if cpu_total_json['status'] == 'success' and len(cpu_total_json['data']['result']) > 0:
            node.cpu = cpu_total_json['data']['result'][0]['value'][1]
        # mem total
        mem_total_url = 'http://{host}/api/v1/query?query=sum(machine_memory_bytes{{kubernetes_io_hostname="{hostname}"}})'
        mem_total_query_url = mem_total_url.format(host=self.host, hostname=node.name)
        mem_total_req = requests.get(mem_total_query_url)
        mem_total_json = mem_total_req.json()
        if mem_total_json['status'] == 'success' and len(mem_total_json['data']['result']) > 0:
            node.mem = int(mem_total_json['data']['result'][0]['value'][1]) / 1024 / 1024 / 1024
        # mem total
        mem_usage_url = 'http://{host}/api/v1/query?query=sum(container_memory_working_set_bytes{{id=\"/\",kubernetes_io_hostname="{hostname}"}})/sum(machine_memory_bytes{{kubernetes_io_hostname="{hostname}"}})'
        mem_usage_query_url = mem_usage_url.format(host=self.host, hostname=node.name)
        mem_usage_req = requests.get(mem_usage_query_url)
        mem_usage_json = mem_usage_req.json()
        if mem_usage_json['status'] == 'success' and len(mem_usage_json['data']['result']) > 0:
            node.mem_usage = mem_usage_json['data']['result'][0]['value'][1]
        return node

    def get_msg_from_loki(self, cluster_name):
        label_url = "http://{host}/loki/api/v1/label/container_name/values"
        label_query_url = label_url.format(host=self.host)
        label_req = requests.get(label_query_url)
        loki_containers = []
        now = time.time()
        # 乘以1000000是loki的要求
        end = int(round(now * 1000 * 1000000))
        start = int(round(now * 1000 - 3600000) * 1000000)
        if label_req.ok:
            label_req_json = label_req.json()
            values = label_req_json.get('values', [])
            for name in values:
                error_count = 0
                prom_url = 'http://{host}/api/prom/query?limit=1000&query={{container_name="{name}"}}&start={start}&end={end}'
                prom_query_url = prom_url.format(host=self.host, name=name, start=start, end=end)
                prom_req = requests.get(prom_query_url)
                if prom_req.ok:
                    prom_req_json = prom_req.json()
                    streams = prom_req_json.get('streams', [])
                    for stream in streams:
                        entries = stream.get('entries', [])
                        for entry in entries:
                            line = entry.get('line', None)
                            if line is not None and 'level=error' in line:
                                error_count = error_count + 1
                if error_count > 0:
                    loki_container = LokiContainer(name, error_count, cluster_name)
                    loki_containers.append(loki_container.__dict__)
        return loki_containers
