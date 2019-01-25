# -*- coding: utf-8 -*-

import os

from celery import Celery, platforms

# set the default Django settings module for the 'celery' program.
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'fit2ansible.settings')
os.environ.setdefault('PYTHONOPTIMIZE', '1')
os.environ.setdefault('C_FORCE_ROOT', '1')

from django.conf import settings


app = Celery('celery_api')
platforms.C_FORCE_ROOT = True

# Using a string here means the worker will not have to
# pickle the object when using Windows.
app.config_from_object('django.conf:settings', namespace='CELERY')
app.autodiscover_tasks(lambda: [app_config.split('.')[0] for app_config in settings.INSTALLED_APPS])


@app.task
def add(x, y):
    print('This is a \033[1;35m test \033[0m!')
    print('\033[1;33;44mThis is a test !\033[0m')
    return x + y
