# -*- coding: utf-8 -*-
#

from rest_framework import serializers
from django.shortcuts import reverse

from ..models import Project
from ..models.mixins import AbstractExecutionModel
from ..ctx import set_current_project, get_current_project


class ReadSerializerMixin(serializers.Serializer):
    project = serializers.SlugRelatedField(queryset=Project.objects.all(), slug_field='name')


class ProjectSerializerMixin(serializers.Serializer):
    project = serializers.HiddenField(default=get_current_project)


class ExecutionSerializerMixin(serializers.Serializer):
    result_summary = serializers.JSONField(read_only=True)
    log_url = serializers.SerializerMethodField()
    log_ws_url = serializers.SerializerMethodField()

    _declare_fields = [
        'id', 'num', 'state', 'timedelta', 'log_url', 'log_ws_url',
        'result_summary', 'date_created', 'date_start', 'date_end',
    ]
    _read_only_fields = [
        'id', 'num', 'state', 'timedelta', 'log_url', 'log_ws_url',
        'result_summary', 'date_created', 'date_start', 'date_end',
    ]

    @staticmethod
    def get_log_url(obj):
        return reverse('celery-api:task-log-api', kwargs={'pk': obj.id})

    @staticmethod
    def get_log_ws_url(obj):
        return '/ws/tasks/{}/log/'.format(obj.id)


