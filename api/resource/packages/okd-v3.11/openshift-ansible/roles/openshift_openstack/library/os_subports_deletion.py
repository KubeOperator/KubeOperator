#!/usr/bin/python
# -*- coding: utf-8 -*-

# Copyright 2018 Red Hat, Inc. and/or its affiliates
# and other contributors as indicated by the @author tags.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# pylint: disable=unused-wildcard-import,wildcard-import,unused-import,redefined-builtin

''' os_subports_deletion '''
import keystoneauth1

from ansible.module_utils.basic import AnsibleModule

try:
    import shade
    HAS_SHADE = True
except ImportError:
    HAS_SHADE = False

DOCUMENTATION = '''
---
module: os_subports_deletion
short_description: Delete subports belonging to a trunk
description:
    - Detach and deletes all the Neutron Subports belonging to a trunk
author:
    - "Luis Tomas Bolivar <ltomasbo@redhat.com>"
'''

RETURN = '''
'''


def main():
    ''' Main module function '''
    module = AnsibleModule(
        argument_spec=dict(
            trunk_name=dict(default=False, type='str'),
        ),
        supports_check_mode=True,
    )

    if not HAS_SHADE:
        module.fail_json(msg='shade is required for this module')

    try:
        cloud = shade.openstack_cloud()
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to connect to the cloud')

    try:
        adapter = keystoneauth1.adapter.Adapter(
            session=cloud.keystone_session,
            service_type=cloud.cloud_config.get_service_type('network'),
            interface=cloud.cloud_config.get_interface('network'),
            endpoint_override=cloud.cloud_config.get_endpoint('network'),
            version=cloud.cloud_config.get_api_version('network'))
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to get an adapter to talk to the Neutron '
                             'API')

    try:
        trunk_response = adapter.get('/trunks')
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to retrieve Neutron trunk information')

    subports = []
    for trunk in trunk_response.json()['trunks']:
        if trunk['name'] == module.params['trunk_name']:
            trunk_id = trunk['id']
            for subport in trunk['sub_ports']:
                subports.append(subport['port_id'])

    data = _get_data(subports)
    try:
        adapter.put('/trunks/' + trunk_id + '/remove_subports',
                    **data)
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to detach subports')

    try:
        for port in subports:
            adapter.delete('/ports/' + port)
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to delete Neutron subports')

    module.exit_json(
        changed=True)


def _get_data(subports):
    ports_list = [{"port_id": port_id.encode('ascii')} for port_id in subports]
    sub_ports = str({"sub_ports": ports_list}).replace('\'', '\"')
    return {'data': sub_ports}


if __name__ == '__main__':
    main()
