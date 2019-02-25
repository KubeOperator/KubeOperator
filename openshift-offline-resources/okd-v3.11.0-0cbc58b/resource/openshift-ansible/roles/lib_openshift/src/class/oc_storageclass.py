# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-instance-attributes
class OCStorageClass(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'storageclass'

    # pylint allows 5
    # pylint: disable=too-many-arguments
    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for OCStorageClass '''
        super(OCStorageClass, self).__init__(None, kubeconfig=config.kubeconfig, verbose=verbose)
        self.config = config
        self.storage_class = None

    def exists(self):
        ''' return whether a storageclass exists'''
        if self.storage_class:
            return True

        return False

    def get(self):
        '''return storageclass '''
        result = self._get(self.kind, self.config.name)
        if result['returncode'] == 0:
            self.storage_class = StorageClass(content=result['results'][0])
        elif '\"%s\" not found' % self.config.name in result['stderr']:
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
        # parameters are currently unable to be updated.  need to delete and recreate
        self.delete()
        # pause here and attempt to wait for delete.
        # Better option would be to poll
        time.sleep(5)
        return self.create()

    def needs_update(self):
        ''' verify an update is needed '''
        # check if params have updated
        if self.storage_class.get_parameters() != self.config.parameters:
            return True

        for anno_key, anno_value in self.storage_class.get_annotations().items():
            if 'is-default-class' in anno_key and anno_value != self.config.default_storage_class:
                return True

        # check if mount options have updated
        if set(self.storage_class.get_mount_options()) != set(self.config.mount_options):
            return True

        # check if reclaim policy has been updated
        if self.storage_class.get_reclaim_policy() != self.config.reclaim_policy:
            return True

        return False

    @staticmethod
    def provisioner_name_qualified(provisioner_name):
        pattern = re.compile(r'^[a-z0-9A-Z-_.]+\/[a-z0-9A-Z-_.]+$')
        return pattern.match(provisioner_name)

    @staticmethod
    # pylint: disable=too-many-return-statements,too-many-branches
    # TODO: This function should be refactored into its individual parts.
    def run_ansible(params, check_mode):
        '''run the oc_storageclass module'''

        # Make sure that the provisioner is fully qualified before using it
        # E.g. if 'aws-efs' is provided as a provisioner, convert it to 'kubernetes.io/aws-efs'
        # but if the name is already qualified  (e.g. 'openshift.org/aws-efs') then leave it be.
        raw_provisioner_name = params['provisioner']
        if OCStorageClass.provisioner_name_qualified(raw_provisioner_name):
            qualified_provisioner_name = raw_provisioner_name
        else:
            qualified_provisioner_name = "kubernetes.io/{}".format(params['provisioner'])

        rconfig = StorageClassConfig(params['name'],
                                     provisioner=qualified_provisioner_name,
                                     parameters=params['parameters'],
                                     annotations=params['annotations'],
                                     api_version="storage.k8s.io/{}".format(params['api_version']),
                                     default_storage_class=params.get('default_storage_class', 'false'),
                                     kubeconfig=params['kubeconfig'],
                                     mount_options=params['mount_options'],
                                     reclaim_policy=params['reclaim_policy']
                                    )

        oc_sc = OCStorageClass(rconfig, verbose=params['debug'])

        state = params['state']

        api_rval = oc_sc.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': 'list'}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_sc.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'Would have performed a delete.'}

                api_rval = oc_sc.delete()

                return {'changed': True, 'results': api_rval, 'state': 'absent'}

            return {'changed': False, 'state': 'absent'}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_sc.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'Would have performed a create.'}

                # Create it here
                api_rval = oc_sc.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_sc.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': 'present'}

            ########
            # Update
            ########
            if oc_sc.needs_update():
                api_rval = oc_sc.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_sc.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': 'present'}

            return {'changed': False, 'results': api_rval, 'state': 'present'}


        return {'failed': True,
                'changed': False,
                'msg': 'Unknown state passed. %s' % state,
                'state': 'unknown'}
