from rest_framework import serializers

from kubeops_api.models.host import Host
from storage.models import NfsStorage

__all__ = [
    'NfsStorageSerializer'
]


class NfsStorageSerializer(serializers.ModelSerializer):
    meta = serializers.JSONField()
    nfs_host = serializers.SlugRelatedField(
        queryset=Host.objects.all(),
        slug_field='name', required=False
    )

    class Meta:
        model = NfsStorage
        read_only_fields = ['id', 'date_created']
        fields = ['id', 'name', 'server', 'allow_ip', 'nfs_host', 'path', 'date_created']
