from rest_framework import serializers
from django.shortcuts import reverse

from cloud_provider.models import Plan, Zone
from kubeops_api.models.auth import AuthTemplate
from kubeops_api.models.credential import Credential
from kubeops_api.models.host import Host
from ansible_api.serializers import GroupSerializer, ProjectSerializer
from ansible_api.serializers import HostSerializer as AnsibleHostSerializer
from ansible_api.serializers.inventory import HostReadSerializer
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.deploy import DeployExecution
from kubeops_api.models.host import Volume
from kubeops_api.models.node import Node
from kubeops_api.models.package import Package
from kubeops_api.models.role import Role
from kubeops_api.models.setting import Setting
from kubeops_api.models.backup_storage import BackupStorage
from kubeops_api.models.backup_strategy import BackupStrategy
from kubeops_api.models.cluster_backup import ClusterBackup

__all__ = [
    'PackageSerializer', 'ClusterSerializer', 'NodeSerializer',
    'RoleSerializer', 'DeployExecutionSerializer', 'SettingSerializer', 'HostSerializer',
    'CredentialSerializer', 'BackupStrategySerializer','BackupStorageSerializer'
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


class SettingSerializer(serializers.ModelSerializer):
    class Meta:
        model = Setting
        fields = ['id', 'name', 'key', 'helper', 'order', 'value']
        read_only_fields = ['id', 'name', 'key', 'helper', 'order']


class PackageSerializer(serializers.ModelSerializer):
    meta = serializers.JSONField()

    class Meta:
        model = Package
        read_only_fields = ['id', 'name', 'meta', 'date_created']
        fields = ['id', 'name', 'meta', 'date_created']


class AuthTemplateSerializer(serializers.ModelSerializer):
    meta = serializers.JSONField()

    class Meta:
        model = AuthTemplate
        read_only_fields = ['id', 'name', 'meta', 'date_created']
        fields = ['id', 'name', 'meta', 'date_created']


class VolumeSerializer(serializers.ModelSerializer):
    class Meta:
        model = Volume
        fields = [
            'id', 'name', 'size',
        ]
        read_only_fields = ['id', 'name', 'size', ]


class HostSerializer(HostReadSerializer):
    credential = serializers.SlugRelatedField(
        queryset=Credential.objects.all(),
        slug_field='name', required=False
    )
    zone = serializers.SlugRelatedField(
        queryset=Zone.objects.all(),
        slug_field='name', required=False
    )
    volumes = VolumeSerializer(required=False, many=True)

    class Meta:
        model = Host
        extra_kwargs = HostReadSerializer.Meta.extra_kwargs
        fields = [
            'id', 'name', 'ip', 'cluster', 'credential', 'memory', 'os', 'os_version', 'cpu_core', 'volumes', 'zone',
            'region'
        ]
        read_only_fields = ['id', 'comment', 'memory', 'os', 'os_version', 'cpu_core', 'volumes', 'zone', 'region']


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
    apps = serializers.JSONField(read_only=True)
    current_execution = DeployExecutionSerializer(read_only=True)

    class Meta:
        model = Cluster
        fields = ['id', 'name', 'package', 'worker_size', 'persistent_storage', 'network_plugin', 'template', 'plan',
                  'comment', 'date_created', 'resource', 'resource_version', 'current_execution', 'status', 'nodes',
                  'apps', 'deploy_type', 'zone', 'region', 'meta', 'zones', 'cloud_provider']
        read_only_fields = ['id', 'date_created', 'current_execution', 'status', 'resource', 'resource_version',
                            'nodes', 'apps', 'zone', 'region', 'meta', 'zones', 'cloud_provider']


class BackupStorageSerializer(serializers.ModelSerializer):
    credentials = serializers.DictField()

    class Meta:
        model = BackupStorage
        fields = ['id', 'name', 'region', 'credentials', 'type', 'date_created', 'status']
        read_only_fields = ['id', 'date_created']

class BackupStrategySerializer(serializers.ModelSerializer):

    class Meta:
        model = BackupStrategy
        fields = ['id', 'cron' ,'save_num', 'project_id','backup_storage_id']

class ClusterBackupSerializer(serializers.ModelSerializer):

    backup_storage = serializers.SlugRelatedField(
        queryset=BackupStorage.objects.all(),
        slug_field='id', required=False
    )

    class Meta:
        model = ClusterBackup
        fields = ['id', 'name', 'size', 'date_created', 'project_id', 'folder', 'backup_storage']
