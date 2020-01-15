from kubeops_api.models.setting import Setting

http_prefix = 'http://'
https_prefix = 'https://'


def get_component_urls(cluster):
    urls = {}
    app_url = cluster.get_config("APP_DOMAIN").get('value')
    if app_url:
        urls = {
            "grafana": http_prefix + "grafana." + app_url,
            "prometheus": http_prefix + "prometheus." + app_url,
            "registry-ui": http_prefix + "registry-ui." + app_url,
            "dashboard": https_prefix + "dashboard." + app_url,
            "traefik": http_prefix + "traefik." + app_url,
            "scope": http_prefix + "scope.weave." + app_url,
            "ceph": http_prefix + "ceph." + app_url,
            "kubeapps-plus": http_prefix + "kubeapps-plus." + app_url
        }
    return urls
