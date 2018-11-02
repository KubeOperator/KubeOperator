#!/usr/bin/python
# -*- coding: utf-8 -*-

# (c) 2015, Tobias Florek <tob@butter.sh>
# Licensed under the terms of the MIT License
"""
An ansible module to query the RPM database. For use, when yum/dnf are not
available.
"""

# pylint: disable=redefined-builtin,wildcard-import,unused-wildcard-import
from ansible.module_utils.basic import *  # noqa: F403

DOCUMENTATION = """
---
module: rpm_q
short_description: Query the RPM database
author: Tobias Florek
options:
  name:
    description:
    - The name of the package to query
    required: true
  state:
    description:
    - Whether the package is supposed to be installed or not
    choices: [present, absent]
    default: present
"""

EXAMPLES = """
- rpm_q: name=ansible state=present
- rpm_q: name=ansible state=absent
"""

RPM_BINARY = '/bin/rpm'


def main():
    """
    Checks rpm -q for the named package and returns the installed packages
    or None if not installed.
    """
    module = AnsibleModule(  # noqa: F405
        argument_spec=dict(
            name=dict(required=True),
            state=dict(default='present', choices=['present', 'absent'])
        ),
        supports_check_mode=True
    )

    name = module.params['name']
    state = module.params['state']

    # pylint: disable=invalid-name
    rc, out, err = module.run_command([RPM_BINARY, '-q', name])

    installed = out.rstrip('\n').split('\n')

    if rc != 0:
        if state == 'present':
            module.fail_json(msg="%s is not installed" % name, stdout=out, stderr=err, rc=rc)
        else:
            module.exit_json(changed=False)
    elif state == 'present':
        module.exit_json(changed=False, installed_versions=installed)
    else:
        module.fail_json(msg="%s is installed", installed_versions=installed)


if __name__ == '__main__':
    main()
