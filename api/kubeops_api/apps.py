

from django.apps import AppConfig


class OpenshiftApiConfig(AppConfig):
    name = 'kubeops_api'

    def ready(self):
        from . import signal_handlers
