# -*- coding: utf-8 -*-
#

from rest_framework import serializers
from django.shortcuts import reverse

from ..models import AdHoc, AdHocExecution
from .mixins import ProjectSerializerMixin

__all__ = [
    'AdHocReadSerializer', 'AdHocSerializer', 'AdHocExecutionSerializer',
]


class AdHocReadSerializer(serializers.ModelSerializer):
    args = serializers.JSONField(required=False, allow_null=True)
    result = serializers.JSONField(read_only=True)
    summary = serializers.JSONField(read_only=True)

    class Meta:
        model = AdHoc
        fields = [
            'id', 'pattern', 'module', 'args', 'project', 'summary', 'result'
        ]
        read_only_fields = (
            'id', 'created_by',
        )


class AdHocSerializer(AdHocReadSerializer, ProjectSerializerMixin):
    pass


class AdHocExecutionSerializer(serializers.ModelSerializer, ProjectSerializerMixin):
    result_summary = serializers.JSONField(read_only=True)
    log_url = serializers.SerializerMethodField()
    log_ws_url = serializers.SerializerMethodField()

    class Meta:
        model = AdHocExecution
        fields = [
            'id', 'num', 'state', 'timedelta', 'log_url', 'log_ws_url',
            'result_summary', 'date_created', 'date_start', 'date_end',
            'adhoc', 'project'
        ]
        read_only_fields = [
            'id', 'num', 'state', 'timedelta', 'log_url', 'log_ws_url',
            'result_summary', 'date_created', 'date_start', 'date_end',
            'project'
        ]

    @staticmethod
    def get_log_url(obj):
        return reverse('celery-api:task-log-api', kwargs={'pk': obj.id})

    @staticmethod
    def get_log_ws_url(obj):
        return '/ws/tasks/{}/log/'.format(obj.id)

