import uuid

from django.conf import settings
from django.db import models

__all__ = ['Setting']


class Setting(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    tab = models.CharField(max_length=128, default="system")
    key = models.CharField(max_length=128, blank=False)
    value = models.CharField(max_length=255, blank=True, default=None, null=True)

    @classmethod
    def get_db_settings(cls):
        sts = cls.objects.all()
        result = {}
        for s in sts:
            result[s.key] = s.value
        return result

    @classmethod
    def get_settings(cls, tab=None):
        sts = cls.objects.all()
        result = {}
        for setting in sts:
            if tab and not setting.tab == tab:
                continue
            r = None
            if hasattr(settings, setting.key):
                r = settings.__getattr__(setting.key)
            result[setting.key] = r
        return result

    @classmethod
    def set_settings(cls, sts, tab=None):
        for k, v in sts.items():
            defaults = {"key": k, "value": v}
            if tab:
                defaults.update({"tab": tab})
            cls.objects.update_or_create(defaults=defaults, key=k)
            Setting.apply_settings()

    @classmethod
    def apply_settings(cls):
        try:
            sts = cls.objects.all()
            for setting in sts:
                settings.__setattr__(setting.key, setting.value)
        except Exception as e:
            pass
