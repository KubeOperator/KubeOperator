# coding: utf-8
import logging
from celery import shared_task

from common.utils import get_object_or_none
from .models import Playbook, AdHoc, Role, PlaybookExecution, AdHocExecution
from .ansible.runner import AdHocRunner
from .ctx import set_current_project, change_to_root
from .inventory import WithHostInfoInventory

logger = logging.getLogger(__file__)


@shared_task
def execute_playbook(tid, **kwargs):
    change_to_root()
    playbook = get_object_or_none(Playbook, id=tid)
    if playbook:
        set_current_project(playbook.project)
        return playbook.execute()
    else:
        msg = "No playbook found: {}".format(tid)
        logger.error(msg)
        return {"error": msg}


@shared_task
def start_playbook_execution(eid, **kwargs):
    change_to_root()
    execution = get_object_or_none(PlaybookExecution, id=eid)
    if execution:
        set_current_project(execution.project)
        return execution.start()
    else:
        msg = "No execution found: {}".format(eid)
        logger.error(msg)
        return {"error": msg}


@shared_task
def execute_adhoc(tid, **kwargs):
    change_to_root()
    adhoc = get_object_or_none(AdHoc, id=tid)
    if adhoc:
        set_current_project(adhoc.project)
        return adhoc.execute()
    else:
        msg = "No adhoc found: {}".format(tid)
        logger.error(msg)
        return {"error": msg}


@shared_task
def start_adhoc_execution(eid, **kwargs):
    change_to_root()
    execution = get_object_or_none(AdHocExecution, id=eid)
    if execution:
        set_current_project(execution.project)
        return execution.start()
    else:
        msg = "No execution found: {}".format(eid)
        logger.error(msg)
        return {"error": msg}


@shared_task
def install_role(tid, **kwargs):
    role = get_object_or_none(Role, id=tid)
    if not role:
        return {"error": "Role {} not found".find(tid)}
    if role.state != Role.STATE_NOT_INSTALL:
        return {"error": "Role {} may be installed".find(role.name)}
    return role.install()


@shared_task
def run_im_adhoc(adhoc_data, inventory_data):
    inventory = WithHostInfoInventory(inventory_data)
    runner = AdHocRunner(inventory)
    pattern = adhoc_data.get('pattern') or ''
    module = adhoc_data.get('module') or 'ping'
    args = adhoc_data.get('args') or ''
    tasks = [{'action': {'module': module, 'args': args}}]
    result = runner.run(tasks, pattern=pattern)
    return result


@shared_task
def hello(name, callback=None):
    print("hello")


@shared_task
def hello_callback(result):
    print("Hello {} :".format(result))
    result += ':'
    return result
