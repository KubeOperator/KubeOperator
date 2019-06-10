# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCPVC(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'pvc'

    # pylint allows 5
    # pylint: disable=too-many-arguments
    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for OCVolume '''
        super(OCPVC, self).__init__(config.namespace, config.kubeconfig)
        self.config = config
        self.namespace = config.namespace
        self._pvc = None

    @property
    def pvc(self):
        ''' property function pvc'''
        if not self._pvc:
            self.get()
        return self._pvc

    @pvc.setter
    def pvc(self, data):
        ''' setter function for yedit var '''
        self._pvc = data

    def bound(self):
        '''return whether the pvc is bound'''
        if self.pvc.get_volume_name():
            return True

        return False

    def exists(self):
        ''' return whether a pvc exists '''
        if self.pvc:
            return True

        return False

    def get(self):
        '''return pvc information '''
        result = self._get(self.kind, self.config.name)
        if result['returncode'] == 0:
            self.pvc = PersistentVolumeClaim(content=result['results'][0])
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
        return self._replace_content(self.kind, self.config.name, self.config.data)

    def needs_update(self):
        ''' verify an update is needed '''
        if self.pvc.get_volume_name() or self.pvc.is_bound():
            return False

        skip = []
        return not Utils.check_def_equal(self.config.data, self.pvc.yaml_dict, skip_keys=skip, debug=True)

    # pylint: disable=too-many-branches,too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_pvc module'''
        pconfig = PersistentVolumeClaimConfig(params['name'],
                                              params['namespace'],
                                              params['kubeconfig'],
                                              params['access_modes'],
                                              params['volume_capacity'],
                                              params['selector'],
                                              params['storage_class_name'],
                                             )
        oc_pvc = OCPVC(pconfig, verbose=params['debug'])

        state = params['state']

        api_rval = oc_pvc.get()
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
            if oc_pvc.exists():

                if check_mode:
                    return {'changed': False, 'msg': 'CHECK_MODE: Would have performed a delete.'}

                api_rval = oc_pvc.delete()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_pvc.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}

                # Create it here
                api_rval = oc_pvc.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_pvc.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if oc_pvc.pvc.is_bound() or oc_pvc.pvc.get_volume_name():
                api_rval['msg'] = '##### - This volume is currently bound.  Will not update - ####'
                return {'changed': False, 'results': api_rval, 'state': state}

            if oc_pvc.needs_update():
                api_rval = oc_pvc.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_pvc.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True, 'msg': 'Unknown state passed. {}'.format(state)}
