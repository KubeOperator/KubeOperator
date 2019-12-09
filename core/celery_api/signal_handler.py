# -*- coding: utf-8 -*-
#
import logging

from celery.signals import after_setup_logger
from celery.utils.log import get_logger
from kombu.utils.encoding import safe_str

from .logger import CeleryTaskFileHandler

safe_str = lambda x: x
logger = get_logger(__file__)


@after_setup_logger.connect
def add_celery_redis_handler(sender=None, logger=None, loglevel=None, format=None, **kwargs):
    if not logger:
        return
    handler = CeleryTaskFileHandler()
    handler.setLevel(loglevel)
    formatter = logging.Formatter(format)
    handler.setFormatter(formatter)
    logger.addHandler(handler)


# @task_failure.connect
# def on_task_failed(sender, task_id, **kwargs):
#     CeleryTask.objects.filter(id=task_id).update(state=CeleryTask.STATE_FAILURE, date_finished=timezone.now())
