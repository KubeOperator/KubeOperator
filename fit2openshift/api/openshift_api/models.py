import os
<<<<<<< HEAD
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
=======

from django.conf import settings
from django.db import models

from ansible_api.models import Project, Host, Group, Playbook
from ansible_api.models.mixins import AbstractProjectResourceModel, AbstractExecutionModel
from ansible_api.signals import pre_execution_start, post_execution_start


class Cluster(Project):
    STAGE1_NAME = "1-offline-package-prepare"
    STAGE2_NAME = "2-prerequisites-check"
    STAGE3_NAME = "3-deploy-cluster"

    class Meta:
        proxy = True

    def create_playbooks(self, version):
        playbooks_data = [
            {
                "name": self.STAGE1_NAME,
                "alias": "playbooks/offline_prepare.yml",
                "comment": "1-离线包准备"
            }, {
                "name": self.STAGE2_NAME,
                "alias": "playbooks/prerequisites.yml",
                "comment": "2-环境准备和依赖检查"
            }, {
                "name": self.STAGE3_NAME,
                "alias": "playbooks/deploy_cluster.yml",
                "comment": "3-部署集群"
            }
        ]
        for data in playbooks_data:
            self.playbook_set.create(
                name=data["name"], alias=data["alias"], type=Playbook.TYPE_GIT,
                git={"repo": os.path.join(settings.BASE_DIR, "data", "openshift-ansible"), "branch": version},
                comment=data["comment"]
            )

    def deploy(self):
        execution = DeployExecution.objects.create(project=self)
        return execution.start()

    def configs(self):
        self.change_to()
        return Role.osev3().vars

    def save(self, force_insert=False, force_update=False, using=None,
             update_fields=None):
        instance = super().save(force_insert=force_insert, force_update=force_update,
                                using=using, update_fields=update_fields)
        Role.create_internal_roles(self)
        Node.create_localhost()
        self.create_playbooks('release-3.10')
        return instance
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce


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

<<<<<<< HEAD
    def get_var(self, key, default):
        return self.vars.get(key, default)


class Role(Group):
=======
    @classmethod
    def create_localhost(cls):
        cls.objects.create(name="localhost", vars={"ansible_connection": "local"})

    def get_var(self, key, default):
        return self.vars.get(key, default)

    def get_node_group_label(self):
        return self.get_var("openshift_node_group_name", "-").split('-')[-1]


class Role(Group):
    NODE_GROUP_MASTER = "node-config-master"
    NODE_GROUP_INFRA = "node-config-infra"
    NODE_GROUP_COMPUTE = "node-config-compute"
    NODE_GROUP_STORAGE = "node-config-compute-storage"
    NODE_GROUP_MASTER_INFRA = "node-config-master-infra"
    NODE_GROUP_ALL_IN_ONE = "node-config-all-in-one"

    ROLE_MASTERS = "masters"
    ROLE_NODES = "nodes"
    ROLE_ETCD = "etcd"
    ROLE_INFRA = "infra"
    ROLE_COMPUTE = "compute"
    ROLE_LB = "lb"
    ROLE_NFS = "nfs"
    ROLE_OSEv3 = "OSEv3"

    ROLE_INTERNAL_NAMES = [
        {"name": ROLE_MASTERS, "comment": "主节点"},
        {"name": ROLE_COMPUTE, "comment": "计算节点"},
        {"name": ROLE_INFRA, "comment": "架构节点"},
        {"name": ROLE_ETCD, "comment": "ETCD节点", "children": (ROLE_MASTERS,)},
        {"name": ROLE_NODES, "comment": "节点", "children": (ROLE_MASTERS, ROLE_INFRA, ROLE_COMPUTE)},
        {
            "name": ROLE_OSEv3, "comment": "OSEv3",
            "children": (ROLE_MASTERS, ROLE_NODES, ROLE_ETCD, ROLE_LB, ROLE_ETCD),
            "vars": {
                "openshift_deployment_type": "origin",
                "openshift_master_identity_providers": [
                    {'name': 'htpasswd_auth', 'login': 'true', 'challenge': 'true', 'kind': 'HTPasswdPasswordIdentityProvider'}
                ],
                "openshift_disable_check": "disk_availability,docker_storage,memory_availability,docker_image_availability"
            }
        },
    ]

>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
    class Meta:
        proxy = True

    @property
    def nodes(self):
        return self.hosts

    @nodes.setter
    def nodes(self, value):
        self.hosts.set(value)

<<<<<<< HEAD
    def __str__(self):
        return "%s %s" % (self.project, self.name)


class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    project = models.ForeignKey('ansible_api.Project', on_delete=models.CASCADE)

    def start(self):
        result = {"raw": {}, "summary": {}}
        pre_deploy_execution_start.send(self.__class__, execution=self)
        for playbook in self.project.playbook_set.all().order_by('name'):
            print("\n>>> Start run {} ".format(playbook.name))
=======
    @classmethod
    def masters(cls):
        return cls.objects.get(name=cls.ROLE_MASTERS)

    @classmethod
    def infra(cls):
        return cls.objects.get(name=cls.ROLE_INFRA)

    @classmethod
    def osev3(cls):
        return cls.objects.get(name=cls.ROLE_OSEv3)

    @classmethod
    def compute(cls):
        return cls.objects.get(name=cls.ROLE_COMPUTE)

    @classmethod
    def create_internal_roles(cls, cluster):
        cluster.change_to()
        for r in cls.ROLE_INTERNAL_NAMES:
            role = cls.objects.create(
                name=r["name"], comment=r.get("comment", ""),
                vars=r.get("vars", {}), project=cluster
            )
            children_names = r.get("children")
            if children_names:
                children = cls.objects.filter(name__in=children_names)
                role.children.set(children)

    def on_nodes_join(self, nodes):
        self.__class__.update_node_group_labels()
        pass

    def on_nodes_leave(self, nodes):
        self.__class__.update_node_group_labels()
        pass

    @classmethod
    def update_node_group_labels(cls):
        if Node.objects.all().count() == 1:
            Node.objects.first().add_vars({"openshift_node_group_name": cls.NODE_GROUP_ALL_IN_ONE})
            return
        infra_role = cls.infra()
        masters_role = cls.masters()
        compute_role = cls.compute()
        if infra_role.nodes.count() == 0:
            tag_name = cls.NODE_GROUP_MASTER_INFRA
        else:
            tag_name = cls.NODE_GROUP_MASTER
        for node in Node.objects.filter(groups=masters_role):
            node.add_vars({"openshift_node_group_name": tag_name})
        for node in Node.objects.filter(groups=compute_role):
            node.add_vars({"openshift_node_group_name": cls.NODE_GROUP_COMPUTE})
        for node in Node.objects.filter(groups=infra_role):
            node.add_vars({"openshift_node_group_name": cls.NODE_GROUP_INFRA})


class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    project = models.ForeignKey('Cluster', on_delete=models.CASCADE)

    def start(self):
        result = {"raw": {}, "summary": {}}
        pre_execution_start.send(self.__class__, execution=self)
        for playbook in self.project.playbook_set.all():
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
            _result = playbook.execute()
            result["summary"].update(_result["summary"])
            if not _result.get('summary', {}).get('success', False):
                break
<<<<<<< HEAD
        post_deploy_execution_start.send(self.__class__, execution=self, result=result)
        return result


=======
        post_execution_start.send(self.__class__, execution=self, result=result)
        return result
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
