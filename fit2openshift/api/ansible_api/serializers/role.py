# -*- coding: utf-8 -*-
#
from rest_framework import serializers
from rest_framework_bulk.serializers import BulkListSerializer

from ..models import Role
from .mixins import ProjectSerializerMixin, ReadSerializerMixin
from .base import GitSerializer


class RoleReadSerializer(ReadSerializerMixin, serializers.ModelSerializer):
    git = GitSerializer(required=False)
    meta = serializers.JSONField(required=False, allow_null=True)
    meta_ext = serializers.JSONField(required=False, allow_null=True)

    class Meta:
        model = Role
        list_serializer_class = BulkListSerializer
        fields = [
            'id', 'name', 'type', 'galaxy_name', 'git', 'url',
            'logo', 'logo_url', 'categories', 'version', 'state', 'comment',
            'meta', 'meta_ext', 'project', 'created_by', 'date_created',
            'date_updated'
        ]
        read_only_fields = [
            'id', 'state', 'meta', 'create_by',
            'date_created', 'date_updated', 'logo_url',
            'created_by'
        ]


class RoleSerializer(ProjectSerializerMixin, RoleReadSerializer):
    pass


class SimpleRoleSerializer(RoleSerializer):
    def get_field_names(self, declared_fields, info):
        return ['name', 'type', 'galaxy_name', 'url', 'git', 'project']
