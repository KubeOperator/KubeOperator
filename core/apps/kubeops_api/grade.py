import logging
from datetime import datetime
from django.core.cache import cache
from validator.validator import run_validate
from kubeops_api.models.cluster import Cluster
from kubeops_api.utils.kubernetes import get_first_master_host

__cache_key = "_grade"

log = logging.getLogger("kubeops")


def validate_cluster(cluster: Cluster):
    token = cluster.get_cluster_token()
    host = get_first_master_host(cluster)
    data = run_validate(host, token)
    cache_cluster_grade(cluster.id, data)
    return data


def query_cluster_grade(cluster: Cluster):
    _from_cache = cache.get("{}{}".format(cluster.id, __cache_key))
    if _from_cache:
        delta = datetime.now() - _from_cache.created_time
        if delta.seconds > 60:
            return validate_cluster(cluster)
        return _from_cache
    else:
        return validate_cluster(cluster)


def cache_cluster_grade(cluster_id, data):
    cache.set("{}{}".format(cluster_id, __cache_key), data)
