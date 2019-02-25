# pylint: skip-file
# flake8: noqa


class PolicyGroupException(Exception):
    ''' PolicyGroup exception'''
    pass


class PolicyGroupConfig(OpenShiftCLIConfig):
    ''' PolicyGroupConfig is a DTO for group related policy.  '''
    def __init__(self, namespace, kubeconfig, policy_options):
        super(PolicyGroupConfig, self).__init__(policy_options['name']['value'],
                                                namespace, kubeconfig, policy_options)
        self.kind = self.get_kind()
        self.namespace = namespace

    def get_kind(self):
        ''' return the kind we are working with '''
        if self.config_options['resource_kind']['value'] == 'role':
            return 'rolebinding'
        elif self.config_options['resource_kind']['value'] == 'cluster-role':
            return 'clusterrolebinding'
        elif self.config_options['resource_kind']['value'] == 'scc':
            return 'scc'

        return None


# pylint: disable=too-many-return-statements
class PolicyGroup(OpenShiftCLI):
    ''' Class to handle attaching policies to users '''


    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for PolicyGroup '''
        super(PolicyGroup, self).__init__(config.namespace, config.kubeconfig, verbose)
        self.config = config
        self.verbose = verbose
        self._rolebinding = None
        self._scc = None
        self._cluster_role_bindings = None
        self._role_bindings = None

    @property
    def rolebindings(self):
        if self._role_bindings is None:
            results = self._get('rolebindings', None)
            if results['returncode'] != 0:
                raise OpenShiftCLIError('Could not retrieve rolebindings')
            self._role_bindings = results['results'][0]['items']

        return self._role_bindings

    @property
    def clusterrolebindings(self):
        if self._cluster_role_bindings is None:
            results = self._get('clusterrolebindings', None)
            if results['returncode'] != 0:
                raise OpenShiftCLIError('Could not retrieve clusterrolebindings')
            self._cluster_role_bindings = results['results'][0]['items']

        return self._cluster_role_bindings

    @property
    def role_binding(self):
        ''' role_binding getter '''
        return self._rolebinding

    @role_binding.setter
    def role_binding(self, binding):
        ''' role_binding setter '''
        self._rolebinding = binding

    @property
    def security_context_constraint(self):
        ''' security_context_constraint getter '''
        return self._scc

    @security_context_constraint.setter
    def security_context_constraint(self, scc):
        ''' security_context_constraint setter '''
        self._scc = scc

    def get(self):
        '''fetch the desired kind'''
        resource_name = self.config.config_options['name']['value']
        if resource_name == 'cluster-reader':
            resource_name += 's'

        # oc adm policy add-... creates policy bindings with the name
        # "[resource_name]-binding", however some bindings in the system
        # simply use "[resource_name]". So try both.

        results = self._get(self.config.kind, resource_name)
        if results['returncode'] == 0:
            return results

        # Now try -binding naming convention
        return self._get(self.config.kind, resource_name + "-binding")

    def exists_role_binding(self):
        ''' return whether role_binding exists '''
        bindings = None
        if self.config.config_options['resource_kind']['value'] == 'cluster-role':
            bindings = self.clusterrolebindings
        else:
            bindings = self.rolebindings

        if bindings is None:
            return False

        for binding in bindings:
            if self.config.config_options['rolebinding_name']['value'] is not None and \
                    binding['metadata']['name'] != self.config.config_options['rolebinding_name']['value']:
                continue
            if binding['roleRef']['name'] == self.config.config_options['name']['value'] and \
                    binding['groupNames'] is not None and \
                    self.config.config_options['group']['value'] in binding['groupNames']:
                self.role_binding = binding
                return True

        return False

    def exists_scc(self):
        ''' return whether scc exists '''
        results = self.get()
        if results['returncode'] == 0:
            self.security_context_constraint = SecurityContextConstraints(results['results'][0])

            if self.security_context_constraint.find_group(self.config.config_options['group']['value']) != None:
                return True

            return False

        return results

    def exists(self):
        '''does the object exist?'''
        if self.config.config_options['resource_kind']['value'] == 'cluster-role':
            return self.exists_role_binding()

        elif self.config.config_options['resource_kind']['value'] == 'role':
            return self.exists_role_binding()

        elif self.config.config_options['resource_kind']['value'] == 'scc':
            return self.exists_scc()

        return False

    def perform(self):
        '''perform action on resource'''
        cmd = ['policy',
               self.config.config_options['action']['value'],
               self.config.config_options['name']['value'],
               self.config.config_options['group']['value']]

        if self.config.config_options['rolebinding_name']['value'] is not None:
            cmd.extend(['--rolebinding-name', self.config.config_options['rolebinding_name']['value']])

        return self.openshift_cmd(cmd, oadm=True)

    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_adm_policy_group module'''

        state = params['state']

        action = None
        if state == 'present':
            action = 'add-' + params['resource_kind'] + '-to-group'
        else:
            action = 'remove-' + params['resource_kind'] + '-from-group'

        nconfig = PolicyGroupConfig(params['namespace'],
                                    params['kubeconfig'],
                                    {'action': {'value': action, 'include': False},
                                     'group': {'value': params['group'], 'include': False},
                                     'resource_kind': {'value': params['resource_kind'], 'include': False},
                                     'name': {'value': params['resource_name'], 'include': False},
                                     'rolebinding_name': {'value': params['rolebinding_name'], 'include': False},
                                    })

        policygroup = PolicyGroup(nconfig, params['debug'])

        # Run the oc adm policy group related command

        ########
        # Delete
        ########
        if state == 'absent':
            if not policygroup.exists():
                return {'changed': False, 'state': 'absent'}

            if check_mode:
                return {'changed': False, 'msg': 'CHECK_MODE: would have performed a delete.'}

            api_rval = policygroup.perform()

            if api_rval['returncode'] != 0:
                return {'msg': api_rval}

            return {'changed': True, 'results' : api_rval, state:'absent'}

        if state == 'present':
            ########
            # Create
            ########
            results = policygroup.exists()
            if isinstance(results, dict) and 'returncode' in results and results['returncode'] != 0:
                return {'msg': results}

            if not results:

                if check_mode:
                    return {'changed': False, 'msg': 'CHECK_MODE: would have performed a create.'}

                api_rval = policygroup.perform()

                if api_rval['returncode'] != 0:
                    return {'msg': api_rval}

                return {'changed': True, 'results': api_rval, state: 'present'}

            return {'changed': False, state: 'present'}

        return {'failed': True, 'changed': False, 'results': 'Unknown state passed. %s' % state, state: 'unknown'}
