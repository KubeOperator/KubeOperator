#!/usr/bin/env python
# pylint: disable=missing-docstring
#     ___ ___ _  _ ___ ___    _ _____ ___ ___
#    / __| __| \| | __| _ \  /_\_   _| __|   \
#   | (_ | _|| .` | _||   / / _ \| | | _|| |) |
#    \___|___|_|\_|___|_|_\/_/_\_\_|_|___|___/_ _____
#   |   \ / _ \  | \| |/ _ \_   _| | __|   \_ _|_   _|
#   | |) | (_) | | .` | (_) || |   | _|| |) | |  | |
#   |___/ \___/  |_|\_|\___/ |_|   |___|___/___| |_|
#
# Copyright 2016 Red Hat, Inc. and/or its affiliates
# and other contributors as indicated by the @author tags.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# -*- -*- -*- Begin included fragment: lib/import.py -*- -*- -*-

# pylint: disable=wrong-import-order,wrong-import-position,unused-import

from __future__ import print_function  # noqa: F401
import copy  # noqa: F401
import fcntl  # noqa: F401
import json   # noqa: F401
import os  # noqa: F401
import re  # noqa: F401
import shutil  # noqa: F401
import tempfile  # noqa: F401
import time  # noqa: F401

try:
    import ruamel.yaml as yaml  # noqa: F401
except ImportError:
    import yaml  # noqa: F401

from ansible.module_utils.basic import AnsibleModule

# -*- -*- -*- End included fragment: lib/import.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: doc/repoquery -*- -*- -*-

DOCUMENTATION = '''
---
module: repoquery
short_description: Query package information from Yum repositories
description:
  - Query package information from Yum repositories.
options:
  state:
    description:
    - The expected state. Currently only supports list.
    required: false
    default: list
    choices: ["list"]
    aliases: []
  name:
    description:
    - The name of the package to query
    required: true
    default: None
    aliases: []
  query_type:
    description:
    - Narrows the packages queried based off of this value.
    - If repos, it narrows the query to repositories defined on the machine.
    - If installed, it narrows the query to only packages installed on the machine.
    - If available, it narrows the query to packages that are available to be installed.
    - If recent, it narrows the query to only recently edited packages.
    - If updates, it narrows the query to only packages that are updates to existing installed packages.
    - If extras, it narrows the query to packages that are not present in any of the available repositories.
    - If all, it queries all of the above.
    required: false
    default: repos
    aliases: []
  verbose:
    description:
    - Shows more detail for the requested query.
    required: false
    default: false
    aliases: []
  show_duplicates:
    description:
    - Shows multiple versions of a package.
    required: false
    default: false
    aliases: []
  match_version:
    description:
    - Match the specific version given to the package.
    required: false
    default: None
    aliases: []
author:
- "Matt Woodson <mwoodson@redhat.com>"
extends_documentation_fragment: []
'''

EXAMPLES = '''
# Example 1: Get bash versions
  - name: Get bash version
    repoquery:
      name: bash
      show_duplicates: True
    register: bash_out

# Results:
#    ok: [localhost] => {
#        "bash_out": {
#            "changed": false,
#            "results": {
#                "cmd": "/usr/bin/repoquery --quiet --pkgnarrow=repos --queryformat=%{version}|%{release}|%{arch}|%{repo}|%{version}-%{release} --show-duplicates bash",
#                "package_found": true,
#                "package_name": "bash",
#                "returncode": 0,
#                "versions": {
#                    "available_versions": [
#                        "4.2.45",
#                        "4.2.45",
#                        "4.2.45",
#                        "4.2.46",
#                        "4.2.46",
#                        "4.2.46",
#                        "4.2.46"
#                    ],
#                    "available_versions_full": [
#                        "4.2.45-5.el7",
#                        "4.2.45-5.el7_0.2",
#                        "4.2.45-5.el7_0.4",
#                        "4.2.46-12.el7",
#                        "4.2.46-19.el7",
#                        "4.2.46-20.el7_2",
#                        "4.2.46-21.el7_3"
#                    ],
#                    "latest": "4.2.46",
#                    "latest_full": "4.2.46-21.el7_3"
#                }
#            },
#            "state": "present"
#        }
#    }



# Example 2: Get bash versions verbosely
  - name: Get bash versions verbosely
    repoquery:
      name: bash
      show_duplicates: True
      verbose: True
    register: bash_out

# Results:
#    ok: [localhost] => {
#        "bash_out": {
#            "changed": false,
#            "results": {
#                "cmd": "/usr/bin/repoquery --quiet --pkgnarrow=repos --queryformat=%{version}|%{release}|%{arch}|%{repo}|%{version}-%{release} --show-duplicates bash",
#                "package_found": true,
#                "package_name": "bash",
#                "raw_versions": {
#                    "4.2.45-5.el7": {
#                        "arch": "x86_64",
#                        "release": "5.el7",
#                        "repo": "rhel-7-server-rpms",
#                        "version": "4.2.45",
#                        "version_release": "4.2.45-5.el7"
#                    },
#                    "4.2.45-5.el7_0.2": {
#                        "arch": "x86_64",
#                        "release": "5.el7_0.2",
#                        "repo": "rhel-7-server-rpms",
#                        "version": "4.2.45",
#                        "version_release": "4.2.45-5.el7_0.2"
#                    },
#                    "4.2.45-5.el7_0.4": {
#                        "arch": "x86_64",
#                        "release": "5.el7_0.4",
#                        "repo": "rhel-7-server-rpms",
#                        "version": "4.2.45",
#                        "version_release": "4.2.45-5.el7_0.4"
#                    },
#                    "4.2.46-12.el7": {
#                        "arch": "x86_64",
#                        "release": "12.el7",
#                        "repo": "rhel-7-server-rpms",
#                        "version": "4.2.46",
#                        "version_release": "4.2.46-12.el7"
#                    },
#                    "4.2.46-19.el7": {
#                        "arch": "x86_64",
#                        "release": "19.el7",
#                        "repo": "rhel-7-server-rpms",
#                        "version": "4.2.46",
#                        "version_release": "4.2.46-19.el7"
#                    },
#                    "4.2.46-20.el7_2": {
#                        "arch": "x86_64",
#                        "release": "20.el7_2",
#                        "repo": "rhel-7-server-rpms",
#                        "version": "4.2.46",
#                        "version_release": "4.2.46-20.el7_2"
#                    },
#                    "4.2.46-21.el7_3": {
#                        "arch": "x86_64",
#                        "release": "21.el7_3",
#                        "repo": "rhel-7-server-rpms",
#                        "version": "4.2.46",
#                        "version_release": "4.2.46-21.el7_3"
#                    }
#                },
#                "results": "4.2.45|5.el7|x86_64|rhel-7-server-rpms|4.2.45-5.el7\n4.2.45|5.el7_0.2|x86_64|rhel-7-server-rpms|4.2.45-5.el7_0.2\n4.2.45|5.el7_0.4|x86_64|rhel-7-server-rpms|4.2.45-5.el7_0.4\n4.2.46|12.el7|x86_64|rhel-7-server-rpms|4.2.46-12.el7\n4.2.46|19.el7|x86_64|rhel-7-server-rpms|4.2.46-19.el7\n4.2.46|20.el7_2|x86_64|rhel-7-server-rpms|4.2.46-20.el7_2\n4.2.46|21.el7_3|x86_64|rhel-7-server-rpms|4.2.46-21.el7_3\n",
#                "returncode": 0,
#                "versions": {
#                    "available_versions": [
#                        "4.2.45",
#                        "4.2.45",
#                        "4.2.45",
#                        "4.2.46",
#                        "4.2.46",
#                        "4.2.46",
#                        "4.2.46"
#                    ],
#                    "available_versions_full": [
#                        "4.2.45-5.el7",
#                        "4.2.45-5.el7_0.2",
#                        "4.2.45-5.el7_0.4",
#                        "4.2.46-12.el7",
#                        "4.2.46-19.el7",
#                        "4.2.46-20.el7_2",
#                        "4.2.46-21.el7_3"
#                    ],
#                    "latest": "4.2.46",
#                    "latest_full": "4.2.46-21.el7_3"
#                }
#            },
#            "state": "present"
#        }
#    }

# Example 3: Match a specific version
  - name: matched versions repoquery test
    repoquery:
      name: atomic-openshift
      show_duplicates: True
      match_version: 3.3
    register: openshift_out

# Result:

#    ok: [localhost] => {
#        "openshift_out": {
#            "changed": false,
#            "results": {
#                "cmd": "/usr/bin/repoquery --quiet --pkgnarrow=repos --queryformat=%{version}|%{release}|%{arch}|%{repo}|%{version}-%{release} --show-duplicates atomic-openshift",
#                "package_found": true,
#                "package_name": "atomic-openshift",
#                "returncode": 0,
#                "versions": {
#                    "available_versions": [
#                        "3.2.0.43",
#                        "3.2.1.23",
#                        "3.3.0.32",
#                        "3.3.0.34",
#                        "3.3.0.35",
#                        "3.3.1.3",
#                        "3.3.1.4",
#                        "3.3.1.5",
#                        "3.3.1.7",
#                        "3.4.0.39"
#                    ],
#                    "available_versions_full": [
#                        "3.2.0.43-1.git.0.672599f.el7",
#                        "3.2.1.23-1.git.0.88a7a1d.el7",
#                        "3.3.0.32-1.git.0.37bd7ea.el7",
#                        "3.3.0.34-1.git.0.83f306f.el7",
#                        "3.3.0.35-1.git.0.d7bd9b6.el7",
#                        "3.3.1.3-1.git.0.86dc49a.el7",
#                        "3.3.1.4-1.git.0.7c8657c.el7",
#                        "3.3.1.5-1.git.0.62700af.el7",
#                        "3.3.1.7-1.git.0.0988966.el7",
#                        "3.4.0.39-1.git.0.5f32f06.el7"
#                    ],
#                    "latest": "3.4.0.39",
#                    "latest_full": "3.4.0.39-1.git.0.5f32f06.el7",
#                    "matched_version_found": true,
#                    "matched_version_full_latest": "3.3.1.7-1.git.0.0988966.el7",
#                    "matched_version_latest": "3.3.1.7",
#                    "matched_versions": [
#                        "3.3.0.32",
#                        "3.3.0.34",
#                        "3.3.0.35",
#                        "3.3.1.3",
#                        "3.3.1.4",
#                        "3.3.1.5",
#                        "3.3.1.7"
#                    ],
#                    "matched_versions_full": [
#                        "3.3.0.32-1.git.0.37bd7ea.el7",
#                        "3.3.0.34-1.git.0.83f306f.el7",
#                        "3.3.0.35-1.git.0.d7bd9b6.el7",
#                        "3.3.1.3-1.git.0.86dc49a.el7",
#                        "3.3.1.4-1.git.0.7c8657c.el7",
#                        "3.3.1.5-1.git.0.62700af.el7",
#                        "3.3.1.7-1.git.0.0988966.el7"
#                    ],
#                    "requested_match_version": "3.3"
#                }
#            },
#            "state": "present"
#        }
#    }

'''

# -*- -*- -*- End included fragment: doc/repoquery -*- -*- -*-

# -*- -*- -*- Begin included fragment: lib/repoquery.py -*- -*- -*-

'''
   class that wraps the repoquery commands in a subprocess
'''

# pylint: disable=too-many-lines,wrong-import-position,wrong-import-order

from collections import defaultdict  # noqa: E402


# pylint: disable=no-name-in-module,import-error
# Reason: pylint errors with "No name 'version' in module 'distutils'".
#         This is a bug: https://github.com/PyCQA/pylint/issues/73
from distutils.version import LooseVersion  # noqa: E402

import subprocess  # noqa: E402


class RepoqueryCLIError(Exception):
    '''Exception class for repoquerycli'''
    pass


def _run(cmds):
    ''' Actually executes the command. This makes mocking easier. '''
    proc = subprocess.Popen(cmds,
                            stdin=subprocess.PIPE,
                            stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE)

    stdout, stderr = proc.communicate()

    return proc.returncode, stdout, stderr


# pylint: disable=too-few-public-methods
class RepoqueryCLI(object):
    ''' Class to wrap the command line tools '''
    def __init__(self,
                 verbose=False):
        ''' Constructor for RepoqueryCLI '''
        self.verbose = verbose
        self.verbose = True

    def _repoquery_cmd(self, cmd, output=False, output_type='json'):
        '''Base command for repoquery '''
        cmds = ['/usr/bin/repoquery', '--plugins', '--quiet']

        cmds.extend(cmd)

        rval = {}
        results = ''
        err = None

        if self.verbose:
            print(' '.join(cmds))

        returncode, stdout, stderr = _run(cmds)

        rval = {
            "returncode": returncode,
            "results": results,
            "cmd": ' '.join(cmds),
        }

        if returncode == 0:
            if output:
                if output_type == 'raw':
                    rval['results'] = stdout

            if self.verbose:
                print(stdout)
                print(stderr)

            if err:
                rval.update({
                    "err": err,
                    "stderr": stderr,
                    "stdout": stdout,
                    "cmd": cmds
                })

        else:
            rval.update({
                "stderr": stderr,
                "stdout": stdout,
                "results": {},
            })

        return rval

# -*- -*- -*- End included fragment: lib/repoquery.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: class/repoquery.py -*- -*- -*-


class Repoquery(RepoqueryCLI):
    ''' Class to wrap the repoquery
    '''
    # pylint: disable=too-many-arguments,too-many-instance-attributes
    def __init__(self, name, query_type, show_duplicates,
                 match_version, ignore_excluders, verbose):
        ''' Constructor for YumList '''
        super(Repoquery, self).__init__(None)
        self.name = name
        self.query_type = query_type
        self.show_duplicates = show_duplicates
        self.match_version = match_version
        self.ignore_excluders = ignore_excluders
        self.verbose = verbose

        if self.match_version:
            self.show_duplicates = True

        self.query_format = "%{version}|%{release}|%{arch}|%{repo}|%{version}-%{release}"

        self.tmp_file = None

    def build_cmd(self):
        ''' build the repoquery cmd options '''

        repo_cmd = []

        repo_cmd.append("--pkgnarrow=" + self.query_type)
        repo_cmd.append("--queryformat=" + self.query_format)

        if self.show_duplicates:
            repo_cmd.append('--show-duplicates')

        if self.ignore_excluders:
            repo_cmd.append('--config=' + self.tmp_file.name)

        repo_cmd.append(self.name)

        return repo_cmd

    @staticmethod
    def process_versions(query_output):
        ''' format the package data into something that can be presented '''

        version_dict = defaultdict(dict)

        for version in query_output.decode().split('\n'):
            pkg_info = version.split("|")

            pkg_version = {}
            pkg_version['version'] = pkg_info[0]
            pkg_version['release'] = pkg_info[1]
            pkg_version['arch'] = pkg_info[2]
            pkg_version['repo'] = pkg_info[3]
            pkg_version['version_release'] = pkg_info[4]

            version_dict[pkg_info[4]] = pkg_version

        return version_dict

    def format_versions(self, formatted_versions):
        ''' Gather and present the versions of each package '''

        versions_dict = {}
        versions_dict['available_versions_full'] = list(formatted_versions.keys())

        # set the match version, if called
        if self.match_version:
            versions_dict['matched_versions_full'] = []
            versions_dict['requested_match_version'] = self.match_version
            versions_dict['matched_versions'] = []

        # get the "full version (version - release)
        versions_dict['available_versions_full'].sort(key=LooseVersion)
        versions_dict['latest_full'] = versions_dict['available_versions_full'][-1]

        # get the "short version (version)
        versions_dict['available_versions'] = []
        for version in versions_dict['available_versions_full']:
            versions_dict['available_versions'].append(formatted_versions[version]['version'])

            if self.match_version:
                if version.startswith(self.match_version):
                    versions_dict['matched_versions_full'].append(version)
                    versions_dict['matched_versions'].append(formatted_versions[version]['version'])

        versions_dict['available_versions'].sort(key=LooseVersion)
        versions_dict['latest'] = versions_dict['available_versions'][-1]

        # finish up the matched version
        if self.match_version:
            if versions_dict['matched_versions_full']:
                versions_dict['matched_version_found'] = True
                versions_dict['matched_versions'].sort(key=LooseVersion)
                versions_dict['matched_version_latest'] = versions_dict['matched_versions'][-1]
                versions_dict['matched_version_full_latest'] = versions_dict['matched_versions_full'][-1]
            else:
                versions_dict['matched_version_found'] = False
                versions_dict['matched_versions'] = []
                versions_dict['matched_version_latest'] = ""
                versions_dict['matched_version_full_latest'] = ""

        return versions_dict

    def repoquery(self):
        '''perform a repoquery '''

        if self.ignore_excluders:
            # Duplicate yum.conf and reset exclude= line to an empty string
            # to clear a list of all excluded packages
            self.tmp_file = tempfile.NamedTemporaryFile()

            with open("/etc/yum.conf", "r") as file_handler:
                yum_conf_lines = file_handler.readlines()

            yum_conf_lines = [l for l in yum_conf_lines if not l.startswith("exclude=")]

            with open(self.tmp_file.name, "w") as file_handler:
                file_handler.writelines(yum_conf_lines)
                file_handler.flush()

        repoquery_cmd = self.build_cmd()

        rval = self._repoquery_cmd(repoquery_cmd, True, 'raw')

        # check to see if there are actual results
        rval['package_name'] = self.name
        if rval['results']:
            processed_versions = Repoquery.process_versions(rval['results'].strip())
            formatted_versions = self.format_versions(processed_versions)

            rval['package_found'] = True
            rval['versions'] = formatted_versions

            if self.verbose:
                rval['raw_versions'] = processed_versions
            else:
                del rval['results']

        # No packages found
        else:
            rval['package_found'] = False

        if self.ignore_excluders:
            self.tmp_file.close()

        return rval

    @staticmethod
    def run_ansible(params, check_mode):
        '''run the ansible idempotent code'''

        repoquery = Repoquery(
            params['name'],
            params['query_type'],
            params['show_duplicates'],
            params['match_version'],
            params['ignore_excluders'],
            params['verbose'],
        )

        state = params['state']

        if state == 'list':
            results = repoquery.repoquery()

            if results['returncode'] != 0:
                return {'failed': True,
                        'msg': results}

            return {'changed': False, 'results': results, 'state': 'list', 'check_mode': check_mode}

        return {'failed': True,
                'changed': False,
                'msg': 'Unknown state passed. %s' % state,
                'state': 'unknown'}

# -*- -*- -*- End included fragment: class/repoquery.py -*- -*- -*-

# -*- -*- -*- Begin included fragment: ansible/repoquery.py -*- -*- -*-


def main():
    '''
    ansible repoquery module
    '''
    module = AnsibleModule(
        argument_spec=dict(
            state=dict(default='list', type='str', choices=['list']),
            name=dict(default=None, required=True, type='str'),
            query_type=dict(default='repos', required=False, type='str',
                            choices=[
                                'installed', 'available', 'recent',
                                'updates', 'extras', 'all', 'repos'
                            ]),
            verbose=dict(default=False, required=False, type='bool'),
            show_duplicates=dict(default=False, required=False, type='bool'),
            match_version=dict(default=None, required=False, type='str'),
            ignore_excluders=dict(default=False, required=False, type='bool'),
            retries=dict(default=4, required=False, type='int'),
            retry_interval=dict(default=5, required=False, type='int'),
        ),
        supports_check_mode=False,
        required_if=[('show_duplicates', True, ['name'])],
    )

    tries = 1
    while True:
        rval = Repoquery.run_ansible(module.params, module.check_mode)
        if 'failed' not in rval:
            module.exit_json(**rval)
        elif tries > module.params['retries']:
            module.fail_json(**rval)
        tries += 1
        time.sleep(module.params['retry_interval'])


if __name__ == "__main__":
    main()

# -*- -*- -*- End included fragment: ansible/repoquery.py -*- -*- -*-
