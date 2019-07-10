from kubernetes import client, config


def list_nodes(config_path):
    config.load_kube_config(config_file=config_path)
    v1 = client.CoreV1Api.connect_get_node_proxy()
    v1.read_node_("master-1.nmss.f2c.com")
