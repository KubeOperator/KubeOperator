import uuid

from django.db import models

class ClusterBackup(models.Model):

    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=255, null=False,unique=True, blank=False)
    size = models.IntegerField(default=0)
    date_created = models.DateTimeField(auto_now_add= True)
    project_id = models.CharField(max_length=256, null=False, blank=False)
    folder = models.CharField(max_length=64,blank=False)
    backup_storage_id = models.CharField(max_length=255,null=False,default="")
