import uuid

from django.db import models
from django.utils.translation import ugettext_lazy as _
from common import models as common_models


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
