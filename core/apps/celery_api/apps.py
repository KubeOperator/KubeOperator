from django.apps import AppConfig


class CeleryApiConfig(AppConfig):
    name = 'celery_api'

    def ready(self):
        from . import signal_handler
