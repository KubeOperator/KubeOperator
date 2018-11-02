# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCClusterRole(OpenShiftCLI):
    ''' Class to manage clusterrole objects'''
    kind = 'clusterrole'

    def __init__(self,
                 name,
                 rules=None,
                 kubeconfig=None,
                 verbose=False):
        ''' Constructor for OCClusterRole '''
        super(OCClusterRole, self).__init__(None, kubeconfig=kubeconfig, verbose=verbose)
        self.verbose = verbose
        self.name = name
        self._clusterrole = None
        self._inc_clusterrole = ClusterRole.builder(name, rules)

    @property
    def clusterrole(self):
        ''' property for clusterrole'''
        if self._clusterrole is None:
            self.get()
        return self._clusterrole

    @clusterrole.setter
    def clusterrole(self, data):
        ''' setter function for clusterrole property'''
        self._clusterrole = data

    @property
    def inc_clusterrole(self):
        ''' property for inc_clusterrole'''
        return self._inc_clusterrole

    @inc_clusterrole.setter
    def inc_clusterrole(self, data):
        ''' setter function for inc_clusterrole property'''
        self._inc_clusterrole = data

    def exists(self):
        ''' return whether a clusterrole exists '''
        if self.clusterrole:
            return True

        return False

    def get(self):
        '''return a clusterrole '''
        result = self._get(self.kind, self.name)

        if result['returncode'] == 0:
            self.clusterrole = ClusterRole(content=result['results'][0])
            result['results'] = self.clusterrole.yaml_dict

        elif '"{}" not found'.format(self.name) in result['stderr']:
            result['returncode'] = 0
            self.clusterrole = None

        return result

    def delete(self):
        '''delete the object'''
        return self._delete(self.kind, self.name)

    def create(self):
        '''create a clusterrole from the proposed incoming clusterrole'''
        return self._create_from_content(self.name, self.inc_clusterrole.yaml_dict)

    def update(self):
        '''update a project'''
        return self._replace_content(self.kind, self.name, self.inc_clusterrole.yaml_dict)

    def needs_update(self):
        ''' verify an update is needed'''
        return not self.clusterrole.compare(self.inc_clusterrole, self.verbose)

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_clusterrole module'''

        oc_clusterrole = OCClusterRole(params['name'],
                                       params['rules'],
                                       params['kubeconfig'],
                                       params['debug'])

        state = params['state']

        api_rval = oc_clusterrole.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval, 'state': state}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_clusterrole.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a delete.'}

                api_rval = oc_clusterrole.delete()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_clusterrole.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}

                # Create it here
                api_rval = oc_clusterrole.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_clusterrole.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if oc_clusterrole.needs_update():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed an update.'}

                api_rval = oc_clusterrole.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_clusterrole.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True,
                'changed': False,
                'msg': 'Unknown state passed. [%s]' % state}
