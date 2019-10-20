from rest_framework import serializers

from kubeops_api.models.host import Host
from storage.models import NfsStorage

__all__ = [
    'NfsStorageSerializer'
]


class NfsStorageSerializer(serializers.ModelSerializer):
    vars = serializers.DictField()

    class Meta:
        model = NfsStorage
        read_only_fields = ['id', 'status', 'date_created']
        fields = ['id', 'name', 'vars', 'status', 'date_created']
