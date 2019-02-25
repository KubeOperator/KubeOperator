# pylint: skip-file
# flake8: noqa

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
