#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/16 
=================================================='''
from django.contrib.auth.models import User
from kubeops_api.models.item import Item


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

    def split_receiver_by_send_type(self, receivers):
        pass
