from django.contrib.auth.models import AbstractUser
from django.db import models


class User(AbstractUser):
    current_item = models.ForeignKey('kubeops_api.Item', on_delete=models.SET_NULL, null=True,
                                     related_name='current_item')

    @property
    def items(self):
        items = []
        for item in self.item_set.all():
            items.append(item.name)
        return items
