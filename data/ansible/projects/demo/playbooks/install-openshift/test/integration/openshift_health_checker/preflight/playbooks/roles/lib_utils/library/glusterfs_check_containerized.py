#!/usr/bin/env python
"""glusterfs_check_containerized module"""
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

import subprocess

from ansible.module_utils.basic import AnsibleModule


DOCUMENTATION = '''
---
module: glusterfs_check_containerized

short_description: Check health of each volume in glusterfs on openshift.

version_added: "2.6"

description:
    - This module attempts to ensure all volumes are in healthy state
      in a glusterfs cluster.  The module is meant to be failure-prone, retries
      should be executed at the ansible level, they are not implemented in
      this module.
      This module by executing the following (roughly):
      oc exec --namespace=<namespace> <podname> -- gluster volume list
      for volume in <volume list>:
        gluster volume heal <volume> info

author:
    - "Michael Gugino <mgugino@redhat.com>"
'''

EXAMPLES = '''
- name: glusterfs volumes check
  glusterfs_check_containerized
    oc_bin: "/usr/bin/oc"
    oc_conf: "/etc/origin/master/admin.kubeconfig"
    oc_namespace: "glusterfs"
    cluster_name: "glusterfs"
'''


def fail(module, err):
    """Fail on error"""
    result = {'failed': True,
              'changed': False,
              'msg': err,
              'state': 'unknown'}
    module.fail_json(**result)


def call_or_fail(module, call_args):
    """Call subprocess.check_output and return utf-8 decoded stdout or fail"""
    try:
        # Must decode as utf-8 for python3 compatibility
        res = subprocess.check_output(call_args).decode('utf-8')
    except subprocess.CalledProcessError as err:
        fail(module, str(err))
    return res


def get_valid_nodes(module, oc_exec, exclude_node):
    """Return a list of nodes that will be used to filter running pods"""
    call_args = oc_exec + ['get', 'nodes']
    res = call_or_fail(module, call_args)
    valid_nodes = []
    for line in res.split('\n'):
        fields = line.split()
        if not fields:
            continue
        if fields[0] != exclude_node and fields[1] == "Ready":
            valid_nodes.append(fields[0])
    if not valid_nodes:
        fail(module,
             'Unable to find suitable node in get nodes output: {}'.format(res))
    return valid_nodes


def select_pod(module, oc_exec, cluster_name, valid_nodes):
    """Select a pod to attempt to run gluster commands on"""
    call_args = oc_exec + ['get', 'pods', '-owide']
    res = call_or_fail(module, call_args)
    # res is returned as a tab/space-separated list with headers.
    res_lines = res.split('\n')
    pod_name = None
    name_search = 'glusterfs-{}'.format(cluster_name)
    res_lines = list(filter(None, res.split('\n')))

    for line in res_lines[1:]:
        fields = line.split()
        if not fields:
            continue
        if name_search in fields[0]:
            if fields[2] == "Running" and fields[6] in valid_nodes:
                pod_name = fields[0]
                break

    if pod_name is None:
        fail(module,
             "Unable to find suitable pod in get pods output: {}".format(res))
    else:
        return pod_name


def get_volume_list(module, oc_exec, pod_name):
    """Retrieve list of active volumes from gluster cluster"""
    call_args = oc_exec + ['exec', pod_name, '--', 'gluster', 'volume', 'list']
    res = call_or_fail(module, call_args)
    # This should always at least return heketidbstorage, so no need to check
    # for empty string.
    return list(filter(None, res.split('\n')))


def check_volume_health_info(module, oc_exec, pod_name, volume):
    """Check health info of gluster volume"""
    call_args = oc_exec + ['exec', pod_name, '--', 'gluster', 'volume', 'heal',
                           volume, 'info']
    res = call_or_fail(module, call_args)
    # Output is not easily parsed
    for line in res.split('\n'):
        if line.startswith('Number of entries:'):
            cols = line.split(':')
            if cols[1].strip() != '0':
                fail(module, 'volume {} is not ready'.format(volume))


def check_volumes(module, oc_exec, pod_name):
    """Check status of all volumes on cluster"""
    volume_list = get_volume_list(module, oc_exec, pod_name)
    for volume in volume_list:
        check_volume_health_info(module, oc_exec, pod_name, volume)


def run_module():
    '''Run this module'''
    module_args = dict(
        oc_bin=dict(type='path', required=True),
        oc_conf=dict(type='path', required=True),
        oc_namespace=dict(type='str', required=True),
        cluster_name=dict(type='str', required=True),
        exclude_node=dict(type='str', required=True),
    )
    module = AnsibleModule(
        supports_check_mode=False,
        argument_spec=module_args
    )
    oc_bin = module.params['oc_bin']
    oc_conf = '--config={}'.format(module.params['oc_conf'])
    oc_namespace = '--namespace={}'.format(module.params['oc_namespace'])
    cluster_name = module.params['cluster_name']
    exclude_node = module.params['exclude_node']

    oc_exec = [oc_bin, oc_conf, oc_namespace]

    # create a nodes to find a pod on; We don't want to try to execute on a
    # pod running on a "NotReady" node or the inventory_hostname node because
    # the pods might not actually be alive.
    valid_nodes = get_valid_nodes(module, [oc_bin, oc_conf], exclude_node)

    # Need to find an alive pod to run gluster commands in.
    pod_name = select_pod(module, oc_exec, cluster_name, valid_nodes)

    check_volumes(module, oc_exec, pod_name)

    result = {'changed': False}
    module.exit_json(**result)


def main():
    """main"""
    run_module()


if __name__ == '__main__':
    main()
