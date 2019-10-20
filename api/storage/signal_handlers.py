import json
import os

from django.db.models.signals import post_save
from django.dispatch import receiver

from storage.models import NfsStorage


@receiver(post_save, sender=NfsStorage)
def on_cluster_save(sender, instance=None, created=True, **kwargs):
    if created and instance:
        instance.on_nfs_save()
