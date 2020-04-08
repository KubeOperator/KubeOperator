from django.apps import AppConfig


class KubeOperatorApiConfig(AppConfig):
    name = 'kubeops_api'

    def ready(self):
        from . import signal_handlers
        from kubeops_api.models.setting import Setting
        try:
            Setting.subscribe_setting_change()
        except Exception:
            pass
