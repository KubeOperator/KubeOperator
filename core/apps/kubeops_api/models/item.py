import uuid

from django.db import models

__all__ = ["Item"]


class Item(models.Model):
    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128, null=False, blank=False, unique=True)
    users = models.ManyToManyField('users.Profile')
    description = models.CharField(max_length=256, null=True)
    date_created = models.DateTimeField(auto_now_add=True)


class ItemRole(models.Model):
    ITEM_ROLE_VIEWER = 'VIEWER'
    ITEM_ROLE_MANAGER = 'MANAGER'
    CLUSTER_STATUS_CHOICES = (
        (ITEM_ROLE_VIEWER, "viewer"),
        (ITEM_ROLE_MANAGER, "manager"),
    )
    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    role = models.CharField(max_length=128, null=False, blank=False, unique=True, choices=CLUSTER_STATUS_CHOICES)
    user = models.ForeignKey('users.Profile', on_delete=models.CASCADE)
    item = models.ForeignKey('Item', on_delete=models.CASCADE)
