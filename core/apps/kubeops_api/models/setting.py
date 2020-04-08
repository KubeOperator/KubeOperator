import logging
import threading
import uuid
from django.conf import settings
from django.db import models

__all__ = ['Setting']

logger = logging.getLogger(__name__)

from kubeops_api.utils.redis import RedisHelper


class Setting(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    tab = models.CharField(max_length=128, default="system")
    key = models.CharField(max_length=128, blank=False)
    value = models.CharField(max_length=255, blank=True, default=None, null=True)
    __helper = RedisHelper()
    __channel_name = "setting_change"

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
            cls.__helper.publish(cls.__channel_name, "1")
            logger.debug("send a message to channel")

    @classmethod
    def __apply_settings(cls):
        sts = Setting.objects.all()
        for setting in sts:
            settings.__setattr__(setting.key, setting.value)

    @classmethod
    def subscribe_setting_change(cls):
        cls.__apply_settings()
        sub = cls.__helper.subscribe(cls.__channel_name)

        def listen(s):
            while True:
                _ = s.parse_response()
                cls.__apply_settings()

        t = threading.Thread(target=listen, args=(sub,))
        t.daemon = True
        t.start()
