# pylint: skip-file
# flake8: noqa


class OCGroup(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'group'

    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for OCGroup '''
        super(OCGroup, self).__init__(config.namespace, config.kubeconfig)
        self.config = config
        self.namespace = config.namespace
        self._group = None

    @property
    def group(self):
        ''' property function service'''
        if not self._group:
            self.get()
        return self._group

    @group.setter
    def group(self, data):
        ''' setter function for yedit var '''
        self._group = data

    def exists(self):
        ''' return whether a group exists '''
        if self.group:
            return True

        return False

    def get(self):
        '''return group information '''
        result = self._get(self.kind, self.config.name)
        if result['returncode'] == 0:
            self.group = Group(content=result['results'][0])
        elif 'groups.user.openshift.io \"{}\" not found'.format(self.config.name) in result['stderr']:
            result['returncode'] = 0
            result['results'] = [{}]

        return result

    def delete(self):
        '''delete the object'''
        return self._delete(self.kind, self.config.name)

    def create(self):
        '''create the object'''
        return self._create_from_content(self.config.name, self.config.data)

    def update(self):
        '''update the object'''
        return self._replace_content(self.kind, self.config.name, self.config.data)

    def needs_update(self):
        ''' verify an update is needed '''
        return not Utils.check_def_equal(self.config.data, self.group.yaml_dict, skip_keys=['users'], debug=True)

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def run_ansible(params, check_mode=False):
        '''run the oc_group module'''

        gconfig = GroupConfig(params['name'],
                              params['namespace'],
                              params['kubeconfig'],
                             )
        oc_group = OCGroup(gconfig, verbose=params['debug'])

        state = params['state']

        api_rval = oc_group.get()

        if api_rval['returncode'] != 0:
            return {'failed': True, 'msg': api_rval}

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': state}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_group.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a delete.'}

                api_rval = oc_group.delete()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_group.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}

                # Create it here
                api_rval = oc_group.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_group.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if oc_group.needs_update():
                api_rval = oc_group.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_group.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True, 'msg': 'Unknown state passed. {}'.format(state)}
