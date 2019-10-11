from rest_framework import serializers

from cloud_provider.models import CloudProviderTemplate, Region, Zone, Plan

__all__ = [
    'CloudProviderTemplateSerializer', 'RegionSerializer', 'ZoneSerializer',
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
        read_only_fields = ['id', 'zone_size', 'cluster_size', 'date_created', 'template', 'comment']
        fields = ['id', 'zone_size', 'cluster_size', 'name', 'vars', 'date_created', 'template', 'comment',
                  'cloud_region']


class ZoneSerializer(serializers.ModelSerializer):
    vars = serializers.DictField(required=False, default={})
    region = serializers.SlugRelatedField(
        queryset=Region.objects.all(),
        slug_field='name', required=True
    )

    class Meta:
        model = Zone
        read_only_fields = ['id', 'status', 'date_created', 'ip_available_size', 'host_size']
        fields = ['id', 'name', 'vars', 'date_created', 'cloud_zone', 'region', 'status', 'ip_available_size',
                  'host_size', 'provider']


class PlanSerializer(serializers.ModelSerializer):
    region = serializers.SlugRelatedField(
        queryset=Region.objects.all(),
        slug_field='name', required=True
    )
    zone = serializers.SlugRelatedField(
        queryset=Zone.objects.all(),
        slug_field='name', required=False
    )
    zones = serializers.SlugRelatedField(
        queryset=Zone.objects.all(),
        slug_field='name', many=True, required=False
    )
    vars = serializers.DictField(required=False, default={})

    class Meta:
        model = Plan
        read_only_fields = ['id', 'date_created']
        fields = ['id', 'name', 'vars', 'date_created', 'zone', 'zones', 'region', 'deploy_template']
