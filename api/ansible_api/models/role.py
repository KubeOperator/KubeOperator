# -*- coding: utf-8 -*-
#
from collections import OrderedDict
import os

import git
from django.db import models
from django.conf import settings
from django.utils.translation import ugettext_lazy as _

from common import models as common_models
from .mixins import AbstractProjectResourceModel
from .utils import name_validator
from ..ansible.galaxy import MyGalaxyRole, MyGalaxyAPI

__all__ = ['Role']


class Role(AbstractProjectResourceModel):
    STATE_NOT_INSTALL = 'uninstalled'
    STATE_INSTALLED = 'installed'
    STATE_INSTALLING = 'installing'
    STATE_FAILED = 'failed'
    STATE_CHOICES = (
        (STATE_NOT_INSTALL, 'UnInstalled'),
        (STATE_INSTALLED, 'Installed'),
        (STATE_INSTALLING, 'Installing'),
        (STATE_FAILED, 'Failed')
    )
    TYPE_GIT = 'git'
    TYPE_HTTP = 'http'
    TYPE_GALAXY = 'galaxy'
    TYPE_FILE = 'file'
    TYPE_CHOICES = (
        (TYPE_GALAXY, 'galaxy'),
        (TYPE_GIT, 'git'),
        (TYPE_HTTP, 'http'),
        (TYPE_FILE, 'file'),
    )

    name = models.CharField(max_length=128, validators=[name_validator])
    type = models.CharField(max_length=16, choices=TYPE_CHOICES, default=TYPE_GALAXY)
    comment = models.CharField(max_length=1024, blank=True, verbose_name=_("Comment"))
    galaxy_name = models.CharField(max_length=128, blank=True, null=True)
    git = common_models.JsonDictCharField(max_length=4096, default={'repo': '', 'branch': 'master'})
    url = models.CharField(max_length=1024, verbose_name=_("Url"), blank=True)
    logo = models.ImageField(verbose_name='Logo', upload_to="logo", null=True)
    categories = models.CharField(max_length=256, verbose_name=_("Tags"), blank=True)
    version = models.CharField(max_length=1024, blank=True, default='master')
    state = models.CharField(default=STATE_NOT_INSTALL, choices=STATE_CHOICES, max_length=16)
    meta = common_models.JsonDictTextField(verbose_name=_("Meta"), blank=True)
    meta_ext = common_models.JsonDictTextField(verbose_name=_("Meta Ext"), blank=True)
    created_by = models.CharField(max_length=128, blank=True, null=True, default='')
    date_created = models.DateTimeField(auto_now_add=True)
    date_updated = models.DateTimeField(auto_now=True)

    class Meta:
        unique_together = ('name', 'project')

    def __str__(self):
        return self.name

    def delete(self, using=None, keep_parents=False):
        role = MyGalaxyRole(self.name, path=self.project.roles_dir)
        role.remove()
        return super().delete(using=using, keep_parents=keep_parents)

    @property
    def _role(self):
        role = MyGalaxyRole(self.name, path=self.project.roles_dir)
        return role

    @property
    def variables(self):
        return self._role.default_variables

    @property
    def role_dir(self):
        return os.path.join(self.project.roles_dir, self.name)

    @property
    def meta_all(self):
        meta = OrderedDict([
            ('name', self.name),
            ('version', self.version),
            ('comment', self.comment),
            ('state', self.get_state_display()),
            ('url', self.url),
            ('type', self.type),
            ('categories', self.categories),
        ])
        if isinstance(self.meta, dict):
            meta.update(self.meta)
        if isinstance(self.meta_ext, dict):
            meta.update(self.meta_ext)
        meta.pop('readme', None)
        meta.pop('readme_html', None)
        galaxy_info = meta.pop('galaxy_info', {})
        for k, v in galaxy_info.items():
            if k == 'platforms':
                v = ' '.join([i['name']+str(i['versions']) for i in v])
            meta[k] = v
        return meta

    def install_from_galaxy(self):
        api = MyGalaxyAPI()
        role = MyGalaxyRole(self.galaxy_name, path=self.project.roles_dir)
        success, error = role.install()
        if success:
            self.comment = api.lookup_role_by_name(self.galaxy_name)['description']
            self.url = api.role_git_url(self.galaxy_name)
            self.version = role.version
            self.meta = role.metadata
            categories = ''
            if self.meta and self.meta['galaxy_info'].get('categories'):
                categories = self.meta['galaxy_info']['categories']
            elif self.meta and self.meta['galaxy_info'].get('galaxy_tags'):
                categories = self.meta['galaxy_info']['galaxy_tags']
            self.categories = ','.join(categories) if isinstance(categories, list) else str(categories)
            os.rename(os.path.join(self.project.roles_dir, self.galaxy_name), self.role_dir)
        return success, error

    def install_from_git(self):
        success, error = True, None
        if not self.git.get('repo'):
            success = False
            error = 'Not repo get'
            return success, error
        print("Install playbook from: {}".format(self.git.get('repo')))
        try:
            if os.path.isdir(os.path.join(self.role_dir, '.git')):
                repo = git.Repo(self.role_dir)
                remote = repo.remote()
                remote.pull()
            else:
                git.Repo.clone_from(
                    self.git['repo'], self.role_dir,
                    branch=self.git.get('branch'), depth=1,
                )
        except Exception as e:
            success = False
            error = e
        return success, error

    def install(self):
        self.state = self.STATE_INSTALLING
        if self.type == self.TYPE_GALAXY:
            success, err = self.install_from_galaxy()
        elif self.type == self.TYPE_GIT:
            success, err = self.install_from_git()
        else:
            success = False
            err = Exception("From {}, using other function".format(self.type))
        if success:
            self.state = self.STATE_INSTALLED
        else:
            self.state = self.STATE_FAILED
        self.save()
        return success, err

    def uninstall(self):
        role = MyGalaxyRole(self.name, path=self.project.roles_dir)
        role.remove()

    @property
    def logo_url(self):
        default = settings.STATIC_URL + "ansible/img/role_logo_default.png"
        if self.logo:
            return self.logo.url
        return default

