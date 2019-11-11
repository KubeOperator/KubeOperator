import uuid
from django.db import models

__all__ = ['Setting']


class Setting(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    key = models.CharField(max_length=128, blank=False)
    value = models.CharField(max_length=255, blank=True, default=None, null=True)

    @classmethod
    def get_setting(cls, key):
        setting = cls.objects.get(key=key)
        return setting

    @classmethod
    def get_settings(cls):
        settings = cls.objects.all()
        result = {}
        for setting in settings:
            result[setting.key] = setting.value
        return result

    @classmethod
    def set_settings(cls, settings):
        for k, v in settings.items():
            defaults = {"key": k, "value": v}
            cls.objects.update_or_create(defaults=defaults, key=k)
