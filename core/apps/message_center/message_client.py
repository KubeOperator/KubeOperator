#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/16 
=================================================='''
import json
from django.contrib.auth.models import User
from kubeops_api.models.item import Item
from .models import Message, UserNotificationConfig, UserReceiver, UserMessage
from kubeops_api.models.setting import Setting
from ko_notification_utils.email_smtp import Email
from .message_thread import EmailThread


class MessageClient():

    def __init__(self):
        pass

    def get_receivers(self, item_id):
        receivers = []
        admin = User.objects.filter(is_superuser=1)
        receivers.extend(list(admin))

        if item_id is not None:
            item = Item.objects.get(id=item_id)
            profiles = item.profiles.all()
            for profile in profiles:
                receivers.append(profile.user)

        return receivers

    def split_receiver_by_send_type(self, receivers, type):
        messageReceivers = []
        setting_email_enable = Setting.objects.get(key='SMTP_STATUS').value = 'ENABLE'
        email_receivers = ''
        for receiver in receivers:
            config = UserNotificationConfig.objects.get(type=type, user_id=receiver.id)
            user_receiver = UserReceiver.objects.get(user_id=receiver.id)
            if config.vars['LOCAL'] == 'ENABLE':
                messageReceivers.append(MessageReceiver(user_id=receiver.id, receive=receiver.username, send_type='LOCAL'))

            if config.vars['DINGTALK'] == 'ENABLE' and user_receiver.vars['DINGTALK'] != '':
                messageReceivers.append(
                    MessageReceiver(user_id=receiver.id, receive=user_receiver.vars['DINGTALK'], send_type='DINGTALK'))
            if config.vars['WORKWEIXIN'] == 'ENABLE' and user_receiver.vars['WORKWEIXIN'] != '':
                messageReceivers.append(
                    MessageReceiver(user_id=receiver.id, receive=user_receiver.vars['WORKWEIXIN'], send_type='WORKWEIXIN'))
            if setting_email_enable and config.vars['EMAIL'] == 'ENABLE' and user_receiver.vars['EMAIL'] != '':
                email_receivers = email_receivers + ',' + receiver.email

        if len(email_receivers) > 0:
            email_receivers[0] = ''
            messageReceivers.append(
                MessageReceiver(user_id=1, receive=email_receivers, send_type='EMAIL'))

        return messageReceivers

    def insert_message(self, message):
        title = message.get('title', None)
        item_id = message.get('item_id', None)
        content = message.get('content', None)
        type = message.get('type', None)
        level = message.get('level', None)
        message = Message.objects.create(title=title, content=json.dumps(content), type=type, level=level)
        message_receivers = self.split_receiver_by_send_type(receivers=self.get_receivers(item_id), type=type)
        user_messages = []
        for message_receiver in message_receivers:
            user_message = UserMessage(receive=message_receiver.receive, user_id=message_receiver.user_id,
                                       send_type=message_receiver.send_type,
                                       read_status=UserMessage.MESSAGE_READ_STATUS_UNREAD,
                                       receive_status=UserMessage.MESSAGE_RECEIVE_STATUS_WAITING, message_id=message.id)
            user_messages.append(user_message)
        UserMessage.objects.bulk_create(user_messages)
        thread = EmailThread(func=MessageClient.send_email,message_id=message.id)
        thread.start()


    def send_email(self, message_id):
        user_message = UserMessage.objects.get(message_id=message_id, send_type=UserMessage.MESSAGE_SEND_TYPE_EMAIL)
        setting_email = Setting.get_settings("email")
        email = Email(address=setting_email['SMTP_ADDRESS'], port=setting_email['SMTP_PORT'],
                      username=setting_email['SMTP_USERNAME'], password=setting_email['SMTP_PASSWORD'])
        res = email.send_message(receiver=user_message.receive,title=user_message.message.title,content=user_message.message.content)
        if res.success:
            user_message.receive_status = UserMessage.MESSAGE_RECEIVE_STATUS_SUCCESS
            user_message.save()


class MessageReceiver():

    def __init__(self, user_id, receive, send_type):
        self.user_id = user_id
        self.receive = receive
        self.send_type = send_type
