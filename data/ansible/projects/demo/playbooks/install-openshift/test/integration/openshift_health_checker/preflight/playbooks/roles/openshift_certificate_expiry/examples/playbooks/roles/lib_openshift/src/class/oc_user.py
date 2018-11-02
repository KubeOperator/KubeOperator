# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-instance-attributes
class OCUser(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'users'

    def __init__(self,
                 config,
                 groups=None,
                 verbose=False):
        ''' Constructor for OCUser '''
        # namespace has no meaning for user operations, hardcode to 'default'
        super(OCUser, self).__init__('default', config.kubeconfig)
        self.config = config
        self.groups = groups
        self._user = None

    @property
    def user(self):
        ''' property function user'''
        if not self._user:
            self.get()
        return self._user

    @user.setter
    def user(self, data):
        ''' setter function for user '''
        self._user = data

    def exists(self):
        ''' return whether a user exists '''
        if self.user:
            return True

        return False

    def get(self):
        ''' return user information '''
        result = self._get(self.kind, self.config.username)
        if result['returncode'] == 0:
            self.user = User(content=result['results'][0])
        elif 'users \"%s\" not found' % self.config.username in result['stderr']:
            result['returncode'] = 0
            result['results'] = [{}]

        return result

    def delete(self):
        ''' delete the object '''
        return self._delete(self.kind, self.config.username)

    def create_group_entries(self):
        ''' make entries for user to the provided group list '''
        if self.groups != None:
            for group in self.groups:
                cmd = ['groups', 'add-users', group, self.config.username]
                rval = self.openshift_cmd(cmd, oadm=True)
                if rval['returncode'] != 0:
                    return rval

                return rval

        return {'returncode': 0}

    def create(self):
        ''' create the object '''
        rval = self.create_group_entries()
        if rval['returncode'] != 0:
            return rval

        return self._create_from_content(self.config.username, self.config.data)

    def group_update(self):
        ''' update group membership '''
        rval = {'returncode': 0}
        cmd = ['get', 'groups', '-o', 'json']
        all_groups = self.openshift_cmd(cmd, output=True)

        # pylint misindentifying all_groups['results']['items'] type
        # pylint: disable=invalid-sequence-index
        for group in all_groups['results']['items']:
            # If we're supposed to be in this group
            if group['metadata']['name'] in self.groups \
               and (group['users'] is None or self.config.username not in group['users']):
                cmd = ['groups', 'add-users', group['metadata']['name'],
                       self.config.username]
                rval = self.openshift_cmd(cmd, oadm=True)
                if rval['returncode'] != 0:
                    return rval
            # else if we're in the group, but aren't supposed to be
            elif group['users'] != None and self.config.username in group['users'] \
                 and group['metadata']['name'] not in self.groups:
                cmd = ['groups', 'remove-users', group['metadata']['name'],
                       self.config.username]
                rval = self.openshift_cmd(cmd, oadm=True)
                if rval['returncode'] != 0:
                    return rval

        return rval

    def update(self):
        ''' update the object '''
        rval = self.group_update()
        if rval['returncode'] != 0:
            return rval

        # need to update the user's info
        return self._replace_content(self.kind, self.config.username, self.config.data, force=True)

    def needs_group_update(self):
        ''' check if there are group membership changes '''
        cmd = ['get', 'groups', '-o', 'json']
        all_groups = self.openshift_cmd(cmd, output=True)

        # pylint misindentifying all_groups['results']['items'] type
        # pylint: disable=invalid-sequence-index
        for group in all_groups['results']['items']:
            # If we're supposed to be in this group
            if group['metadata']['name'] in self.groups \
               and (group['users'] is None or self.config.username not in group['users']):
                return True
            # else if we're in the group, but aren't supposed to be
            elif group['users'] != None and self.config.username in group['users'] \
                 and group['metadata']['name'] not in self.groups:
                return True

        return False

    def needs_update(self):
        ''' verify an update is needed '''
        skip = []
        if self.needs_group_update():
            return True

        return not Utils.check_def_equal(self.config.data, self.user.yaml_dict, skip_keys=skip, debug=True)

    # pylint: disable=too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode=False):
        ''' run the oc_user module

            params comes from the ansible portion of this module
            check_mode: does the module support check mode. (module.check_mode)
        '''

        uconfig = UserConfig(params['kubeconfig'],
                             params['username'],
                             params['full_name'],
                            )

        oc_user = OCUser(uconfig, params['groups'],
                         verbose=params['debug'])
        state = params['state']

        api_rval = oc_user.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': "list"}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_user.exists():

                if check_mode:
                    return {'changed': False, 'msg': 'Would have performed a delete.'}

                api_rval = oc_user.delete()

                return {'changed': True, 'results': api_rval, 'state': "absent"}
            return {'changed': False, 'state': "absent"}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_user.exists():

                if check_mode:
                    return {'changed': False, 'msg': 'Would have performed a create.'}

                # Create it here
                api_rval = oc_user.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_user.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': "present"}

            ########
            # Update
            ########
            if oc_user.needs_update():
                api_rval = oc_user.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                orig_cmd = api_rval['cmd']
                # return the created object
                api_rval = oc_user.get()
                # overwrite the get/list cmd
                api_rval['cmd'] = orig_cmd

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': "present"}

            return {'changed': False, 'results': api_rval, 'state': "present"}

        return {'failed': True,
                'changed': False,
                'results': 'Unknown state passed. %s' % state,
                'state': "unknown"}
