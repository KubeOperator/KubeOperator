import uuid

from django.db import models

__all__ = ["Script", "ScriptExecution"]

from ansible_api.ansible import AdHocRunner
from ansible_api.models.mixins import AbstractExecutionModel
from ops.signals import pre_script_execution_start, post_script_execution_start


class Script(models.Model):
    SCRIPT_TYPE_PYTHON = "python"
    SCRIPT_TYPE_SHELL = "shell"
    SCRIPT_TYPE_CHOICES = (
        (SCRIPT_TYPE_PYTHON, "python"),
        (SCRIPT_TYPE_SHELL, "shell")
    )
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.SlugField(max_length=128, allow_unicode=True, unique=True)
    type = models.CharField(max_length=64, default=SCRIPT_TYPE_SHELL, choices=SCRIPT_TYPE_CHOICES)
    content = models.TextField(default="")
    date_created = models.DateTimeField(auto_now_add=True)


class ScriptExecution(AbstractExecutionModel):
    script = models.ForeignKey("Script", on_delete=models.CASCADE)
    cluster = models.ForeignKey("kubeops_api.Cluster", on_delete=models.CASCADE)
    targets = models.ManyToManyField("kubeops_api.Host")

    def start(self):
        pre_script_execution_start.send(self.__class__, execution=self)
        runner = AdHocRunner(self.cluster)
        tasks = [{'action': {'module': "shell", 'args': self.script.content}}]
        result = runner.run(tasks, pattern=self.targets)
        post_script_execution_start.send(self.__class__, execution=self, result=result)
