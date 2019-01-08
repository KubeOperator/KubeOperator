from rest_framework import serializers
from django.shortcuts import reverse

from ansible_api.serializers import HostSerializer, GroupSerializer, ProjectSerializer
<<<<<<< HEAD
from .models import Cluster, Node, Role, DeployExecution, Package


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


class ClusterSerializer(ProjectSerializer):
    package = serializers.SlugRelatedField(
        queryset=Package.objects.all(),
        slug_field='name', required=False
    )

    class Meta:
        model = Cluster
        fields = ['id', 'name', 'package', 'template', 'comment', 'date_created']
        read_only_fields = ['id', 'date_created']
=======
from .models import Cluster, Node, Role, DeployExecution


class ClusterSerializer(ProjectSerializer):
    configs = serializers.JSONField(read_only=True)

    class Meta:
        model = Cluster
        fields = ['id', 'name', 'configs', 'comment']
        read_only_fields = ['id']
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce


class NodeSerializer(HostSerializer):
    roles = serializers.SlugRelatedField(
<<<<<<< HEAD
        many=True, queryset=Role.objects.all(),
        slug_field='name', required=False
    )
    meta = serializers.JSONField()
=======
        many=True, read_only=False, queryset=Role.objects.all(),
        slug_field='name', required=False,
    )

    class Meta:
        model = Node
        extra_kwargs = HostSerializer.Meta.extra_kwargs
        read_only_fields = list(filter(lambda x: x not in ('groups',), HostSerializer.Meta.read_only_fields))
        fields = list(filter(lambda x: x not in ('groups',), HostSerializer.Meta.fields))
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce

    def get_field_names(self, declared_fields, info):
        names = super().get_field_names(declared_fields, info)
        names.append('roles')
        return names

<<<<<<< HEAD
    def save(self, **kwargs):
        self.validated_data['groups'] = self.validated_data.pop('roles', [])
        return super().save(**kwargs)

    class Meta:
        model = Node
        extra_kwargs = HostSerializer.Meta.extra_kwargs
        fields = [
            'id', 'name', 'ip', 'username', 'password', 'vars', 'comment',
            'roles'
        ]
        read_only_fields = ['id']
=======
    def create(self, validated_data):
        validated_data['groups'] = validated_data.pop('roles', [])
        return super().create(validated_data)
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce


class RoleSerializer(GroupSerializer):
    nodes = serializers.SlugRelatedField(
<<<<<<< HEAD
        many=True,  queryset=Node.objects.all(),
        slug_field='name', required=False
    )
    meta = serializers.JSONField()

    class Meta:
        model = Role
        fields = ['id', 'name', 'nodes', 'children', 'vars', 'meta', 'comment']
        read_only_fields = ['id']


class DeployExecutionSerializer(serializers.ModelSerializer):
=======
        many=True, read_only=False, queryset=Node.objects.all(),
        slug_field='name', required=False
    )

    class Meta:
        model = Role
        fields = ["id", "name", "nodes", "children", "vars", "comment"]
        read_only_fields = ["id", "children", "vars"]


class DeployReadExecutionSerializer(serializers.ModelSerializer):
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
    result_summary = serializers.JSONField(read_only=True)
    log_url = serializers.SerializerMethodField()
    log_ws_url = serializers.SerializerMethodField()

    class Meta:
        model = DeployExecution
<<<<<<< HEAD
        fields = '__all__'
        read_only_fields = [
            'id', 'state', 'num', 'result_summary', 'result_raw',
            'date_created', 'date_start', 'date_end', 'project', 'timedelta'
=======
        fields = [
            'id', 'state', 'num', 'result_summary', 'result_raw',
            'date_created', 'date_start', 'date_end',
        ]
        read_only_fields = [
            'id', 'state', 'num', 'result_summary', 'result_raw',
            'date_created', 'date_start', 'date_end'
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
        ]

    @staticmethod
    def get_log_url(obj):
        return reverse('celery-api:task-log-api', kwargs={'pk': obj.id})

    @staticmethod
    def get_log_ws_url(obj):
        return '/ws/tasks/{}/log/'.format(obj.id)
