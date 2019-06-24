import os

import yaml
from django.db import models

from common.models import JsonTextField
from fit2ansible import settings
from django.utils.translation import ugettext_lazy as _

__all__ = ['AuthTemplate']


class AuthTemplate(models.Model):
    name = models.CharField(max_length=128, verbose_name='名称')
    meta = JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    auth_dir = os.path.join(settings.BASE_DIR, 'resource', 'auth')

    def __str__(self):
        return self.name

    @property
    def path(self):
        return os.path.join(self.auth_dir, self.name)

    @classmethod
    def lookup(cls):
        for d in os.listdir(cls.auth_dir):
            full_path = os.path.join(cls.auth_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            cls.objects.update_or_create(defaults=defaults, name=d)
