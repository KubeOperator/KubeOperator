import os
import sys

from django.apps import AppConfig
from django.db.backends.signals import connection_created
from django.dispatch import receiver


class OpenshiftApiConfig(AppConfig):
    name = 'openshift_api'

    @receiver(connection_created, dispatch_uid="my_unique_identifier")
    def on_db_connection_ready(sender, **kwargs):
        from .signals import django_ready
        from openshift_api.models import Setting
        hostname = Setting.objects.filter(key='hostname').first()
        if hostname:
            os.putenv("REGISTORY_HOSTNAME", hostname.value)
            django_ready.send(OpenshiftApiConfig)

    def ready(self):
        from . import signal_handlers
