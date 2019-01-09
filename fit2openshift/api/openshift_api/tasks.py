from celery import shared_task

from common.utils import get_object_or_none
from ansible_api.ctx import change_to_root

from .models import DeployExecution


@shared_task
def start_deploy_execution(eid, **kwargs):
    change_to_root()
    execution = get_object_or_none(DeployExecution, id=eid)
    if execution:
        execution.project.change_to()
        return execution.start()
    else:
        msg = "No execution found: {}".format(eid)
        print(msg)
        return {"error": msg}
