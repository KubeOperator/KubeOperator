from rest_framework import serializers

from storage.models import NfsStorage

__all__ = [
    'NfsStorageSerializer'
]


class NfsStorageSerializer(serializers.ModelSerializer):
    meta = serializers.JSONField()

    class Meta:
        model = NfsStorage
        read_only_fields = ['id', 'name', 'server', 'path', 'date_created']
        fields = ['id', 'name', 'date_created']
