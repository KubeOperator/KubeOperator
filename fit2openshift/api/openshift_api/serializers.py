from rest_framework import serializers
from django.shortcuts import reverse

from ansible_api.serializers import GroupSerializer, ProjectSerializer
from ansible_api.serializers import HostSerializer as AnsibleHostSerializer
from ansible_api.serializers.inventory import HostReadSerializer
from .models import Cluster, Node, Role, DeployExecution, Package, Host

__all__ = [
    'PackageSerializer', 'ClusterSerializer', 'NodeSerializer',
    'RoleSerializer', 'DeployExecutionSerializer',
]


class PackageSerializer(serializers.ModelSerializer):
    meta = serializers.JSONField()

    class Meta:
        model = Package
        read_only_fields = ['id', 'name', 'meta', 'date_created']
        fields = ['id', 'name', 'meta', 'date_created']


class HostSerializer(HostReadSerializer):
    meta = serializers.JSONField()

    class Meta:
        model = Host
        extra_kwargs = HostReadSerializer.Meta.extra_kwargs
        fields = [
            'id', 'name', 'ip', 'username', 'password', 'comment', 'memory', 'os', 'os_version', 'cpu_core', 'comment',
            'cluster'
        ]
        read_only_fields = ['id', 'memory', 'os', 'os_version', 'cpu_core', 'comment', 'cluster']


class ClusterSerializer(ProjectSerializer):
    package = serializers.SlugRelatedField(
        queryset=Package.objects.all(),
        slug_field='name', required=False
    )

    class Meta:
        model = Cluster
        fields = ['id', 'name', 'package', 'template', 'comment', 'current_task_id', 'state', 'date_created', ]
        read_only_fields = ['id', 'date_created', 'current_task_id', 'state']


class ClusterConfigSerializer(serializers.Serializer):
    key = serializers.CharField(max_length=128)
    value = serializers.JSONField()


class NodeSerializer(AnsibleHostSerializer):
    roles = serializers.SlugRelatedField(
        many=True, queryset=Role.objects.all(),
        slug_field='name', required=False
    )

    meta = serializers.JSONField()

    def get_field_names(self, declared_fields, info):
        names = super().get_field_names(declared_fields, info)
        names.append('roles')
        return names

    def save(self, **kwargs):
        self.validated_data['groups'] = self.validated_data.pop('roles', [])
        return super().save(**kwargs)

    class Meta:
        model = Node
        extra_kwargs = AnsibleHostSerializer.Meta.extra_kwargs

        fields = [
            'id', 'name', 'ip', 'vars', 'roles', 'host', 'host_memory', 'host_cpu_core', 'host_os', 'host_os_version'
        ]
        read_only_fields = ['id', 'host_memory', 'host_cpu_core', 'host_os', 'host_os_version', 'ip']


class RoleSerializer(GroupSerializer):
    nodes = serializers.SlugRelatedField(
        many=True, queryset=Node.objects.all(),
        slug_field='name', required=False
    )
    meta = serializers.JSONField()

    class Meta:
        model = Role
        fields = ['id', 'name', 'nodes', 'children', 'vars', 'meta', 'comment']
        read_only_fields = ['id']


class DeployExecutionSerializer(serializers.ModelSerializer):
    result_summary = serializers.JSONField(read_only=True)
    log_url = serializers.SerializerMethodField()
    log_ws_url = serializers.SerializerMethodField()
    progress_ws_url = serializers.SerializerMethodField()

    class Meta:
        model = DeployExecution
        fields = '__all__'
        read_only_fields = [
            'id', 'state', 'num', 'result_summary', 'result_raw',
            'date_created', 'date_start', 'date_end', 'project', 'timedelta', 'current_task', 'progress'
        ]

    @staticmethod
    def get_log_url(obj):
        return reverse('celery-api:task-log-api', kwargs={'pk': obj.id})

    @staticmethod
    def get_log_ws_url(obj):
        return '/ws/tasks/{}/log/'.format(obj.id)

    @staticmethod
    def get_progress_ws_url(obj):
        return '/ws/progress/{}/'.format(obj.id)
