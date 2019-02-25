# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCEnv(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''

    container_path = {"pod": "spec.containers[0].env",
                      "dc":  "spec.template.spec.containers[0].env",
                      "rc":  "spec.template.spec.containers[0].env",
                     }

    # pylint allows 5. we need 6
    # pylint: disable=too-many-arguments
    def __init__(self,
                 namespace,
                 kind,
                 env_vars,
                 resource_name=None,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 verbose=False):
        ''' Constructor for OpenshiftOC '''
        super(OCEnv, self).__init__(namespace, kubeconfig=kubeconfig, verbose=verbose)
        self.kind = kind
        self.name = resource_name
        self.env_vars = env_vars
        self._resource = None

    @property
    def resource(self):
        ''' property function for resource var'''
        if not self._resource:
            self.get()
        return self._resource

    @resource.setter
    def resource(self, data):
        ''' setter function for resource var'''
        self._resource = data

    def key_value_exists(self, key, value):
        ''' return whether a key, value  pair exists '''
        return self.resource.exists_env_value(key, value)

    def key_exists(self, key):
        ''' return whether a key exists '''
        return self.resource.exists_env_key(key)

    def get(self):
        '''return environment variables '''
        result = self._get(self.kind, self.name)
        if result['returncode'] == 0:
            if self.kind == 'dc':
                self.resource = DeploymentConfig(content=result['results'][0])
                result['results'] = self.resource.get(OCEnv.container_path[self.kind]) or []
        return result

    def delete(self):
        ''' delete environment variables '''
        if self.resource.delete_env_var(self.env_vars.keys()):
            return self._replace_content(self.kind, self.name, self.resource.yaml_dict)

        return {'returncode': 0, 'changed': False}

    def put(self):
        '''place env vars into dc '''
        for update_key, update_value in self.env_vars.items():
            self.resource.update_env_var(update_key, update_value)

        return self._replace_content(self.kind, self.name, self.resource.yaml_dict)

    # pylint: disable=too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_env module'''

        ocenv = OCEnv(params['namespace'],
                      params['kind'],
                      params['env_vars'],
                      resource_name=params['name'],
                      kubeconfig=params['kubeconfig'],
                      verbose=params['debug'])

        state = params['state']

        api_rval = ocenv.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': "list"}

        ########
        # Delete
        ########
        if state == 'absent':
            for key in params.get('env_vars', {}).keys():
                if ocenv.resource.exists_env_key(key):

                    if check_mode:
                        return {'changed': False,
                                'msg': 'CHECK_MODE: Would have performed a delete.'}

                    api_rval = ocenv.delete()

                    return {'changed': True, 'state': 'absent'}

            return {'changed': False, 'state': 'absent'}

        if state == 'present':
            ########
            # Create
            ########
            for key, value in params.get('env_vars', {}).items():
                if not ocenv.key_value_exists(key, value):

                    if check_mode:
                        return {'changed': False,
                                'msg': 'CHECK_MODE: Would have performed a create.'}

                    # Create it here
                    api_rval = ocenv.put()

                    if api_rval['returncode'] != 0:
                        return {'failed': True, 'msg': api_rval}

                    # return the created object
                    api_rval = ocenv.get()

                    if api_rval['returncode'] != 0:
                        return {'failed': True, 'msg': api_rval}

                    return {'changed': True, 'results': api_rval['results'], 'state': 'present'}

            return {'changed': False, 'results': api_rval['results'], 'state': 'present'}


        return {'failed': True,
                'changed': False,
                'msg': 'Unknown state passed. %s' % state}
