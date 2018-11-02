# coding: utf-8
import sys
import os
import logging
from celery import shared_task, subtask

from common.utils import get_object_or_none
from .models import Playbook, AdHoc, Role
from .ctx import set_current_project, change_to_root

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
def execute_adhoc(tid, passwords=None, **kwargs):
    change_to_root()
    adhoc = get_object_or_none(AdHoc, id=tid)
    if adhoc:
        set_current_project(adhoc.project)
        return adhoc.execute(passwords)
    else:
        msg = "No adhoc found: {}".format(tid)
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
def hello(name, callback=None):
    print("hello")


@shared_task
def hello_callback(result):
    print("Hello {} :".format(result))
    result += ':'
    return result
