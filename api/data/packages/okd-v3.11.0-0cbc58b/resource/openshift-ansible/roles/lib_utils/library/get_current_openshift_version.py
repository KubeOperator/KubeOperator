#!/usr/bin/env python
# pylint: disable=missing-docstring
#
# Copyright 2017 Red Hat, Inc. and/or its affiliates
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

import os

from ansible.module_utils.basic import AnsibleModule


DOCUMENTATION = '''
---
module: get_current_openshift_version

short_description: Discovers installed openshift version on masters and nodes

version_added: "2.4"

description:
    - This module checks various files and program outputs to get the
      currently installed openshfit version

options:
    deployment_type:
        description:
            - openshift_deployment_type
        required: true


author:
    - "Michael Gugino <mgugino@redhat.com>"
'''

EXAMPLES = '''
- name: Set openshift_current_version
  get_current_openshift_version:
    deployment_type: openshift_deployment_type
'''


def chomp_commit_offset(version):
    """Chomp any "+git.foo" commit offset string from the given `version`
    and return the modified version string.

Ex:
- chomp_commit_offset(None)                 => None
- chomp_commit_offset(1337)                 => "1337"
- chomp_commit_offset("v3.4.0.15+git.derp") => "v3.4.0.15"
- chomp_commit_offset("v3.4.0.15")          => "v3.4.0.15"
- chomp_commit_offset("v1.3.0+52492b4")     => "v1.3.0"
    """
    if version is None:
        return version
    else:
        # Stringify, just in case it's a Number type. Split by '+' and
        # return the first split. No concerns about strings without a
        # '+', .split() returns an array of the original string.
        return str(version).split('+')[0]


def get_container_openshift_version(deployment_type):
    """
    If containerized, see if we can determine the installed version via the
    systemd environment files.
    """
    service_type_dict = {'origin': 'origin',
                         'openshift-enterprise': 'atomic-openshift'}
    service_type = service_type_dict[deployment_type]

    for filename in ['/etc/sysconfig/%s-master-controllers', '/etc/sysconfig/%s-node']:
        env_path = filename % service_type
        if not os.path.exists(env_path):
            continue

        with open(env_path) as env_file:
            for line in env_file:
                if line.startswith("IMAGE_VERSION="):
                    tag = line[len("IMAGE_VERSION="):].strip()
                    # Remove leading "v" and any trailing release info, we just want
                    # a version number here:
                    no_v_version = tag[1:] if tag[0] == 'v' else tag
                    version = no_v_version.split("-")[0]
                    return version
    return None


def parse_openshift_version(output):
    """ Apply provider facts to supplied facts dict

        Args:
            string: output of 'openshift version'
        Returns:
            string: the version number
    """
    versions = dict(e.split(' v') for e in output.splitlines() if ' v' in e)
    ver = versions.get('openshift', '')
    # Remove trailing build number and commit hash from older versions, we need to return a straight
    # w.x.y.z version here for use as openshift_version throughout the playbooks/roles. (i.e. 3.1.1.6-64-g80b61da)
    ver = ver.split('-')[0]
    return ver


def get_openshift_version(module, deployment_type):
    """ Get current version of openshift on the host.

        Checks a variety of ways ranging from fastest to slowest.

        Args:
            facts (dict): existing facts

        Returns:
            version: the current openshift version
    """
    version = None

    if os.path.isfile('/usr/bin/openshift'):
        _, output, _ = module.run_command(['/usr/bin/openshift', 'version'])  # noqa: F405
        version = parse_openshift_version(output)
    else:
        version = get_container_openshift_version(deployment_type)

    return chomp_commit_offset(version)


def run_module():
    '''Run this module'''
    module_args = dict(
        deployment_type=dict(type='str', required=True)
    )

    module = AnsibleModule(
        argument_spec=module_args,
        supports_check_mode=False
    )

    # First, create our dest dir if necessary
    deployment_type = module.params['deployment_type']
    changed = False
    ansible_facts = {}

    current_version = get_openshift_version(module, deployment_type)
    if current_version is not None:
        ansible_facts = {'openshift_current_version': current_version}

    # Passing back ansible_facts will set_fact the values.
    result = {'changed': changed, 'ansible_facts': ansible_facts}

    module.exit_json(**result)


def main():
    run_module()


if __name__ == '__main__':
    main()
