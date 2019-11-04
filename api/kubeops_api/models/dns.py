import uuid

from django.db import models
from django.utils.translation import ugettext_lazy as _


class DNS(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    dns1 = models.CharField(max_length=128, null=True)
    dns2 = models.CharField(max_length=128, null=True)
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))

    class Meta:
        get_latest_by = 'date_created'
        ordering = ('-date_created',)
