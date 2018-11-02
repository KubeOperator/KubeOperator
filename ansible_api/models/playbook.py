# -*- coding: utf-8 -*-
#

import os
import json
import yaml
import uuid
import git
import requests

from django.db import models
from django.utils.translation import ugettext_lazy as _
from django_celery_beat.models import PeriodicTask

from common import models as common_models
from celery_api.utils import (
    delete_celery_periodic_task, disable_celery_periodic_task,
    create_or_update_periodic_task, get_celery_task_log_path
)

from .mixins import AbstractProjectResourceModel
from .role import Role
from .utils import format_result_as_list
from ..utils import get_logger
from ..ansible.runner import PlayBookRunner
from ..signals import pre_playbook_exec, post_playbook_exec

logger = get_logger(__file__)

__all__ = ['Play', 'Playbook', 'PlaybookExecution']


class PlaybookQuerySet(models.QuerySet):

    def delete(self):
        for obj in self:
            obj.delete()
        return super().delete()

    def update(self, **kwargs):
        queryset = super().update(**kwargs)
        for obj in self:
            obj.save()
        return queryset


class Play(AbstractProjectResourceModel):
    name = models.CharField(max_length=128, verbose_name=_('Name'), blank=True, null=True)
    pattern = models.CharField(max_length=1024, default='all', verbose_name=_('Pattern'))
    gather_facts = models.BooleanField(default=False)
    vars = common_models.JsonDictTextField(verbose_name=_('Vars'), blank=True, null=True)
    tasks = common_models.JsonListTextField(verbose_name=_('Tasks'), blank=True, null=True)
    roles = common_models.JsonListTextField(verbose_name=_('Roles'), blank=True, null=True)

    @staticmethod
    def format_data(data, fmt='py'):
        if fmt == 'yaml':
            return yaml.safe_dump(data, default_flow_style=False)
        elif fmt == 'json':
            return json.dumps(data, indent=4)
        else:
            return data

    def get_play_data(self, fmt='py'):
        data = {
            'hosts': self.pattern,
            'gather_facts': self.gather_facts,
            'vars': self.vars or [],
            'tasks': self.tasks or [],
            'roles': self.roles or [],
        }
        return self.format_data(data, fmt=fmt)

    @classmethod
    def get_plays_data(cls, plays, fmt='py'):
        data = []
        for play in plays:
            data.append(play.get_play_data())
        return cls.format_data(data, fmt=fmt)

    def get_play_roles_names(self):
        names = []
        for role in self.roles or []:
            name = role['role'] if isinstance(role, dict) else role
            names.append(name)
        return names

    def check_role(self):
        for role_name in self.get_play_roles_names():
            try:
                role = self.project.role_set.get(name=role_name)
            except Role.DoesNotExist:
                error = "- Role not exist in project: {}".format(role_name)
                logger.error(error)
                return False, error
            if role.state != Role.STATE_INSTALLED:
                success, error = role.install()
                if not success:
                    msg = "- Install role failed {}: {}".format(role_name, error)
                    logger.error(msg)
                return False, error
        return True, None

    @classmethod
    def get_plays_roles_names(cls, plays):
        names = []
        for play in plays:
            names.extend(play.get_play_roles_names())
        return names


class Playbook(AbstractProjectResourceModel):
    TYPE_JSON, TYPE_TEXT, TYPE_FILE, TYPE_GIT, TYPE_HTTP = ('json', 'text', 'file', 'git', 'http')
    TYPE_CHOICES = (
        (TYPE_JSON, TYPE_JSON),
        (TYPE_TEXT, TYPE_TEXT),
        (TYPE_FILE, TYPE_FILE),
        (TYPE_GIT, TYPE_GIT),
        (TYPE_HTTP, TYPE_HTTP),
    )
    name = models.SlugField(max_length=128, allow_unicode=True, verbose_name=_('Name'))
    type = models.CharField(choices=TYPE_CHOICES, default=TYPE_JSON, max_length=16)
    plays = models.ManyToManyField('Play', verbose_name='Plays')
    git = common_models.JsonDictCharField(max_length=4096, default={'repo': '', 'branch': 'master'})
    url = models.URLField(verbose_name=_("http url"), blank=True)
    rel_path = models.CharField(max_length=128, blank=True, default='site.yml')

    # Extra schedule content
    is_periodic = models.BooleanField(default=False, verbose_name=_("Enable"))
    interval = models.CharField(verbose_name=_("Interval"), null=True, blank=True, max_length=128, help_text=_("s/m/d"))
    crontab = models.CharField(verbose_name=_("Crontab"), null=True, blank=True, max_length=128, help_text=_("5 * * * *"))
    meta = common_models.JsonDictTextField(blank=True, verbose_name=_("Meta"))

    times = models.IntegerField(default=0)
    comment = models.TextField(blank=True, verbose_name=_("Comment"))
    is_active = models.BooleanField(default=True, verbose_name=_("Active"))
    created_by = models.CharField(max_length=128, blank=True, null=True)
    date_created = models.DateTimeField(auto_now_add=True)

    class Meta:
        unique_together = ["name", "project"]

    def __str__(self):
        return '{}-{}'.format(self.project, self.name)

    @property
    def playbook_dir(self):
        path = os.path.join(self.project.playbooks_dir, str(self.name))
        os.makedirs(path, exist_ok=True)
        return path

    @property
    def latest_execute_history(self):
        try:
            return self.history.all().latest()
        except PlaybookExecution.DoesNotExist:
            return None

    @property
    def playbook_path(self):
        # path = 'playbooks/prerequisites.yml'
        path = os.path.join(self.playbook_dir, self.rel_path)
        return path

    def get_plays_data(self, fmt='py'):
        return Play.get_plays_data(self.plays.all(), fmt=fmt)

    def install_from_git(self):
        success, error = True, None
        if not self.git.get('repo'):
            success, error = False, 'Not repo get'
            return success, error
        try:
            print("Install playbook from: {}".format(self.git.get('repo')))
            if os.path.isdir(os.path.join(self.playbook_dir, '.git')):
                repo = git.Repo(self.playbook_dir)
                remote = repo.remote()
                remote.pull()
            else:
                git.Repo.clone_from(
                    self.git['repo'], self.playbook_dir,
                    branch=self.git.get('branch'), depth=1,
                )
        except Exception as e:
            success, error = False, e
        return success, error

    def install_from_http(self):
        if os.listdir(self.playbook_dir):
            os.removedirs(self.playbook_dir)
        r = requests.get(self.url)
        tmp_file_path = os.path.join(self.playbook_dir, 'tmp')
        with open(tmp_file_path, 'wb') as f:
            f.write(r.content)
        # TODO: compress it

    def install_from_plays(self):
        for play in self.plays.all():
            success, error = play.check_role()
            if not success:
                return success, error
        with open(self.playbook_path, 'w') as f:
            f.write(self.get_plays_data(fmt='yaml'))
        return True, None

    def install(self):
        if self.type == self.TYPE_JSON:
            return self.install_from_plays()
        elif self.type == self.TYPE_GIT:
            return self.install_from_git()
        else:
            return False, 'Not support {}'.format(self.type)

    def pre_check(self):
        return True, None

    def execute(self, save_history=True):
        result = {"raw": {}, "summary": {}}
        success, err = self.install()
        if not success:
            result["summary"] = {"error": str(err)}
        os.chdir(self.playbook_dir)
        try:
            pre_playbook_exec.send(self.__class__, playbook=self, save_history=save_history)
            runner = PlayBookRunner(
                self.project.inventory_obj,
                options=self.project.cleaned_options,
            )
            result = runner.run(self.playbook_path)
        except IndexError as e:
            result["summary"] = {'error': str(e)}
        finally:
            post_playbook_exec.send(self.__class__, playbook=self, save_history=save_history, result=result)
        return format_result_as_list(result.get('summary', {}))

    def create_period_task(self):
        from ..tasks import execute_playbook
        tasks = {
            self.__str__(): {
                "task": execute_playbook.name,
                "interval": self.interval or None,
                "crontab": self.crontab or None,
                "args": (str(self.id),),
                "kwargs": {"name": self.__str__()},
                "enabled": True,
            }
        }
        create_or_update_periodic_task(tasks)

    def disable_period_task(self):
        disable_celery_periodic_task(self.__str__())

    def remove_period_task(self):
        if self.is_periodic:
            delete_celery_periodic_task(self.__str__())

    @property
    def period_task(self):
        try:
            return PeriodicTask.objects.get(name=self.__str__())
        except PeriodicTask.DoesNotExist:
            return None

    @staticmethod
    def test_tasks():
        return [
            {
                "name": "Test ping",
                "ping": ""
            },
            {
                "name": "Ifconfig",
                "command": "ifconfig"
            }
        ]

    @staticmethod
    def test_roles():
        return [
            {
                "role": "bennojoy.memcached",
                "memcached_port": 11244,
                "memcached_cache_size": 512
            }
        ]


class PlaybookExecution(AbstractProjectResourceModel):
    id = models.CharField(primary_key=True, default=uuid.uuid4, max_length=36)
    playbook = models.ForeignKey(Playbook, related_name='history', on_delete=models.SET_NULL, null=True)
    num = models.IntegerField(default=1)
    timedelta = models.FloatField(default=0.0, verbose_name=_('Time'), null=True)
    is_finished = models.BooleanField(default=False, verbose_name=_('Is finished'))
    is_success = models.BooleanField(default=False, verbose_name=_('Is success'))
    raw = common_models.JsonDictTextField(blank=True, null=True, default='{}', verbose_name=_('Adhoc raw result'))
    summary = common_models.JsonDictTextField(blank=True, null=True, default='{}', verbose_name=_('Adhoc summary'))
    date_start = models.DateTimeField(auto_now_add=True, verbose_name=_('Start time'))
    date_finished = models.DateTimeField(blank=True, null=True, verbose_name=_('End time'))

    def __str__(self):
        return "{} run at {}".format(self.playbook.__str__(), self.date_start)

    @property
    def log_path(self):
        return get_celery_task_log_path(self.id)

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

    class Meta:
        get_latest_by = 'date_start'
