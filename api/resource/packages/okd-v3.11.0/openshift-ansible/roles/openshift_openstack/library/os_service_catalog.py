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

''' os_service_catalog_facts '''

from ansible.module_utils.basic import AnsibleModule

try:
    import shade
    HAS_SHADE = True
except ImportError:
    HAS_SHADE = False

DOCUMENTATION = '''
---
module: os_service_catalog_facts
short_description: Retrieve OpenStack service catalog facts
description:
    - Retrieves all the available OpenStack services
notes:
    - This module creates a new top-level C(openstack_service_catalog) fact
      which contains a dictionary of OpenStack service endpoints like
      network and load-balancers.
author:
    - "Antoni Segura Puimedon <antoni@redhat.com>"
'''

RETURN = '''
openstack_service_catalog:
    description: OpenStack available services.
    type: dict
    returned: always
    sample:
      alarming:
      - adminURL: http://172.16.0.9:8042
        id: 2c40b50da0bb44178db91c8a9a29a46e
        internalURL: http://172.16.0.9:8042
        publicURL: https://mycloud.org:13042
        region: regionOne
      cloudformation:
      - adminURL: http://172.16.0.9:8000/v1
        id: 46648eded04e463281a9cba7ddcc45cb
        internalURL: http://172.16.0.9:8000/v1
        publicURL: https://mycloud.org:13005/v1
        region: regionOne
      compute:
      - adminURL: http://172.16.0.9:8774/v2.1
        id: bff1bc5dd92842c281b2358a6d15c5bc
        internalURL: http://172.16.0.9:8774/v2.1
        publicURL: https://mycloud.org:13774/v2.1
        region: regionOne
      event:
      - adminURL: http://172.16.0.9:8779
        id: 608ac3666ef24f2e8f240785b8612efb
        internalURL: http://172.16.0.9:8779
        publicURL: https://mycloud.org:13779
        region: regionOne
      identity:
      - adminURL: https://mycloud.org:35357
        id: 4d07689ce46b4d51a01cc873bc772c80
        internalURL: http://172.16.0.9:5000
        publicURL: https://mycloud.org:13000
        region: regionOne
      image:
      - adminURL: http://172.16.0.9:9292
        id: 1850105115ea493eb65f3f704d421291
        internalURL: http://172.16.0.9:9292
        publicURL: https://mycloud.org:13292
        region: regionOne
      metering:
      - adminURL: http://172.16.0.9:8777
        id: 4cae4dcabe0a4914a6ec6dabd62490ba
        internalURL: http://172.16.0.9:8777
        publicURL: https://mycloud.org:13777
        region: regionOne
      metric:
      - adminURL: http://172.16.0.9:8041
        id: 29bcecf9a06f40f782f19dd7492af352
        internalURL: http://172.16.0.9:8041
        publicURL: https://mycloud.org:13041
        region: regionOne
      network:
      - adminURL: http://172.16.0.9:9696
        id: 5d5785c9b8174c21bfb19dc3b16c87fa
        internalURL: http://172.16.0.9:9696
        publicURL: https://mycloud.org:13696
        region: regionOne
      object-store:
      - adminURL: http://172.17.0.9:8080
        id: 031f1e342fdf4f25b6099d1f3b0847e3
        internalURL: http://172.17.0.9:8080/v1/AUTH_6d2847d6a6414308a67644eefc7b98c7
        publicURL: https://mycloud.org:13808/v1/AUTH_6d2847d6a6414308a67644eefc7b98c7
        region: regionOne
      orchestration:
      - adminURL: http://172.16.0.9:8004/v1/6d2847d6a6414308a67644eefc7b98c7
        id: 1e6cecbd15b3413d9411052c52b9d433
        internalURL: http://172.16.0.9:8004/v1/6d2847d6a6414308a67644eefc7b98c7
        publicURL: https://mycloud.org:13004/v1/6d2847d6a6414308a67644eefc7b98c7
        region: regionOne
      placement:
      - adminURL: http://172.16.0.9:8778/placement
        id: 1f2551e5450c4bd6a9f716f92e93a154
        internalURL: http://172.16.0.9:8778/placement
        publicURL: https://mycloud.org:13778/placement
        region: regionOne
      volume:
      - adminURL: http://172.16.0.9:8776/v1/6d2847d6a6414308a67644eefc7b98c7
        id: 38e369a0e17346fe8e37a20146e005ef
        internalURL: http://172.16.0.9:8776/v1/6d2847d6a6414308a67644eefc7b98c7
        publicURL: https://mycloud.org:13776/v1/6d2847d6a6414308a67644eefc7b98c7
        region: regionOne
      volumev2:
      - adminURL: http://172.16.0.9:8776/v2/6d2847d6a6414308a67644eefc7b98c7
        id: 113a0bff9f2347b6b8774407a1c8d572
        internalURL: http://172.16.0.9:8776/v2/6d2847d6a6414308a67644eefc7b98c7
        publicURL: https://mycloud.org:13776/v2/6d2847d6a6414308a67644eefc7b98c7
        region: regionOne
      volumev3:
      - adminURL: http://172.16.0.9:8776/v3/6d2847d6a6414308a67644eefc7b98c7
        id: 9982c0afd28941a19feb1ffb13b91daf
        internalURL: http://172.16.0.9:8776/v3/6d2847d6a6414308a67644eefc7b98c7
        publicURL: https://mycloud.org:13776/v3/6d2847d6a6414308a67644eefc7b98c7
        region: regionOne
'''


def main():
    ''' Main module function '''
    module = AnsibleModule(argument_spec={}, supports_check_mode=True)

    if not HAS_SHADE:
        module.fail_json(msg='shade is required for this module')

    try:
        cloud = shade.openstack_cloud()
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to connect to the cloud')

    try:
        service_catalog = cloud.cloud_config.get_service_catalog()
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to retrieve the service catalog')

    try:
        endpoints = service_catalog.get_endpoints()
    # pylint: disable=broad-except
    except Exception:
        module.fail_json(msg='Failed to retrieve the service catalog '
                         'endpoints')

    module.exit_json(
        changed=False,
        ansible_facts={'openstack_service_catalog': endpoints})


if __name__ == '__main__':
    main()
