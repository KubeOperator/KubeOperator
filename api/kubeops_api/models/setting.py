import uuid
from django.db import models

__all__ = ['Setting']


class Setting(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    key = models.CharField(max_length=128, blank=False)
    value = models.CharField(max_length=255, blank=True, default=None, null=True)
    name = models.CharField(max_length=128, blank=False)
    helper = models.CharField(max_length=255, blank=True)
    order = models.IntegerField(default=0)