from rest_framework import serializers

from .models import Cluster, Node, Role, DeployExecution
from ansible_api.serializers import HostSerializer, GroupSerializer, ProjectSerializer


class ClusterSerializer(ProjectSerializer):
    class Meta:
        model = Cluster
        fields = ['id', 'name', 'comment']
        read_only_fields = ['id']


class NodeSerializer(HostSerializer):
    roles = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=Role.objects.all(),
        slug_field='name', required=False,
    )

    class Meta:
        model = Node
        extra_kwargs = HostSerializer.Meta.extra_kwargs
        read_only_fields = list(filter(lambda x: x not in ('vars', 'groups'), HostSerializer.Meta.read_only_fields))
        fields = list(filter(lambda x: x not in ('vars', 'groups'), HostSerializer.Meta.fields))

    def get_field_names(self, declared_fields, info):
        names = super().get_field_names(declared_fields, info)
        names.append('roles')
        return names


class RoleSerializer(GroupSerializer):
    nodes = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=Node.objects.all(),
        slug_field='name', required=False
    )

    class Meta:
        model = Role
        fields = ["id", "name", "nodes"]
        read_only_fields = ["id"]


class DeployExecutionSerializer(serializers.Serializer):
    pass


class DeployReadExecutionSerializer(serializers.ModelSerializer):

    class Meta:
        model = DeployExecution
        fields = ['id', 'state', 'num', 'result_summary', 'result_raw', 'date_created', 'date_start', 'date_end',]
        read_only_fields = ['id', 'state', 'num', 'result_summary', 'result_raw', 'date_created', 'date_start', 'date_end']
