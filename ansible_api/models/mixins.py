# -*- coding: utf-8 -*-
#
import uuid

from django.db import models
from django.utils.translation import ugettext_lazy as _

from common import models as common_models
from celery_api.utils import get_celery_task_log_path
from ..ctx import current_project, set_current_project


class ProjectResourceManager(models.Manager):
    def get_queryset(self):
        queryset = super(ProjectResourceManager, self).get_queryset()
        if not current_project:
            return queryset
        if current_project.is_real():
            queryset = queryset.filter(project=current_project.id)
        return queryset

    def create(self, **kwargs):
        if 'project' not in kwargs and current_project.is_real():
            kwargs['project'] = current_project._get_current_object()
        return super().create(**kwargs)

    def all(self):
        if current_project:
            return super().all()
        else:
            return self

    def set_current_org(self, project):
        set_current_project(project)
        return self


class AbstractProjectResourceModel(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    project = models.ForeignKey('Project', on_delete=models.CASCADE)
    objects = ProjectResourceManager()

    name = 'Not-Sure'

    class Meta:
        abstract = True

    def __str__(self):
        return '{}: {}'.format(self.project, self.name)

    def save(self, force_insert=False, force_update=False, using=None,
             update_fields=None):
        if not hasattr(self, 'project') and current_project.is_real():
            self.project = current_project._get_current_object()
        return super().save(force_insert=force_insert, force_update=force_update,
                            using=using, update_fields=update_fields)


class AbstractExecutionModel(models.Model):
    STATE_PENDING = 'PENDING'
    STATE_STARTED = 'STARTED'
    STATE_FAILURE = 'FAILURE'
    STATE_SUCCESS = 'SUCCESS'
    STATE_RETRY = 'RETRY'

    STATUS_CHOICES = (
        (STATE_PENDING, _('Pending')),
        (STATE_STARTED, _('Started')),
        (STATE_SUCCESS, _('Success')),
        (STATE_FAILURE, _('Failure')),
        (STATE_RETRY, _('Retry')),
    )

    timedelta = models.FloatField(default=0.0, verbose_name=_('Time'), null=True)
    state = models.CharField(choices=STATUS_CHOICES, default=STATE_PENDING, max_length=16)
    num = models.IntegerField(default=1)
    result_summary = common_models.JsonDictTextField(blank=True, null=True, default={}, verbose_name=_('Result summary'))
    result_raw = common_models.JsonDictTextField(blank=True, null=True, default={}, verbose_name=_('Result raw'))
    date_created = models.DateTimeField(auto_now_add=True, null=True, verbose_name=_('Create time'))
    date_start = models.DateTimeField(null=True, verbose_name=_('Start time'))
    date_end = models.DateTimeField(null=True, verbose_name=_('End time'))

    class Meta:
        abstract = True

    @property
    def stdout(self):
        with open(self.log_path, 'r') as f:
            data = f.read()
        return data

    @property
    def success_hosts(self):
        return self.result_summary.get('contacted', []) if self.result_summary else []

    @property
    def failed_hosts(self):
        return self.result_summary.get('dark', {}) if self.result_summary else []

    @property
    def log_path(self):
        return get_celery_task_log_path(self.id)
