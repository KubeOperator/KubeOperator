from celery import shared_task

from common.utils import get_object_or_none
from ops.models.script import ScriptExecution


@shared_task
def start_script_execution(eid, **kwargs):
    execution = get_object_or_none(ScriptExecution, id=eid)
    if execution:
        return execution.start
    else:
        msg = "No execution found: {}".format(eid)
        return {"error": msg}
