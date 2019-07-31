from rest_framework import serializers

from cloud_provider.models import CloudProviderTemplate, Region

__all__ = [
    'CloudProviderTemplateSerializer', 'RegionSerializer'
]


class CloudProviderTemplateSerializer(serializers.ModelSerializer):
    meta = serializers.JSONField()

    class Meta:
        model = CloudProviderTemplate
        read_only_fields = ['id', 'name', 'meta', 'date_created']
        fields = ['id', 'name', 'meta', 'date_created']


class RegionSerializer(serializers.ModelSerializer):
    vars = serializers.DictField(required=False, default={})
    template = serializers.SlugRelatedField(
        queryset=CloudProviderTemplate.objects.all(),
        slug_field='name', required=False
    )

    class Meta:
        model = Region
        read_only_fields = ['id', 'date_created', 'template', 'comment']
        fields = ['id', 'name', 'vars', 'date_created', 'template', 'comment', 'cloud_region']
