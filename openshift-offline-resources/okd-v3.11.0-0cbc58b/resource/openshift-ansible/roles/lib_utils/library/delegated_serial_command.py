#!/usr/bin/python
# -*- coding: utf-8 -*-

# (c) 2012, Michael DeHaan <michael.dehaan@gmail.com>, and others
# (c) 2016, Andrew Butcher <abutcher@redhat.com>
#
# This module is derrived from the Ansible command module.
#
# Ansible is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Ansible is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Ansible.  If not, see <http://www.gnu.org/licenses/>.


# pylint: disable=unused-wildcard-import,wildcard-import,unused-import,redefined-builtin

''' delegated_serial_command '''

import datetime
import errno
import glob
import shlex
import os
import fcntl
import time

DOCUMENTATION = '''
---
module: delegated_serial_command
short_description: Executes a command on a remote node
version_added: historical
description:
     - The M(command) module takes the command name followed by a list
       of space-delimited arguments.
     - The given command will be executed on all selected nodes. It
       will not be processed through the shell, so variables like
       C($HOME) and operations like C("<"), C(">"), C("|"), and C("&")
       will not work (use the M(shell) module if you need these
       features).
     - Creates and maintains a lockfile such that this module will
       wait for other invocations to proceed.
options:
  command:
    description:
      - the command to run
    required: true
    default: null
  creates:
    description:
      - a filename or (since 2.0) glob pattern, when it already
        exists, this step will B(not) be run.
    required: no
    default: null
  removes:
    description:
      - a filename or (since 2.0) glob pattern, when it does not
        exist, this step will B(not) be run.
    version_added: "0.8"
    required: no
    default: null
  chdir:
    description:
      - cd into this directory before running the command
    version_added: "0.6"
    required: false
    default: null
  executable:
    description:
      - change the shell used to execute the command. Should be an
        absolute path to the executable.
    required: false
    default: null
    version_added: "0.9"
  warn:
    version_added: "1.8"
    default: yes
    description:
      - if command warnings are on in ansible.cfg, do not warn about
        this particular line if set to no/false.
    required: false
  lockfile:
    default: yes
    description:
      - the lockfile that will be created
  timeout:
    default: yes
    description:
      - time in milliseconds to wait to obtain the lock
notes:
    -  If you want to run a command through the shell (say you are using C(<),
       C(>), C(|), etc), you actually want the M(shell) module instead. The
       M(command) module is much more secure as it's not affected by the user's
       environment.
    - " C(creates), C(removes), and C(chdir) can be specified after
       the command. For instance, if you only want to run a command if
       a certain file does not exist, use this."
author:
    - Ansible Core Team
    - Michael DeHaan
    - Andrew Butcher
'''

EXAMPLES = '''
# Example from Ansible Playbooks.
- delegated_serial_command:
    command: /sbin/shutdown -t now

# Run the command if the specified file does not exist.
- delegated_serial_command:
    command: /usr/bin/make_database.sh arg1 arg2
    creates: /path/to/database
'''

# Dict of options and their defaults
OPTIONS = {'chdir': None,
           'creates': None,
           'command': None,
           'executable': None,
           'NO_LOG': None,
           'removes': None,
           'warn': True,
           'lockfile': None,
           'timeout': None}


def check_command(commandline):
    ''' Check provided command '''
    arguments = {'chown': 'owner', 'chmod': 'mode', 'chgrp': 'group',
                 'ln': 'state=link', 'mkdir': 'state=directory',
                 'rmdir': 'state=absent', 'rm': 'state=absent', 'touch': 'state=touch'}
    commands = {'git': 'git', 'hg': 'hg', 'curl': 'get_url or uri', 'wget': 'get_url or uri',
                'svn': 'subversion', 'service': 'service',
                'mount': 'mount', 'rpm': 'yum, dnf or zypper', 'yum': 'yum', 'apt-get': 'apt',
                'tar': 'unarchive', 'unzip': 'unarchive', 'sed': 'template or lineinfile',
                'rsync': 'synchronize', 'dnf': 'dnf', 'zypper': 'zypper'}
    become = ['sudo', 'su', 'pbrun', 'pfexec', 'runas']
    warnings = list()
    command = os.path.basename(commandline.split()[0])
    # pylint: disable=line-too-long
    if command in arguments:
        warnings.append("Consider using file module with {0} rather than running {1}".format(arguments[command], command))
    if command in commands:
        warnings.append("Consider using {0} module rather than running {1}".format(commands[command], command))
    if command in become:
        warnings.append(
            "Consider using 'become', 'become_method', and 'become_user' rather than running {0}".format(command,))
    return warnings


# pylint: disable=too-many-statements,too-many-branches,too-many-locals
def main():
    ''' Main module function '''
    module = AnsibleModule(  # noqa: F405
        argument_spec=dict(
            _uses_shell=dict(type='bool', default=False),
            command=dict(required=True),
            chdir=dict(),
            executable=dict(),
            creates=dict(),
            removes=dict(),
            warn=dict(type='bool', default=True),
            lockfile=dict(default='/tmp/delegated_serial_command.lock'),
            timeout=dict(type='int', default=30)
        )
    )

    shell = module.params['_uses_shell']
    chdir = module.params['chdir']
    executable = module.params['executable']
    command = module.params['command']
    creates = module.params['creates']
    removes = module.params['removes']
    warn = module.params['warn']
    lockfile = module.params['lockfile']
    timeout = module.params['timeout']

    if command.strip() == '':
        module.fail_json(rc=256, msg="no command given")

    iterated = 0
    lockfd = open(lockfile, 'w+')
    while iterated < timeout:
        try:
            fcntl.flock(lockfd, fcntl.LOCK_EX | fcntl.LOCK_NB)
            break
        # pylint: disable=invalid-name
        except IOError as e:
            if e.errno != errno.EAGAIN:
                module.fail_json(msg="I/O Error {0}: {1}".format(e.errno, e.strerror))
            else:
                iterated += 1
                time.sleep(0.1)

    if chdir:
        chdir = os.path.abspath(os.path.expanduser(chdir))
        os.chdir(chdir)

    if creates:
        # do not run the command if the line contains creates=filename
        # and the filename already exists.  This allows idempotence
        # of command executions.
        path = os.path.expanduser(creates)
        if glob.glob(path):
            module.exit_json(
                cmd=command,
                stdout="skipped, since %s exists" % path,
                changed=False,
                stderr=False,
                rc=0
            )

    if removes:
        # do not run the command if the line contains removes=filename
        # and the filename does not exist.  This allows idempotence
        # of command executions.
        path = os.path.expanduser(removes)
        if not glob.glob(path):
            module.exit_json(
                cmd=command,
                stdout="skipped, since %s does not exist" % path,
                changed=False,
                stderr=False,
                rc=0
            )

    warnings = list()
    if warn:
        warnings = check_command(command)

    if not shell:
        command = shlex.split(command)
    startd = datetime.datetime.now()

    # pylint: disable=invalid-name
    rc, out, err = module.run_command(command, executable=executable, use_unsafe_shell=shell)

    fcntl.flock(lockfd, fcntl.LOCK_UN)
    lockfd.close()

    endd = datetime.datetime.now()
    delta = endd - startd

    if out is None:
        out = ''
    if err is None:
        err = ''

    module.exit_json(
        cmd=command,
        stdout=out.rstrip("\r\n"),
        stderr=err.rstrip("\r\n"),
        rc=rc,
        start=str(startd),
        end=str(endd),
        delta=str(delta),
        changed=True,
        warnings=warnings,
        iterated=iterated
    )


# import module snippets
# pylint: disable=wrong-import-position
from ansible.module_utils.basic import *  # noqa: F402,F403

main()
