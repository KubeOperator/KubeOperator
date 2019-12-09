# -*- coding: utf-8 -*-
#

from rest_framework import serializers
from rest_framework_bulk import BulkListSerializer
from django.db import transaction

from .mixins import ProjectSerializerMixin, ReadSerializerMixin
from ..models import Inventory, ClusterGroup, ClusterHost, Group, Host
from ..ctx import current_project


__all__ = [
    'ClusterHostSerializer', 'ClusterGroupSerializer',
    'HostSerializer', 'GroupSerializer',
    'InventorySerializer',
]


class ClusterHostSerializer(serializers.ModelSerializer):
    vars = serializers.DictField(required=False, allow_null=True, default={})
    groups = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=ClusterGroup.objects.all(),
        slug_field='name', required=False,
    )

    class Meta:
        model = ClusterHost
        extra_kwargs = {
            'password': {'write_only': True},
            'private_key': {'write_only': True},
        }
        read_only_fields = ['id']
        fields = [
            'id', 'name', 'ip', 'port', 'username', 'password',
            'private_key', 'groups', 'vars',
        ]


class HostReadSerializer(ReadSerializerMixin, serializers.ModelSerializer):
    vars = serializers.DictField(required=False, default={})
    groups = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=ClusterGroup.objects.all(),
        slug_field='name', required=False,
    )

    class Meta:
        model = Host
        extra_kwargs = {
            'password': {'write_only': True},
            'username': {'write_only': True},
            'private_key': {'write_only': True},
        }
        read_only_fields = ['id']
        fields = [
            'id', 'name', 'ip', 'port', 'username', 'password',
            'private_key', 'groups', 'vars', 'project', 'comment'
        ]


class HostSerializer(HostReadSerializer, ProjectSerializerMixin):
    pass


class ClusterGroupSerializer(serializers.ModelSerializer):
    children = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=ClusterGroup.objects.all(),
        slug_field='name', required=False,
    )
    hosts = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=ClusterHost.objects.all(),
        slug_field='name', required=False
    )
    vars = serializers.DictField(required=False, default={})

    class Meta:
        model = ClusterGroup
        list_serializer_class = BulkListSerializer
        read_only_fields = ['id']
        fields = ['id', 'name', 'hosts', 'children', 'vars', ]


class GroupReadSerializer(ReadSerializerMixin, serializers.ModelSerializer):
    children = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=Group.objects.all(),
        slug_field='name', required=False,
    )
    hosts = serializers.SlugRelatedField(
        many=True, read_only=False, queryset=Host.objects.all(),
        slug_field='name', required=False
    )
    vars = serializers.DictField(required=False, default={})

    class Meta:
        model = Group
        list_serializer_class = BulkListSerializer
        read_only_fields = ['id']
        fields = ['id', 'name', 'hosts', 'children', 'vars', 'project', 'comment']


class GroupSerializer(GroupReadSerializer, ProjectSerializerMixin):
    pass


class InventorySerializer(serializers.Serializer):
    hosts = HostSerializer(many=True, read_only=False, required=False)
    groups = GroupSerializer(many=True, read_only=False, required=False)

    hosts_groups_map = {}
    groups_children_map = {}
    groups_hosts_map = {}
    _save_point = None

    def clean_hosts_data(self, initial_data):
        __hosts = initial_data.get('hosts')
        cleaned_hosts = []

        for host in __hosts:
            groups = host.pop('groups', [])
            self.hosts_groups_map[host['name']] = groups
            cleaned_hosts.append(host)
        return cleaned_hosts

    def clean_groups_data(self, initial_data):
        __groups = initial_data.get('groups', [])
        cleaned_groups = []

        for group in __groups:
            group_hosts = group.pop('hosts', [])
            group_children = group.pop('children', [])

            not_exits = set(group_hosts) - set(self.hosts_groups_map.keys())
            if not_exits:
                msg = 'Group `{}` hosts `{}` not exist'.format(group['name'], not_exits)
                raise serializers.ValidationError(msg)
            self.groups_children_map[group['name']] = group_children
            self.groups_hosts_map[group['name']] = group_hosts
            cleaned_groups.append(group)

    def is_valid(self, raise_exception=False):
        if not self.initial_data:
            raise serializers.ValidationError({"inventory": "inventory empty"})
        self._save_point = transaction.savepoint()
        try:
            current_project.clear_inventory()
            self.clean_inventory_data(self.initial_data)
            valid = super().is_valid(raise_exception=raise_exception)
        except serializers.ValidationError as e:
            transaction.savepoint_rollback(self._save_point)
            raise e
        if not valid:
            transaction.savepoint_rollback(self._save_point)
        return valid

    def clean_inventory_data(self, initial_data):
        self.clean_hosts_data(initial_data)
        self.clean_groups_data(initial_data)

    def add_groups(self):
        _groups = self.validated_data.get('groups')
        if not _groups:
            return []
        serializer = GroupSerializer(data=_groups, many=True)
        if serializer.is_valid():
            groups = serializer.save()
        else:
            raise serializers.ValidationError(
                "Groups is not valid: {}".format(serializer.errors)
            )
        return groups

    def add_hosts(self):
        _hosts = self.validated_data.get('hosts')
        if not _hosts:
            return []
        serializer = HostSerializer(
            data=_hosts, many=True,
        )
        if serializer.is_valid():
            hosts = serializer.save()
        else:
            raise serializers.ValidationError(
                "Hosts is not valid: {}".format(serializer.errors)
            )
        return hosts

    def set_host_groups(self):
        for host_name, group_names in self.hosts_groups_map.items():
            host = Host.objects.get(name=host_name)
            host.set_groups(group_names)

    def set_group_children(self):
        for group_name, children in self.groups_children_map.items():
            group = Group.objects.get(name=group_name)
            group.add_children(children)

    def set_group_hosts(self):
        for group_name, hosts_name in self.groups_hosts_map.items():
            group = Group.objects.get(name=group_name)
            hosts = Host.objects.filter(name__in=hosts_name)
            group.hosts.add(*hosts)

    def update(self, instance, validated_data):
        current_project.clear_inventory()
        return self.create(validated_data)

    def create(self, validated_data):
        try:
            hosts = self.add_hosts()
            groups = self.add_groups()
            self.set_host_groups()
            self.set_group_children()
            self.set_group_hosts()
        except serializers.ValidationError as e:
            transaction.savepoint_rollback(self._save_point)
            raise e
        transaction.savepoint_commit(self._save_point)
        return Inventory(hosts=hosts, groups=groups)

