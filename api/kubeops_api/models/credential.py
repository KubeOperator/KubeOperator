import os
import uuid

import paramiko
from django.db import models
from django.utils.translation import ugettext_lazy as _
from common import models as common_models
from common.utils import ssh_key_string_to_obj
from fit2ansible import settings
from hashlib import md5

__all__ = ["Credential"]


class Credential(models.Model):
    CREDENTIAL_TYPE_PASSWORD = "password"
    CREDENTIAL_TYPE_PRIVATE_KEY = "privateKey"
    CREDENTIAL_TYPE_CHOICES = (
        (CREDENTIAL_TYPE_PASSWORD, "password"),
        (CREDENTIAL_TYPE_PRIVATE_KEY, "privateKey")
    )
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.SlugField(max_length=128, allow_unicode=True, unique=True, verbose_name=_('Name'))
    username = models.CharField(max_length=256, default='root')
    password = common_models.EncryptCharField(max_length=4096, blank=True, null=True)
    private_key = common_models.EncryptCharField(max_length=8192, blank=True, null=True)
    type = models.CharField(max_length=128, choices=CREDENTIAL_TYPE_CHOICES, default=CREDENTIAL_TYPE_PASSWORD)
    date_created = models.DateTimeField(auto_now_add=True)

    @property
    def private_key_obj(self):
        return ssh_key_string_to_obj(self.private_key, self.password)

    @property
    def private_key_path(self):
        if not self.type == 'privateKey':
            return None
        tmp_dir = os.path.join(settings.BASE_DIR, 'data', 'tmp')
        if not os.path.isdir(tmp_dir):
            os.makedirs(tmp_dir)
        key_name = '.' + md5(self.private_key.encode('utf-8')).hexdigest()
        key_path = os.path.join(tmp_dir, key_name)
        if not os.path.exists(key_path):
            self.private_key_obj.write_private_key_file(key_path)
            os.chmod(key_path, 0o400)
        return key_path
