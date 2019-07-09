import logging

from celery import shared_task

from celery_api.utils import register_as_period_task
from common.utils import get_object_or_none
from ansible_api.ctx import change_to_root
from kubeops_api.models.deploy import DeployExecution

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
        print(msg)
        return {"error": msg}


@shared_task
@register_as_period_task(interval=20)
def test():
    logger.info('test')
