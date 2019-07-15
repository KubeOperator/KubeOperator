import requests


def generate_grafana_urls(cluster):
    app_url = cluster.get_config("APP_DOMAIN").get('value')
    db_url = 'http://' + 'grafana.' + app_url
    urls = list_grafana_dbs(db_url)
    return urls


def list_grafana_dbs(db_url):
    urls = {}
    etcd_title = 'Etcd by Prometheus'
    apps_title = 'Kubernetes Apps'
    cluster_health_title = 'Kubernetes Cluster Health'
    nodes_title = 'Kubernetes Nodes (prometheus)'
    cluster_title = 'Cluster Monitoring for Kubernetes'
    try:
        res = requests.get(db_url + "/api/search")
        dbs = res.json()
        for db in dbs:
            url = db_url + db.get('url')
            if db.get('title') == etcd_title:
                urls['etcd_grafana'] = url
            elif db.get('title') == apps_title:
                urls['apps_grafana'] = url
            elif db.get('title') == cluster_health_title:
                urls['cluster_health_grafana'] = url
            elif db.get('title') == nodes_title:
                urls['nodes_grafana'] = url
            elif db.get('title') == cluster_title:
                urls['cluster_grafana'] = url
    except Exception as e:
        urls = {}
    return urls
