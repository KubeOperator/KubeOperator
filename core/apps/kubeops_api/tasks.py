import logging
import json
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

logger = logging.getLogger("kubeops")


@shared_task
def sync_host_info(host_id):
    host = Host.objects.get(id=host_id)
    try:
        print("gather info ...")
        host.gather_info()
    except Exception as e:
        logger.error("refresh host {} status error: {}".format(host.name, e.args))


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


@periodic_task(run_every=crontab(minute=0, hour=1), name='task.cluster_backup', time_limit=1200)
def cluster_backup():
    try:
        kubeops_api.cluster_backup_utils.cluster_backup()
    except Exception as e:
        logger.error("cluster_backup error")


@periodic_task(run_every=crontab(minute='*/5'), name='task.save_cluster_data', time_limit=360)
def save_cluster_data():
    try:
        kubeops_api.cluster_monitor.put_cluster_data_to_redis()
    except Exception as e:
        logger.error("save_cluster_data error")


@periodic_task(run_every=crontab(minute=0, hour='*/1'), name='task.get_loki_data_hour', time_limit=1200)
def get_loki_data_hour():
    try:
        kubeops_api.cluster_monitor.put_loki_data_to_redis()
    except Exception as e:
        logger.error("get_loki_data_hour error")


@periodic_task(run_every=crontab(minute='*/5'), name='task.save_cluster_event', time_limit=300)
def save_cluster_event():
    try:
        kubeops_api.cluster_monitor.put_event_data_to_es()
    except Exception as e:
        logger.error("save_cluster_event error" + json.dumps(e))


@periodic_task(run_every=crontab(minute="*/5"), name='task.host_health_check')
def host_health_check():
    for host in Host.objects.all():
        logger.info("start host: {} health check".format(host.name))
        try:
            host.health_check()
        except Exception as e:
            logger.error("host {} health check error".format(host.name))


@periodic_task(run_every=crontab(minute="*/5"), name='task.node_health_check')
def node_health_check():
    for cluster in Cluster.objects.filter(status=Cluster.CLUSTER_STATUS_RUNNING):
        logger.info("start cluster: {} health check".format(cluster.name))
        try:
            cluster.node_health_check()
        except Exception as e:
            logger.error("cluster {} health check error".format(cluster.name))
