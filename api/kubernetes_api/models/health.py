from django.db import models

from ansible_api.models import Project
from ansible_api.models.mixins import AbstractProjectResourceModel
from common import models as common_models
from kubernetes import client, config, watch


class HealthCheck(models.Model, AbstractProjectResourceModel):
    msg = models.TextField(null=True, blank=True, default='No message')
    health = models.BooleanField(default=True)
    date_created = models.DateTimeField(auto_now_add=True)
    meta = common_models.JsonDictTextField(blank=True, null=True)

