import uuid

from django.contrib.auth.models import User
from django.db import models

__all__ = ["Profile"]


class Profile(models.Model):
    USER_SOURCE_LOCAL = "local"
    USER_SOURCE_LDAP = "ldap"
    USER_SOURCE_CHOICES = (
        (USER_SOURCE_LDAP, 'ldap'),
        (USER_SOURCE_LOCAL, 'local'),
    )
    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='profile')
    source = models.CharField(max_length=128, choices=USER_SOURCE_CHOICES, default=USER_SOURCE_LOCAL)

    @property
    def items(self):
        return self.item_set.all()

    @property
    def item_role_mappings(self):
        return self.itemrolemapping_set.all()
