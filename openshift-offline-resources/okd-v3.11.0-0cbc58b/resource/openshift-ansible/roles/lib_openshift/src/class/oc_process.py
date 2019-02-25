# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCProcess(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''

    # pylint allows 5. we need 6
    # pylint: disable=too-many-arguments
    def __init__(self,
                 namespace,
                 tname=None,
                 params=None,
                 create=False,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 tdata=None,
                 verbose=False):
        ''' Constructor for OpenshiftOC '''
        super(OCProcess, self).__init__(namespace, kubeconfig=kubeconfig, verbose=verbose)
        self.name = tname
        self.data = tdata
        self.params = params
        self.create = create
        self._template = None

    @property
    def template(self):
        '''template property'''
        if self._template is None:
            results = self._process(self.name, False, self.params, self.data)
            if results['returncode'] != 0:
                raise OpenShiftCLIError('Error processing template [%s]: %s' %(self.name, results))
            self._template = results['results']['items']

        return self._template

    def get(self):
        '''get the template'''
        results = self._get('template', self.name)
        if results['returncode'] != 0:
            # Does the template exist??
            if 'not found' in results['stderr']:
                results['returncode'] = 0
                results['exists'] = False
                results['results'] = []

        return results

    def delete(self, obj):
        '''delete a resource'''
        return self._delete(obj['kind'], obj['metadata']['name'])

    def create_obj(self, obj):
        '''create a resource'''
        return self._create_from_content(obj['metadata']['name'], obj)

    def process(self, create=None):
        '''process a template'''
        do_create = False
        if create != None:
            do_create = create
        else:
            do_create = self.create

        return self._process(self.name, do_create, self.params, self.data)

    def exists(self):
        '''return whether the template exists'''
        # Always return true if we're being passed template data
        if self.data:
            return True
        t_results = self._get('template', self.name)

        if t_results['returncode'] != 0:
            # Does the template exist??
            if 'not found' in t_results['stderr']:
                return False
            else:
                raise OpenShiftCLIError('Something went wrong. %s' % t_results)

        return True

    def needs_update(self):
        '''attempt to process the template and return it for comparison with oc objects'''
        obj_results = []
        for obj in self.template:

            # build a list of types to skip
            skip = []

            if obj['kind'] == 'ServiceAccount':
                skip.extend(['secrets', 'imagePullSecrets'])
            if obj['kind'] == 'BuildConfig':
                skip.extend(['lastTriggeredImageID'])
            if obj['kind'] == 'ImageStream':
                skip.extend(['generation'])
            if obj['kind'] == 'DeploymentConfig':
                skip.extend(['lastTriggeredImage'])

            # fetch the current object
            curr_obj_results = self._get(obj['kind'], obj['metadata']['name'])
            if curr_obj_results['returncode'] != 0:
                # Does the template exist??
                if 'not found' in curr_obj_results['stderr']:
                    obj_results.append((obj, True))
                    continue

            # check the generated object against the existing object
            if not Utils.check_def_equal(obj, curr_obj_results['results'][0], skip_keys=skip):
                obj_results.append((obj, True))
                continue

            obj_results.append((obj, False))

        return obj_results

    # pylint: disable=too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_process module'''

        ocprocess = OCProcess(params['namespace'],
                              params['template_name'],
                              params['params'],
                              params['create'],
                              kubeconfig=params['kubeconfig'],
                              tdata=params['content'],
                              verbose=params['debug'])

        state = params['state']

        api_rval = ocprocess.get()

        if state == 'list':
            if api_rval['returncode'] != 0:
                return {"failed": True, "msg" : api_rval}

            return {"changed" : False, "results": api_rval, "state": state}

        elif state == 'present':
            if check_mode and params['create']:
                return {"changed": True, 'msg': "CHECK_MODE: Would have processed template."}

            if not ocprocess.exists() or not params['reconcile']:
            #FIXME: this code will never get run in a way that succeeds when
            #       module.params['reconcile'] is true. Because oc_process doesn't
            #       create the actual template, the check of ocprocess.exists()
            #       is meaningless. Either it's already here and this code
            #       won't be run, or this code will fail because there is no
            #       template available for oc process to use. Have we conflated
            #       the template's existence with the existence of the objects
            #       it describes?

            # Create it here
                api_rval = ocprocess.process()
                if api_rval['returncode'] != 0:
                    return {"failed": True, "msg": api_rval}

                if params['create']:
                    return {"changed": True, "results": api_rval, "state": state}

                return {"changed": False, "results": api_rval, "state": state}

        # verify results
        update = False
        rval = []
        all_results = ocprocess.needs_update()
        for obj, status in all_results:
            if status:
                ocprocess.delete(obj)
                results = ocprocess.create_obj(obj)
                results['kind'] = obj['kind']
                rval.append(results)
                update = True

        if not update:
            return {"changed": update, "results": api_rval, "state": state}

        for cmd in rval:
            if cmd['returncode'] != 0:
                return {"failed": True, "changed": update, "msg": rval, "state": state}

        return {"changed": update, "results": rval, "state": state}
