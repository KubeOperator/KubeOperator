#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/20 
=================================================='''
import datetime
import json
from uuid import UUID

class DateEncoder(json.JSONEncoder):

    def default(self, obj):
        if isinstance(obj, UUID):
            return str(obj)
        if isinstance(obj,datetime.datetime):
            return obj.strftime("%Y-%m-%d %H:%M:%S")
        else:
            return json.JSONEncoder.default(self,obj)
