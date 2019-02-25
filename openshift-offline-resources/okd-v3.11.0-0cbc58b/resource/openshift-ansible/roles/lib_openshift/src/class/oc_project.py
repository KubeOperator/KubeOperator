# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCProject(OpenShiftCLI):
    ''' Project Class to manage project/namespace objects'''
    kind = 'namespace'

    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for OCProject '''
        super(OCProject, self).__init__(None, config.kubeconfig)
        self.config = config
        self._project = None

    @property
    def project(self):
        ''' property for project'''
        if not self._project:
            self.get()
        return self._project

    @project.setter
    def project(self, data):
        ''' setter function for project propeorty'''
        self._project = data

    def exists(self):
        ''' return whether a project exists '''
        if self.project:
            return True

        return False

    def get(self):
        '''return project '''
        result = self._get(self.kind, self.config.name)

        if result['returncode'] == 0:
            self.project = Project(content=result['results'][0])
            result['results'] = self.project.yaml_dict

        elif 'namespaces "%s" not found' % self.config.name in result['stderr']:
            result = {'results': [], 'returncode': 0}

        return result

    def delete(self):
        '''delete the object'''
        return self._delete(self.kind, self.config.name)

    def create(self):
        '''create a project '''
        cmd = ['new-project', self.config.name]
        cmd.extend(self.config.to_option_list())

        return self.openshift_cmd(cmd, oadm=True)

    def update(self):
        '''update a project '''

        if self.config.config_options['display_name']['value'] is not None:
            self.project.update_annotation('display-name', self.config.config_options['display_name']['value'])

        if self.config.config_options['description']['value'] is not None:
            self.project.update_annotation('description', self.config.config_options['description']['value'])

        # work around for immutable project field
        if self.config.config_options['node_selector']['value'] is not None:
            self.project.update_annotation('node-selector', self.config.config_options['node_selector']['value'])

        return self._replace_content(self.kind, self.config.name, self.project.yaml_dict)

    def needs_update(self):
        ''' verify an update is needed '''
        if self.config.config_options['display_name']['value'] is not None:
            result = self.project.find_annotation("display-name")
            if result != self.config.config_options['display_name']['value']:
                return True

        if self.config.config_options['description']['value'] is not None:
            result = self.project.find_annotation("description")
            if result != self.config.config_options['description']['value']:
                return True

        if self.config.config_options['node_selector']['value'] is not None:
            result = self.project.find_annotation("node-selector")
            if result != self.config.config_options['node_selector']['value']:
                return True

        return False

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_project module'''

        node_selector = None
        if params['node_selector'] is not None:
            node_selector = ','.join(params['node_selector'])

        pconfig = ProjectConfig(
            params['name'],
            'None',
            params['kubeconfig'],
            {
                'admin': {'value': params['admin'], 'include': True},
                'admin_role': {'value': params['admin_role'], 'include': True},
                'description': {'value': params['description'], 'include': True},
                'display_name': {'value': params['display_name'], 'include': True},
                'node_selector': {'value': node_selector, 'include': True},
            },
        )

        oadm_project = OCProject(pconfig, verbose=params['debug'])

        state = params['state']

        api_rval = oadm_project.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': state}

        ########
        # Delete
        ########
        if state == 'absent':
            if oadm_project.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a delete.'}

                api_rval = oadm_project.delete()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not oadm_project.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}

                # Create it here
                api_rval = oadm_project.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oadm_project.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if oadm_project.needs_update():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed an update.'}

                api_rval = oadm_project.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oadm_project.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True,
                'changed': False,
                'msg': 'Unknown state passed. [%s]' % state}
