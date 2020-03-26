import kubernetes

from kubeops_api.models.cluster import Cluster


def get_api_client(cluster: Cluster):
    configuration = kubernetes.client.Configuration()
    configuration.api_key_prefix['authorization'] = 'Bearer'
    configuration.api_key['authorization'] = cluster.get_cluster_token()
    configuration.debug = True
    configuration.host = "https://{}:6443".format(get_first_master_host(cluster))
    configuration.verify_ssl = False
    return kubernetes.client.AppsV1Api(configuration)


def get_first_master_host(cluster: Cluster):
    cluster.change_to()
    master = cluster.group_set.get(name='master').hosts.first()
    if master and master.ip:
        return master.ip
