from django.db import models

# Create your models here.
from ansible_api.models import Project


class NfsStorage(Project):
    NFS_STATUS_CREATING = 'CREATING'
    NFS_STATUS_RUNNING = 'RUNNING'

    NFS_STATUS_CHOICES = (
        (NFS_STATUS_CREATING, 'CREATING'),
        (NFS_STATUS_RUNNING, 'RUNNING')
    )
    server = models.CharField(max_length=128, null=True)
    path = models.CharField(max_length=128, null=True)
    status = models.CharField(max_length=128, choices=NFS_STATUS_CHOICES, null=True)
