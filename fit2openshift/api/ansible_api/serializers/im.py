# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from rest_framework import serializers

from .inventory import InventorySerializer
from .role import SimpleRoleSerializer
from .playbook import PlaySerializer
from .adhoc import AdHocSerializer
from ..models import Playbook


__all__ = ['IMPlaybookSerializer', 'IMAdHocSerializer']


class IMBaseSerializer(serializers.Serializer):
    inventory = InventorySerializer(required=False)

    project = None
    inv_serializer = None

    def check_inventory(self):
        hosts = self.initial_data.get("inventory", {}).get("hosts")
        if not hosts:
            raise serializers.ValidationError("hosts is null")

        for host in hosts:
            if not host.get('name'):
                raise serializers.ValidationError({"hosts", "name is null"})

    def is_valid(self, raise_exception=False):
        self.check_inventory()
        return super().is_valid(raise_exception=raise_exception)


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
        pass

    def is_valid(self, raise_exception=False):
        adhoc_data = self.initial_data.get('adhoc')
        if not adhoc_data.get("pattern"):
            raise serializers.ValidationError("pattern is null")
        elif not adhoc_data.get("module"):
            raise serializers.ValidationError("module is null")
        return super().is_valid()

    def update(self, instance, validated_data):
        pass
