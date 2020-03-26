#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/23 
=================================================='''
import threading


class EmailThread(threading.Thread):

    def __init__(self, func, message_id):
        threading.Thread.__init__(self)
        self.func = func
        self.message_id = message_id

    def run(self):
        self.func(self.message_id)
