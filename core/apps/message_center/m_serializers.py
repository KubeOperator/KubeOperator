#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/17 
=================================================='''

from rest_framework import serializers
from message_center.models import UserNotificationConfig, UserReceiver, UserMessage, Message

__all__ = ["UserNotificationConfigSerializer","MessageSerializer"]


class UserNotificationConfigSerializer(serializers.ModelSerializer):
    vars = serializers.DictField(required=False)

    class Meta:
        model = UserNotificationConfig
        fields = [
            'id', 'user_id', 'type', 'vars'
        ]


class UserReceiverSerializer(serializers.ModelSerializer):
    vars = serializers.DictField(required=False)

    class Meta:
        model = UserReceiver
        fields = [
            'id', 'vars', 'user_id'
        ]



class UserMessageSerializer(serializers.ModelSerializer):
    message = serializers.SlugRelatedField(
        queryset=Message.objects.all(),
        slug_field='id', required=False
    )
    message_detail  = serializers.DictField(required=False)

    class Meta:
        model = UserMessage
        fields = [
            'id', 'receive', 'user_id', 'send_type', 'receive_status', 'read_status', 'date_created', 'message_id',
            'message', 'message_detail'
        ]

class MessageSerializer(serializers.ModelSerializer):
    class Meta:
        model = Message
        fields = [
            'id', 'title', 'sender', 'contend', 'type', 'read_status', 'date_created', 'level',
        ]


