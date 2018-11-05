# -*- coding: utf-8 -*-
#

import logging
import uuid

from django.dispatch import receiver
from django.db.models.signals import pre_delete, post_save, post_delete
from django.utils import timezone
from celery import current_task

from .models import Playbook, Role, PlaybookExecution
from .signals import pre_execution_start, post_execution_start

logger = logging.getLogger(__file__)


@receiver(post_save, sender=Playbook)
def on_playbook_create_or_update(sender, instance, created, **kwargs):
    if instance.is_periodic and instance.is_active:
        if created:
            instance.create_period_task()
    else:
        instance.disable_period_task()


@receiver(pre_delete, sender=Playbook)
def on_playbook_delete(sender, instance, **kwargs):
    instance.remove_period_task()


@receiver(post_save, sender=Role)
def on_role_create_or_update(sender, instance, created, **kwargs):
    from .tasks import install_role
    if created:
        install_role.delay(instance.id)


@receiver(pre_execution_start)
def on_execution_start(sender, execution, **kwargs):
    execution.date_start = timezone.now()
    execution.state = execution.STATE_STARTED
    execution.save()


@receiver(post_execution_start)
def on_execution_end(sender, execution, result, **kwargs):
    date_finished = timezone.now()
    timedelta = (timezone.now() - execution.date_start).seconds
    state = execution.STATE_FAILURE
    if result.get('summary', {}).get("success", False):
        state = execution.STATE_SUCCESS
    execution.result_summary = result.get('summary', {})
    execution.result_raw = result.get('raw', {})
    execution.state = state
    execution.date_finished = date_finished
    execution.timedelta = timedelta
    execution.save()
