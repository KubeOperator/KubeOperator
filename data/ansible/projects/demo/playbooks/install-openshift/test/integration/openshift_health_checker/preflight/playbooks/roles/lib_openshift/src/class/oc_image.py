# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-arguments
class OCImage(OpenShiftCLI):
    ''' Class to import and create an imagestream object'''
    def __init__(self,
                 namespace,
                 registry_url,
                 image_name,
                 image_tag,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 verbose=False):
        ''' Constructor for OCImage'''
        super(OCImage, self).__init__(namespace, kubeconfig)
        self.registry_url = registry_url
        self.image_name = image_name
        self.image_tag = image_tag
        self.verbose = verbose

    def get(self):
        '''return a image by name '''
        results = self._get('imagestream', self.image_name)
        results['exists'] = False
        if results['returncode'] == 0 and results['results'][0]:
            results['exists'] = True

        if results['returncode'] != 0 and '"{}" not found'.format(self.image_name) in results['stderr']:
            results['returncode'] = 0

        return results

    def create(self, url=None, name=None, tag=None):
        '''Create an image '''
        return self._import_image(url, name, tag)


    # pylint: disable=too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode):
        ''' run the oc_image module'''

        ocimage = OCImage(params['namespace'],
                          params['registry_url'],
                          params['image_name'],
                          params['image_tag'],
                          kubeconfig=params['kubeconfig'],
                          verbose=params['debug'])

        state = params['state']

        api_rval = ocimage.get()

        #####
        # Get
        #####
        if state == 'list':
            if api_rval['returncode'] != 0:
                return {"failed": True, "msg": api_rval}
            return {"changed": False, "results": api_rval, "state": "list"}

        ########
        # Create
        ########
        if state == 'present':

            if not Utils.exists(api_rval['results'], params['image_name']):

                if check_mode:
                    return {"changed": False, "msg": 'CHECK_MODE: Would have performed a create'}

                api_rval = ocimage.create(params['registry_url'],
                                          params['image_name'],
                                          params['image_tag'])

                if api_rval['returncode'] != 0:
                    return {"failed": True, "msg": api_rval}

                # return the newly created object
                api_rval = ocimage.get()

                if api_rval['returncode'] != 0:
                    return {"failed": True, "msg": api_rval}

                return {"changed": True, "results": api_rval, "state": "present"}

            # image exists, no change
            return {"changed": False, "results": api_rval, "state": "present"}

        return {"failed": True, "changed": False, "msg": "Unknown state passed. {0}".format(state)}
