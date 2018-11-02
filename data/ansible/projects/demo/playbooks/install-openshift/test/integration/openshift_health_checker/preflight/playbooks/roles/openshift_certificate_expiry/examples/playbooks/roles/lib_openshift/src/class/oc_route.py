# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCRoute(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'route'

    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for OCVolume '''
        super(OCRoute, self).__init__(config.namespace, kubeconfig=config.kubeconfig, verbose=verbose)
        self.config = config
        self._route = None

    @property
    def route(self):
        ''' property function for route'''
        if not self._route:
            self.get()
        return self._route

    @route.setter
    def route(self, data):
        ''' setter function for route '''
        self._route = data

    def exists(self):
        ''' return whether a route exists '''
        if self.route:
            return True

        return False

    def get(self):
        '''return route information '''
        result = self._get(self.kind, self.config.name)
        if result['returncode'] == 0:
            self.route = Route(content=result['results'][0])
        elif 'routes \"%s\" not found' % self.config.name in result['stderr']:
            result['returncode'] = 0
            result['results'] = [{}]
        elif 'namespaces \"%s\" not found' % self.config.namespace in result['stderr']:
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
        return self._replace_content(self.kind,
                                     self.config.name,
                                     self.config.data,
                                     force=(self.config.host != self.route.get_host()))

    def needs_update(self):
        ''' verify an update is needed '''
        skip = []
        return not Utils.check_def_equal(self.config.data, self.route.yaml_dict, skip_keys=skip, debug=self.verbose)

    @staticmethod
    def get_cert_data(path, content):
        '''get the data for a particular value'''
        rval = None
        if path and os.path.exists(path) and os.access(path, os.R_OK):
            rval = open(path).read()
        elif content:
            rval = content

        return rval

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def run_ansible(params, check_mode=False):
        ''' run the oc_route module

            params comes from the ansible portion for this module
            files: a dictionary for the certificates
                   {'cert': {'path': '',
                             'content': '',
                             'value': ''
                            }
                   }
            check_mode: does the module support check mode.  (module.check_mode)
        '''
        files = {'destcacert': {'path': params['dest_cacert_path'],
                                'content': params['dest_cacert_content'],
                                'value': None, },
                 'cacert': {'path': params['cacert_path'],
                            'content': params['cacert_content'],
                            'value': None, },
                 'cert': {'path': params['cert_path'],
                          'content': params['cert_content'],
                          'value': None, },
                 'key': {'path': params['key_path'],
                         'content': params['key_content'],
                         'value': None, }, }

        if params['tls_termination'] and params['tls_termination'].lower() != 'passthrough':  # E501

            for key, option in files.items():
                if not option['path'] and not option['content']:
                    continue

                option['value'] = OCRoute.get_cert_data(option['path'], option['content'])  # E501

                if not option['value']:
                    return {'failed': True,
                            'msg': 'Verify that you pass a correct value for %s' % key}

        rconfig = RouteConfig(params['name'],
                              params['namespace'],
                              params['kubeconfig'],
                              params['labels'],
                              files['destcacert']['value'],
                              files['cacert']['value'],
                              files['cert']['value'],
                              files['key']['value'],
                              params['host'],
                              params['tls_termination'],
                              params['service_name'],
                              params['wildcard_policy'],
                              params['weight'],
                              params['port'])

        oc_route = OCRoute(rconfig, verbose=params['debug'])

        state = params['state']

        api_rval = oc_route.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False,
                    'results': api_rval['results'],
                    'state': 'list'}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_route.exists():

                if check_mode:
                    return {'changed': False, 'msg': 'CHECK_MODE: Would have performed a delete.'}  # noqa: E501

                api_rval = oc_route.delete()

                return {'changed': True, 'results': api_rval, 'state': "absent"}  # noqa: E501
            return {'changed': False, 'state': 'absent'}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_route.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}  # noqa: E501

                # Create it here
                api_rval = oc_route.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval, 'state': "present"}  # noqa: E501

                # return the created object
                api_rval = oc_route.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval, 'state': "present"}  # noqa: E501

                return {'changed': True, 'results': api_rval, 'state': "present"}  # noqa: E501

            ########
            # Update
            ########
            if oc_route.needs_update():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed an update.'}  # noqa: E501

                api_rval = oc_route.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval, 'state': "present"}  # noqa: E501

                # return the created object
                api_rval = oc_route.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval, 'state': "present"}  # noqa: E501

                return {'changed': True, 'results': api_rval, 'state': "present"}  # noqa: E501

            return {'changed': False, 'results': api_rval, 'state': "present"}

        # catch all
        return {'failed': True, 'msg': "Unknown State passed"}
