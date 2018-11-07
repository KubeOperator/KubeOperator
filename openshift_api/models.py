import os

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

    def create_internal_roles(self):
        names = ['OSEv3', 'masters', 'nodes', 'etcd', "infra", "compute"]
        for name in names:
            Role.objects.get_or_create(name=name, project=self)

    def save(self, force_insert=False, force_update=False, using=None,
             update_fields=None):
        instance = super().save(force_insert=force_insert, force_update=force_update,
                                using=using, update_fields=update_fields)
        self.create_internal_roles()
        return instance


class Node(Host):
    class Meta:
        proxy = True

    @property
    def roles(self):
        return self.groups


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
    ROLE_OSEv3 = "OSEv3"

    ROLE_INTERNAL_NAMES = (
        (ROLE_MASTERS, "主节点"),
        (ROLE_COMPUTE, "计算节点"),
        (ROLE_INFRA, "架构节点"),
        (ROLE_ETCD, "ETCD节点"),
        (ROLE_NODES, "节点"),
        (ROLE_OSEv3, "OSEv3"),
    )

    class Meta:
        proxy = True

    @property
    def nodes(self):
        return self.hosts

    @nodes.setter
    def nodes(self, value):
        self.hosts.set(value)

    @classmethod
    def create_internal_roles(cls):
        for name, comment in cls.ROLE_INTERNAL_NAMES:
            pass
        for name in names:
            Role.objects.get_or_create(name=name, project=self)


    @classmethod
    def update_nodes_group_label(cls):
        pass



class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    project = models.ForeignKey('Cluster', on_delete=models.CASCADE)

    def start(self):
        result = {"raw": {}, "summary": {}}
        pre_execution_start.send(self.__class__, execution=self)
        for playbook in self.project.playbook_set.all():
            _result = playbook.execute()
            result["summary"].update(_result["result"])
            if not _result.get('summary', {}).get('success', False):
                break
        post_execution_start.send(self.__class__, execution=self, result=result)
        return result
