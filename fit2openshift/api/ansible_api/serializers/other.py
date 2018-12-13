# -*- coding: utf-8 -*-
#

from rest_framework import serializers


__all__ = [
    'ModuleSerializer', 'TaskSerializer',
]


class ModuleSerializer(serializers.Serializer):
    pass


class TaskSerializer(serializers.Serializer):
    task = serializers.CharField(max_length=1024)
