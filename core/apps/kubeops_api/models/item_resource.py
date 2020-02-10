import uuid

from django.db import models

__all__ = ["ItemResource"]


class ItemResource(models.Model):
    RESOURCE_TYPE_CLUSTER = 'CLUSTER'
    RESOURCE_TYPE_HOST = 'HOST'
    RESOURCE_TYPE_PLAN = 'PLAN'
    RESOURCE_TYPE_BACKUP_STORAGE = 'BACKUP_STORAGE'
    RESOURCE_TYPE_STORAGE = 'STORAGE'

    RESOURCE_TYPE_CHOICES  = (
        (RESOURCE_TYPE_CLUSTER,'CLUSTER'),
        (RESOURCE_TYPE_HOST,'HOST'),
        (RESOURCE_TYPE_PLAN,'PLAN'),
        (RESOURCE_TYPE_BACKUP_STORAGE,'BACKUP_STORAGE'),
        (RESOURCE_TYPE_STORAGE,'STORAGE')
    )

    item_id = models.UUIDField(max_length=255, default=uuid.uuid4)
    resource_id = models.UUIDField(max_length=255, default=uuid.uuid4)
    resource_type = models.CharField(max_length=64,choices=RESOURCE_TYPE_CHOICES)

