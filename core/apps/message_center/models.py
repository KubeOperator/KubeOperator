import uuid
import json

from uuid import UUID
from django.db import models
from django.contrib.auth.models import User
from common import models as common_models
from datetime import date
from datetime import datetime


# Create your models here.

class Message(models.Model):
    MESSAGE_TYPE_SYSTEM = 'SYSTEM'
    MESSAGE_TYPE_CLUSTER = 'CLUSTER'

    MESSAGE_TYPE_CHOICES = (
        (MESSAGE_TYPE_SYSTEM, 'SYSTEM'),
        (MESSAGE_TYPE_CLUSTER, 'CLUSTER'),
    )

    MESSAGE_LEVEL_WARNING = 'WARNING'
    MESSAGE_LEVEL_ERROR = 'ERROR'
    MESSAGE_LEVEL_INFO = 'INFO'

    MESSAGE_LEVEL_CHOICES = (
        (MESSAGE_LEVEL_WARNING, 'WARNING'),
        (MESSAGE_LEVEL_ERROR, 'ERROR'),
        (MESSAGE_LEVEL_INFO, 'INFO'),
    )

    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    title = models.CharField(max_length=255, null=False, unique=False, blank=False)
    sender = models.CharField(max_length=255, null=False, unique=False, blank=False)
    content = models.TextField(max_length=65535)
    date_created = models.DateTimeField(auto_now_add=True)
    type = models.CharField(max_length=128, choices=MESSAGE_TYPE_CHOICES, default=MESSAGE_TYPE_SYSTEM, null=False)
    level = models.CharField(max_length=64, choices=MESSAGE_LEVEL_CHOICES, default=MESSAGE_LEVEL_INFO, null=False)


class UserMessage(models.Model):
    MESSAGE_SEND_TYPE_EMAIL = 'EMAIL'
    MESSAGE_SEND_TYPE_LOCAL = 'LOCAL'
    MESSAGE_SEND_TYPE_DINGTALK = 'DINGTALK'
    MESSAGE_SEND_TYPE_WORKWEIXIN = 'WORKWEIXIN'

    MESSAGE_SEND_TYPE_CHOICES = (
        (MESSAGE_SEND_TYPE_EMAIL, 'EMAIL'),
        (MESSAGE_SEND_TYPE_LOCAL, 'LOCAL'),
        (MESSAGE_SEND_TYPE_DINGTALK, 'DINGTALK'),
        (MESSAGE_SEND_TYPE_WORKWEIXIN, 'WORKWEIXIN')
    )

    MESSAGE_RECEIVE_STATUS_SUCCESS = 'SUCCESS'
    MESSAGE_RECEIVE_STATUS_WAITING = 'WAITING'
    MESSAGE_RECEIVE_STATUS_FAILED = 'FAILED'

    MESSAGE_RECEIVE_STATUS_CHOICES = (
        (MESSAGE_RECEIVE_STATUS_SUCCESS, 'SUCCESS'),
        (MESSAGE_RECEIVE_STATUS_WAITING, 'WAITING'),
        (MESSAGE_RECEIVE_STATUS_FAILED, 'FAILED'),
    )

    MESSAGE_READ_STATUS_READ = 'READ'
    MESSAGE_READ_STATUS_UNREAD = 'UNREAD'

    MESSAGE_READ_STATUS_CHOICES = (
        (MESSAGE_READ_STATUS_READ, 'READ'),
        (MESSAGE_READ_STATUS_UNREAD, 'UNREAD'),
    )

    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    receive = models.CharField(max_length=64, null=True, unique=False, blank=False)
    user_id = models.CharField(max_length=64, null=True, unique=False, blank=False)
    message = models.ForeignKey(Message, on_delete=models.CASCADE, to_field='id')
    send_type = models.CharField(max_length=64, choices=MESSAGE_SEND_TYPE_CHOICES, default=MESSAGE_SEND_TYPE_LOCAL,
                                 null=False)
    receive_status = models.CharField(max_length=64, choices=MESSAGE_RECEIVE_STATUS_CHOICES,
                                      default=MESSAGE_RECEIVE_STATUS_WAITING)
    read_status = models.CharField(max_length=64, choices=MESSAGE_READ_STATUS_CHOICES,
                                   default=MESSAGE_READ_STATUS_UNREAD)
    date_created = models.DateTimeField(auto_now_add=True)

    @property
    def message_detail(self):
        message = Message.objects.get(id=self.message_id)

        detail = {
            "title":str(message.title),
            "sender":message.sender,
            "content":str(message.content),
            "date_created":message.date_created.strftime("%Y-%m-%d %H:%M:%S"),
            "type":message.type,
            "level":message.level
        }

        return detail

    class Meta:
        ordering = ('-date_created',)

class UserNotificationConfig(models.Model):
    user = models.ForeignKey(User, on_delete=models.CASCADE, to_field='id')
    vars = common_models.JsonDictTextField(default={})
    type = models.CharField(max_length=128, choices=Message.MESSAGE_TYPE_CHOICES, default=Message.MESSAGE_TYPE_SYSTEM,
                            null=False)


class UserReceiver(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE, to_field='id')
    vars = common_models.JsonDictTextField(default={})

