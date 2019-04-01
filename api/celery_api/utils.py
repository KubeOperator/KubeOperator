# -*- coding: utf-8 -*-
#
import os
import json
from functools import wraps
import sys

from django.conf import settings
from django.db.utils import ProgrammingError, OperationalError
from django.core.cache import cache
from django_celery_beat.models import PeriodicTask, IntervalSchedule, CrontabSchedule

from .contrib import NoStripLoggingProxy

INTERVAL_UNIT_MAP = {
    's': IntervalSchedule.SECONDS,
    'm': IntervalSchedule.MINUTES,
    'h': IntervalSchedule.HOURS,
    'd': IntervalSchedule.DAYS,
}


def add_register_period_task(name):
    key = "__REGISTER_PERIODIC_TASKS"
    value = cache.get(key, [])
    value.append(name)
    cache.set(key, value)


def get_register_period_tasks():
    key = "__REGISTER_PERIODIC_TASKS"
    return cache.get(key, [])


def add_after_app_shutdown_clean_task(name):
    key = "__AFTER_APP_SHUTDOWN_CLEAN_TASKS"
    value = cache.get(key, [])
    value.append(name)
    cache.set(key, value)


def get_after_app_shutdown_clean_tasks():
    key = "__AFTER_APP_SHUTDOWN_CLEAN_TASKS"
    return cache.get(key, [])


def add_after_app_ready_task(name):
    key = "__AFTER_APP_READY_RUN_TASKS"
    value = cache.get(key, [])
    value.append(name)
    cache.set(key, value)


def get_after_app_ready_tasks():
    key = "__AFTER_APP_READY_RUN_TASKS"
    return cache.get(key, [])


def create_or_update_periodic_task(tasks, pk=None):
    """
    :param tasks: {
        'add-every-monday-morning': {
            'task': 'tasks.add' # A registered celery task,
            'interval': 30,
            'crontab': "30 7 * * *",
            'args': (16, 16),
            'kwargs': {},
            'enabled': False,
        },
    }
    :return:
    """
    # Todo: check task valid, task and callback must be a celery task
    for name, detail in tasks.items():
        interval = None
        crontab = None
        try:
            IntervalSchedule.objects.all().count()
        except (ProgrammingError, OperationalError):
            return None

        _interval = detail.get("interval")
        if isinstance(_interval, int):
            _interval = '{}s'.format(_interval)

        if isinstance(_interval, str) and _interval and _interval[:-1].isdigit():
            val = int(_interval[:-1])
            unit = INTERVAL_UNIT_MAP.get(_interval[-1].lower())
            if not unit:
                raise SyntaxError("Schedule is not valid: {}".format(_interval))
            intervals = IntervalSchedule.objects.filter(
                every=val, period=unit
            )
            if intervals:
                interval = intervals[0]
            else:
                interval = IntervalSchedule.objects.create(
                    every=val,
                    period=unit,
                )
        elif isinstance(detail.get("crontab"), str):
            try:
                minute, hour, day, month, week = detail["crontab"].split()
            except ValueError:
                raise SyntaxError("crontab is not valid")
            kwargs = dict(
                minute=minute, hour=hour, day_of_week=week,
                day_of_month=day, month_of_year=month,
            )
            contabs = CrontabSchedule.objects.filter(
                **kwargs
            )
            if contabs:
                crontab = contabs[0]
            else:
                crontab = CrontabSchedule.objects.create(**kwargs)
        else:
            raise SyntaxError("Schedule is not valid: {}".format(_interval))

        defaults = dict(
            interval=interval,
            crontab=crontab,
            name=name,
            task=detail['task'],
            args=json.dumps(detail.get('args', [])),
            kwargs=json.dumps(detail.get('kwargs', {})),
            enabled=detail.get('enabled', True),
        )

        task = PeriodicTask.objects.update_or_create(
            defaults=defaults, id=pk,
        )
        return task


def disable_celery_periodic_task(period_task_name):
    from django_celery_beat.models import PeriodicTask
    PeriodicTask.objects.filter(name=period_task_name).update(enabled=False)


def delete_celery_periodic_task(period_task_name):
    from django_celery_beat.models import PeriodicTask
    PeriodicTask.objects.filter(name=period_task_name).delete()


def register_as_period_task(crontab=None, interval=None):
    """
    Warning: Task must be have not any args and kwargs
    :param crontab:  "* * * * *"
    :param interval:  60*60*60
    :return:
    """
    if crontab is None and interval is None:
        raise SyntaxError("Must set crontab or interval one")

    def decorate(func):
        if crontab is None and interval is None:
            raise SyntaxError("Interval and crontab must set one")

        # Because when this decorator run, the task was not created,
        # So we can't use func.name
        name = '{func.__module__}.{func.__name__}'.format(func=func)
        if name not in get_register_period_tasks():
            create_or_update_periodic_task({
                name: {
                    'task': name,
                    'interval': interval,
                    'crontab': crontab,
                    'args': (),
                    'enabled': True,
                }
            })
            add_register_period_task(name)

        @wraps(func)
        def wrapper(*args, **kwargs):
            return func(*args, **kwargs)
        return wrapper
    return decorate


def after_app_ready_start(func):
    # Because when this decorator run, the task was not created,
    # So we can't use func.name
    name = '{func.__module__}.{func.__name__}'.format(func=func)
    if name not in get_after_app_ready_tasks():
        add_after_app_ready_task(name)

    @wraps(func)
    def decorate(*args, **kwargs):
        return func(*args, **kwargs)
    return decorate


def after_app_shutdown_clean(func):
    # Because when this decorator run, the task was not created,
    # So we can't use func.name
    name = '{func.__module__}.{func.__name__}'.format(func=func)
    if name not in get_after_app_shutdown_clean_tasks():
        add_after_app_shutdown_clean_task(name)

    @wraps(func)
    def decorate(*args, **kwargs):
        return func(*args, **kwargs)
    return decorate


def get_celery_task_log_path(task_id):
    task_id = str(task_id)
    rel_path = os.path.join(task_id[0], task_id[1], task_id + '.log')
    path = os.path.join(settings.CELERY_LOG_DIR, rel_path)
    os.makedirs(os.path.dirname(path), exist_ok=True)
    return path

