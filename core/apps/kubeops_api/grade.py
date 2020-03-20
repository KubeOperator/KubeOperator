from validator.validator import run_validate
from kubeops_api.models.cluster import Cluster
from kubeops_api.utils.kubernetes import get_first_master_host


def validate_cluster(cluster: Cluster):
    token = cluster.get_cluster_token()
    host = get_first_master_host(cluster)
    return run_validate(host, token)
