# -*- coding: utf-8 -*-
#

from rest_framework import serializers

from ..models import Playbook, PlaybookExecution, Play
from .mixins import ProjectSerializerMixin, ReadSerializerMixin
from .base import GitSerializer

__all__ = [
    'PlaybookReadSerializer', 'PlaybookSerializer', 'PlayReadSerializer',
    'PlaybookExecutionSerializer', 'PlaySerializer',
]


class PlayReadSerializer(ReadSerializerMixin, serializers.ModelSerializer):
    vars = serializers.DictField(required=False, allow_null=True)
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


class PlaySerializer(ProjectSerializerMixin, PlayReadSerializer):
    pass


class PlaybookReadSerializer(ReadSerializerMixin, serializers.ModelSerializer):
    plays = PlayReadSerializer(many=True, required=False)
    git = GitSerializer(required=False)

    class Meta:
        model = Playbook
        fields = [
            'id', 'name', 'project', 'type', 'plays', 'git',
            'rel_path',
            'is_periodic', 'interval', 'crontab', 'is_active',
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


class PlaybookSerializer(ProjectSerializerMixin, PlaybookReadSerializer):
    plays = PlaySerializer(many=True, required=False)


class PlaybookExecutionSerializer(ProjectSerializerMixin, serializers.ModelSerializer):
    summary = serializers.JSONField(read_only=True)
    raw = serializers.JSONField(read_only=True)

    class Meta:
        model = PlaybookExecution
        fields = [
            'id', 'playbook', 'is_finished', 'is_success', 'timedelta',
            'raw', 'summary', 'date_start', 'date_finished',
        ]

