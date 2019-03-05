# -*- coding: utf-8 -*-
#

from django.shortcuts import reverse
from rest_framework import serializers

from ..models import Playbook, PlaybookExecution, Play
from .mixins import ProjectSerializerMixin, ReadSerializerMixin, ExecutionSerializerMixin
from .base import GitSerializer

__all__ = [
    'PlaybookReadSerializer', 'PlaybookSerializer', 'PlayReadSerializer',
    'PlaybookExecutionSerializer', 'PlaySerializer',
]


class PlayReadSerializer(ReadSerializerMixin, serializers.ModelSerializer):
    vars = serializers.DictField(required=False, default={})
    tasks = serializers.ListField(child=serializers.DictField(), required=False, allow_null=True)
    roles = serializers.ListField(child=serializers.DictField(), required=False, allow_null=True)

    class Meta:
        model = Play
        read_only_fields = ['id']
        fields = ['id', 'name', 'pattern', 'vars', 'tasks', 'roles', 'project']

    def validate(self, data):
        if not data.get("tasks") and not data.get('roles'):
            raise serializers.ValidationError(
                {"tasks": "tasks or roles require one"}
            )
        return data


class PlaySerializer(PlayReadSerializer, ProjectSerializerMixin):
    pass


class PlaybookReadSerializer(ReadSerializerMixin, serializers.ModelSerializer):
    plays = PlayReadSerializer(many=True, required=False)
    extra_vars = serializers.DictField(required=False, default={})
    git = GitSerializer(required=False)

    class Meta:
        model = Playbook
        fields = [
            'id', 'name', 'alias', 'project', 'type', 'plays', 'git',
            'extra_vars', 'is_periodic', 'interval', 'crontab', 'is_active',
            'comment', 'created_by', 'date_created', 'project'
        ]
        read_only_fields = ['id', 'created_by', 'date_created']

    def save(self, **kwargs):
        plays_data = self.validated_data.pop('plays', None)
        playbook = super().save(**kwargs)
        if plays_data is not None:
            play_serializer = self.fields.get('plays')
            play_serializer.initial_data = plays_data
            play_serializer.is_valid(raise_exception=True)
            plays = play_serializer.save()
            playbook.plays.all().delete()
            playbook.plays.set(plays)
        return playbook


class PlaybookSerializer(PlaybookReadSerializer, ProjectSerializerMixin):
    plays = PlaySerializer(many=True, required=False)


class PlaybookExecutionSerializer(serializers.ModelSerializer, ProjectSerializerMixin):
    extra_vars = serializers.DictField(required=False, default={})
    result_summary = serializers.JSONField(read_only=True)
    log_url = serializers.SerializerMethodField()
    log_ws_url = serializers.SerializerMethodField()

    class Meta:
        model = PlaybookExecution
        fields = [
            'id', 'extra_vars', 'num', 'state', 'timedelta', 'log_url',
            'log_ws_url', 'result_summary', 'date_created', 'date_start',
            'date_end', 'playbook', 'project'
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


