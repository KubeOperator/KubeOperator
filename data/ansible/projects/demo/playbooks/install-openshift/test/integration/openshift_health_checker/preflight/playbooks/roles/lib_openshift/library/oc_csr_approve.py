#!/usr/bin/env python
'''oc_csr_approve module'''
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

import base64
import json

from ansible.module_utils.basic import AnsibleModule

try:
    from json.decoder import JSONDecodeError
except ImportError:
    JSONDecodeError = ValueError

DOCUMENTATION = '''
---
module: oc_csr_approve

short_description: Retrieve, approve, and verify node client csrs

version_added: "2.4"

description:
    - Runs various commands to list csrs, approve csrs, and verify nodes are
      ready.

author:
    - "Michael Gugino <mgugino@redhat.com>"
'''

EXAMPLES = '''
# Pass in a message
- name: Place credentials in file
  oc_csr_approve:
    oc_bin: "/usr/bin/oc"
    oc_conf: "/etc/origin/master/admin.kubeconfig"
    node_list: ['node1.example.com', 'node2.example.com']
'''

CERT_MODE = {'client': 'client auth', 'server': 'server auth'}


def parse_subject_cn(subject_str):
    '''parse output of openssl req -noout -subject to retrieve CN.
       example input:
         'subject=/C=US/CN=test.io/L=Raleigh/O=Red Hat/ST=North Carolina/OU=OpenShift\n'
         or
         'subject=C = US, CN = test.io, L = City, O = Company, ST = State, OU = Dept\n'
       example output: 'test.io'
    '''
    stripped_string = subject_str[len('subject='):].strip()
    kv_strings = [x.strip() for x in stripped_string.split(',')]
    if len(kv_strings) == 1:
        kv_strings = [x.strip() for x in stripped_string.split('/')][1:]
    for item in kv_strings:
        item_parts = [x.strip() for x in item.split('=')]
        if item_parts[0] == 'CN':
            return item_parts[1]


class CSRapprove(object):
    """Approves csr requests"""

    def __init__(self, module, oc_bin, oc_conf, node_list):
        '''init method'''
        self.module = module
        self.oc_bin = oc_bin
        self.oc_conf = oc_conf
        self.node_list = node_list
        self.all_subjects_found = []
        self.unwanted_csrs = []
        # Build a dictionary to hold all of our output information so nothing
        # is lost when we fail.
        self.result = {'changed': False, 'rc': 0,
                       'oc_get_nodes': None,
                       'client_csrs': None,
                       'server_csrs': None,
                       'all_subjects_found': self.all_subjects_found,
                       'client_approve_results': [],
                       'server_approve_results': [],
                       'unwanted_csrs': self.unwanted_csrs}

    def run_command(self, command, rc_opts=None):
        '''Run a command using AnsibleModule.run_command, or fail'''
        if rc_opts is None:
            rc_opts = {}
        rtnc, stdout, err = self.module.run_command(command, **rc_opts)
        if rtnc:
            self.result['failed'] = True
            self.result['msg'] = str(err)
            self.result['state'] = 'unknown'
            self.module.fail_json(**self.result)
        return stdout

    def get_nodes(self):
        '''Get all nodes via oc get nodes -ojson'''
        # json output is necessary for consistency here.
        command = "{} {} get nodes -ojson".format(self.oc_bin, self.oc_conf)
        stdout = self.run_command(command)
        try:
            data = json.loads(stdout)
        except JSONDecodeError as err:
            self.result['failed'] = True
            self.result['msg'] = str(err)
            self.result['state'] = 'unknown'
            self.module.fail_json(**self.result)
        self.result['oc_get_nodes'] = data
        return [node['metadata']['name'] for node in data['items']]

    def get_csrs(self):
        '''Retrieve csrs from cluster using oc get csr -ojson'''
        command = "{} {} get csr -ojson".format(self.oc_bin, self.oc_conf)
        stdout = self.run_command(command)
        try:
            data = json.loads(stdout)
        except JSONDecodeError as err:
            self.result['failed'] = True
            self.result['msg'] = str(err)
            self.result['state'] = 'unknown'
            self.module.fail_json(**self.result)
        return data['items']

    def process_csrs(self, csrs, mode):
        '''Return a dictionary of pending csrs where the format of the dict is
           k=csr name, v=Subject Common Name'''
        csr_dict = {}
        for item in csrs:
            name = item['metadata']['name']
            request_data = base64.b64decode(item['spec']['request'])
            command = "openssl req -noout -subject"
            # ansible's module.run_command accepts data to pipe via stdin as
            # as 'data' kwarg.
            rc_opts = {'data': request_data, 'binary_data': True}
            stdout = self.run_command(command, rc_opts=rc_opts)
            self.all_subjects_found.append(stdout)

            status = item['status'].get('conditions')
            if status:
                # If status is not an empty dictionary, cert is not pending.
                self.unwanted_csrs.append(item)
                continue
            if CERT_MODE[mode] not in item['spec']['usages']:
                self.unwanted_csrs.append(item)
                continue
            # parse common_name from subject string.
            common_name = parse_subject_cn(stdout)
            if common_name and common_name.startswith('system:node:'):
                # common name is typically prepended with system:node:.
                common_name = common_name.split('system:node:')[1]
            # we only want to approve csrs from nodes we know about.
            if common_name in self.node_list:
                csr_dict[name] = common_name
            else:
                self.unwanted_csrs.append(item)

        return csr_dict

    def confirm_needed_requests_present(self, not_ready_nodes, csr_dict):
        '''Ensure all non-Ready nodes have a csr, or fail'''
        nodes_needed = set(not_ready_nodes)
        for _, val in csr_dict.items():
            nodes_needed.discard(val)

        # check that we found all of our needed nodes
        if nodes_needed:
            missing_nodes = ', '.join(nodes_needed)
            self.result['failed'] = True
            self.result['msg'] = "Could not find csr for nodes: {}".format(missing_nodes)
            self.result['state'] = 'unknown'
            self.module.fail_json(**self.result)

    def approve_csrs(self, csr_pending_list, mode):
        '''Loop through csr_pending_list and call:
           oc adm certificate approve <item>'''
        res_mode = "{}_approve_results".format(mode)
        base_command = "{} {} adm certificate approve {}"
        approve_results = []
        for csr in csr_pending_list:
            command = base_command.format(self.oc_bin, self.oc_conf, csr)
            rtnc, stdout, err = self.module.run_command(command)
            approve_results.append(stdout)
            if rtnc:
                self.result['failed'] = True
                self.result['msg'] = str(err)
                self.result[res_mode] = approve_results
                self.result['state'] = 'unknown'
                self.module.fail_json(**self.result)
        self.result[res_mode] = approve_results
        # We set changed for approved client or server csrs.
        self.result['changed'] = bool(approve_results) or bool(self.result['changed'])

    def get_ready_nodes_server(self, nodes_list):
        '''Determine which nodes have working server certificates'''
        ready_nodes_server = []
        base_command = "{} {} get --raw /api/v1/nodes/{}/proxy/healthz"
        for node in nodes_list:
            # need this to look like /api/v1/nodes/<node>/proxy/healthz
            command = base_command.format(self.oc_bin, self.oc_conf, node)
            rtnc, _, _ = self.module.run_command(command)
            if not rtnc:
                # if we can hit that api endpoint, the node has a valid server
                # cert.
                ready_nodes_server.append(node)
        return ready_nodes_server

    def verify_server_csrs(self):
        '''We approved some server csrs, now we need to validate they are working.
           This function will attempt to retry 10 times in case of failure.'''
        # Attempt to try node endpoints a few times.
        attempts = 0
        # Find not_ready_nodes for server-side again
        nodes_server_ready = self.get_ready_nodes_server(self.node_list)
        # Create list of nodes that still aren't ready.
        not_ready_nodes_server = set([item for item in self.node_list if item not in nodes_server_ready])
        while not_ready_nodes_server:
            nodes_server_ready = self.get_ready_nodes_server(not_ready_nodes_server)

            # if we have same number of nodes_server_ready now, all of the previous
            # not_ready_nodes are now ready.
            if not len(not_ready_nodes_server - set(nodes_server_ready)):
                break
            attempts += 1
            if attempts > 9:
                self.result['failed'] = True
                self.result['rc'] = 1
                missing_nodes = not_ready_nodes_server - set(nodes_server_ready)
                msg = "Some nodes still not ready after approving server certs: {}"
                msg = msg.format(", ".join(missing_nodes))
                self.result['msg'] = msg
                self.module.fail_json(**self.result)

    def run(self):
        '''execute the csr approval process'''
        all_nodes = self.get_nodes()
        # don't need to check nodes that have already joined the cluster because
        # client csr needs to be approved for now to show in output of
        # oc get nodes.
        not_found_nodes = [item for item in self.node_list
                           if item not in all_nodes]

        # Get all csrs, no good way to filter on pending.
        client_csrs = self.get_csrs()
        # process data in csrs and build a dictionary of client requests
        client_csr_dict = self.process_csrs(client_csrs, "client")
        self.result['client_csrs'] = client_csr_dict

        # This method is fail-happy and expects all not found nodes have available
        # csrs.  Handle failure for this method via ansible retry/until.
        self.confirm_needed_requests_present(not_found_nodes,
                                             client_csr_dict)
        # If for some reason a node is found in oc get nodes but it still needs
        # a client csr approved, this method will approve all outstanding
        # client csrs for any node in our self.node_list.
        self.approve_csrs(client_csr_dict, 'client')

        # # Server Cert Section # #
        # Find not_ready_nodes for server-side
        nodes_server_ready = self.get_ready_nodes_server(self.node_list)
        # Create list of nodes that definitely need a server cert approved.
        not_ready_nodes_server = [item for item in self.node_list
                                  if item not in nodes_server_ready]

        # Get all csrs again, no good way to filter on pending.
        server_csrs = self.get_csrs()
        # process data in csrs and build a dictionary of server requests
        server_csr_dict = self.process_csrs(server_csrs, "server")
        self.result['server_csrs'] = server_csr_dict

        # This will fail if all server csrs are not present, but probably shouldn't
        # at this point since we spent some time hitting the api to see if the
        # nodes are already responding.
        self.confirm_needed_requests_present(not_ready_nodes_server,
                                             server_csr_dict)
        self.approve_csrs(server_csr_dict, 'server')

        self.verify_server_csrs()

        # We made it here, everything was successful, cleanup some debug info
        # so we don't spam logs.
        for key in ('client_csrs', 'server_csrs', 'unwanted_csrs'):
            self.result.pop(key)
        self.module.exit_json(**self.result)


def run_module():
    '''Run this module'''
    module_args = dict(
        oc_bin=dict(type='path', required=False, default='oc'),
        oc_conf=dict(type='path', required=False, default='/etc/origin/master/admin.kubeconfig'),
        node_list=dict(type='list', required=True),
    )
    module = AnsibleModule(
        supports_check_mode=False,
        argument_spec=module_args
    )
    oc_bin = module.params['oc_bin']
    oc_conf = '--config={}'.format(module.params['oc_conf'])
    node_list = module.params['node_list']

    approver = CSRapprove(module, oc_bin, oc_conf, node_list)
    approver.run()


def main():
    '''main'''
    run_module()


if __name__ == '__main__':
    main()
