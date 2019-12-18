from django.apps import AppConfig


class KubeOperatorApiConfig(AppConfig):
    name = 'kubeops_api'

    def ready(self):
        from . import signal_handlers
