#!/usr/bin/python
# -*- coding: utf-8 -*-
''' kubeclient_ca ansible module '''

import base64
import yaml
from ansible.module_utils.basic import AnsibleModule


DOCUMENTATION = '''
---
module: kubeclient_ca
short_description: Modify kubeclient certificate-authority-data
author: Andrew Butcher
requirements: [ ]
'''
EXAMPLES = '''
- kubeclient_ca:
    client_path: /etc/origin/master/admin.kubeconfig
    ca_path: /etc/origin/master/ca-bundle.crt

- slurp:
    src: /etc/origin/master/ca-bundle.crt
  register: ca_data
- kubeclient_ca:
    client_path: /etc/origin/master/admin.kubeconfig
    ca_data: "{{ ca_data.content }}"
'''


def main():
    ''' Modify kubeconfig located at `client_path`, setting the
        certificate authority data to specified `ca_data` or contents of
        `ca_path`.
    '''

    module = AnsibleModule(  # noqa: F405
        argument_spec=dict(
            client_path=dict(required=True),
            ca_data=dict(required=False, default=None),
            ca_path=dict(required=False, default=None),
            backup=dict(required=False, default=True, type='bool'),
        ),
        supports_check_mode=True,
        mutually_exclusive=[['ca_data', 'ca_path']],
        required_one_of=[['ca_data', 'ca_path']]
    )

    client_path = module.params['client_path']
    ca_data = module.params['ca_data']
    ca_path = module.params['ca_path']
    backup = module.params['backup']

    try:
        with open(client_path) as client_config_file:
            client_config_data = yaml.safe_load(client_config_file.read())

        if ca_data is None:
            with open(ca_path) as ca_file:
                ca_data = base64.standard_b64encode(ca_file.read())

        changes = []
        # Naively update the CA information for each cluster in the
        # kubeconfig.
        for cluster in client_config_data['clusters']:
            if cluster['cluster']['certificate-authority-data'] != ca_data:
                cluster['cluster']['certificate-authority-data'] = ca_data
                changes.append(cluster['name'])

        if not module.check_mode:
            if len(changes) > 0 and backup:
                module.backup_local(client_path)

            with open(client_path, 'w') as client_config_file:
                client_config_string = yaml.dump(client_config_data, default_flow_style=False)
                client_config_string = client_config_string.replace('\'\'', '""')
                client_config_file.write(client_config_string)

        return module.exit_json(changed=(len(changes) > 0))

    # ignore broad-except error to avoid stack trace to ansible user
    # pylint: disable=broad-except
    except Exception as error:
        return module.fail_json(msg=str(error))


if __name__ == '__main__':
    main()
