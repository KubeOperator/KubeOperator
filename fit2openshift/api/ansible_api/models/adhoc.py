# -*- coding: utf-8 -*-
#

from celery import current_task
from django.db import models
from django.utils.translation import ugettext_lazy as _

from common import models as common_models
from ..ansible.runner import AdHocRunner
from ..signals import pre_execution_start, post_execution_start
from .mixins import AbstractProjectResourceModel, AbstractExecutionModel


__all__ = ['AdHoc', 'AdHocExecution']


class AdHoc(AbstractProjectResourceModel):
    pattern = models.CharField(max_length=1024, default='all', verbose_name=_('Pattern'))
    module = models.CharField(max_length=128, default='command', verbose_name=_("Module"))
    args = common_models.JsonTextField(verbose_name=_("Args"))

    execute_times = models.IntegerField(default=0)
    created_by = models.CharField(max_length=128, blank=True, null=True, default='')

    def __str__(self):
        return "{}: {}".format(self.module, self.args)

    def execute(self):
        pk = None
        if current_task:
            pk = current_task.request.id
        execution = AdHocExecution(adhoc=self, pk=pk)
        execution.save()
        result = execution.start()
        return result

    @property
    def tasks(self):
        return [{
            "name": self.__str__(),
            "action": {
                "module": self.module,
                "args": self.args,
            },
        }]

    @property
    def inventory(self):
        return self.project.inventory_obj

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


class AdHocExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    adhoc = models.ForeignKey('AdHoc', on_delete=models.SET_NULL, null=True)

    def start(self):
        result = {"raw": {}, "summary": {}}
        try:
            pre_execution_start.send(self.__class__, execution=self)
            runner = AdHocRunner(self.adhoc.inventory, options=self.project.cleaned_options)
            result = runner.run(self.adhoc.tasks, pattern=self.adhoc.pattern)
        except Exception as e:
            result['summary'] = {'error': str(e)}
        finally:
            post_execution_start.send(self.__class__, execution=self, result=result)
        return result

