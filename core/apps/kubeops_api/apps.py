from django.apps import AppConfig


class KubeOperatorApiConfig(AppConfig):
    name = 'kubeops_api'

    def ready(self):
        from . import signal_handlers
        from kubeops_api.models.setting import Setting
        Setting.apply_settings()
