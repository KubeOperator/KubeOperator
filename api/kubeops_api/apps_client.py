import requests


class AppsClient():
    def __init__(self, cluster):
        self.cluster = cluster

    def get(self, app, url):
        app_domain = app + "." + self.cluster.get_config('APP_DOMAIN')
        header = {
            "Host": app_domain
        }
        host_ip = self.cluster.get_first_master().ip
        req_url = str(url).replace(app_domain, host_ip)
        return requests.get(headers=header, url=req_url)
