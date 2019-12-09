# -*- coding: utf-8 -*-
#
from rest_framework import serializers

__all__ = ['TaskResultSerializer']


class TaskResultSerializer(serializers.Serializer):
    id = serializers.UUIDField()
    result = serializers.JSONField()
    state = serializers.CharField(max_length=16)
