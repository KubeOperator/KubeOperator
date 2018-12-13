from django.apps import AppConfig


class AnsibleApiConfig(AppConfig):
    name = 'ansible_api'

    def ready(self):
        import ansible_api.signal_handlers
