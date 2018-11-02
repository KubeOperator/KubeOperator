# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-instance-attributes
class OCServiceAccount(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'sa'

    # pylint allows 5
    # pylint: disable=too-many-arguments
    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for OCVolume '''
        super(OCServiceAccount, self).__init__(config.namespace, kubeconfig=config.kubeconfig, verbose=verbose)
        self.config = config
        self.service_account = None

    def exists(self):
        ''' return whether a volume exists '''
        if self.service_account:
            return True

        return False

    def get(self):
        '''return volume information '''
        result = self._get(self.kind, self.config.name)
        if result['returncode'] == 0:
            self.service_account = ServiceAccount(content=result['results'][0])
        elif '\"%s\" not found' % self.config.name in result['stderr']:
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
        # need to update the tls information and the service name
        for secret in self.config.secrets:
            result = self.service_account.find_secret(secret)
            if not result:
                self.service_account.add_secret(secret)

        for secret in self.config.image_pull_secrets:
            result = self.service_account.find_image_pull_secret(secret)
            if not result:
                self.service_account.add_image_pull_secret(secret)

        return self._replace_content(self.kind, self.config.name, self.config.data)

    def needs_update(self):
        ''' verify an update is needed '''
        # since creating an service account generates secrets and imagepullsecrets
        # check_def_equal will not work
        # Instead, verify all secrets passed are in the list
        for secret in self.config.secrets:
            result = self.service_account.find_secret(secret)
            if not result:
                return True

        for secret in self.config.image_pull_secrets:
            result = self.service_account.find_image_pull_secret(secret)
            if not result:
                return True

        return False

    @staticmethod
    # pylint: disable=too-many-return-statements,too-many-branches
    # TODO: This function should be refactored into its individual parts.
    def run_ansible(params, check_mode):
        '''run the oc_serviceaccount module'''

        rconfig = ServiceAccountConfig(params['name'],
                                       params['namespace'],
                                       params['kubeconfig'],
                                       params['secrets'],
                                       params['image_pull_secrets'],
                                      )

        oc_sa = OCServiceAccount(rconfig,
                                 verbose=params['debug'])

        state = params['state']

        api_rval = oc_sa.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': 'list'}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_sa.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'Would have performed a delete.'}

                api_rval = oc_sa.delete()

                return {'changed': True, 'results': api_rval, 'state': 'absent'}

            return {'changed': False, 'state': 'absent'}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_sa.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'Would have performed a create.'}

                # Create it here
                api_rval = oc_sa.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_sa.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': 'present'}

            ########
            # Update
            ########
            if oc_sa.needs_update():
                api_rval = oc_sa.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_sa.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': 'present'}

            return {'changed': False, 'results': api_rval, 'state': 'present'}


        return {'failed': True,
                'changed': False,
                'msg': 'Unknown state passed. %s' % state,
                'state': 'unknown'}
