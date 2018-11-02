# -*- coding: utf-8 -*-
#

import uuid
import datetime
import os

from django.db import models
from django.utils.translation import ugettext_lazy as _

from common import models as common_models
from ..ansible.runner import AdHocRunner
from ..signals import pre_adhoc_exec, post_adhoc_exec
from .mixins import AbstractProjectResourceModel
from .utils import format_result_as_list


__all__ = ['AdHoc', 'AdHocExecution']


class AdHoc(AbstractProjectResourceModel):
    pattern = models.CharField(max_length=1024, default='all', verbose_name=_('Pattern'))
    module = models.CharField(max_length=128, default='command', verbose_name=_("Module"))
    args = common_models.JsonTextField(verbose_name=_("Args"))
    created_by = models.CharField(max_length=128, blank=True, null=True, default='')

    COMMAND_MODULES = ('shell', 'command', 'raw')

    def __str__(self):
        return "{}: {}".format(self.module, self.clean_args)

    @property
    def clean_args(self):
        if self.module in self.COMMAND_MODULES and isinstance(self.args, str):
            if self.args.startswith('executable='):
                _args = self.args.split(' ')
                executable, command = _args[0].split('=')[1], ' '.join(_args[1:])
                args = {'executable': executable, '_raw_params':  command}
            else:
                args = {'_raw_params':  self.args}
            return args
        else:
            return self.args

    @property
    def tasks(self):
        return [{
            "name": self.__str__(),
            "action": {
                "module": self.module,
                "args": self.clean_args,
            },
        }]

    @property
    def inventory(self):
        return self.project.inventory_obj

    def execute(self, save_history=True):
        result = {"raw": {}, "summary": {}}
        try:
            pre_adhoc_exec.send(self.__class__, adhoc=self, save_history=save_history)
            runner = AdHocRunner(self.inventory, options=self.project.cleaned_options)
            result = runner.run(self.tasks, pattern=self.pattern)
        except Exception as e:
            result['summary'] = {'error': str(e)}
        finally:
            post_adhoc_exec.send(self.__class__, adhoc=self, save_history=save_history, result=result)
        return format_result_as_list(result.get('summary', {}))

    @staticmethod
    def test_tasks():
        return [
            {
                "name": "Test ping",
                "actions": [{
                    "module": "ping",
                    "args": ""
                }],
            },
        ]


class AdHocExecution(AbstractProjectResourceModel):
    id = models.CharField(primary_key=True, default=uuid.uuid4, max_length=36)
    adhoc = models.ForeignKey('AdHoc', on_delete=models.SET_NULL, null=True)
    is_finished = models.BooleanField(default=False, verbose_name=_('Is finished'))
    is_success = models.BooleanField(default=False, verbose_name=_('Is success'))
    timedelta = models.FloatField(default=0.0, verbose_name=_('Time'), null=True)
    raw = common_models.JsonDictTextField(blank=True, null=True, default='{}', verbose_name=_('Adhoc raw result'))
    summary = common_models.JsonDictTextField(blank=True, null=True, default='{}', verbose_name=_('Adhoc summary'))
    date_start = models.DateTimeField(auto_now_add=True, verbose_name=_('Start time'))
    date_finished = models.DateTimeField(blank=True, null=True, verbose_name=_('End time'))

    @property
    def summary_as_list(self):
        summary = {"contacted": [], "dark": []}
        for status, result in self.summary.items():
            for hostname, _tasks in result.items():
                tasks = []
                for task_name, detail in _tasks.items():
                    detail["task"] = task_name
                    tasks.append(detail)
                summary[status].append({"hostname": hostname, "tasks": tasks})
        return summary

    @property
    def stdout(self):
        with open(self.log_path, 'r') as f:
            data = f.read()
        return data

    @property
    def success_hosts(self):
        return self.summary.get('contacted', []) if self.summary else []

    @property
    def failed_hosts(self):
        return self.summary.get('dark', {}) if self.summary else []

    @property
    def log_path(self):
        dt = datetime.datetime.now().strftime('%Y-%m-%d')
        log_dir = os.path.join(self.project.adhoc_dir, dt)
        if not os.path.exists(log_dir):
            os.makedirs(log_dir)
        return os.path.join(log_dir, str(self.id) + '.log')