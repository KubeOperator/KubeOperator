#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/23 
=================================================='''
import threading

class MessageThread(threading.Thread):

    def __init__(self, func, user_message):
        threading.Thread.__init__(self)
        self.func = func
        self.user_message = user_message

    def run(self):
        self.func(self.user_message)
