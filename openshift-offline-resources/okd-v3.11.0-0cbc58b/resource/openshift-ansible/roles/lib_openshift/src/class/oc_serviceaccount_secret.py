# pylint: skip-file
# flake8: noqa

class OCServiceAccountSecret(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''

    kind = 'sa'
    def __init__(self, config, verbose=False):
        ''' Constructor for OpenshiftOC '''
        super(OCServiceAccountSecret, self).__init__(config.namespace, kubeconfig=config.kubeconfig, verbose=verbose)
        self.config = config
        self.verbose = verbose
        self._service_account = None

    @property
    def service_account(self):
        ''' Property for the service account '''
        if not self._service_account:
            self.get()
        return self._service_account

    @service_account.setter
    def service_account(self, data):
        ''' setter for the service account '''
        self._service_account = data

    def exists(self, in_secret):
        ''' verifies if secret exists in the service account '''
        result = self.service_account.find_secret(in_secret)
        if not result:
            return False
        return True

    def get(self):
        ''' get the service account definition from the master '''
        sao = self._get(OCServiceAccountSecret.kind, self.config.name)
        if sao['returncode'] == 0:
            self.service_account = ServiceAccount(content=sao['results'][0])
            sao['results'] = self.service_account.get('secrets')
        return sao

    def delete(self):
        ''' delete secrets '''

        modified = []
        for rem_secret in self.config.secrets:
            modified.append(self.service_account.delete_secret(rem_secret))

        if any(modified):
            return self._replace_content(OCServiceAccountSecret.kind, self.config.name, self.service_account.yaml_dict)

        return {'returncode': 0, 'changed': False}

    def put(self):
        ''' place secrets into sa '''
        modified = False
        for add_secret in self.config.secrets:
            if not self.service_account.find_secret(add_secret):
                self.service_account.add_secret(add_secret)
                modified = True

        if modified:
            return self._replace_content(OCServiceAccountSecret.kind, self.config.name, self.service_account.yaml_dict)

        return {'returncode': 0, 'changed': False}


    @staticmethod
    # pylint: disable=too-many-return-statements,too-many-branches
    # TODO: This function should be refactored into its individual parts.
    def run_ansible(params, check_mode):
        ''' run the oc_serviceaccount_secret module'''

        sconfig = ServiceAccountConfig(params['service_account'],
                                       params['namespace'],
                                       params['kubeconfig'],
                                       [params['secret']],
                                       None)

        oc_sa_sec = OCServiceAccountSecret(sconfig, verbose=params['debug'])

        state = params['state']

        api_rval = oc_sa_sec.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': "list"}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_sa_sec.exists(params['secret']):

                if check_mode:
                    return {'changed': True, 'msg': 'Would have removed the " + \
                            "secret from the service account.'}

                api_rval = oc_sa_sec.delete()

                return {'changed': True, 'results': api_rval, 'state': "absent"}

            return {'changed': False, 'state': "absent"}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_sa_sec.exists(params['secret']):

                if check_mode:
                    return {'changed': True, 'msg': 'Would have added the ' + \
                            'secret to the service account.'}

                # Create it here
                api_rval = oc_sa_sec.put()
                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_sa_sec.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': "present"}


            return {'changed': False, 'results': api_rval, 'state': "present"}


        return {'failed': True,
                'changed': False,
                'msg': 'Unknown state passed. %s' % state,
                'state': 'unknown'}
