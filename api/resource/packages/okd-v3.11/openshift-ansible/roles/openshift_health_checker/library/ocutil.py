#!/usr/bin/python
"""Interface to OpenShift oc command"""

import os
import shlex
import shutil
import subprocess

from ansible.module_utils.basic import AnsibleModule


ADDITIONAL_PATH_LOOKUPS = ['/usr/local/bin', os.path.expanduser('~/bin')]


def locate_oc_binary():
    """Find and return oc binary file"""
    # https://github.com/openshift/openshift-ansible/issues/3410
    # oc can be in /usr/local/bin in some cases, but that may not
    # be in $PATH due to ansible/sudo
    paths = os.environ.get("PATH", os.defpath).split(os.pathsep) + ADDITIONAL_PATH_LOOKUPS

    oc_binary = 'oc'

    # Use shutil.which if it is available, otherwise fallback to a naive path search
    try:
        which_result = shutil.which(oc_binary, path=os.pathsep.join(paths))
        if which_result is not None:
            oc_binary = which_result
    except AttributeError:
        for path in paths:
            if os.path.exists(os.path.join(path, oc_binary)):
                oc_binary = os.path.join(path, oc_binary)
                break

    return oc_binary


def main():
    """Module that executes commands on a remote OpenShift cluster"""

    module = AnsibleModule(
        argument_spec=dict(
            namespace=dict(type="str", required=False),
            config_file=dict(type="str", required=True),
            cmd=dict(type="str", required=True),
            extra_args=dict(type="list", default=[]),
        ),
    )

    cmd = [locate_oc_binary(), '--config', module.params["config_file"]]
    if module.params["namespace"]:
        cmd += ['-n', module.params["namespace"]]
    cmd += shlex.split(module.params["cmd"]) + module.params["extra_args"]

    failed = True
    try:
        cmd_result = subprocess.check_output(list(cmd), stderr=subprocess.STDOUT)
        failed = False
    except subprocess.CalledProcessError as exc:
        cmd_result = '[rc {}] {}\n{}'.format(exc.returncode, ' '.join(exc.cmd), exc.output)
    except OSError as exc:
        # we get this when 'oc' is not there
        cmd_result = str(exc)

    module.exit_json(
        changed=False,
        failed=failed,
        result=cmd_result,
    )


if __name__ == '__main__':
    main()
