import logging
import os
import uuid

import yaml
from django.db import models
from common import models as common_models
from fit2ansible import settings
from django.utils.translation import ugettext_lazy as _
from kubeops_api.adhoc import storage_health_check
from kubeops_api.models.host import Host

logger = logging.getLogger(__name__)


class StorageTemplate(models.Model):
    name = models.CharField(max_length=128, unique=True, verbose_name='名称')
    meta = common_models.JsonDictTextField(blank=True, null=True)
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    templates_dir = os.path.join(settings.BASE_DIR, 'resource', 'storage')

    def __str__(self):
        return self.name

    @property
    def path(self):
        return os.path.join(self.templates_dir, self.name)

    @classmethod
    def lookup(cls):
        for d in os.listdir(cls.templates_dir):
            full_path = os.path.join(cls.templates_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            cls.objects.update_or_create(defaults=defaults, name=d)


class Storage(models.Model):
    STORAGE_STATUS_INVALID = 'invalid'
    STORAGE_STATUS_VALID = 'valid'
    STORAGE_STATUS_UNKNOWN = 'unknown'

    STORATE_STATUS_CHOICES = (
        (STORAGE_STATUS_INVALID, 'invalid'),
        (STORAGE_STATUS_VALID, 'valid'),
        (STORAGE_STATUS_UNKNOWN, 'unknown')

    )
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128, unique=True)
    template = models.ForeignKey("StorageTemplate", null=True, on_delete=models.SET_NULL)
    vars = common_models.JsonDictTextField(default={}, blank=True, null=True, verbose_name=_('Vars'))
    status = models.CharField(max_length=128, choices=STORATE_STATUS_CHOICES, default=STORAGE_STATUS_UNKNOWN)
    date_created = models.DateTimeField(auto_now_add=True)
    comment = models.CharField(max_length=128, blank=True, null=True, verbose_name=_("Comment"))

    def health_check(self):
        meta = self.template.meta.get('health_check')
        module = meta.get('module')
        command = meta.get('command')
        real_command = self.replace_vars(command)
        host = Host(name='localhost', vars={"ansible_connection": "local"}, )
        logger.info('execute command: ' + real_command)
        if storage_health_check(host, module, real_command):
            self.status = Storage.STORAGE_STATUS_VALID
        else:
            self.status = Storage.STORAGE_STATUS_INVALID
        self.save()

    def replace_vars(self, command):
        for k, v in self.vars.items():
            if k in command:
                command = command.replace('$' + k, v)
        return command
