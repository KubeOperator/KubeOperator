"""
Ansible action plugin to check for usage of deprecated variables in Openshift Ansible inventory files.
"""

from ansible.plugins.action import ActionBase


# pylint: disable=too-few-public-methods
class ActionModule(ActionBase):
    """Action plugin to execute deprecated var checks."""

    def run(self, tmp=None, task_vars=None):
        result = super(ActionModule, self).run(tmp, task_vars)

        # pylint: disable=line-too-long
        deprecation_header = "[DEPRECATION WARNING]: The following are deprecated variables and will be no longer be used in the next minor release. Please update your inventory accordingly."

        facts = self._task.args.get('facts', [])
        dep_var_list = self._task.args.get('vars', [])
        dep_header = self._task.args.get('header', deprecation_header)

        deprecation_message = "No deprecations found"
        is_changed = False

        for param in dep_var_list:
            if param in facts:
                if not is_changed:
                    deprecation_message = dep_header
                    is_changed = True

                deprecation_message = deprecation_message + "\n\t" + param

        result['changed'] = is_changed
        result['msg'] = deprecation_message

        return result
