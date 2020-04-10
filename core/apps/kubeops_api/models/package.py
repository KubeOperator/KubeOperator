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
    PACKAGE_STATE_OFFLINE = "offline"
    PACKAGE_STATE_ONLINE = "online"
    PACKAGE_STATE_CHOICES = (
        (PACKAGE_STATE_OFFLINE, 'offline'),
        (PACKAGE_STATE_ONLINE, 'online'),
    )

    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    meta = JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    state = models.CharField(max_length=32, default=PACKAGE_STATE_OFFLINE, choices=PACKAGE_STATE_CHOICES)
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
            cls.check_package_health()
            thread = threading.Thread(target=cls.start_container(instance))
            thread.daemon = True
            thread.start()

    @classmethod
    def start_container(cls, package):
        cls.check_package_health()
        if not is_package_container_exists(package.name):
            create_package_container(package)
        if not is_package_container_start(package.name):
            start_package_container(package)

    @classmethod
    def check_package_health(cls):
        ps = cls.objects.all()
        cs = list_package_containers()
        for p in ps:
            if p not in cs:
                p.state = Package.PACKAGE_STATE_OFFLINE
            else:
                p.state = Package.PACKAGE_STATE_ONLINE
            p.save()
