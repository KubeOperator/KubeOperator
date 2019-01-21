import os
import uuid
import yaml

from django.conf import settings
from django.db import models
from django.utils.translation import ugettext_lazy as _

from common.models import JsonTextField
from ansible_api.models import Project, Host, Group, Playbook
from ansible_api.models.mixins import (
    AbstractProjectResourceModel, AbstractExecutionModel
)
from .signals import pre_deploy_execution_start, post_deploy_execution_start


__all__ = ['Package', 'Cluster', 'Node', 'Role', 'DeployExecution']


# 离线包的model
class Package(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    meta = JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))

    packages_dir = os.path.join(settings.BASE_DIR, 'data', 'packages')

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = _('Package')

    @property
    def path(self):
        return os.path.join(self.packages_dir, self.name)

    @classmethod
    def lookup(cls):
        for d in os.listdir(cls.packages_dir):
            full_path = os.path.join(cls.packages_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            cls.objects.update_or_create(defaults=defaults, name=d)


class Cluster(Project):
    package = models.ForeignKey("Package", null=True, on_delete=models.SET_NULL)
    template = models.CharField(max_length=64, blank=True, default='')

    def create_roles(self):
        _roles = {}
        for role in self.package.meta.get('roles', []):
            _roles[role['name']] = role
        template = None
        for tmp in self.package.meta.get('templates', []):
            if tmp['name'] == self.template:
                template = tmp
                break
        for role in template.get('roles', []):
            _roles[role['name']] = role
        roles_data = [role for role in _roles.values()]

        children_data = {}
        for data in roles_data:
            children_data[data['name']] = data.pop('children', [])
            Role.objects.update_or_create(defaults=data, name=data['name'])

        for name, children_name in children_data.items():
            try:
                role = Role.objects.get(name=name)
                children = Role.objects.filter(name__in=children_name)
                role.children.set(children)
            except Role.DoesNotExist:
                pass

    def create_node_localhost(self):
        Node.objects.create(
            name="localhost", vars={"ansible_connection": "local"},
            project=self, meta={"hidden": True}
        )

    def create_install_playbooks(self):
        for data in self.package.meta.get('install_playbooks', []):
            url = 'file:///{}'.format(os.path.join(self.package.path))
            Playbook.objects.create(
                name=data['name'], alias=data['alias'],
                type=Playbook.TYPE_LOCAL, url=url, project=self,
            )

    def create_playbooks(self):
        self.create_install_playbooks()

    def on_cluster_create(self):
        self.change_to()
        self.create_roles()
        self.create_node_localhost()
        self.create_install_playbooks()

    def configs(self, tp='list'):
        self.change_to()
        role = Role.objects.get(name='OSEv3')
        configs = role.vars
        if tp == 'list':
            configs = [{'key': k, 'value': v} for k, v in configs.items()]
        return configs

    def set_config(self, k, v):
        self.change_to()
        role = Role.objects.get(name='OSEv3')
        _vars = role.vars
        _vars[k] = v
        role.vars = _vars
        role.save()

    def get_config(self, k):
        v = self.configs(tp='dict').get(k)
        return {'key': k, 'value': v}

    def del_config(self, k):
        self.change_to()
        role = Role.objects.get(name='OSEv3')
        _vars = role.vars
        _vars.pop(k, None)
        role.vars = _vars
        role.save()


class ClusterConfig(models.Model):
    key = models.CharField(max_length=1024)
    value = JsonTextField()

    class Meta:
        abstract = True


class Node(Host):
    class Meta:
        proxy = True

    @property
    def roles(self):
        return self.groups

    @roles.setter
    def roles(self, value):
        self.groups.set(value)

    def add_vars(self, _vars):
        __vars = {k: v for k, v in self.vars.items()}
        __vars.update(_vars)
        if self.vars != __vars:
            self.vars = __vars
            self.save()

    def remove_var(self, key):
        __vars = self.vars
        if key in __vars:
            del __vars[key]
            self.vars = __vars
            self.save()

    def get_var(self, key, default):
        return self.vars.get(key, default)


class Role(Group):
    class Meta:
        proxy = True

    @property
    def nodes(self):
        return self.hosts

    @nodes.setter
    def nodes(self, value):
        self.hosts.set(value)

    def __str__(self):
        return "%s %s" % (self.project, self.name)


class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    project = models.ForeignKey('ansible_api.Project', on_delete=models.CASCADE)

    def start(self):
        result = {"raw": {}, "summary": {}}
        pre_deploy_execution_start.send(self.__class__, execution=self)
        for playbook in self.project.playbook_set.all().order_by('name'):
            print("\n>>> Start run {} ".format(playbook.name))
            _result = playbook.execute()
            result["summary"].update(_result["summary"])
            if not _result.get('summary', {}).get('success', False):
                break
        post_deploy_execution_start.send(self.__class__, execution=self, result=result)
        return result

    class Meta:
        get_latest_by = 'date_created'
        ordering = ('-date_created', )



