# -*- coding: utf-8 -*-
from __future__ import unicode_literals
import uuid

from rest_framework import serializers

from .inventory import InventorySerializer
from .role import SimpleRoleSerializer
from .playbook import PlaySerializer
from .adhoc import AdHocSerializer
from ..ctx import set_current_project
from ..models import Playbook, Project


__all__ = ['IMPlaybookSerializer', 'IMAdHocSerializer']


class IMBaseSerializer(serializers.Serializer):
    inventory = InventorySerializer(required=False)

    project = None
    inv_serializer = None

    def is_valid(self, raise_exception=False):
        self.project = Project.objects.create(
            name=str(uuid.uuid4()), comment='#IM#'
        )
        set_current_project(self.project)
        self.inv_serializer = InventorySerializer(
            data=self.initial_data.pop('inventory', {})
        )
        self.inv_serializer.is_valid(raise_exception=True)
        return super().is_valid(raise_exception=raise_exception)

    def create_inventory(self):
        self.inv_serializer.save()


class IMPlaybookSerializer(IMBaseSerializer):
    roles = SimpleRoleSerializer(many=True, required=False, allow_null=True)
    plays = PlaySerializer(many=True)

    def create(self, validated_data):
        self.create_inventory()
        self.create_roles()
        plays = self.create_plays()
        playbook = Playbook.objects.create(name=self.project.name, project=self.project)
        playbook.plays.set(plays)
        return playbook

    def update(self, instance, validated_data):
        pass

    def create_roles(self):
        roles_data = self.validated_data.get('roles')
        if not roles_data:
            return
        serializer = SimpleRoleSerializer(
            data=roles_data, many=True,
        )
        serializer.is_valid(raise_exception=True)
        serializer.save()

    def create_plays(self):
        serializer = PlaySerializer(
            data=self.validated_data.get('plays'),
            many=True,
        )
        serializer.is_valid(raise_exception=True)
        return serializer.save()


class IMAdHocSerializer(IMBaseSerializer):
    adhoc = AdHocSerializer(required=True)

    project = None
    inv_serializer = None

    def create(self, validated_data):
        self.create_inventory()
        adhoc = self.create_adhoc()
        return adhoc

    def update(self, instance, validated_data):
        pass

    def create_adhoc(self):
        serializer = AdHocSerializer(
            data=self.validated_data.get('adhoc'),
        )
        serializer.is_valid(raise_exception=True)
        return serializer.save()
