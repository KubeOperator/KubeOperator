# -*- coding: utf-8 -*-
#

import logging
import uuid

from django.dispatch import receiver
from django.db.models.signals import pre_delete, post_save, post_delete
from django.utils import timezone
from celery import current_task

from .models import Playbook, Role, PlaybookExecution
from .signals import pre_playbook_exec, post_playbook_exec

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


@receiver(pre_playbook_exec)
def on_playbook_start(sender, playbook, save_history, **kwargs):
    if not save_history:
        return
    pk = current_task.request.id if current_task else str(uuid.uuid4())
    history = PlaybookExecution.objects.create(
        id=pk, playbook=playbook, project=playbook.project
    )
    playbook._history = history


@receiver(post_playbook_exec)
def on_playbook_end(sender, playbook, save_history, result, **kwargs):
    if not hasattr(playbook, '_history'):
        return
    if not save_history:
        return
    playbook.times += 1
    playbook.save()
    date_finished = timezone.now()
    timedelta = (timezone.now() - playbook._history.date_start).seconds
    data = {
        'raw': result.get('raw'),
        'summary': result.get('summary'),
        'num': playbook.times,
        'is_success': result.get('summary', {}).get("success", False),
        'is_finished': True,
        'date_finished': date_finished,
        'timedelta': timedelta
    }
    PlaybookExecution.objects.filter(pk=playbook._history.id).update(**data)



# @receiver(post_save, sender=AdHoc)
# def on_adhoc_create_or_update(sender, instance, created, **kwargs):
#     from .tasks import run_adhoc
#     if created:
#         run_adhoc.delay(instance.id)
