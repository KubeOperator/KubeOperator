# -*- coding: utf-8 -*-
#

import os
import shutil
import json
import yaml
import git
import requests
import shutil

from celery import current_task
from django.db import models
from django.utils.translation import ugettext_lazy as _
from django_celery_beat.models import PeriodicTask

from common import models as common_models
from celery_api.utils import (
    delete_celery_periodic_task, disable_celery_periodic_task,
    create_or_update_periodic_task
)

from .mixins import AbstractProjectResourceModel, AbstractExecutionModel
from .role import Role
from ..utils import get_logger
from ..ansible.runner import PlayBookRunner
from ..signals import pre_execution_start, post_execution_start

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


class Playbook(AbstractProjectResourceModel):
    TYPE_JSON, TYPE_TEXT, TYPE_FILE, TYPE_GIT, TYPE_HTTP, TYPE_LOCAL = (
        'json', 'text', 'file', 'git', 'http', 'local',
    )
    TYPE_CHOICES = (
        (TYPE_JSON, TYPE_JSON),
        (TYPE_TEXT, TYPE_TEXT),
        (TYPE_FILE, TYPE_FILE),
        (TYPE_GIT, TYPE_GIT),
        (TYPE_HTTP, TYPE_HTTP),
        (TYPE_LOCAL, TYPE_LOCAL),
    )
    UPDATE_POLICY_ALWAYS, UPDATE_POLICY_IF_NOT_PRESENT, UPDATE_POLICY_NEVER = ('always', 'if_not_present', 'never')
    UPDATE_POLICY_CHOICES = (
        (UPDATE_POLICY_IF_NOT_PRESENT, _('Always')),
        (UPDATE_POLICY_ALWAYS, _("If not present")),
        (UPDATE_POLICY_NEVER, _("Never")),
    )
    name = models.SlugField(max_length=128, allow_unicode=True, verbose_name=_('Name'))
    alias = models.CharField(max_length=128, blank=True, default='site.yml')
    type = models.CharField(choices=TYPE_CHOICES, default=TYPE_JSON, max_length=16)
    plays = models.ManyToManyField('Play', verbose_name='Plays')
    git = common_models.JsonDictCharField(max_length=4096, default={'repo': '', 'branch': 'master'})
    url = models.URLField(verbose_name=_("http url"), blank=True)
    update_policy = models.CharField(choices=UPDATE_POLICY_CHOICES, max_length=16, default=UPDATE_POLICY_IF_NOT_PRESENT)

    # Extra schedule content
    is_periodic = models.BooleanField(default=False, verbose_name=_("Enable"))
    interval = models.CharField(verbose_name=_("Interval"), null=True, blank=True, max_length=128, help_text=_("s/m/d"))
    crontab = models.CharField(verbose_name=_("Crontab"), null=True, blank=True, max_length=128, help_text=_("5 * * * *"))
    meta = common_models.JsonDictTextField(blank=True, verbose_name=_("Meta"))

    execute_times = models.IntegerField(default=0)
    comment = models.TextField(blank=True, verbose_name=_("Comment"))
    is_active = models.BooleanField(default=True, verbose_name=_("Active"))
    created_by = models.CharField(max_length=128, blank=True, null=True)
    date_created = models.DateTimeField(auto_now_add=True)

    class Meta:
        unique_together = ["name", "project"]

    def __str__(self):
        return '{}-{}'.format(self.project, self.name)

    def playbook_dir(self, auto_create=True):
        path = os.path.join(self.project.playbooks_dir, str(self.name))
        if not os.path.isdir(path) and auto_create:
            os.makedirs(path, exist_ok=True)
        return path

    @property
    def playbook_path(self):
        path = os.path.join(self.playbook_dir(), self.alias)
        return path

    @property
    def latest_execution(self):
        try:
            return self.executions.all().latest()
        except PlaybookExecution.DoesNotExist:
            return None

    def get_plays_data(self, fmt='py'):
        return Play.get_plays_data(self.plays.all(), fmt=fmt)

    def install_from_git(self):
        success, error = True, None
        if not self.git.get('repo'):
            success, error = False, 'Not repo get'
            return success, error
        try:
            if os.path.isdir(os.path.join(self.playbook_dir(), '.git')):
                if self.update_policy == self.UPDATE_POLICY_ALWAYS:
                    print("Update playbook from: {}".format(self.git.get('repo')))
                    repo = git.Repo(self.playbook_dir())
                    remote = repo.remote()
                    remote.pull()
            else:
                print("Install playbook from: {}".format(self.git.get('repo')))
                git.Repo.clone_from(
                    self.git['repo'], self.playbook_dir(),
                    branch=self.git.get('branch'), depth=1,
                )
        except Exception as e:
            success, error = False, e
        return success, error

    def install_from_http(self):
        if os.listdir(self.playbook_dir()):
            os.removedirs(self.playbook_dir())
        r = requests.get(self.url)
        tmp_file_path = os.path.join(self.playbook_dir(), 'tmp')
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

    def install_from_local(self):
        playbook_dir = self.playbook_dir(auto_create=False)
        if self.update_policy == self.UPDATE_POLICY_NEVER:
            return True, None
        if os.path.isfile(self.playbook_path) and \
                self.update_policy == self.UPDATE_POLICY_IF_NOT_PRESENT:
            return True, None
        shutil.rmtree(playbook_dir, ignore_errors=True)
        url = self.url
        if self.url.startswith('file://'):
            url = self.url.replace('file://', '')
        try:
            shutil.copytree(url, playbook_dir)
        except Exception as e:
            return False, e
        return True, None

    def install(self):
        if self.type == self.TYPE_JSON:
            return self.install_from_plays()
        elif self.type == self.TYPE_GIT:
            return self.install_from_git()
        elif self.type == self.TYPE_LOCAL:
            return self.install_from_local()
        else:
            return False, 'Not support {}'.format(self.type)

    def execute(self):
        pk = current_task.request.id if current_task else None
        execution = PlaybookExecution(playbook=self, pk=pk)
        execution.save()
        result = execution.start()
        return result

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

    def cleanup(self):
        self.remove_period_task()
        shutil.rmtree(self.playbook_dir(), ignore_errors=True)


class PlaybookExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    playbook = models.ForeignKey(Playbook, related_name='executions', on_delete=models.SET_NULL, null=True)

    class Meta:
        get_latest_by = 'date_start'

    def __str__(self):
        return "{} run at {}".format(self.playbook.__str__(), self.date_start)

    def start(self):
        pre_execution_start.send(self.__class__, execution=self)
        result = {"raw": {}, "summary": {}}
        success, err = self.playbook.install()
        if not success:
            result["summary"] = {"error": str(err)}
            post_execution_start.send(self.__class__, execution=self, result=result)
            return result
        os.chdir(self.playbook.playbook_dir())
        try:
            runner = PlayBookRunner(
                self.project.inventory_obj,
                options=self.project.cleaned_options,
            )
            result = runner.run(self.playbook.playbook_path)
        except IndexError as e:
            result["summary"] = {'error': str(e)}
        finally:
            post_execution_start.send(self.__class__, execution=self, result=result)
        return result
