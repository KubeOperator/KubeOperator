import logging
import os

from django.db import models

import kubeops_api
from ansible_api.models import Project, Playbook
from fit2ansible.settings import ANSIBLE_PROJECTS_DIR
from kubeops_api.adhoc import fetch_cluster_config
from kubeops_api.models.auth import AuthTemplate
from kubeops_api.models.node import Node
from kubeops_api.models.role import Role
from django.db.models import Q

logger = logging.getLogger(__name__)


class Cluster(Project):
    CLUSTER_STATUS_READY = 'READY'
    CLUSTER_STATUS_RUNNING = 'RUNNING'
    CLUSTER_STATUS_ERROR = 'ERROR'
    CLUSTER_STATUS_WARNING = 'WARNING'
    CLUSTER_STATUS_INSTALLING = 'INSTALLING'
    CLUSTER_STATUS_DELETING = 'DELETING'

    CLUSTER_STATUS_CHOICES = (
        (CLUSTER_STATUS_RUNNING, 'running'),
        (CLUSTER_STATUS_INSTALLING, 'installing'),
        (CLUSTER_STATUS_DELETING, 'deleting'),
        (CLUSTER_STATUS_READY, 'ready'),
        (CLUSTER_STATUS_ERROR, 'error'),
        (CLUSTER_STATUS_WARNING, 'warning')
    )

    package = models.ForeignKey("Package", null=True, on_delete=models.SET_NULL)
    persistent_storage = models.ForeignKey('Storage', null=True, on_delete=models.SET_NULL)
    auth_template = models.ForeignKey('kubeops_api.AuthTemplate', null=True, on_delete=models.SET_NULL)
    template = models.CharField(max_length=64, blank=True, default='')
    config_path = models.CharField(max_length=128, blank=True, null=True, default=None)
    status = models.CharField(max_length=128, choices=CLUSTER_STATUS_CHOICES, default=CLUSTER_STATUS_READY)

    @property
    def current_execution(self):
        current = kubeops_api.models.deploy.DeployExecution.objects.filter(project=self).first()
        return current

    @property
    def resource(self):
        return self.package.meta['resource']

    @property
    def operations(selfs):
        return selfs.package.meta['operations']

    @property
    def resource_version(self):
        return self.package.meta['version']

    @property
    def nodes(self):
        self.change_to()
        nodes = Node.objects.all().filter(~Q(name__in=['::1', '127.0.0.1', 'localhost']))
        n = []
        for node in nodes:
            n.append(node.name)
        return n

    def create_storage(self):
        if self.persistent_storage:
            print(self.persistent_storage.vars)
            self.set_config_storage(self.persistent_storage.vars)

    def get_template_meta(self):
        for template in self.package.meta.get('templates', []):
            if template['name'] == self.template:
                return template['name']

    def create_playbooks(self):
        for playbook in self.package.meta.get('playbooks', []):
            url = 'file:///{}'.format(os.path.join(self.package.path))
            Playbook.objects.create(
                name=playbook['name'], alias=playbook['alias'],
                type=Playbook.TYPE_LOCAL, url=url, project=self
            )

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
        config_role = Role.objects.get(name='config')
        private_var = template['private_vars']
        role_vars = config_role.vars
        role_vars.update(private_var)
        config_role.vars = role_vars
        config_role.save()

    def configs(self, tp='list'):
        self.change_to()
        role = Role.objects.get(name='config')
        configs = role.vars
        if tp == 'list':
            configs = [{'key': k, 'value': v} for k, v in configs.items()]
        return configs

    def set_config(self, k, v):
        self.change_to()
        role = Role.objects.select_for_update().get(name='config')
        _vars = role.vars
        if isinstance(v, str):
            v = v.strip()
        _vars[k] = v
        role.vars = _vars
        role.save()

    def set_config_storage(self, vars):
        self.change_to()
        config_role = Role.objects.get(name='config')
        role_vars = config_role.vars
        role_vars.update(vars)
        config_role.vars = role_vars
        config_role.save()
        config_role.vars = role_vars
        config_role.save()

    def get_config(self, k):
        v = self.configs(tp='dict').get(k)
        return {'key': k, 'value': v}

    def del_config(self, k):
        self.change_to()
        role = Role.objects.get(name='config')
        _vars = role.vars
        _vars.pop(k, None)
        role.vars = _vars
        role.save()

    def create_node_localhost(self):
        local_nodes = ['localhost', '127.0.0.1', '::1']
        for name in local_nodes:
            node = Node.objects.create(
                name=name, vars={"ansible_connection": "local"},
                project=self, meta={"hidden": True},
            )
            node.set_groups(group_names=['config'])

    def fetch_config(self):
        self.change_to()
        master = self.group_set.get(name='master').hosts.first()
        dest = fetch_cluster_config(master, os.path.join(ANSIBLE_PROJECTS_DIR, self.name))
        self.config_path = dest
        self.save()

    def on_cluster_create(self):
        self.change_to()
        self.create_roles()
        self.create_playbooks()
        self.create_node_localhost()
        self.create_storage()
