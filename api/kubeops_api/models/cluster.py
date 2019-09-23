import logging
import os
import shutil

import yaml
from django.db import models

import kubeops_api
from ansible_api.models import Project, Playbook
from fit2ansible.settings import ANSIBLE_PROJECTS_DIR, CLUSTER_CONFIG_DIR, KUBEEASZ_DIR
from kubeops_api.adhoc import fetch_cluster_config, get_cluster_token
from kubeops_api.cloud_provider import create_hosts, delete_hosts
from common import models as common_models
from kubeops_api.components import get_component_urls
from kubeops_api.models.auth import AuthTemplate
from kubeops_api.models.node import Node
from kubeops_api.models.role import Role
from django.db.models import Q

logger = logging.getLogger(__name__)
__all__ = ["Cluster"]


class Cluster(Project):
    CLUSTER_STATUS_READY = 'READY'
    CLUSTER_STATUS_RUNNING = 'RUNNING'
    CLUSTER_STATUS_ERROR = 'ERROR'
    CLUSTER_STATUS_WARNING = 'WARNING'
    CLUSTER_STATUS_INSTALLING = 'INSTALLING'
    CLUSTER_STATUS_DELETING = 'DELETING'
    CLUSTER_DEPLOY_TYPE_MANUAL = 'MANUAL'
    CLUSTER_DEPLOY_TYPE_AUTOMATIC = 'AUTOMATIC'

    CLUSTER_STATUS_CHOICES = (
        (CLUSTER_STATUS_RUNNING, 'running'),
        (CLUSTER_STATUS_INSTALLING, 'installing'),
        (CLUSTER_STATUS_DELETING, 'deleting'),
        (CLUSTER_STATUS_READY, 'ready'),
        (CLUSTER_STATUS_ERROR, 'error'),
        (CLUSTER_STATUS_WARNING, 'warning')
    )

    CLUSTER_DEPLOY_TYPE_CHOICES = (
        (CLUSTER_DEPLOY_TYPE_MANUAL, 'manual'),
        (CLUSTER_DEPLOY_TYPE_AUTOMATIC, 'automatic'),
    )

    package = models.ForeignKey("Package", null=True, on_delete=models.SET_NULL)
    persistent_storage = models.CharField(max_length=128, null=True, blank=True)
    network_plugin = models.CharField(max_length=128, null=True, blank=True)
    auth_template = models.ForeignKey('kubeops_api.AuthTemplate', null=True, on_delete=models.SET_NULL)
    template = models.CharField(max_length=64, blank=True, default='')
    plan = models.ForeignKey('cloud_provider.Plan', on_delete=models.SET_NULL, null=True)
    worker_size = models.IntegerField(default=3)
    status = models.CharField(max_length=128, choices=CLUSTER_STATUS_CHOICES, default=CLUSTER_STATUS_READY)
    deploy_type = models.CharField(max_length=128, choices=CLUSTER_DEPLOY_TYPE_CHOICES,
                                   default=CLUSTER_DEPLOY_TYPE_MANUAL)
    configs = common_models.JsonDictTextField(default={})

    @property
    def region(self):
        if self.plan:
            return self.plan.region.name

    @property
    def zone(self):
        if self.plan:
            return self.plan.zone.name

    @property
    def zones(self):
        if self.plan.zones:
            zones = []
            for zone in self.plan.zones.all():
                zones.append(zone.name)
            return zones

    @property
    def cloud_provider(self):
        if self.plan:
            return self.plan.region.vars['provider']

    @property
    def current_execution(self):
        current = kubeops_api.models.deploy.DeployExecution.objects.filter(project=self).first()
        return current

    @property
    def resource(self):
        return self.package.meta['resource']

    @property
    def apps(self):
        return get_component_urls(self)

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

    @property
    def node_size(self):
        self.change_to()
        nodes = Node.objects.all().filter(~Q(name__in=['::1', '127.0.0.1', 'localhost']))
        return len(nodes)

    @property
    def current_workers(selfs):
        selfs.change_to()
        return Node.objects.filter(groups__name__in=['worker'])

    def change_status(self, status):
        self.status = status
        self.save()

    def get_steps(self, opt):
        for operation in self.package.meta['operations']:
            if operation['name'] == opt:
                return operation['steps']

    def create_network_plugin(self):
        cluster_configs = self.load_config_file()
        if self.network_plugin:
            networks = cluster_configs.get('networks', [])
            vars = {}
            for net in networks:
                if net["name"] == self.network_plugin:
                    vars = net.get('vars', {})
            self.set_config_unlock(vars)

    def create_storage(self):
        cluster_configs = self.load_config_file()
        if self.persistent_storage:
            storages = cluster_configs.get('storages', [])
            vars = {}
            for storage in storages:
                if storage['name'] == self.persistent_storage:
                    vars = storage.get('vars', {})
            self.set_config_unlock(vars)

    def set_package_configs(self):
        self.configs.update(self.package.meta['vars'])
        self.save()

    def get_template_meta(self):
        for template in self.package.meta.get('templates', []):
            if template['name'] == self.template:
                return template['name']

    def create_playbooks(self):
        for playbook in self.package.meta.get('playbooks', []):
            url = 'file:///{}'.format(os.path.join(KUBEEASZ_DIR))
            Playbook.objects.create(
                name=playbook['name'], alias=playbook['alias'],
                type=Playbook.TYPE_LOCAL, url=url, project=self
            )

    @staticmethod
    def load_config_file():
        with open(os.path.join(CLUSTER_CONFIG_DIR, "config.yml")) as f:
            return yaml.load(f.read())

    def create_roles(self):
        config_file = self.load_config_file()
        _roles = {}
        for role in config_file.get('roles', []):
            _roles[role['name']] = role
        template = None
        for tmp in config_file.get('templates', []):
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

    def set_config(self, k, v):
        cluster = Cluster.objects.select_for_update().get(name=self.name)
        _vars = cluster.configs
        if isinstance(v, str):
            v = v.strip()
        _vars[k] = v
        cluster.configs = _vars
        cluster.save()

    def get_config(self, k):
        v = self.configs.get(k)
        return {'key': k, 'value': v}

    def get_configs(self):
        configs = [{'key': k, 'value': v} for k, v in self.configs.items()]
        return configs

    def del_config(self, k):
        _vars = self.vars
        _vars.pop(k, None)
        self.vars = _vars
        self.save()

    def set_config_unlock(self, vars):
        configs = self.configs
        configs.update(vars)
        self.configs = configs
        self.save()

    def create_node_localhost(self):
        local_nodes = ['localhost', '127.0.0.1', '::1']
        for name in local_nodes:
            node = Node.objects.create(
                name=name, vars={"ansible_connection": "local"},
                project=self, meta={"hidden": True},
            )
            node.set_groups(group_names=['config'])

    def create_node(self, role, host):
        node = Node.objects.create(
            name=host.name,
            host=host,
            project=self
        )
        node.set_groups(group_names=[role])

    def create_resource(self):
        create_hosts(self)

    def destroy_resource(self):
        delete_hosts(self)

    def fetch_config(self):
        path = None
        if self.status == Cluster.CLUSTER_STATUS_RUNNING:
            self.change_to()
            master = self.group_set.get(name='master').hosts.first()
            dest = fetch_cluster_config(master, os.path.join(ANSIBLE_PROJECTS_DIR, self.name))
            path = dest
        return path

    def get_cluster_token(self):
        token = None
        if self.status == Cluster.CLUSTER_STATUS_RUNNING:
            self.change_to()
            master = self.group_set.get(name='master').hosts.first()
            token = get_cluster_token(master)
        return token

    def delete_data(self):
        path = os.path.join(ANSIBLE_PROJECTS_DIR, self.name)
        if os.path.exists(path):
            shutil.rmtree(path)

    def set_plan_configs(self):
        if self.plan and self.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            self.set_config_unlock(self.plan.mixed_vars)

    def on_cluster_create(self):
        self.change_to()
        self.create_roles()
        self.create_playbooks()
        self.create_node_localhost()
        self.create_network_plugin()
        self.set_package_configs()
        self.create_storage()
        self.set_plan_configs()

    def on_cluster_delete(self):
        self.delete_data()
