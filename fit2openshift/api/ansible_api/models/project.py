# -*- coding: utf-8 -*-
#
import os
import uuid

from django.utils.translation import ugettext_lazy as _
from django.db import models
from django.conf import settings

from common import models as common_models
from .inventory import Inventory
from ..ansible import BaseInventory
from ..ctx import set_current_project


__all__ = ['Project']


class Project(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.SlugField(max_length=128, allow_unicode=True, unique=True, verbose_name=_('Name'))
    # Run full_options, ex: forks,
    options = common_models.JsonCharField(max_length=1024, blank=True, null=True, verbose_name=_('Run options'))
    comment = models.CharField(max_length=128, blank=True, null=True, verbose_name=_("Comment"))
    meta = common_models.JsonDictTextField(blank=True, null=True)
    created_by = models.CharField(max_length=128, blank=True, null=True, default='')
    date_created = models.DateTimeField(auto_now_add=True)

    __root_id = '00000000-0000-0000-0000-000000000000'
    __public_id = '00000000-0000-0000-0000-000000000001'

    def __str__(self):
        return self.name

    @property
    def inventory(self):
        return Inventory(self.host_set.all(), self.group_set.all())

    @property
    def inventory_file_path(self):
        return os.path.join(self.project_dir, 'hosts.yaml')

    def refresh_inventory_file(self):
        with open(self.inventory_file_path, 'w') as f:
            f.write(self.inventory.get_data(fmt='yaml'))

    @property
    def roles_dir(self):
        roles_dir = os.path.join(self.project_dir, 'roles')
        if not os.path.isdir(roles_dir):
            os.makedirs(roles_dir, exist_ok=True)
        return roles_dir

    @property
    def project_dir(self):
        project_dir = os.path.join(settings.ANSIBLE_PROJECTS_DIR, self.name)
        if not os.path.isdir(project_dir):
            os.makedirs(project_dir, exist_ok=True)
        return project_dir

    @property
    def playbooks_dir(self):
        playbooks_dir = os.path.join(self.project_dir, 'playbooks')
        if not os.path.isdir(playbooks_dir):
            os.makedirs(playbooks_dir, exist_ok=True)
        return playbooks_dir

    @property
    def adhoc_dir(self):
        adhoc_dir = os.path.join(self.project_dir, 'adhoc')
        if not os.path.isdir(adhoc_dir):
            os.makedirs(adhoc_dir, exist_ok=True)
        return adhoc_dir

    @classmethod
    def root_project(cls):
        return cls(id=cls.__root_id, name='ROOT')

    @classmethod
    def public_project(cls):
        return cls(id=cls.__public_id, name='Public')

    def is_real(self):
        return self.id not in [self.__root_id, self.__public_id]

    @property
    def inventory_obj(self):
        return self.inventory.as_object()

    def get_inventory_data(self):
        return self.inventory.get_data(fmt='py')

    def change_to(self):
        set_current_project(self)

    def clear_inventory(self):
        self.group_set.all().delete()
        self.host_set.all().delete()

    @property
    def cleaned_options(self):
        options = self.options or {}
        options['roles_path'] = [self.roles_dir]
        return options

    @staticmethod
    def test_inventory():
        data = {
            "hosts": [
                {
                    "hostname": "192.168.244.128",
                    "vars": {
                        "ansible_ssh_user": "root",
                        "ansible_ssh_pass": "redhat"
                    }
                },
                {
                    "hostname": "gaga",
                    "vars": {
                        "ansible_ssh_host": "192.168.1.1"
                    }
                }
            ],
            "groups": [
                {"name": "apache", "hosts": ["gaga"]},
                {"name": "web", "hosts": ["192.168.244.128"],
                 "vars": {"hello": "world"}, "children": ["apache"]},
            ]
        }
        return data

    @classmethod
    def get_test_inventory(cls):
        return BaseInventory(cls.test_inventory())
