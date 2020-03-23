import logging
from datetime import datetime
from django.core.cache import cache
from validator.validator import run_validate
from kubeops_api.models.cluster import Cluster

__cache_key = "_grade"

log = logging.getLogger("kubeops")


def validate_cluster(cluster: Cluster):
    # token = cluster.get_cluster_token()
    # host = get_first_master_host(cluster)
    host = "172.16.10.239"
    token = "eyJhbGciOiJSUzI1NiIsImtpZCI6Im1ndnhLdGt0S2laYTdCdUxWcFc3SU1RN2FOcE5Tbmpia2JIYXpqNWJGcVUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrdWJlcm5ldGVzLWRhc2hib2FyZC10b2tlbi05d2dzdiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrdWJlcm5ldGVzLWRhc2hib2FyZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjNjZmUwNTkxLTYwM2ItNDVhNi04ZDAxLTVhZDQ2MjRlZWEyMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprdWJlcm5ldGVzLWRhc2hib2FyZCJ9.PVDQrpD_0bLIwDXbmPsL1EzEfwYtuRqsTUbpK81Su5k_V-v40UfE1ENAFjEyUdeGyzmGI4BpCZpaDNUZKQ4XdZYodYhB981zaA6GM93VSMrwdi2dl5Krjcfj5WcmqcbARcYZu2-9PHWh4UXQkYhaWLcZiM4VhiqFveXC4nMsC_AaALdAoiWYZ743RGdrs1w64rVSguLzZaVXDrFHRXz8cIOHNtzDuoznaXHD0k6g1Lz2cmAwi8dyHcy0LWiOMCZcIxcCXlNl-DyBSwF_plKtRgprjUJgY0zh4nBosz3y3zqKfW2K_hay4XazhvjzbFqQxGGy3RQXdsfCUaPq4HqQrw"
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
