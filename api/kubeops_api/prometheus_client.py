import requests
from kubeops_api.models.host import Host


class PrometheusClient():

    def __init__(self, config):
        self.host = config.get("host", None)
        self.table_name = config.get("table_name", None)
        self.param = config.get("param", None)
        self.start = config.get("start", None)
        self.end = config.get("end", None)

    def query(self):
        url = "http://{host}/api/v1/query?query={table_name}{param}&start={start}&end={end}"
        query_url = url.format(host=self.host, table_name=self.table_name, param=self.param, start=self.start,
                               end=self.end)
        req = requests.get(query_url)
        return req.json()

    def targets(self):
        url = "http://{host}/api/v1/targets"
        query_url = url.format(host=self.host)
        req = requests.get(query_url)
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
        # cpu usage
        cpu_usage_url = 'http://{host}/api/v1/query?query=sum(rate(container_cpu_usage_seconds_total{{id=\"/\",kubernetes_io_hostname="{hostname}"}}[5m]))/sum(machine_cpu_cores{{kubernetes_io_hostname="{hostname}"}})'
        cpu_usage_query_url = cpu_usage_url.format(host=self.host, hostname=node.name)
        req = requests.get(cpu_usage_query_url)
        cpu_usage_json = req.json()
        if cpu_usage_json['status'] == 'success' and len(cpu_usage_json['data']['result']) > 0:
            node.cpu_usage = cpu_usage_json['data']['result'][0]['value'][1]
        # cpu total
        cpu_total_url = 'http://{host}/api/v1/query?query=sum(machine_cpu_cores{{kubernetes_io_hostname="{hostname}"}})'
        cpu_total_query_url = cpu_total_url.format(host=self.host,hostname=node.name)
        cpu_total_req = requests.get(cpu_total_query_url)
        cpu_total_json = cpu_total_req.json()
        if cpu_total_json['status'] == 'success' and len(cpu_total_json['data']['result']) > 0:
            node.cpu = cpu_total_json['data']['result'][0]['value'][1]
        # mem total
        mem_total_url = 'http://{host}/api/v1/query?query=sum(machine_memory_bytes{{kubernetes_io_hostname="{hostname}"}})'
        mem_total_query_url = mem_total_url.format(host=self.host,hostname=node.name)
        mem_total_req = requests.get(mem_total_query_url)
        mem_total_json = mem_total_req.json()
        if mem_total_json['status'] == 'success' and len(mem_total_json['data']['result']) > 0:
            node.mem = int(mem_total_json['data']['result'][0]['value'][1]) / 1024 / 1024 / 1024
        # mem total
        mem_usage_url = 'http://{host}/api/v1/query?query=sum(container_memory_working_set_bytes{{id=\"/\",kubernetes_io_hostname="{hostname}"}})/sum(machine_memory_bytes{{kubernetes_io_hostname="{hostname}"}})'
        mem_usage_query_url = mem_usage_url.format(host=self.host,hostname=node.name)
        mem_usage_req = requests.get(mem_usage_query_url)
        mem_usage_json = mem_usage_req.json()
        if mem_usage_json['status'] == 'success' and len(mem_usage_json['data']['result']) > 0:
            node.mem_usage = mem_usage_json['data']['result'][0]['value'][1]
        return node
