# -*- coding: utf-8 -*-
#
import uuid
import yaml

from django.db import models, transaction

from common import models as common_models
from .utils import name_validator
from .mixins import AbstractProjectResourceModel
from ..inventory import LocalModelInventory


__all__ = ['ClusterHost', 'ClusterGroup', 'Host', 'Group', 'Inventory']


class BaseHost(models.Model):
    name = models.CharField(max_length=1024, validators=[name_validator])
    ip = models.GenericIPAddressField(null=True)
    port = models.IntegerField(default=22)
    username = models.CharField(max_length=1024, default='root')
    password = common_models.EncryptCharField(max_length=4096, blank=True, null=True)
    private_key = common_models.EncryptCharField(max_length=8192, blank=True, null=True)
    vars = common_models.JsonDictTextField(default={})
    meta = common_models.JsonDictTextField(default={})
    comment = models.TextField(blank=True)

    class Meta:
        abstract = True

    @property
    def ansible_vars(self):
        host_vars = {k: v for k, v in self.vars.items()}
        host_vars['ansible_ssh_host'] = self.ip
        host_vars['ansible_ssh_port'] = self.port
        host_vars['ansible_ssh_user'] = self.username
        host_vars['ansible_ssh_pass'] = self.password
        return host_vars

    def add_to_groups(self, group_names, auto_create=True):
        with transaction.atomic():
            for name in group_names:
                group = self.groups.model.get_group(name=name, auto_create=auto_create)
                group.hosts.add(self)

    def set_groups(self, group_names, auto_create=True):
        with transaction.atomic():
            groups = []
            for name in group_names:
                group = self.groups.model.get_group(name=name, auto_create=auto_create)
                groups.append(group)
            self.groups.set(groups)


class ClusterHost(BaseHost):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    projects = models.ManyToManyField('Project', related_name='cluster_hosts')

    @property
    def group_cls(self):
        return ClusterGroup

    class Meta:
        unique_together = ('name',)

    def __str__(self):
        return self.name


class Host(AbstractProjectResourceModel, BaseHost):
    is_project_private = True

    class Meta:
        unique_together = ('name', 'project')

    def __str__(self):
        return '{}: {}'.format(self.project, self.name)


class BaseGroup(models.Model):
    GROUP_NAME = (
        ('master', 'master'),
        ('node', 'node'),
        ('etcd', 'etcd'),

    )
    name = models.CharField(choices=GROUP_NAME, max_length=64, validators=[name_validator])
    vars = common_models.JsonDictTextField(default={})
    hosts = models.ManyToManyField('BaseHost', related_name='groups')
    children = models.ManyToManyField('BaseGroup', related_name='parents', blank=True)
    meta = common_models.JsonDictTextField(default={})
    comment = models.TextField(blank=True)

    class Meta:
        abstract = True

    @property
    def children_names(self):
        return [child.name for child in self.children.all()]

    @property
    def hosts_names(self):
        return [host.name for host in self.hosts.all()]

    @classmethod
    def get_group(cls, name, auto_create=True):
        try:
            group = cls.objects.get(name=name)
        except cls.DoesNotExist as e:
            if auto_create:
                group = cls(name=name)
                group.save()
            else:
                raise e
        return group

    def add_children(self, group_names, auto_create=True):
        with transaction.atomic():
            groups = []
            for name in group_names:
                group = self.children.model.get_group(name, auto_create=auto_create)
                groups.append(group)
            self.children.add(*groups)

    def add_hosts(self, host_names):
        with transaction.atomic():
            hosts = []
            for host_name in host_names:
                host = self.hosts.model.objects.get(name=host_name)
                hosts.append(host)
            self.hosts.add(*hosts)


class ClusterGroup(BaseGroup):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    hosts = models.ManyToManyField('ClusterHost', related_name='groups')
    children = models.ManyToManyField('ClusterGroup', related_name='parents', blank=True)
    projects = models.ManyToManyField('Project', related_name='cluster_groups')

    class Meta:
        unique_together = ('name',)

    def __str__(self):
        return self.name


class Group(AbstractProjectResourceModel, BaseGroup):
    hosts = models.ManyToManyField('Host', related_name='groups')
    children = models.ManyToManyField('Group', related_name='parents', blank=True)

    is_project_private = True

    class Meta:
        unique_together = ('name', 'project')

    def __str__(self):
        return '{}: {}'.format(self.project, self.name)


class Inventory:
    def __init__(self, hosts=None, groups=None):
        if hosts is None and groups is None:
            raise OSError('Hosts or groups at least one required')
        self._hosts = hosts
        self._groups = groups

    @property
    def hosts(self):
        hosts = list(self._hosts)
        return hosts

    @property
    def groups(self):
        groups = list(self._groups)
        return groups

    def as_object(self):
        return LocalModelInventory(self)

    def get_data_yaml(self):
        return self.get_data(fmt='yaml')

    def get_data(self, fmt='py'):
        data = {}
        group_all_hosts = {}
        group_all_data = {'hosts': group_all_hosts}
        data['all'] = group_all_data

        for host in self.hosts:
            group_all_hosts[host.name] = host.ansible_vars

        for group in self.groups:
            group_data = {}
            children = {child: {} for child in group.children_names}
            hosts = {host: {} for host in group.hosts_names}
            if group.vars:
                group_data["vars"] = group.vars
            if children:
                group_data['children'] = children
            if hosts:
                group_data['hosts'] = hosts
            data[group.name] = group_data

        if fmt in ['py']:
            return data
        elif fmt in ['yaml', 'yaml']:
            return yaml.safe_dump(data, default_flow_style=False)
        else:
            return None
