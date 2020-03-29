#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/13 
=================================================='''
import json
import redis
import kubeoperator.settings

from rest_framework.views import APIView
from ko_notification_utils.email_smtp import Email
from ko_notification_utils.work_weixin import WorkWeiXin
from rest_framework.response import Response
from rest_framework import status
from rest_framework.viewsets import ModelViewSet
from .m_serializers import UserNotificationConfigSerializer, UserReceiverSerializer, UserMessageSerializer
from .models import UserNotificationConfig, UserReceiver, UserMessage, Message
from django.core.paginator import Paginator, EmptyPage, PageNotAnInteger
from django.http import JsonResponse
from django.template import Template, Context ,loader
from message_center.message_client import send_email

class EmailCheckView(APIView):

    def post(self, request, *args, **kwargs):
        email_config = request.data
        email = Email(address=email_config['SMTP_ADDRESS'], port=email_config['SMTP_PORT'],
                      username=email_config['SMTP_USERNAME'], password=email_config['SMTP_PASSWORD'])

        result = email.login()
        if result.success:
            return Response(data={'msg': '校验成功！'}, status=status.HTTP_200_OK)
        else:
            return Response(data={'msg': '校验失败！'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)


class WorkWeixinCheckView(APIView):

    def post(self, request, *args, **kwargs):
        weixin_config = request.data
        weixin = WorkWeiXin(corp_id=weixin_config['WEIXIN_CORP_ID'], corp_secret=weixin_config['WEIXIN_CORP_SECRET'],
                            agent_id=weixin_config['WEIXIN_AGENT_ID'])
        result = weixin.get_token()
        if result.success:
            redis_cli = redis.StrictRedis(host=kubeoperator.settings.REDIS_HOST,
                                               port=kubeoperator.settings.REDIS_PORT)
            redis_cli.set('WORK_WEIXIN_TOKEN',result.data['access_token'],result.data['expires_in'])

            return Response(data={'msg': '校验成功！'}, status=status.HTTP_200_OK)
        else:
            return Response(data={'msg': '校验失败！' + json.dumps(result.data)},
                            status=status.HTTP_500_INTERNAL_SERVER_ERROR)


class SubscribeViewSet(ModelViewSet):
    serializer_class = UserNotificationConfigSerializer
    queryset = UserNotificationConfig.objects.all()

    http_method_names = ['post', 'get', 'head', 'options']

    lookup_field = 'id'
    lookup_url_kwarg = 'id'

    def list(self, request, *args, **kwargs):
        user = request.user
        self.queryset = UserNotificationConfig.objects.filter(user_id=user.id)
        return super().list(self, request, *args, **kwargs)

    def post(self, request, *args, **kwargs):
        config = UserNotificationConfig.objects.get(id=kwargs['id'])
        config.vars = request.data['vars']
        config.save()
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        headers = self.get_success_headers(serializer.data)
        return Response(serializer.data, status=status.HTTP_201_CREATED, headers=headers)


class UserReceiverViewSet(ModelViewSet):
    serializer_class = UserReceiverSerializer

    lookup_field = 'id'
    lookup_url_kwarg = 'id'

    def list(self, request, *args, **kwargs):
        user = request.user
        self.queryset = UserReceiver.objects.filter(user_id=user.id)
        return super().list(self, request, *args, **kwargs)

    def post(self, request, *args, **kwargs):
        receiver = UserReceiver.objects.get(id=kwargs['id'])
        receiver.vars = request.data['vars']
        receiver.save()
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        headers = self.get_success_headers(serializer.data)
        return Response(serializer.data, status=status.HTTP_201_CREATED, headers=headers)


class UserMessageView(ModelViewSet):
    serializer_class = UserMessageSerializer

    lookup_field = 'id'
    lookup_url_kwarg = 'id'

    def list(self, request, *args, **kwargs):
        user = request.user
        limit = request.query_params.get('limit')
        page = request.query_params.get('page')
        type = request.query_params.get('type')
        level = request.query_params.get('level')
        readStatus = request.query_params.get('readStatus')
        user_messages = UserMessage.objects.filter(user_id=user.id,
                                                   send_type=UserMessage.MESSAGE_SEND_TYPE_LOCAL).values()
        if readStatus != 'ALL':
            user_messages = user_messages.filter(read_status=readStatus)
        if type != 'ALL':
            ids = Message.objects.filter(type=type).values_list('id')
            user_messages = user_messages.filter(message_id__in=ids)
        if level != 'ALL':
            l_ids = Message.objects.filter(level=level).values_list('id')
            user_messages = user_messages.filter(message_id__in=l_ids)
        for user_message in user_messages:
            user_message['message_detail'] = Message.objects.filter(id=user_message['message_id']).values()[0]
        paginator = Paginator(user_messages, limit)
        try:
            user_messages = paginator.page(page)
        except PageNotAnInteger:
            user_messages = paginator.page(1)
        except EmptyPage:
            user_messages = paginator.page(paginator.num_pages)

        return JsonResponse(data={"total": paginator.count,
                                  "page_num": paginator.num_pages,
                                  "data": list(user_messages.object_list)})

    def post(self, request, *args, **kwargs):
        if kwargs['id'] == 'ALL':
            user = request.user
            UserMessage.objects.filter(user_id=user.id, send_type=UserMessage.MESSAGE_SEND_TYPE_LOCAL).update(
                read_status=UserMessage.MESSAGE_READ_STATUS_READ)
            return Response({"msg": "更新成功！"}, status=status.HTTP_200_OK)
        else:
            user_msg = UserMessage.objects.get(id=kwargs['id'])
            user_msg.read_status = UserMessage.MESSAGE_READ_STATUS_READ
            user_msg.save()
            serializer = self.get_serializer(user_msg, data=request.data)
            serializer.is_valid(raise_exception=True)
            headers = self.get_success_headers(serializer.data)
            return Response(serializer.data, status=status.HTTP_200_OK, headers=headers)


class UnReadMessageView(APIView):

    def get(self, request, *args, **kwargs):
        user = request.user
        info_ids = Message.objects.filter(level=Message.MESSAGE_LEVEL_INFO).values_list('id')
        info_message = UserMessage.objects.filter(user_id=user.id, send_type=UserMessage.MESSAGE_SEND_TYPE_LOCAL,
                                                  read_status=UserMessage.MESSAGE_READ_STATUS_UNREAD).filter(
            message_id__in=info_ids)
        warn_ids = Message.objects.filter(level=Message.MESSAGE_LEVEL_WARNING).values_list('id')
        warn_message = UserMessage.objects.filter(user_id=user.id, send_type=UserMessage.MESSAGE_SEND_TYPE_LOCAL,
                                                  read_status=UserMessage.MESSAGE_READ_STATUS_UNREAD).filter(
            message_id__in=warn_ids)
        result = {
            "info": len(info_message),
            "warning": len(warn_message)
        }
        return Response(result)
