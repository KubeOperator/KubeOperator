import os
import threading
import uuid

import yaml
from django.db import models
from common.models import JsonTextField
from django.utils.translation import ugettext_lazy as _
from kubeoperator.settings import PACKAGE_DIR
from kubeops_api.package_manage import *

logger = logging.getLogger('kubeops')
__all__ = ['Package']


class Package(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    meta = JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    packages_dir = PACKAGE_DIR

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = _('Package')

    @property
    def path(self):
        return os.path.join(self.packages_dir, self.name)

    @property
    def repo_port(self):
        return self.meta['vars']['repo_port']

    @property
    def registry_port(self):
        return self.meta['vars']['registry_port']

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
            instance, _ = cls.objects.update_or_create(defaults=defaults, name=d)
            thread = threading.Thread(target=cls.start_container(instance))
            thread.start()

    @classmethod
    def start_container(cls, package):
        if not is_package_container_exists(package.name):
            create_package_container(package)
            return
        if not is_package_container_start(package.name):
            start_package_container(package)
