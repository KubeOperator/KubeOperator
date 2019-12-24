import logging

from celery.task import periodic_task
from celery.schedules import crontab
import threading

from kubeops_api.models.cluster import Cluster
from kubeops_api.models.host import Host

logger = logging.getLogger("health_check")


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
