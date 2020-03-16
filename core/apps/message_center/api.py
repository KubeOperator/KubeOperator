#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/13 
=================================================='''
import json

from rest_framework.views import APIView
from ko_notification_utils.email_smtp import Email
from ko_notification_utils.work_weixin import WorkWinXin
from rest_framework.response import Response
from rest_framework import viewsets, status


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
        weixin = WorkWinXin(corp_id=weixin_config['WEIXIN_CORP_ID'], corp_secret=weixin_config['WEIXIN_CORP_SECRET'],
                            agent_id=weixin_config['WEIXIN_AGENT_ID'])
        result = weixin.get_token()
        if result.success:
            return Response(data={'msg': '校验成功！'}, status=status.HTTP_200_OK)
        else:
            return Response(data={'msg': '校验失败！' + json.dumps(result.data)},
                            status=status.HTTP_500_INTERNAL_SERVER_ERROR)
