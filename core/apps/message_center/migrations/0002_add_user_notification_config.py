#!/usr/bin/env python
# -*- coding: UTF-8 -*-
'''=================================================
@Author ：zk.wang
@Date   ：2020/3/17 
=================================================='''

from django.db import migrations


class Migration(migrations.Migration):
    dependencies = [
        ('message_center', '0001_initial')
    ]

    def forwards_func(apps, schema_editor):
        User = apps.get_model('auth', 'user')
        UserNotificationConfig = apps.get_model('message_center', 'UserNotificationConfig')
        db_alias = schema_editor.connection.alias
        users = User.objects.using(db_alias).all()
        user_notification_configs = []
        vars = {
            "LOCAL": "ENABLE",
            "EMAIL": "DISABLE",
            "DINGTALK": "DISABLE",
            "WORKWEIXIN": "DISABLE",
        }
        for user in users:
            user_notification_configs.append(UserNotificationConfig(vars=vars, user=user, type='CLUSTER'))
            user_notification_configs.append(UserNotificationConfig(vars=vars, user=user, type='SYSTEM'))

        UserNotificationConfig.objects.using(db_alias).bulk_create(user_notification_configs)

    operations = [
        migrations.RunPython(forwards_func),
    ]
