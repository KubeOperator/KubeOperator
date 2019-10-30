from kubeops_api.models.cluster import Cluster
from kubeops_api.models.setting import Setting
from kubeops_api.models.node import Node
from kubeops_api.prometheus_client import PrometheusClient
from kubeops_api.models.cluster_health_history import ClusterHealthHistory

import datetime


def get_cluster_health_msg_hour():
    clusters = Cluster.objects.all()
    for cluster in clusters:
        cluster.change_to()
        nodes = Node.objects.all()
        host = "prometheus.apps"
        for node in nodes:
            print('---------------')
            role_names=[]
            for role in node.roles.all():
                role_names.append(role.name)
            print(role_names)
            if 'master' in role_names:
                host = host + node.name[7:]
        config = {
            'host': host
        }
        prometheus_client = PrometheusClient(config)
        result = prometheus_client.handle_targets_message(prometheus_client.targets())
        if result['success']:
            month = datetime.datetime.now().strftime('%Y-%m')
            clusterHealthHistory = ClusterHealthHistory(project_id=cluster.id,available_rate=result['rate'],
                                                        date_type=ClusterHealthHistory.CLUSTER_HEALTH_HISTORY_DATE_TYPE_HOUR,
                                                        month=month)
            clusterHealthHistory.save()

def handle_cluster_health_msg_day():
    hour_msg = ClusterHealthHistory.objects.filter(date_type=ClusterHealthHistory.CLUSTER_HEALTH_HISTORY_DATE_TYPE_HOUR)
    cluster_rate = 0
    for hour in hour_msg:
        cluster_rate = cluster_rate+hour.available_rate
    cluster_rate = cluster_rate / len(hour_msg)
    month = datetime.datetime.now().strftime('%Y-%m')
    clusterHealthHistory = ClusterHealthHistory(project_id=hour_msg[0].project_id,available_rate=cluster_rate,
                                                date_type=ClusterHealthHistory.CLUSTER_HEALTH_HISTORY_DATE_TYPE_DAY,
                                                month=month)
    clusterHealthHistory.save()
    ClusterHealthHistory.objects.filter(ClusterHealthHistory.CLUSTER_HEALTH_HISTORY_DATE_TYPE_HOUR).delete()
