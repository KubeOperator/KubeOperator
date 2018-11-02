# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-instance-attributes
class OCScale(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''

    # pylint allows 5
    # pylint: disable=too-many-arguments
    def __init__(self,
                 resource_name,
                 namespace,
                 replicas,
                 kind,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 verbose=False):
        ''' Constructor for OCScale '''
        super(OCScale, self).__init__(namespace, kubeconfig=kubeconfig, verbose=verbose)
        self.kind = kind
        self.replicas = replicas
        self.name = resource_name
        self._resource = None

    @property
    def resource(self):
        ''' property function for resource var '''
        if not self._resource:
            self.get()
        return self._resource

    @resource.setter
    def resource(self, data):
        ''' setter function for resource var '''
        self._resource = data

    def get(self):
        '''return replicas information '''
        vol = self._get(self.kind, self.name)
        if vol['returncode'] == 0:
            if self.kind == 'dc':
                # The resource returned from a query could be an rc or dc.
                # pylint: disable=redefined-variable-type
                self.resource = DeploymentConfig(content=vol['results'][0])
                vol['results'] = [self.resource.get_replicas()]
            if self.kind == 'rc':
                # The resource returned from a query could be an rc or dc.
                # pylint: disable=redefined-variable-type
                self.resource = ReplicationController(content=vol['results'][0])
                vol['results'] = [self.resource.get_replicas()]

        return vol

    def put(self):
        '''update replicas into dc '''
        self.resource.update_replicas(self.replicas)
        return self._replace_content(self.kind, self.name, self.resource.yaml_dict)

    def needs_update(self):
        ''' verify whether an update is needed '''
        return self.resource.needs_update_replicas(self.replicas)

    # pylint: disable=too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_scale module'''

        oc_scale = OCScale(params['name'],
                           params['namespace'],
                           params['replicas'],
                           params['kind'],
                           params['kubeconfig'],
                           verbose=params['debug'])

        state = params['state']

        api_rval = oc_scale.get()
        if api_rval['returncode'] != 0:
            return {'failed': True, 'msg': api_rval}

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'result': api_rval['results'], 'state': 'list'}  # noqa: E501

        elif state == 'present':
            ########
            # Update
            ########
            if oc_scale.needs_update():
                if check_mode:
                    return {'changed': True, 'result': 'CHECK_MODE: Would have updated.'}  # noqa: E501
                api_rval = oc_scale.put()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_scale.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'result': api_rval['results'], 'state': 'present'}  # noqa: E501

            return {'changed': False, 'result': api_rval['results'], 'state': 'present'}  # noqa: E501

        return {'failed': True, 'msg': 'Unknown state passed. [{}]'.format(state)}
