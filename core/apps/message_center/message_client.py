#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/16 
=================================================='''
import json
import logging
import kubeoperator.settings
import redis

from django.contrib.auth.models import User
from kubeops_api.models.item import Item
from .models import Message, UserNotificationConfig, UserReceiver, UserMessage
from kubeops_api.models.setting import Setting
from ko_notification_utils.email_smtp import Email
from ko_notification_utils.ding_talk import DingTalk
from ko_notification_utils.work_weixin import WorkWeiXin
from .message_thread import MessageThread
from django.template import Template, Context, loader

logger = logging.getLogger('kubeops')


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
        setting_email_enable = False
        email_receivers = ''
        if len(Setting.objects.filter(key='SMTP_STATUS')) > 0 and Setting.objects.get(key='SMTP_STATUS').value == 'ENABLE':
            setting_email_enable = True
        send_ding_talk_enable = False
        ding_talk_receivers = ''
        if len(Setting.objects.filter(key='DINGTALK_STATUS')) > 0 and Setting.objects.get(key='DINGTALK_STATUS').value == 'ENABLE':
            send_ding_talk_enable = True
        send_weixin_enable = False
        weixin_receivers = ''
        if len(Setting.objects.filter(key='WEIXIN_STATUS')) > 0 and Setting.objects.get(key='WEIXIN_STATUS').value == 'ENABLE':
            send_weixin_enable = True

        for receiver in receivers:
            config = UserNotificationConfig.objects.get(type=type, user_id=receiver.id)
            user_receiver = UserReceiver.objects.get(user_id=receiver.id)
            if config.vars['LOCAL'] == 'ENABLE':
                messageReceivers.append(
                    MessageReceiver(user_id=receiver.id, receive=receiver.username, send_type='LOCAL'))

            if setting_email_enable and config.vars['EMAIL'] == 'ENABLE' and user_receiver.vars['EMAIL'] != '':
                if email_receivers != '':
                    email_receivers = email_receivers + ',' + receiver.email
                else:
                    email_receivers = receiver.email

            if send_ding_talk_enable and config.vars['DINGTALK'] == 'ENABLE' and user_receiver.vars['DINGTALK'] != '':
                if ding_talk_receivers != '':
                    ding_talk_receivers = ding_talk_receivers + ',' + user_receiver.vars['DINGTALK']
                else:
                    ding_talk_receivers = user_receiver.vars['DINGTALK']

            if send_weixin_enable and config.vars['WORKWEIXIN'] == 'ENABLE' and user_receiver.vars['WORKWEIXIN'] != '':
                if weixin_receivers != '':
                    weixin_receivers = weixin_receivers + '|' + user_receiver.vars['WORKWEIXIN']
                else:
                    weixin_receivers = user_receiver.vars['WORKWEIXIN']

        if len(email_receivers) > 0:
            messageReceivers.append(
                MessageReceiver(user_id=1, receive=email_receivers, send_type='EMAIL'))

        if len(ding_talk_receivers) > 0:
            messageReceivers.append(
                MessageReceiver(user_id=1, receive=ding_talk_receivers, send_type='DINGTALK')
            )
        if len(weixin_receivers) > 0:
            messageReceivers.append(
                MessageReceiver(user_id=1, receive=weixin_receivers, send_type='WORKWEIXIN')
            )

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
        email_messages = UserMessage.objects.filter(message_id=message.id,
                                                    send_type=UserMessage.MESSAGE_SEND_TYPE_EMAIL,
                                                    user_id=1)
        if len(email_messages) > 0:
            thread = MessageThread(func=send_email, user_message=email_messages[0])
            thread.start()
        ding_talk_messages = UserMessage.objects.filter(message_id=message.id,
                                                        send_type=UserMessage.MESSAGE_SEND_TYPE_DINGTALK,
                                                        user_id=1)
        if len(ding_talk_messages) > 0:
            thread2 = MessageThread(func=send_ding_talk_msg, user_message=ding_talk_messages[0])
            thread2.start()

        work_weixin_messages = UserMessage.objects.filter(message_id=message.id,
                                                          send_type=UserMessage.MESSAGE_SEND_TYPE_WORKWEIXIN,
                                                          user_id=1)
        if len(work_weixin_messages) > 0:
            thread3 = MessageThread(func=send_work_weixin_msg, user_message=work_weixin_messages[0])
            thread3.start()


def send_email(user_message):
    setting_email = Setting.get_settings("email")
    email = Email(address=setting_email['SMTP_ADDRESS'], port=setting_email['SMTP_PORT'],
                  username=setting_email['SMTP_USERNAME'], password=setting_email['SMTP_PASSWORD'])
    res = email.send_html_mail(receiver=user_message.receive, title=user_message.message.title,
                               content=get_email_content(user_message))
    if res.success:
        user_message.receive_status = UserMessage.MESSAGE_RECEIVE_STATUS_SUCCESS
        user_message.save()
    else:
        logger.error(msg="send email error message_id=" + str(user_message.message_id) + "reason:" + str(res.data),
                     exc_info=True)


def get_email_content(user_message):
    content = json.loads(user_message.message.content)
    try:
        template = loader.get_template(get_email_template(content['resource_type']))
        content['detail'] = json.loads(content['detail'])
        content['title'] = user_message.message.title
        content['date'] = user_message.message.date_created.strftime("%Y-%m-%d %H:%M:%S")
        email_content = template.render(content)
        return email_content
    except Exception as e:
        logger.error(msg="get email content error", exc_info=True)
        return ''


def get_email_template(type):
    templates = {
        "CLUSTER": "cluster.html",
        "CLUSTER_EVENT": "cluster-event.html",
        "CLUSTER_USAGE": "cluster-usage.html",
    }
    return templates[type]


def send_ding_talk_msg(user_message):
    setting_dingTalk = Setting.get_settings("dingTalk")
    ding_talk = DingTalk(webhook=setting_dingTalk['DINGTALK_WEBHOOK'], secret=setting_dingTalk['DINGTALK_SECRET'])

    text = get_msg_content(user_message)
    content = {"title": user_message.message.title, "text": text}

    res = ding_talk.send_markdown_msg(receivers=user_message.receive.split(','), content=content)
    if res.success:
        print("send ding talk success")
        user_message.receive_status = UserMessage.MESSAGE_RECEIVE_STATUS_SUCCESS
        user_message.save()
    else:
        logger.error(msg="send dingtalk error message_id=" + str(user_message.message_id) + "reason:" + str(res.data),
                     exc_info=True)


def send_work_weixin_msg(user_message):
    workWeixin = Setting.get_settings("workWeixin")
    weixin = WorkWeiXin(corp_id=workWeixin['WEIXIN_CORP_ID'], corp_secret=workWeixin['WEIXIN_CORP_SECRET'],
                        agent_id=workWeixin['WEIXIN_AGENT_ID'])
    text = get_msg_content(user_message)
    content = {'content': text}
    token = get_work_weixin_token()

    res = weixin.send_markdown_msg(receivers=user_message.receive, content=content, token=token)
    if res.success:
        print("send work weixin success")
        user_message.receive_status = UserMessage.MESSAGE_RECEIVE_STATUS_SUCCESS
        user_message.save()
    else:
        logger.error(msg="send workweixin error message_id=" + str(user_message.message_id) + "reason:" + str(res.data),
                     exc_info=True)


def get_msg_content(user_message):
    content = json.loads(user_message.message.content)
    type = content['resource_type']
    content['detail'] = json.loads(content['detail'])
    text = ''
    if type == 'CLUSTER_EVENT':
        text = "### " + user_message.message.title + " \n\n " + \
               "> **项目**:" + content['item_name'] + " \n\n " + \
               "> **集群**:" + content['resource_name'] + " \n\n" + \
               "> **名称**:" + content['detail']['name'] + " \n\n " + \
               "> **类别**:" + content['detail']['type'] + " \n\n " + \
               "> **原因**:" + content['detail']['reason'] + " \n\n " + \
               "> **组件**:" + content['detail']['component'] + " \n\n " + \
               "> **NameSpace**:" + content['detail']['namespace'] + " \n\n " + \
               "> **主机**:" + content['detail']['host'] + " \n\n " + \
               "> **告警时间**:" + content['detail']['last_timestamp'] + " \n\n " + \
               "> **详情**:" + content['detail']['message'] + " \n\n " + \
               "<font color=\"info\">本消息由KubeOperator自动发送</font>"

    if type == 'CLUSTER':
        text = "### " + user_message.message.title + "\n\n" + \
               "> **项目**:" + content['item_name'] + "\n\n" + \
               "> **集群**:" + content['resource_name'] + "\n\n" + \
               "> **信息**:" + content['detail']['message'] + "\n\n" + \
               "<font color=\"info\">本消息由KubeOperator自动发送</font>"

    if type == 'CLUSTER_USAGE':
        text = "### " + user_message.message.title + "\n\n" + \
               "> **项目**:" + content['item_name'] + "\n\n" + \
               "> **集群**:" + content['resource_name'] + "\n\n" + \
               "> **详情**:" + content['detail']['message'] + "\n\n" + \
               "<font color=\"info\">本消息由KubeOperator自动发送</font>"
    return text


def get_work_weixin_token():
    redis_cli = redis.StrictRedis(host=kubeoperator.settings.REDIS_HOST, port=kubeoperator.settings.REDIS_PORT)
    if redis_cli.exists('WORK_WEIXIN_TOKEN'):
        return redis_cli.get('WORK_WEIXIN_TOKEN')
    else:
        workWeixin = Setting.get_settings("workWeixin")
        weixin = WorkWeiXin(corp_id=workWeixin['WEIXIN_CORP_ID'], corp_secret=workWeixin['WEIXIN_CORP_SECRET'],
                            agent_id=workWeixin['WEIXIN_AGENT_ID'])
        result = weixin.get_token()
        redis_cli.set('WORK_WEIXIN_TOKEN', result.data['access_token'], result.data['expires_in'])
        return result.data['access_token']


class MessageReceiver():

    def __init__(self, user_id, receive, send_type):
        self.user_id = user_id
        self.receive = receive
        self.send_type = send_type
