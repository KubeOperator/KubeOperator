"""
Health checks for OpenShift clusters.
"""

import json
import operator
import os
import re
import time
import collections

from abc import ABCMeta, abstractmethod, abstractproperty
from importlib import import_module

from ansible.module_utils import six
from ansible.module_utils.six.moves import reduce  # pylint: disable=import-error,redefined-builtin
from ansible.module_utils.six import string_types
from ansible.plugins.filter.core import to_bool as ansible_to_bool


class OpenShiftCheckException(Exception):
    """Raised when a check encounters a failure condition."""

    def __init__(self, name, msg=None):
        # msg is for the message the user will see when this is raised.
        # name is for test code to identify the error without looking at msg text.
        if msg is None:  # for parameter backward compatibility
            msg = name
            name = self.__class__.__name__
        self.name = name
        super(OpenShiftCheckException, self).__init__(msg)


class OpenShiftCheckExceptionList(OpenShiftCheckException):
    """A container for multiple errors that may be detected in one check."""
    def __init__(self, errors):
        self.errors = errors
        super(OpenShiftCheckExceptionList, self).__init__(
            'OpenShiftCheckExceptionList',
            '\n'.join(str(msg) for msg in errors)
        )

    # make iterable
    def __getitem__(self, index):
        return self.errors[index]


FileToSave = collections.namedtuple("FileToSave", "filename contents remote_filename")


# pylint: disable=too-many-instance-attributes; all represent significantly different state.
# Arguably they could be separated into two hashes, one for storing parameters, and one for
# storing result state; but that smells more like clutter than clarity.
@six.add_metaclass(ABCMeta)
class OpenShiftCheck(object):
    """A base class for defining checks for an OpenShift cluster environment.

    Optional init params: method execute_module, dict task_vars, and string tmp
    execute_module is expected to have a signature compatible with _execute_module
    from ansible plugins/action/__init__.py, e.g.:
    def execute_module(module_name=None, module_args=None, tmp=None, task_vars=None, *args):
    This is stored so that it can be invoked in subclasses via check.execute_module("name", args)
    which provides the check's stored task_vars and tmp.

    Optional init param: want_full_results
    If the check can gather logs, tarballs, etc., do so when True; but no need to spend
    the time if they're not wanted (won't be written to output directory).
    """
    # pylint: disable=too-many-arguments
    def __init__(self, execute_module=None, task_vars=None, tmp=None, want_full_results=False,
                 templar=None):
        # store a method for executing ansible modules from the check
        self._execute_module = execute_module
        # the task variables and tmpdir passed into the health checker task
        self.task_vars = task_vars or {}
        # We may need to template some task_vars
        self._templar = templar
        self.tmp = tmp
        # a boolean for disabling the gathering of results (files, computations) that won't
        # actually be recorded/used
        self.want_full_results = want_full_results

        # mainly for testing purposes; see execute_module_with_retries
        self._module_retries = 3
        self._module_retry_interval = 5  # seconds

        # state to be recorded for inspection after the check runs:
        #
        # set to True when the check changes the host, for accurate total "changed" count
        self.changed = False
        # list of OpenShiftCheckException for check to report (alternative to returning a failed result)
        self.failures = []
        # list of FileToSave - files the check specifies to be written locally if so configured
        self.files_to_save = []
        # log messages for the check - tuples of (description, msg) where msg is serializable.
        # These are intended to be a sequential record of what the check observed and determined.
        self.logs = []

    def template_var(self, var_to_template):
        """Return a templated variable if self._templar is not None, else
           just return the variable as-is"""
        if self._templar is not None:
            return self._templar.template(var_to_template)
        return var_to_template

    @abstractproperty
    def name(self):
        """The name of this check, usually derived from the class name."""
        return "openshift_check"

    @property
    def tags(self):
        """A list of tags that this check satisfy.

        Tags are used to reference multiple checks with a single '@tagname'
        special check name.
        """
        return []

    @staticmethod
    def is_active():
        """Returns true if this check applies to the ansible-playbook run."""
        return True

    def is_first_master(self):
        """Determine if running on first master. Returns: bool"""
        masters = self.get_var("groups", "oo_first_master", default=None) or [None]
        return masters[0] == self.get_var("ansible_host")

    @abstractmethod
    def run(self):
        """Executes a check against a host and returns a result hash similar to Ansible modules.

        Actually the direction ahead is to record state in the attributes and
        not bother building a result hash. Instead, return an empty hash and let
        the action plugin fill it in. Or raise an OpenShiftCheckException.
        Returning a hash may become deprecated if it does not prove necessary.
        """
        return {}

    @classmethod
    def subclasses(cls):
        """Returns a generator of subclasses of this class and its subclasses."""
        # AUDIT: no-member makes sense due to this having a metaclass
        for subclass in cls.__subclasses__():  # pylint: disable=no-member
            yield subclass
            for subclass in subclass.subclasses():
                yield subclass

    def register_failure(self, error):
        """Record in the check that a failure occurred.

        Recorded failures are merged into the result hash for now. They are also saved to output directory
        (if provided) <check>.failures.json and registered as a log entry for context <check>.log.json.
        """
        # It should be an exception; make it one if not
        if not isinstance(error, OpenShiftCheckException):
            error = OpenShiftCheckException(str(error))
        self.failures.append(error)
        # duplicate it in the logs so it can be seen in the context of any
        # information that led to the failure
        self.register_log("failure: " + error.name, str(error))

    def register_log(self, context, msg):
        """Record an entry for the check log.

        Notes are intended to serve as context of the whole sequence of what the check observed.
        They are be saved as an ordered list in a local check log file.
        They are not to included in the result or in the ansible log; it's just for the record.
        """
        self.logs.append([context, msg])

    def register_file(self, filename, contents=None, remote_filename=""):
        """Record a file that a check makes available to be saved individually to output directory.

        Either file contents should be passed in, or a file to be copied from the remote host
        should be specified. Contents that are not a string are to be serialized as JSON.

        NOTE: When copying a file from remote host, it is slurped into memory as base64, meaning
        you should avoid using this on huge files (more than say 10M).
        """
        if contents is None and not remote_filename:
            raise OpenShiftCheckException("File data/source not specified; this is a bug in the check.")
        self.files_to_save.append(FileToSave(filename, contents, remote_filename))

    def execute_module(self, module_name=None, module_args=None, save_as_name=None, register=True):
        """Invoke an Ansible module from a check.

        Invoke stored _execute_module, normally copied from the action
        plugin, with its params and the task_vars and tmp given at
        check initialization. No positional parameters beyond these
        are specified. If it's necessary to specify any of the other
        parameters to _execute_module then that should just be invoked
        directly (with awareness of changes in method signature per
        Ansible version).

        So e.g. check.execute_module("foo", dict(arg1=...))

        save_as_name specifies a file name for saving the result to an output directory,
        if needed, and is intended to uniquely identify the result of invoking execute_module.
        If not provided, the module name will be used.
        If register is set False, then the result won't be registered in logs or files to save.

        Return: result hash from module execution.
        """
        if self._execute_module is None:
            raise NotImplementedError(
                self.__class__.__name__ +
                " invoked execute_module without providing the method at initialization."
            )
        result = self._execute_module(module_name, module_args, self.tmp, self.task_vars)
        if result.get("changed"):
            self.changed = True
        for output in ["result", "stdout"]:
            # output is often JSON; attempt to decode
            try:
                result[output + "_json"] = json.loads(result[output])
            except (KeyError, ValueError):
                pass

        if register:
            self.register_log("execute_module: " + module_name, result)
            self.register_file(save_as_name or module_name + ".json", result)
        return result

    def execute_module_with_retries(self, module_name, module_args):
        """Run execute_module and retry on failure."""
        result = {}
        tries = 0
        while True:
            res = self.execute_module(module_name, module_args)
            if tries > self._module_retries or not res.get("failed"):
                result.update(res)
                return result
            result["last_failed"] = res
            tries += 1
            time.sleep(self._module_retry_interval)

    def get_var(self, *keys, **kwargs):
        """Get deeply nested values from task_vars.

        Ansible task_vars structures are Python dicts, often mapping strings to
        other dicts. This helper makes it easier to get a nested value, raising
        OpenShiftCheckException when a key is not found.

        Keyword args:
          default:
            On missing key, return this as default value instead of raising exception.
          convert:
            Supply a function to apply to normalize the value before returning it.
            None is the default (return as-is).
            This function should raise ValueError if the user has provided a value
            that cannot be converted, or OpenShiftCheckException if some other
            problem needs to be described to the user.
        """
        if len(keys) == 1:
            keys = keys[0].split(".")

        try:
            value = reduce(operator.getitem, keys, self.task_vars)
        except (KeyError, TypeError):
            if "default" not in kwargs:
                raise OpenShiftCheckException(
                    "This check expects the '{}' inventory variable to be defined\n"
                    "in order to proceed, but it is undefined. There may be a bug\n"
                    "in Ansible, the checks, or their dependencies."
                    "".format(".".join(map(str, keys)))
                )
            value = kwargs["default"]

        convert = kwargs.get("convert", None)
        try:
            if convert is None:
                return value
            elif convert is bool:  # interpret bool as Ansible does, instead of python truthiness
                return ansible_to_bool(value)
            else:
                return convert(value)

        except ValueError as error:  # user error in specifying value
            raise OpenShiftCheckException(
                'Cannot convert inventory variable to expected type:\n'
                '  "{var}={value}"\n'
                '{error}'.format(var=".".join(keys), value=value, error=error)
            )

        except OpenShiftCheckException:  # some other check-specific problem
            raise

        except Exception as error:  # probably a bug in the function
            raise OpenShiftCheckException(
                'There is a bug in this check. While trying to convert variable \n'
                '  "{var}={value}"\n'
                'the given converter cannot be used or failed unexpectedly:\n'
                '{type}: {error}'.format(
                    var=".".join(keys),
                    value=value,
                    type=error.__class__.__name__,
                    error=error
                ))

    @staticmethod
    def normalize(name_list):
        """Return a clean list of names.

        The input may be a comma-separated string or a sequence. Leading and
        trailing whitespace characters are removed. Empty items are discarded.
        """
        if isinstance(name_list, string_types):
            name_list = name_list.split(',')
        return [name.strip() for name in name_list if name.strip()]

    def get_major_minor_version(self, openshift_image_tag=None):
        """Parse and return the deployed version of OpenShift as a tuple."""

        version = openshift_image_tag or self.get_var("openshift_image_tag")
        components = [int(component) for component in re.findall(r'\d+', version)]

        if len(components) < 2:
            msg = "An invalid version of OpenShift was found for this host: {}"
            raise OpenShiftCheckException(msg.format(version))

        # map major release version across releases to OCP major version
        components[0] = {1: 3}.get(components[0], components[0])

        return tuple(int(x) for x in components[:2])

    def get_required_version(self, name, version_map):
        """Return the correct required version(s) for the current (or nearest) OpenShift version."""
        openshift_version = self.get_major_minor_version()

        earliest = min(version_map)
        latest = max(version_map)
        if openshift_version < earliest:
            return version_map[earliest]
        if openshift_version > latest:
            return version_map[latest]

        required_version = version_map.get(openshift_version)
        if not required_version:
            msg = "There is no recommended version of {} for the current version of OpenShift ({})"
            raise OpenShiftCheckException(msg.format(name, ".".join(str(comp) for comp in openshift_version)))

        return required_version

    def find_ansible_mount(self, path):
        """Return the mount point for path from ansible_mounts."""

        # reorganize list of mounts into dict by path
        mount_for_path = {
            mount['mount']: mount
            for mount
            in self.get_var('ansible_mounts')
        }

        # NOTE: including base cases '/' and '' to ensure the loop ends
        mount_targets = set(mount_for_path.keys()) | {'/', ''}
        mount_point = path
        while mount_point not in mount_targets:
            mount_point = os.path.dirname(mount_point)

        try:
            mount = mount_for_path[mount_point]
            self.register_log("mount point for " + path, mount)
            return mount
        except KeyError:
            known_mounts = ', '.join('"{}"'.format(mount) for mount in sorted(mount_for_path))
            raise OpenShiftCheckException(
                'Unable to determine mount point for path "{}".\n'
                'Known mount points: {}.'.format(path, known_mounts or 'none')
            )


LOADER_EXCLUDES = (
    "__init__.py",
    "mixins.py",
    "logging.py",
)


def load_checks(path=None, subpkg=""):
    """Dynamically import all check modules for the side effect of registering checks."""
    if path is None:
        path = os.path.dirname(__file__)

    modules = []

    for name in os.listdir(path):
        if os.path.isdir(os.path.join(path, name)):
            modules = modules + load_checks(os.path.join(path, name), subpkg + "." + name)
            continue

        if name.endswith(".py") and name not in LOADER_EXCLUDES:
            modules.append(import_module(__package__ + subpkg + "." + name[:-3]))

    return modules
