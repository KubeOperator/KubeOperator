import logging
import os
import threading

from django.db import models
from common import models as common_models

# Create your models here.
from ansible_api.models import Project, Group, Playbook
from fit2ansible.settings import KUBEEASZ_DIR
from kubeops_api.models.node import Node
from kubeops_api.models.host import Host

Logger = logging.getLogger(__name__)


class NfsStorage(Project):
    NFS_STATUS_CREATING = 'CREATING'
    NFS_STATUS_RUNNING = 'RUNNING'
    NFS_STATUS_ERROR = 'ERROR'

    NFS_STATUS_CHOICES = (
        (NFS_STATUS_CREATING, 'CREATING'),
        (NFS_STATUS_RUNNING, 'RUNNING'),
        (NFS_STATUS_ERROR, 'ERROR')
    )

    NFS_OPTION_NEW = 'NEW'
    NFS_OPTION_EXISTS = 'EXISTS'

    NFS_OPTION_CHOICES = (
        (NFS_OPTION_NEW, 'NEW'),
        (NFS_OPTION_EXISTS, 'EXISTS')
    )
    status = models.CharField(max_length=128, choices=NFS_STATUS_CHOICES, default=NFS_STATUS_RUNNING, null=True)
    vars = common_models.JsonDictTextField(default={"allow_ip": "0.0.0.0/0", "storage_nfs_server_path": "/exports"})

    def create_group_node(self):
        host = Host.objects.get(name=self.vars['host'])
        node = Node.objects.create(
            name=host.name,
            host=host,
            project=self
        )
        node.set_groups(['nfs'])
        self.vars['storage_nfs_server'] = node.host.ip
        self.save()

    def create_playbooks(self):
        self.change_to()
        url = 'file:///{}'.format(os.path.join(KUBEEASZ_DIR))
        Playbook.objects.create(
            name='nfs', alias="nfs.yml",
            type=Playbook.TYPE_LOCAL, url=url, project=self
        )

    def deploy_nfs(self):
        playbook = self.playbook_set.get(name="nfs")
        Logger.info('开始部署 NFS 服务')
        thread = threading.Thread(target=self.execute_playbook, args=(playbook, self.vars))
        thread.start()

    def execute_playbook(self, playbook, extra_vars):
        self.change_to()
        self.status = NfsStorage.NFS_STATUS_CREATING
        _result = playbook.execute(extra_vars=extra_vars)
        if _result.get('summary', {}).get("success", False):
            self.status = NfsStorage.NFS_STATUS_RUNNING
        else:
            self.status = NfsStorage.NFS_STATUS_ERROR
        self.save()

    def on_nfs_save(self):
        if self.vars['option'].upper() == NfsStorage.NFS_OPTION_NEW:
            self.change_to()
            self.create_group_node()
            self.create_playbooks()
            self.deploy_nfs()
