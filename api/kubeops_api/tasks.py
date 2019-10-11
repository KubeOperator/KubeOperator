import logging

from celery import shared_task
from celery.task import periodic_task

from common.utils import get_object_or_none
from ansible_api.ctx import change_to_root
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.deploy import DeployExecution
import kubeops_api.cluster_backup_utils
from celery.schedules import crontab


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


@shared_task
def test():
    cluster = Cluster.objects.first()
    cluster.create_resource()


def test_task():
    test.apply_async(
        task_id=str(123)
    )

@periodic_task(run_every=crontab(minute='*/15'),name='task.cluster_backup')
def cluster_backup():
    kubeops_api.cluster_backup_utils.cluster_backup()