from django.db import models
from ansible_api.models import Host as Ansible_Host


class Node(Ansible_Host):
    host = models.ForeignKey('Host', related_name='host', default=None, null=True, on_delete=models.CASCADE)

    @property
    def roles(self):
        return self.groups

    @property
    def host_memory(self):
        return self.host.info.memory

    @property
    def host_cpu_core(self):
        return self.host.info.cpu_core

    @property
    def host_os(self):
        return self.host.info.os

    @property
    def host_os_version(self):
        return self.host.info.os_version

    @property
    def host_volumes(self):
        return self.host.info.volumes

    @roles.setter
    def roles(self, value):
        self.groups.set(value)

    def on_node_save(self):
        self.ip = self.host.ip
        self.username = self.host.username
        self.password = self.host.password
        self.private_key = self.host.private_key
        self.host.node_id = self.id
        self.host.save()
        self.save()

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

