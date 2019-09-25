import uuid


from django.db import models


class BackupStrategy(models.Model):

    id = models.UUIDField(primary_key=True,default=uuid.uuid4)
    cron = models.IntegerField(max_length=64,null=False)
    save_num = models.IntegerField(max_length=64,null=False)
    backup_storage = models.ForeignKey('BackupStorage',on_delete=models.SET_NULL,null=True)
    # project = models.ForeignKey('ansible_api.project',on_delete=models.CASCADE)
    project_id = models.CharField(max_length=64)

