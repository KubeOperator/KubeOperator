import os

from django.apps import AppConfig


class OpenshiftApiConfig(AppConfig):
    name = 'openshift_api'

    def ready(self):
        from openshift_api.models import Setting
        hostname = Setting.objects.filter(key='hostname').first()
        if hostname:
            os.putenv("REGISTORY_HOSTNAME", hostname.value)
