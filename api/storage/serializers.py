from rest_framework import serializers

from storage.models import NfsStorage, CephStorage

__all__ = [
    'NfsStorageSerializer', 'CephStorageSerializer'
]


class NfsStorageSerializer(serializers.ModelSerializer):
    vars = serializers.DictField()

    class Meta:
        model = NfsStorage
        read_only_fields = ['id', 'status', 'date_created']
        fields = ['id', 'name', 'vars', 'status', 'date_created']


class CephStorageSerializer(serializers.ModelSerializer):
    vars = serializers.DictField()

    class Meta:
        model = CephStorage
        read_only_fields = ['id', 'name', 'date_created']
        fields = ['id', 'name', 'vars', 'date_created']
