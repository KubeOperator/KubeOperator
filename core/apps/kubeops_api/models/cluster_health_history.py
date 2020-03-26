import uuid

from django.db import models


class ClusterHealthHistory(models.Model):
    CLUSTER_HEALTH_HISTORY_DATE_TYPE_HOUR = 'HOUR'
    CLUSTER_HEALTH_HISTORY_DATE_TYPE_DAY = 'DAY'

    CLUSTER_HEALTH_HISTORY_DATE_TYPE_CHOICES = (
        (CLUSTER_HEALTH_HISTORY_DATE_TYPE_HOUR, 'HOUR'),
        (CLUSTER_HEALTH_HISTORY_DATE_TYPE_DAY, 'DAY')
    )

    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    project_id = models.CharField(max_length=255, null=False, blank=False)
    available_rate = models.IntegerField(default=0)
    date_type = models.CharField(max_length=255, choices=CLUSTER_HEALTH_HISTORY_DATE_TYPE_CHOICES,
                                 default=CLUSTER_HEALTH_HISTORY_DATE_TYPE_HOUR)
    month = models.CharField(max_length=255, null=True, blank=True)
    date_created = models.DateTimeField(auto_now_add=True)
