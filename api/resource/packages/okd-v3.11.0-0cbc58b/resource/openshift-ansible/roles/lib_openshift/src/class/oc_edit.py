# pylint: skip-file
# flake8: noqa

class Edit(OpenShiftCLI):
    ''' Class to wrap the oc command line tools
    '''
    # pylint: disable=too-many-arguments
    def __init__(self,
                 kind,
                 namespace,
                 resource_name=None,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 separator='.',
                 verbose=False):
        ''' Constructor for OpenshiftOC '''
        super(Edit, self).__init__(namespace, kubeconfig=kubeconfig, verbose=verbose)
        self.kind = kind
        self.name = resource_name
        self.separator = separator

    def get(self):
        '''return a secret by name '''
        return self._get(self.kind, self.name)

    def update(self, file_name, content, edits, force=False, content_type='yaml'):
        '''run update '''
        if file_name:
            if content_type == 'yaml':
                data = yaml.load(open(file_name))
            elif content_type == 'json':
                data = json.loads(open(file_name).read())

            yed = Yedit(filename=file_name, content=data, separator=self.separator)
            # Keep this for compatibility
            if content is not None:
                changes = []

                for key, value in content.items():
                    changes.append(yed.put(key, value))

                if any([not change[0] for change in changes]):
                    return {'returncode': 0, 'updated': False}

            elif edits is not None:
                results = Yedit.process_edits(edits, yed)

                if not results['changed']:
                    return results

            yed.write()

            atexit.register(Utils.cleanup, [file_name])

            return self._replace(file_name, force=force)

        return self._replace_content(self.kind, self.name, content, edits, force=force, sep=self.separator)

    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_edit module'''

        ocedit = Edit(params['kind'],
                      params['namespace'],
                      params['name'],
                      kubeconfig=params['kubeconfig'],
                      separator=params['separator'],
                      verbose=params['debug'])

        api_rval = ocedit.get()

        ########
        # Create
        ########
        if not Utils.exists(api_rval['results'], params['name']):
            return {"failed": True, 'msg': api_rval}

        ########
        # Update
        ########
        if check_mode:
            return {'changed': True, 'msg': 'CHECK_MODE: Would have performed edit'}

        api_rval = ocedit.update(params['file_name'],
                                 params['content'],
                                 params['edits'],
                                 params['force'],
                                 params['file_format'])

        if api_rval['returncode'] != 0:
            return {"failed": True, 'msg': api_rval}

        if 'updated' in api_rval and not api_rval['updated']:
            return {"changed": False, 'results': api_rval, 'state': 'present'}

        # return the created object
        api_rval = ocedit.get()

        if api_rval['returncode'] != 0:
            return {"failed": True, 'msg': api_rval}

        return {"changed": True, 'results': api_rval, 'state': 'present'}
