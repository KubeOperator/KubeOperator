from django.apps import AppConfig


class CloudProviderConfig(AppConfig):
    name = 'cloud_provider'

    def ready(self):
        from . import signal_handlers
