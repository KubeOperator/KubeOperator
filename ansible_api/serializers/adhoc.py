# -*- coding: utf-8 -*-
#

from rest_framework import serializers

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
            'id', 'pattern', 'module', 'args', 'project',
        ]
        read_only_fields = (
            'id', 'created_by',
        )


class AdHocSerializer(ProjectSerializerMixin, AdHocReadSerializer):
    pass


class AdHocExecutionSerializer(serializers.ModelSerializer):
    summary = serializers.JSONField(read_only=True)
    result = serializers.JSONField(read_only=True)

    class Meta:
        model = AdHocExecution
        fields = [
            'id', 'adhoc', 'is_finished', 'is_success', 'timedelta',
            'raw', 'summary', 'date_start', 'date_finished',
        ]

