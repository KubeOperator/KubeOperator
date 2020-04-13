import uuid

from django.db import models

__all__ = ["Script"]


class Script(models.Model):
    SCRIPT_TYPE_PYTHON = "python"
    SCRIPT_TYPE_SHELL = "shell"
    SCRIPT_TYPE_CHOICES = (
        (SCRIPT_TYPE_PYTHON, "python"),
        (SCRIPT_TYPE_SHELL, "shell")
    )
    DEFAULT_SHELL_INTERPRETER = "/usr/bin/bash"
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.SlugField(max_length=128, allow_unicode=True, unique=True)
    type = models.CharField(max_length=64, default=SCRIPT_TYPE_SHELL, choices=SCRIPT_TYPE_CHOICES)
    content = models.TextField(default="")
    date_created = models.DateTimeField(auto_now_add=True)
