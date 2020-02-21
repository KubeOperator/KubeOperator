# -*- coding: utf-8 -*-
#
from django.contrib.auth.hashers import make_password
from django.db import migrations, models


def add_default_admin(apps, schema_editor):
    user_model = apps.get_model("users", "User")
    db_alias = schema_editor.connection.alias
    user_model.objects.using(db_alias).create(
        username="admin",
        email="admin@mycomany.com",
        password=make_password("kubeoperator@admin123"),
        is_superuser=True,
        is_staff=True
    )


class Migration(migrations.Migration):
    initial = True

    dependencies = [
        ('users', '0001_initial'),
    ]

    operations = [
        migrations.RunPython(add_default_admin),
    ]
