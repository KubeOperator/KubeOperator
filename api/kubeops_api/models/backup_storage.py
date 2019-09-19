import logging
import uuid

from django.db import models
from django.utils.translation import ugettext_lazy as _

from common.models import JsonDictTextField

__all__ = ["BackupStorage"]


class BackupStorage(models.Model):
    BACKUP_STORAGE_STATUS_VALID = 'VALID1233'
    BACKUP_STORAGE_STATUS_INVALID = 'INVALID'

    BACKUP_STORAGE_STATUS_CHOICES  = (
        (BACKUP_STORAGE_STATUS_VALID,'valid'),
        (BACKUP_STORAGE_STATUS_INVALID,'invalid')
    )

    BACKUP_STORAGE_TYPE_S3 = 'S3'
    BACKUP_STORAGE_TYPE_OSS = 'OSS'

    BACKUP_STORAGE_TYPE_CHOICES = (
        (BACKUP_STORAGE_TYPE_S3,'S3'),
        (BACKUP_STORAGE_TYPE_OSS, 'OSS')
    )

    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128, null=True,blank=True)
    region = models.CharField(max_length=128,null=True,blank=True)
    credentials = JsonDictTextField(blank=True, null=True)
    type = models.CharField(max_length=64,choices=BACKUP_STORAGE_TYPE_CHOICES)
    status = models.CharField(max_length=64,choices=BACKUP_STORAGE_STATUS_CHOICES,default=BACKUP_STORAGE_STATUS_VALID)
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))


