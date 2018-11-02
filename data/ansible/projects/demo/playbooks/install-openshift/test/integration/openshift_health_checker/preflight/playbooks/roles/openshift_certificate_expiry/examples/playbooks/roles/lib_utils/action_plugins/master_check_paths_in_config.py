"""
Ansible action plugin to ensure inventory variables are set
appropriately and no conflicting options have been provided.
"""
import collections
import six

from ansible.plugins.action import ActionBase
from ansible import errors


FAIL_MSG = """A string value that appears to be a file path located outside of
{} has been found in /etc/origin/master/master-config.yaml.
In 3.10 and newer, all files needed by the master must reside inside of
those directories or a subdirectory or it will not be readable by the
master process. Please migrate all files needed by the master into
one of {} or a subdirectory and update your master configs before
proceeding. The string found was: {}
***********************
NOTE: the following items do not need to be migrated, they will be migrated
for you: {}"""


ITEMS_TO_POP = (
    ('oauthConfig', 'identityProviders'),
)
# Create csv string of dot-separated dictionary keys:
# eg: 'oathConfig.identityProviders, something.else.here'
MIGRATED_ITEMS = ", ".join([".".join(x) for x in ITEMS_TO_POP])

ALLOWED_DIRS = (
    '/etc/origin/master/',
    '/var/lib/origin',
    '/etc/origin/cloudprovider',
    '/etc/origin/kubelet-plugins',
    '/usr/libexec/kubernetes/kubelet-plugins',
)

ALLOWED_DIRS_STRING = ', '.join(ALLOWED_DIRS)


def pop_migrated_fields(mastercfg):
    """Some fields do not need to be searched because they will be migrated
    for users automatically"""
    # Walk down the tree and pop the specific item we migrate / don't care about
    for item in ITEMS_TO_POP:
        field = mastercfg
        for sub_field in item:
            parent_field = field
            field = field[sub_field]
        parent_field.pop(item[len(item) - 1])


def do_item_check(val, strings_to_check):
    """Check type of val, append to strings_to_check if string, otherwise if
    it's a dictionary-like object call walk_mapping, if it's a list-like
    object call walk_sequence, else ignore."""
    if isinstance(val, six.string_types):
        strings_to_check.append(val)
    elif isinstance(val, collections.Sequence):
        # A list-like object
        walk_sequence(val, strings_to_check)
    elif isinstance(val, collections.Mapping):
        # A dictionary-like object
        walk_mapping(val, strings_to_check)
    # If it's not a string, list, or dictionary, we're not interested.


def walk_sequence(items, strings_to_check):
    """Walk recursively through a list, items"""
    for item in items:
        do_item_check(item, strings_to_check)


def walk_mapping(map_to_walk, strings_to_check):
    """Walk recursively through map_to_walk dictionary and add strings to
    strings_to_check"""
    for _, val in map_to_walk.items():
        do_item_check(val, strings_to_check)


def check_strings(strings_to_check):
    """Check the strings we found to see if they look like file paths and if
    they are, fail if not start with /etc/origin/master"""
    for item in strings_to_check:
        if item.startswith('/') or item.startswith('../'):
            matches = 0
            for allowed in ALLOWED_DIRS:
                if item.startswith(allowed):
                    matches += 1
            if matches == 0:
                raise errors.AnsibleModuleError(
                    FAIL_MSG.format(ALLOWED_DIRS_STRING,
                                    ALLOWED_DIRS_STRING,
                                    item, MIGRATED_ITEMS))


# pylint: disable=R0903
class ActionModule(ActionBase):
    """Action plugin to validate no files are needed by master that reside
    outside of /etc/origin/master as masters will now run as pods and cannot
    utilize files outside of that path as they will not be mounted inside the
    containers."""
    def run(self, tmp=None, task_vars=None):
        """Run this action module"""
        result = super(ActionModule, self).run(tmp, task_vars)

        # self.task_vars holds all in-scope variables.
        # Ignore settting self.task_vars outside of init.
        # pylint: disable=W0201
        self.task_vars = task_vars or {}

        # mastercfg should be a dictionary from scraping an existing master's
        # config yaml file.
        mastercfg = self._task.args.get('mastercfg')

        # We migrate some paths for users automatically, so we pop those.
        pop_migrated_fields(mastercfg)

        # Create an empty list to append strings from our config file to to check
        # later.
        strings_to_check = []

        walk_mapping(mastercfg, strings_to_check)

        check_strings(strings_to_check)

        result["changed"] = False
        result["failed"] = False
        result["msg"] = "Aight, configs looking good"
        return result
