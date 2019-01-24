# -*- coding: utf-8 -*-
#
import uuid
import os

from django.conf import settings
from django.db import models
from django.utils.translation import ugettext_lazy as _

from common import models as common_models

__all__ = ['CeleryTask']


class CeleryTask(models.Model):
    STATE_PENDING = 'PENDING'
    STATE_STARTED = 'STARTED'
    STATE_FAILURE = 'FAILURE'
    STATE_SUCCESS = 'SUCCESS'
    STATE_RETRY = 'RETRY'

    LOG_DIR = os.path.join(settings.BASE_DIR, 'data', 'celery')
    STATUS_CHOICES = (
        (STATE_PENDING, _('Pending')),
        (STATE_STARTED, _('Started')),
        (STATE_SUCCESS, _('Success')),
        (STATE_FAILURE, _('Failure')),
        (STATE_RETRY, _('Retry')),
    )
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    root_id = models.UUIDField()
    # root_id = models.UUIDField()
    # parent_id = models.UUIDField(null=True)
    name = models.CharField(max_length=256)
    state = models.CharField(max_length=16, choices=STATUS_CHOICES, default=STATE_PENDING)
    result = common_models.JsonTextField(null=True)
    date_start = models.DateTimeField(auto_now_add=True)
    date_finished = models.DateTimeField(null=True)

    def __str__(self):
        return "{}: {}".format(self.name, self.id)

    @property
    def log_path(self):
        dt = self.date_start.strftime('%Y-%m-%d')
        log_dir = os.path.join(self.LOG_DIR, dt)
        if not os.path.isdir(log_dir):
            os.mkdir(log_dir)
        return os.path.join(log_dir, '{}.log'.format(self.root_id))

    @property
    def is_finished(self):
        return self.state in (self.STATE_SUCCESS, self.STATE_FAILURE)

