#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/16 
=================================================='''
from django.contrib.auth.models import User
from kubeops_api.models.item import Item
from .models import Message, UserNotificationConfig, UserReceiver, UserMessage


class MessageClient():

    def __init__(self):
        pass

    def get_receivers(self, item_id):
        receivers = []
        admin = User.objects.filter(is_superuser=1)
        receivers.append(admin)

        if item_id is not None:
            item = Item.objects.get(id=item_id)
            profiles = item.profiles
            for profile in profiles:
                receivers.append(profile.user)

        return receivers

    def split_receiver_by_send_type(self, receivers, type):
        messageReceivers = []
        for receiver in receivers:
            config = UserNotificationConfig.objects.get(type=type, user_id=receiver.id)
            user_receiver = UserReceiver.objects.get(user_id=receiver.id)
            if config.vars['LOCAL'] == 'ENABLE':
                messageReceivers.append(MessageReceiver(user=receiver, receive=receiver.username, send_type='LOCAL'))
            if config.vars['EMAIL'] == 'ENABLE' and user_receiver.vars['EMAIL'] != '':
                messageReceivers.append(
                    MessageReceiver(user=receiver, receive=user_receiver.vars['EMAIL'], send_type='EMAIL'))
            if config.vars['DINGTALK'] == 'ENABLE' and user_receiver.vars['DINGTALK'] != '':
                messageReceivers.append(
                    MessageReceiver(user=receiver, receive=user_receiver.vars['DINGTALK'], send_type='DINGTALK'))
            if config.vars['WORKWEIXIN'] == 'ENABLE' and user_receiver.vars['WORKWEIXIN'] != '':
                messageReceivers.append(
                    MessageReceiver(user=receiver, receive=user_receiver.vars['WORKWEIXIN'], send_type='WORKWEIXIN'))
        return messageReceivers

    def insert_message(self, title, content, type, level, item_id):
        message = Message.objects.create(title=title, content=content, type=type, level=level)
        message_receivers = self.split_receiver_by_send_type(receivers=self.get_receivers(item_id), type=type)
        user_messages = []
        for message_receiver in message_receivers:
            user_message = UserMessage(receive=message_receiver.receive, user_id=message_receiver.user_id,
                                       send_type=message_receiver.send_type,
                                       read_status=UserMessage.MESSAGE_READ_STATUS_UNREAD,
                                       receive_status=UserMessage.MESSAGE_RECEIVE_STATUS_SUCCESS, message_id=message.id)
            user_messages.append(user_message)
        return UserMessage.objects.bulk_create(user_messages)


class MessageReceiver():

    def __init__(self, user, receive, send_type):
        self.user_id = user.id
        self.receive = receive
        self.send_type = send_type
