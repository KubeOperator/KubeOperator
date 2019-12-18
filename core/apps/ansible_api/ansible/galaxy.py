# -*- coding: utf-8 -*-
#
import os
import yaml
from collections import namedtuple

import markdown2
from ansible.galaxy import Galaxy
from ansible import constants as C
from ansible.errors import AnsibleError
from ansible.galaxy.api import GalaxyAPI
from ansible.galaxy.role import GalaxyRole


GALAXY_OPTION = namedtuple('Option', ['api_server', 'ignore_certs', 'force'])


class MyGalaxy(Galaxy):
    def __init__(self):
        option = GALAXY_OPTION(
            api_server=C.GALAXY_SERVER,
            ignore_certs=True,
            force=False,
        )
        super().__init__(options=option)


class MyGalaxyAPI(GalaxyAPI):
    def __init__(self):
        galaxy = MyGalaxy()
        self._role_info = {}
        super().__init__(galaxy)

    def lookup_role_by_name(self, role_name, notify=True):
        if self._role_info.get(role_name):
            return self._role_info[role_name]
        info = super().lookup_role_by_name(role_name, notify=notify)
        self._role_info[role_name] = info
        return info

    def role_git_url(self, role_name):
        info = self.lookup_role_by_name(role_name)
        if info:
            return 'https://github.com/{github_user}/{github_repo}'.format(**info)
        return None


class MyGalaxyRole(GalaxyRole):
    def __init__(self, name, src=None, version=None, scm=None, path=None):
        galaxy = MyGalaxy()
        self._meta_with_readme = None
        self._default_variables = None
        super().__init__(galaxy, name, src=src, version=version, scm=scm, path=path)

    def force_install(self):
        self.options = self.options._replace(force=True)
        self.install()

    def install(self):
        try:
            super().install()
            return True, None
        except AnsibleError as e:
            return False, e

    @property
    def metadata(self):
        metadata = super().metadata
        if self._meta_with_readme:
            return self._meta_with_readme

        readme_path = os.path.join(self.path, 'README.md')
        if os.path.isfile(readme_path):
            with open(readme_path, 'r') as f:
                readme_data = f.read()
                metadata['readme'] = readme_data
                metadata['readme_html'] = markdown2.markdown(readme_data)
                self._meta_with_readme = metadata
        return self._meta_with_readme

    @property
    def default_variables(self):
        if self._default_variables is None:
            defaults_path = os.path.join(self.path, 'defaults', 'main.yml')
            if os.path.isfile(defaults_path):
                with open(defaults_path, 'r') as f:
                    self._default_variables = yaml.safe_load(f)
        return self._default_variables
