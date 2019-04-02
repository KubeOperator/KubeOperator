"""
Ansible action plugin to help with setting facts conditionally based on other facts.
"""

from ansible.plugins.action import ActionBase


DOCUMENTATION = '''
---
action_plugin: conditional_set_fact

short_description: This will set a fact if the value is defined

description:
    - "To avoid constant set_fact & when conditions for each var we can use this"

author:
    - Eric Wolinetz ewolinet@redhat.com
'''


EXAMPLES = '''
- name: Conditionally set fact
  conditional_set_fact:
    fact1: not_defined_variable

- name: Conditionally set fact
  conditional_set_fact:
    fact1: not_defined_variable
    fact2: defined_variable

- name: Conditionally set fact falling back on default
  conditional_set_fact:
    fact1: not_defined_var | defined_variable

'''


# pylint: disable=too-few-public-methods
class ActionModule(ActionBase):
    """Action plugin to execute deprecated var checks."""

    def run(self, tmp=None, task_vars=None):
        result = super(ActionModule, self).run(tmp, task_vars)
        result['changed'] = False

        facts = self._task.args.get('facts', [])
        var_list = self._task.args.get('vars', [])

        local_facts = dict()

        for param in var_list:
            other_vars = var_list[param].replace(" ", "")

            for other_var in other_vars.split('|'):
                if other_var in facts:
                    local_facts[param] = facts[other_var]
                    break

        if local_facts:
            result['changed'] = True

        result['ansible_facts'] = local_facts

        return result
