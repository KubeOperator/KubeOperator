#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/4/20 
=================================================='''
import logging
import uuid

from django.db import models
from django.utils.translation import ugettext_lazy as _

from common.models import JsonDictTextField, JsonListTextField

__all__ = ["CisLog"]


class CisLog(models.Model):
    CIS_STATUS_FAILED = 'FAILED'
    CIS_STATUS_SUCCESS = 'SUCCESS'

    CIS_STATUS_CHOICES = (
        (CIS_STATUS_FAILED, 'FAILED'),
        (CIS_STATUS_SUCCESS, 'SUCCESS')
    )

    id = models.UUIDField(max_length=255, primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=126, null=False, unique=True)
    cluster_id = models.UUIDField(max_length=255, default=uuid.uuid4)
    detail = JsonListTextField(default={})
    result = JsonDictTextField(default={})
    status = models.CharField(max_length=64, choices=CIS_STATUS_CHOICES, default=CIS_STATUS_FAILED)
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))

    class Meta:
        ordering = ('-date_created',)
