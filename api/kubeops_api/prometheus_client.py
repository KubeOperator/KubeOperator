import requests
from kubeops_api.models.host import Host


class PrometheusClient():

    def __init__(self,config):
        self.host = config.get("host",None)
        self.table_name = config.get("table_name",None)
        self.param = config.get("param",None)
        self.start = config.get("start",None)
        self.end = config.get("end",None)


    def query(self):
        url = "http://{host}/api/v1/query?query={table_name}{param}&start={start}&end={end}"
        query_url = url.format(host=self.host,table_name=self.table_name,param=self.param,start=self.start,end=self.end)
        req = requests.get(query_url)
        return req.json()

    def targets(self):
        url = "http://{host}/api/v1/targets"
        query_url = url.format(host=self.host)
        req = requests.get(query_url)
        return req.json()

    def handle_targets_message(self,json):
        result = {
            'success': True,
            'data': [],
            'rate': 0
        }
        if json.get('status') and json.get('status') == 'success':
            keys = ['kubernetes-control-manager','etcd','kubernetes-nodes','kubernetes-schedule','kubernetes-apiservers']
            for key in keys:
                result['data'].append({
                    'job':key,
                    'data': [],
                    'rate': 0
                })
            active_targets =  json.get('data').get('activeTargets')
            for target in active_targets:
                if target.get('labels').get('job') in keys:
                    index = keys.index(target.get('labels').get('job'))
                    instance_address = target.get('discoveredLabels').get('__address__').split(':')[0]
                    hostName = Host.objects.get(ip=instance_address).name
                    health = target.get('health')
                    result['data'][index]['data'].append({
                        'key':hostName,
                        'value':health
                    })
        else:
             result['success'] = False
        self.calculate_available_rate(result)
        return result

    def calculate_available_rate(self,result):
        service_up = 0
        for res in result['data']:
            job_up = 0
            for job in res['data']:
                if job['value'] == 'up':
                    job_up = job_up +1
            res['rate'] = job_up / len(res['data']) * 100 if len(res['data']) > 0 else 0
            if res['rate'] == 100:
                service_up = service_up+1
        result['rate'] = service_up / len(result['data']) * 100 if len(result['data']) > 0 else 0