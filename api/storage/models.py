import os
import threading

from django.db import models

# Create your models here.
from ansible_api.models import Project, Group, Playbook
from fit2ansible.settings import KUBEEASZ_DIR
from kubeops_api.models.node import Node
from kubeops_api.models.host import Host


class NfsStorage(Project):
    NFS_STATUS_CREATING = 'CREATING'
    NFS_STATUS_RUNNING = 'RUNNING'

    NFS_STATUS_CHOICES = (
        (NFS_STATUS_CREATING, 'CREATING'),
        (NFS_STATUS_RUNNING, 'RUNNING')
    )

    NFS_OPTION_NEW = 'NEW'
    NFS_OPTION_EXISTS = 'EXISTS'

    NFS_OPTION_CHOICES = (
        (NFS_OPTION_NEW, 'NEW'),
        (NFS_OPTION_EXISTS, 'EXISTS')
    )

    server = models.CharField(max_length=128, null=True)
    path = models.CharField(max_length=128, null=True)
    status = models.CharField(max_length=128, choices=NFS_STATUS_CHOICES, null=True)
    nfs_host = models.ForeignKey(Host, on_delete=models.SET_NULL, null=True)
    allow_ip = models.CharField(max_length=128, default='0.0.0.0/0')
    option = models.CharField(choices=NFS_OPTION_CHOICES, max_length=128, null=True)

    def create_group_node(self):
        if self.option == NfsStorage.NFS_OPTION_NEW:
            group = Group.objects.create(name='nfs')
            node = Node.objects.create(
                name=self.nfs_host.name,
                host=self.nfs_host,
                project=self
            )
            node.set_groups([group.name])
            self.server = node.host.ip

    def create_playbooks(self):
        url = 'file:///{}'.format(os.path.join(KUBEEASZ_DIR))
        Playbook.objects.create(
            name='nfs', alias="nfs.yml",
            type=Playbook.TYPE_LOCAL, url=url, project=self
        )

    def deploy_nfs(self):
        playbook = self.playbook_set.get(name="nfs")
        extra_vars = {}
        thread = threading.Thread(target=self.execute_playbook, args=(playbook, extra_vars))
        thread.start()

    def execute_playbook(self, playbook, extra_vars):
        self.status = NfsStorage.NFS_STATUS_CREATING
        _result = playbook.execute(extra_vars=extra_vars)
        if _result.get('summary', {}).get("success", False):
            self.status = NfsStorage.NFS_STATUS_RUNNING

    def on_nfs_save(self):
        self.change_to()
        self.create_group_node()
        self.create_playbooks()
        self.deploy_nfs()
