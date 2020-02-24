import uuid

from django.contrib.auth.models import User
from django.db import models

__all__ = ["Profile"]


class Profile(models.Model):
    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='profile')
    current_item = models.ForeignKey('kubeops_api.Item', related_name='user_item', null=True,on_delete=models.SET_NULL)
