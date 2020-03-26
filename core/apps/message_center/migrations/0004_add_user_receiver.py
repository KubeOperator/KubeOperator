#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/18 
=================================================='''
from django.db import migrations


class Migration(migrations.Migration):
    dependencies = [
        ('message_center', '0003_auto_20200318_0306')
    ]

    def forwards_func(apps, schema_editor):
        User = apps.get_model('auth', 'user')
        UserReceiver = apps.get_model('message_center', 'UserReceiver')
        db_alias = schema_editor.connection.alias
        users = User.objects.using(db_alias).all()
        user_receivers = []

        for user in users:
            vars = {
                "EMAIL": user.email,
                "DINGTALK": "",
                "WORKWEIXIN": "",
            }
            user_receivers.append(UserReceiver(vars=vars, user=user))

        UserReceiver.objects.using(db_alias).bulk_create(user_receivers)

    operations = [
        migrations.RunPython(forwards_func),
    ]
