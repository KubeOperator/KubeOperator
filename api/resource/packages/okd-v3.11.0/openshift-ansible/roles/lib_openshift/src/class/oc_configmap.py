# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-arguments
class OCConfigMap(OpenShiftCLI):
    ''' Openshift ConfigMap Class

        ConfigMaps are a way to store data inside of objects
    '''
    def __init__(self,
                 name,
                 from_file,
                 from_literal,
                 state,
                 namespace,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 verbose=False):
        ''' Constructor for OpenshiftOC '''
        super(OCConfigMap, self).__init__(namespace, kubeconfig=kubeconfig, verbose=verbose)
        self.name = name
        self.state = state
        self._configmap = None
        self._inc_configmap = None
        self.from_file = from_file if from_file is not None else {}
        self.from_literal = from_literal if from_literal is not None else {}

    @property
    def configmap(self):
        if self._configmap is None:
            self._configmap = self.get()

        return self._configmap

    @configmap.setter
    def configmap(self, inc_map):
        self._configmap = inc_map

    @property
    def inc_configmap(self):
        if self._inc_configmap is None:
            results = self.create(dryrun=True, output=True)
            self._inc_configmap = results['results']

        return self._inc_configmap

    @inc_configmap.setter
    def inc_configmap(self, inc_map):
        self._inc_configmap = inc_map

    def from_file_to_params(self):
        '''return from_files in a string ready for cli'''
        return ["--from-file={}={}".format(key, value) for key, value in self.from_file.items()]

    def from_literal_to_params(self):
        '''return from_literal in a string ready for cli'''
        return ["--from-literal={}={}".format(key, value) for key, value in self.from_literal.items()]

    def get(self):
        '''return a configmap by name '''
        results = self._get('configmap', self.name)
        if results['returncode'] == 0 and results['results'][0]:
            self.configmap = results['results'][0]

        if results['returncode'] != 0 and '"{}" not found'.format(self.name) in results['stderr']:
            results['returncode'] = 0

        return results

    def delete(self):
        '''delete a configmap by name'''
        return self._delete('configmap', self.name)

    def create(self, dryrun=False, output=False):
        '''Create a configmap

           :dryrun: Product what you would have done. default: False
           :output: Whether to parse output. default: False
        '''

        cmd = ['create', 'configmap', self.name]
        if self.from_literal is not None:
            cmd.extend(self.from_literal_to_params())

        if self.from_file is not None:
            cmd.extend(self.from_file_to_params())

        if dryrun:
            cmd.extend(['--dry-run', '-ojson'])

        results = self.openshift_cmd(cmd, output=output)

        return results

    def update(self):
        '''run update configmap '''
        return self._replace_content('configmap', self.name, self.inc_configmap)

    def needs_update(self):
        '''compare the current configmap with the proposed and return if they are equal'''
        return not Utils.check_def_equal(self.inc_configmap, self.configmap, debug=self.verbose)

    @staticmethod
    # pylint: disable=too-many-return-statements,too-many-branches
    # TODO: This function should be refactored into its individual parts.
    def run_ansible(params, check_mode):
        '''run the oc_configmap module'''

        oc_cm = OCConfigMap(params['name'],
                            params['from_file'],
                            params['from_literal'],
                            params['state'],
                            params['namespace'],
                            kubeconfig=params['kubeconfig'],
                            verbose=params['debug'])

        state = params['state']

        api_rval = oc_cm.get()

        if 'failed' in api_rval:
            return {'failed': True, 'msg': api_rval}

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval, 'state': state}

        if not params['name']:
            return {'failed': True,
                    'msg': 'Please specify a name when state is absent|present.'}

        ########
        # Delete
        ########
        if state == 'absent':
            if not Utils.exists(api_rval['results'], params['name']):
                return {'changed': False, 'state': 'absent'}

            if check_mode:
                return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a delete.'}

            api_rval = oc_cm.delete()

            if api_rval['returncode'] != 0:
                return {'failed': True, 'msg': api_rval}

            return {'changed': True, 'results': api_rval, 'state': state}

        ########
        # Create
        ########
        if state == 'present':
            if not Utils.exists(api_rval['results'], params['name']):

                if check_mode:
                    return {'changed': True, 'msg': 'Would have performed a create.'}

                api_rval = oc_cm.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                api_rval = oc_cm.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if oc_cm.needs_update():

                api_rval = oc_cm.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                api_rval = oc_cm.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True, 'msg': 'Unknown state passed. {}'.format(state)}
