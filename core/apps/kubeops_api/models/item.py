import uuid

from django.db import models

__all__ = ["Item"]


class Item(models.Model):
    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128, null=False, blank=False, unique=True)
    users = models.ManyToManyField('users.Profile')
    description = models.CharField(max_length=256, null=True)
    date_created = models.DateTimeField(auto_now_add=True)
