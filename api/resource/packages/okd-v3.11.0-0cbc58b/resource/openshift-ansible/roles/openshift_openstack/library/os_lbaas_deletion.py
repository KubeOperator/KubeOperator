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

from oslo_serialization import jsonutils

from ansible.module_utils.basic import AnsibleModule

try:
    import shade
    HAS_SHADE = True
except ImportError:
    HAS_SHADE = False

DOCUMENTATION = '''
---
module: os_lbaas_deletion
short_description: Delete LBaaS created by Kuryr
description:
    - Delete all the LBaaS created by Kuryr with the cascade flag
author:
    - "Luis Tomas Bolivar <ltomasbo@redhat.com>"
'''

RETURN = '''
'''


def main():
    ''' Main module function '''
    module = AnsibleModule(
        argument_spec=dict(
            lbaas_annotation=dict(default=False, type='dict'),
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
        lbaas_state = (
            module.params['lbaas_annotation'][
                'openstack.org/kuryr-lbaas-state'])
    # pylint: disable=broad-except
    except Exception:
        module.exit_json(change=True, msg='No information about the lbaas to '
                         'delete')

    lbaas_data = jsonutils.loads(lbaas_state)['versioned_object.data'][
        'loadbalancer']
    lbaas_id = lbaas_data['versioned_object.data']['id']

    try:
        adapter.delete(
            '/v2.0/lbaas/loadbalancers/' + lbaas_id + '?cascade=True')
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to delete Octavia LBaaS with cascade '
                         'flag')

    module.exit_json(
        changed=True)


if __name__ == '__main__':
    main()
