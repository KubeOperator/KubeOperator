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
    # host = "172.16.10.99"
    # token = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkJaQ2dLU0FZYWJfdXhEWEVDQjFMcVFHTXc2UHctV1V4Nk5TRVY4NXRNX0kifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJ0aWxsZXItdG9rZW4tcDY5bTYiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoidGlsbGVyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiNTcyN2Y4NWItMmY0ZS00ZWU5LWIyZGMtYTI1YTI1ZWEwNjE0Iiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Omt1YmUtc3lzdGVtOnRpbGxlciJ9.eA0pjpCVcneAWzk7cFb4ymIhf5dkKSJ0O4j3klRidwX9jFz_rDkfWFLOu39ohXXqaXVib-ZjIG8963pz0bGZYaWVHtjgPkD5OYXsdAyIqtOwlFdlK7lSGmVwXhOqfVfTEEFyZG-UoxCaoRnD6ZCDH6U6pHSnFb-XJhPWyNk-FGCH6mvvC3Zy42aL97rF-7sX8kFigelEAYwc8BMWYPDE_i3w3Lmyi4ldFaK0YfrL_pR8j6dlIaMLrPpg1Wif3wsIRbakqHPbWzHywLuMQCzdKruvQeoab59U7POl034XNJpxX2Diz_7HoYnN_gD1Fv82PCMqrgJykejBRHucelu9pg"
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
