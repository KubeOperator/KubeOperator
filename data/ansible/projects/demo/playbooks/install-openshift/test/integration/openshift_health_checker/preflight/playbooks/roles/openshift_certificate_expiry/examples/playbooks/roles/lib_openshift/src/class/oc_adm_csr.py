# pylint: skip-file
# flake8: noqa


class OCcsr(OpenShiftCLI):
    ''' Class to wrap the oc adm certificate command line'''
    kind = 'csr'

    # pylint: disable=too-many-arguments
    def __init__(self,
                 nodes=None,
                 approve_all=False,
                 service_account=None,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 verbose=False):
        ''' Constructor for oc adm certificate '''
        super(OCcsr, self).__init__(None, kubeconfig, verbose)
        self.service_account = service_account
        self.nodes = self.create_nodes(nodes)
        self._csrs = []
        self.approve_all = approve_all
        self.verbose = verbose

    @property
    def csrs(self):
        '''property for managing csrs'''
        # any processing needed??
        self._csrs = self._get(resource=self.kind)['results'][0]['items']
        return self._csrs

    def create_nodes(self, nodes):
        '''create a node object to track csr signing status'''
        nodes_list = []

        if nodes is None:
            return nodes_list

        results = self._get(resource='nodes')['results'][0]['items']

        for node in nodes:
            nodes_list.append(dict(name=node, csrs={}, server_accepted=False, client_accepted=False, denied=False))

            # Ready nodes have already been accepted. Mark client and server as accepted.
            for ocnode in results:
                if ocnode['metadata']['name'] == node:
                    for condition in ocnode['status']['conditions']:
                        if condition['type'] == 'Ready' and condition['status'] == 'True':
                            nodes_list[-1]['server_accepted'] = True
                            nodes_list[-1]['client_accepted'] = True

        return nodes_list

    def get(self):
        '''get the current certificate signing requests'''
        return self.csrs

    @staticmethod
    def action_needed(csr, action):
        '''check to see if csr is in desired state'''
        if csr['status'] == {}:
            return True

        state = csr['status']['conditions'][0]['type']

        if action == 'approve' and state != 'Approved':
            return True

        elif action == 'deny' and state != 'Denied':
            return True

        return False

    def get_csr_request(self, request):
        '''base64 decode the request object and call openssl to determine the
           subject and specifically the CN: from the request

           Output:
           (0, '...
                Subject: O=system:nodes, CN=system:node:ip-172-31-54-54.ec2.internal
                ...')
        '''
        import base64
        return self._run(['openssl', 'req', '-noout', '-text'], base64.b64decode(request))[1]

    def match_node(self, csr):
        '''match an inc csr to a node in self.nodes'''
        for node in self.nodes:
            # we need to match based upon the csr's request certificate's CN
            if node['name'] in self.get_csr_request(csr['spec']['request']):
                node['csrs'][csr['metadata']['name']] = csr

                # client certs may come in as either the service_account or as the node during upgrade
                # server certs always come in as the node
                if ((node['name'] in csr['spec']['username'] or
                     csr['spec']['username'] in [self.service_account, 'system:admin']) and
                        csr['status'] and csr['status']['conditions'][0]['type'] == 'Approved'):
                    if 'server auth' in csr['spec']['usages']:
                        node['server_accepted'] = True
                    if 'client auth' in csr['spec']['usages']:
                        node['client_accepted'] = True
                # check type is 'Denied' and mark node as such
                if csr['status'] and csr['status']['conditions'][0]['type'] == 'Denied':
                    node['denied'] = True
                return node
        return None

    def finished(self):
        '''determine if there are more csrs to sign'''
        # if nodes is set and we have nodes then return if all nodes are 'accepted'
        if self.nodes is not None and len(self.nodes) > 0:
            return all([(node['server_accepted'] and node['client_accepted']) or node['denied'] for node in self.nodes])

        # we are approving everything or we still have nodes outstanding
        return False

    def manage(self, action):
        '''run openshift oc adm ca create-server-cert cmd and store results into self.nodes

           we attempt to verify if the node is one that was given to us to accept.

           action - (allow | deny)
        '''

        results = []
        # There are 2 types of requests:
        # - node-bootstrapper-client-ip-172-31-51-246-ec2-internal
        #   The client request allows the client to talk to the api/controller
        # - node-bootstrapper-server-ip-172-31-51-246-ec2-internal
        #   The server request allows the server to join the cluster
        # Here we need to determine how to approve/deny
        # we should query the csrs and verify they are from the nodes we thought
        for csr in self.csrs:
            node = self.match_node(csr)
            # oc adm certificate <approve|deny> csr
            # there are 3 known states: Denied, Approved, {}
            # verify something is needed by OCcsr.action_needed
            # if approve_all, then do it
            # if you passed in nodes, you must have a node that matches
            if self.approve_all or (node and OCcsr.action_needed(csr, action)):
                result = self.openshift_cmd(['certificate', action, csr['metadata']['name']], oadm=True)
                # if we successfully approved
                if result['returncode'] == 0:
                    # client should have service account name in username field
                    # server should have node name in username field
                    if node and csr['metadata']['name'] not in node['csrs']:
                        node['csrs'][csr['metadata']['name']] = csr

                    # mark node as accepted in our list of nodes
                    # we will use {client,server}_accepted fields to determine if we're finished
                    if (node['name'] in csr['spec']['username'] or
                            csr['spec']['username'] in [self.service_account, 'system:admin']):
                        if 'server auth' in csr['spec']['usages']:
                            node['server_accepted'] = True
                        if 'client auth' in csr['spec']['usages']:
                            node['client_accepted'] = True

                results.append(result)

        return results

    @staticmethod
    def run_ansible(params, check_mode=False):
        '''run the oc_adm_csr module'''

        client = OCcsr(params['nodes'],
                       params['approve_all'],
                       params['service_account'],
                       params['kubeconfig'],
                       params['debug'])

        state = params['state']

        api_rval = client.get()

        if state == 'list':
            return {'changed': False, 'results': api_rval, 'state': state}

        if state in ['approve', 'deny']:
            if check_mode:
                return {'changed': True,
                        'msg': "CHECK_MODE: Would have {} the certificate.".format(params['state']),
                        'state': state}

            all_results = []
            finished = False
            timeout = False
            # loop for timeout or block until all nodes pass
            ctr = 0
            while True:

                all_results.extend(client.manage(params['state']))
                if client.finished():
                    finished = True
                    break

                if params['timeout'] == 0:
                    if not params['approve_all']:
                        ctr = 0

                if ctr * 2 > params['timeout']:
                    timeout = True
                    break

                # This provides time for the nodes to send their csr requests between approvals
                time.sleep(2)

                ctr += 1

            for result in all_results:
                if result['returncode'] != 0:
                    return {'failed': True, 'msg': all_results, 'timeout': timeout}

            return dict(changed=len(all_results) > 0,
                        results=all_results,
                        nodes=client.nodes,
                        state=state,
                        finished=finished,
                        timeout=timeout)

        return {'failed': True,
                'msg': 'Unknown state passed. %s' % state}
