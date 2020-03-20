from rest_framework import serializers
from django.shortcuts import reverse

from cloud_provider.models import Plan, Zone
from kubeops_api.models import Condition
from kubeops_api.models.credential import Credential
from kubeops_api.models.host import Host, GPU
from ansible_api.serializers import GroupSerializer, ProjectSerializer
from ansible_api.serializers import HostSerializer as AnsibleHostSerializer
from ansible_api.serializers.inventory import HostReadSerializer
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.deploy import DeployExecution
from kubeops_api.models.host import Volume
from kubeops_api.models.node import Node
from kubeops_api.models.package import Package
from kubeops_api.models.role import Role
from kubeops_api.models.backup_storage import BackupStorage
from kubeops_api.models.backup_strategy import BackupStrategy
from kubeops_api.models.cluster_backup import ClusterBackup
from kubeops_api.models.cluster_health_history import ClusterHealthHistory
from kubeops_api.models.item import Item
from kubeops_api.serializers.host import ConditionSerializer

__all__ = [
    'PackageSerializer', 'ClusterSerializer', 'NodeSerializer', 'RoleSerializer', 'DeployExecutionSerializer',
    'CredentialSerializer', 'BackupStrategySerializer', 'BackupStorageSerializer'
]


class CredentialSerializer(serializers.ModelSerializer):
    class Meta:
        model = Credential
        extra_kwargs = {
            'password': {'write_only': True},
            'private_key': {'write_only': True},
        }
        fields = ['id', 'name', 'username', 'password', 'private_key', 'date_created', 'type']
        read_only_fields = ['id', 'date_created']


class PackageSerializer(serializers.ModelSerializer):
    meta = serializers.JSONField()

    class Meta:
        model = Package
        read_only_fields = ['id', 'name', 'meta', 'date_created']
        fields = ['id', 'name', 'meta', 'date_created']


class ClusterConfigSerializer(serializers.Serializer):
    key = serializers.CharField(max_length=128)
    value = serializers.JSONField()


class NodeSerializer(AnsibleHostSerializer):
    roles = serializers.SlugRelatedField(
        many=True, queryset=Role.objects.all(),
        slug_field='name', required=False
    )
    conditions = ConditionSerializer(required=False, many=True)
    meta = serializers.JSONField()
    info = serializers.DictField()

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
            'id', 'name', 'ip', 'port', 'vars', 'roles', 'host', 'host_memory', 'host_cpu_core', 'host_os',
            'host_os_version',
            'status', 'conditions', "info"
        ]
        read_only_fields = ['id', 'host_memory', 'host_cpu_core', 'host_os', 'host_os_version', 'ip', 'status',
                            'conditions', 'info']


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
    steps = serializers.ListField(required=False, read_only=True)
    params = serializers.DictField(required=False, read_only=False)

    class Meta:
        model = DeployExecution
        fields = '__all__'
        read_only_fields = [
            'id', 'state', 'num', 'result_summary', 'result_raw',
            'date_created', 'date_start', 'date_end', 'project', 'timedelta', 'steps',
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


class ClusterSerializer(ProjectSerializer):
    package = serializers.SlugRelatedField(
        queryset=Package.objects.all(),
        slug_field='name', required=False
    )
    plan = serializers.SlugRelatedField(
        queryset=Plan.objects.all(),
        slug_field='name', required=False
    )
    meta = serializers.DictField(required=False)
    configs = serializers.DictField(required=False)
    apps = serializers.JSONField(read_only=True)
    current_execution = DeployExecutionSerializer(read_only=True)

    class Meta:
        model = Cluster
        fields = ['id', 'name', 'package', 'worker_size', 'persistent_storage', 'network_plugin', 'template',
                  'plan',
                  'comment', 'date_created', 'resource', 'resource_version', 'current_execution', 'status', 'nodes',
                  'apps', 'deploy_type', 'zone', 'region', 'meta', 'zones', 'cloud_provider', 'configs',
                  'cluster_doamin_suffix', 'item_name']
        read_only_fields = ['id', 'date_created', 'current_execution', 'resource', 'resource_version',
                            'nodes', 'apps', 'zone', 'region', 'meta', 'zones', 'cloud_provider', 'item_name']


class BackupStorageSerializer(serializers.ModelSerializer):
    credentials = serializers.DictField()

    class Meta:
        model = BackupStorage
        fields = ['id', 'name', 'region', 'credentials', 'type', 'date_created', 'status']
        read_only_fields = ['id', 'date_created']


class BackupStrategySerializer(serializers.ModelSerializer):
    class Meta:
        model = BackupStrategy
        fields = ['id', 'cron', 'save_num', 'project_id', 'backup_storage_id', 'status']


class ClusterBackupSerializer(serializers.ModelSerializer):
    backup_storage = serializers.SlugRelatedField(
        queryset=BackupStorage.objects.all(),
        slug_field='id', required=False
    )

    class Meta:
        model = ClusterBackup
        fields = ['id', 'name', 'size', 'date_created', 'project_id', 'folder', 'backup_storage']


class ClusterHealthSerializer(serializers.Serializer):
    type = serializers.CharField()
    data = serializers.CharField()
    status = serializers.CharField()


class ClusterHeathHistorySerializer(serializers.ModelSerializer):
    class Meta:
        model = ClusterHealthHistory
        fields = ['id', 'project_id', 'available_rate', 'date_type', 'date_created', 'month']


class ItemSerializer(serializers.ModelSerializer):
    class Meta:
        model = Item
        fields = ['id', 'name', 'description', 'date_created']
