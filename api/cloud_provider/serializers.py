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
        read_only_fields = ['id', 'cluster_size', 'plan_size', 'status', 'date_created']
        fields = ['id', 'name', 'cluster_size', 'plan_size', 'vars', 'date_created', 'cloud_zone', 'region', 'status']


class PlanSerializer(serializers.ModelSerializer):
    region = serializers.SlugRelatedField(
        queryset=Region.objects.all(),
        slug_field='name', required=True
    )
    zone = serializers.SlugRelatedField(
        queryset=Zone.objects.all(),
        slug_field='name', required=True,
    )
    vars = serializers.DictField(required=False, default={})

    class Meta:
        model = Plan
        read_only_fields = ['id', 'date_created']
        fields = ['id', 'name', 'vars', 'date_created', 'zone', 'region']
