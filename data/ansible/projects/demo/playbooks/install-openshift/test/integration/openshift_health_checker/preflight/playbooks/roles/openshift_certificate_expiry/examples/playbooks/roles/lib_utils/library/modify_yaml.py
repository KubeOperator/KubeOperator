#!/usr/bin/python
# -*- coding: utf-8 -*-
''' modify_yaml ansible module '''

import yaml

# ignore pylint errors related to the module_utils import
# pylint: disable=redefined-builtin, unused-wildcard-import, wildcard-import
from ansible.module_utils.basic import *  # noqa: F402,F403


DOCUMENTATION = '''
---
module: modify_yaml
short_description: Modify yaml key value pairs
author: Andrew Butcher
requirements: [ ]
'''
EXAMPLES = '''
- modify_yaml:
    dest: /etc/origin/master/master-config.yaml
    yaml_key: 'kubernetesMasterConfig.masterCount'
    yaml_value: 2
'''


def set_key(yaml_data, yaml_key, yaml_value):
    ''' Updates a parsed yaml structure setting a key to a value.

        :param yaml_data: yaml structure to modify.
        :type yaml_data: dict
        :param yaml_key: Key to modify.
        :type yaml_key: mixed
        :param yaml_value: Value use for yaml_key.
        :type yaml_value: mixed
        :returns: Changes to the yaml_data structure
        :rtype: dict(tuple())
    '''
    changes = []
    ptr = yaml_data
    final_key = yaml_key.split('.')[-1]
    for key in yaml_key.split('.'):
        # Key isn't present and we're not on the final key. Set to empty dictionary.
        if key not in ptr and key != final_key:
            ptr[key] = {}
            ptr = ptr[key]
        # Current key is the final key. Update value.
        elif key == final_key:
            if (key in ptr and module.safe_eval(ptr[key]) != yaml_value) or (key not in ptr):  # noqa: F405
                ptr[key] = yaml_value
                changes.append((yaml_key, yaml_value))
        else:
            # Next value is None and we're not on the final key.
            # Turn value into an empty dictionary.
            if ptr[key] is None and key != final_key:
                ptr[key] = {}
            ptr = ptr[key]
    return changes


def main():
    ''' Modify key (supplied in jinja2 dot notation) in yaml file, setting
        the key to the desired value.
    '''

    # disabling pylint errors for global-variable-undefined and invalid-name
    # for 'global module' usage, since it is required to use ansible_facts
    # pylint: disable=global-variable-undefined, invalid-name,
    # redefined-outer-name
    global module

    module = AnsibleModule(  # noqa: F405
        argument_spec=dict(
            dest=dict(required=True),
            yaml_key=dict(required=True),
            yaml_value=dict(required=True),
            backup=dict(required=False, default=True, type='bool'),
        ),
        supports_check_mode=True,
    )

    dest = module.params['dest']
    yaml_key = module.params['yaml_key']
    yaml_value = module.safe_eval(module.params['yaml_value'])
    backup = module.params['backup']

    # Represent null values as an empty string.
    # pylint: disable=missing-docstring, unused-argument
    def none_representer(dumper, data):
        return yaml.ScalarNode(tag=u'tag:yaml.org,2002:null', value=u'')

    yaml.add_representer(type(None), none_representer)

    try:
        with open(dest) as yaml_file:
            yaml_data = yaml.safe_load(yaml_file.read())

        changes = set_key(yaml_data, yaml_key, yaml_value)

        if len(changes) > 0:
            if backup:
                module.backup_local(dest)
            with open(dest, 'w') as yaml_file:
                yaml_string = yaml.dump(yaml_data, default_flow_style=False)
                yaml_string = yaml_string.replace('\'\'', '""')
                yaml_file.write(yaml_string)

        return module.exit_json(changed=(len(changes) > 0), changes=changes)

    # ignore broad-except error to avoid stack trace to ansible user
    # pylint: disable=broad-except
    except Exception as error:
        return module.fail_json(msg=str(error))


if __name__ == '__main__':
    main()
