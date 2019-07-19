import requests

http_prefix = 'http://'
https_prefix = 'https://'


def get_component_urls(cluster):
    urls = {}
    app_url = cluster.get_config("APP_DOMAIN").get('value')
    if app_url:
        grafana_urls = generate_grafana_urls(cluster, app_url)
        urls.update(grafana_urls)
        prometheus_urls = generate_prometheus_url(cluster, app_url)
        urls.update(prometheus_urls)
        registry_urls = generate_registry_url(cluster, app_url)
        urls.update(registry_urls)
        dashboard_urls = generate_dashboard_url(cluster, app_url)
        urls.update(dashboard_urls)
    return urls


def generate_grafana_urls(cluster, app_url):
    db_url = http_prefix + 'grafana.' + app_url
    urls = {"grafana": db_url}
    if cluster.status == 'RUNNING':
        urls.update(list_grafana_dbs(db_url))
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


def generate_prometheus_url(cluster, app_url):
    prometheus_url = http_prefix + "prometheus." + app_url
    return {
        "prometheus": prometheus_url
    }


def generate_registry_url(cluster, app_url):
    registry_url = http_prefix + "registry." + app_url
    registry_ui_url = http_prefix + "registry-ui." + app_url
    return {
        "registry": registry_url,
        "registry-ui": registry_ui_url
    }


def generate_dashboard_url(cluster, app_url):
    dashboard_url = https_prefix + "dashboard." + app_url
    return {
        "dashboard": dashboard_url
    }
