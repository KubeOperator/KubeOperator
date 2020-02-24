import uuid

from django.db import models

__all__ = ["Item", "ItemRoleMapping"]


class Item(models.Model):
    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128, null=False, blank=False, unique=True)
    profiles = models.ManyToManyField('users.Profile')
    description = models.CharField(max_length=256, null=True)
    date_created = models.DateTimeField(auto_now_add=True)


class ItemRoleMapping(models.Model):
    ITEM_ROLE_VIEWER = 'VIEWER'
    ITEM_ROLE_MANAGER = 'MANAGER'
    ITEM_ROLE_CHOICES = (
        (ITEM_ROLE_VIEWER, "VIEWER"),
        (ITEM_ROLE_MANAGER, "MANAGER"),
    )
    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    role = models.CharField(max_length=128, null=False, blank=False, choices=ITEM_ROLE_CHOICES,
                            default=ITEM_ROLE_VIEWER)
    profile = models.ForeignKey('users.Profile', on_delete=models.CASCADE)
    item = models.ForeignKey('Item', on_delete=models.CASCADE)

    @property
    def item_name(self):
        return self.item.name
