import uuid

from django.db import models
from common import models as common_models
from django.utils.translation import ugettext_lazy as _

__all__ = ['Package']


class Package(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    endpoint = models.CharField(max_length=255, null=True, default='')
    meta = common_models.JsonDictTextField(default={})
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))

    class Meta:
        verbose_name = _('Package')
