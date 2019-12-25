import logging
import threading

from celery import shared_task
from celery.task import periodic_task

from common.utils import get_object_or_none
from ansible_api.ctx import change_to_root
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.deploy import DeployExecution
import kubeops_api.cluster_backup_utils
import kubeops_api.cluster_health_utils
from celery.schedules import crontab
import kubeops_api.cluster_monitor
from kubeops_api.models.host import Host

logger = logging.getLogger(__name__)


@shared_task
def start_deploy_execution(eid, **kwargs):
    change_to_root()
    execution = get_object_or_none(DeployExecution, id=eid)
    if execution:
        execution.project.change_to()
        return execution.start
    else:
        msg = "No execution found: {}".format(eid)
        return {"error": msg}


@periodic_task(run_every=crontab(minute=0, hour=1), name='task.cluster_backup')
def cluster_backup():
    kubeops_api.cluster_backup_utils.cluster_backup()


@periodic_task(run_every=crontab(minute='*/5'), name='task.save_cluster_data')
def save_cluster_data():
    kubeops_api.cluster_monitor.put_cluster_data_to_redis()


@periodic_task(run_every=crontab(minute=0, hour='*/1'), name='task.get_loki_data_hour')
def get_loki_data_hour():
    kubeops_api.cluster_monitor.put_loki_data_to_redis()


@periodic_task(run_every=crontab(minute='*/3'), name='task.save_cluster_event')
def save_cluster_event():
    kubeops_api.cluster_monitor.put_event_data_to_es()


@periodic_task(run_every=crontab(minute="*/5"), name='task.host_health_check')
def host_health_check():
    for host in Host.objects.all():
        logger.info("start host: {} health check".format(host.name))
        t = threading.Thread(target=host.health_check())
        t.start()


@periodic_task(run_every=crontab(minute="*/5"), name='task.node_health_check')
def node_health_check():
    for cluster in Cluster.objects.all():
        logger.info("start cluster: {} health check".format(cluster.name))
        t = threading.Thread(target=cluster.node_health_check())
        t.start()
