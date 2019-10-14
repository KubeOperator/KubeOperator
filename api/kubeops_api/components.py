from kubeops_api.models.setting import Setting

http_prefix = 'http://'
https_prefix = 'https://'


def get_component_urls(cluster):
    urls = {}
    domain_suffix = Setting.objects.get(key="domain_suffix")
    app_url = "apps.{}.{}".format(cluster.name, domain_suffix.value)
    if app_url:
        urls = {
            "grafana": http_prefix + "grafana." + app_url,
            "prometheus": http_prefix + "prometheus." + app_url,
            "registry-ui": http_prefix + "registry-ui." + app_url,
            "dashboard": https_prefix + "dashboard." + app_url,
            "traefik": http_prefix + "traefik." + app_url,
            "scope": http_prefix + "scope.weave." + app_url
        }
    return urls
