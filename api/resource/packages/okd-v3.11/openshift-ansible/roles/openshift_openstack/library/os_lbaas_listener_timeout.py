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

''' os_lbaas_deletion '''
import keystoneauth1

from ansible.module_utils.basic import AnsibleModule

try:
    import shade
    HAS_SHADE = True
except ImportError:
    HAS_SHADE = False

DOCUMENTATION = '''
---
module: os_lbaas_listener_timeout
short_description: Modify Octavia listener connection timeouts
description:
    - Set the client and member data timeouts to the specified value (ms)
author:
    - "Luis Tomas Bolivar <ltomasbo@redhat.com>"
'''

RETURN = '''
'''


def main():
    ''' Main module function '''
    module = AnsibleModule(
        argument_spec=dict(
            timeout=dict(default=50000, type='int'),
            listener_name=dict(required=True, type='str'),
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
            service_type=cloud.cloud_config.get_service_type('load-balancer'),
            interface=cloud.cloud_config.get_interface('load-balancer'),
            endpoint_override=cloud.cloud_config.get_endpoint('load-balancer'),
            version=cloud.cloud_config.get_api_version('load-balancer'))
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to get an adapter to talk to the Octavia '
                             'API')
    try:
        listeners = adapter.get(
            'v2.0/lbaas/listeners?name=' + module.params['listener_name'])
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to retrive listeners')

    listener_id = listeners.json()['listeners'][0]['id']
    timeout_data = {'json': {"listener": {
        "timeout_client_data": module.params['timeout'],
        "timeout_member_data": module.params['timeout']}}}
    try:
        adapter.put(
            '/v2.0/lbaas/listeners/' + listener_id, **timeout_data)
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to increate listener timeout')

    module.exit_json(
        changed=True)


if __name__ == '__main__':
    main()
