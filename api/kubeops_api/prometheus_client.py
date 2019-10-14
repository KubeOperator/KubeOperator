import requests


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