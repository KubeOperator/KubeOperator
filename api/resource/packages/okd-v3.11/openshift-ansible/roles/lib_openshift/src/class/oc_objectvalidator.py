# pylint: skip-file
# flake8: noqa

# pylint: disable=too-many-instance-attributes
class OCObjectValidator(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''

    def __init__(self, kubeconfig):
        ''' Constructor for OCObjectValidator '''
        # namespace has no meaning for object validation, hardcode to 'default'
        super(OCObjectValidator, self).__init__('default', kubeconfig)

    def get_invalid(self, kind, invalid_filter):
        ''' return invalid object information '''

        rval = self._get(kind)
        if rval['returncode'] != 0:
            return False, rval, []

        return True, rval, list(filter(invalid_filter, rval['results'][0]['items']))  # wrap filter with list for py3

    # pylint: disable=too-many-return-statements
    @staticmethod
    def run_ansible(params):
        ''' run the oc_objectvalidator module

            params comes from the ansible portion of this module
        '''

        objectvalidator = OCObjectValidator(params['kubeconfig'])
        all_invalid = {}
        failed = False

        def _is_invalid_namespace(namespace):
            # check if it uses a reserved name
            name = namespace['metadata']['name']
            if not any((name == 'kube',
                        name == 'kubernetes',
                        name == 'openshift',
                        name.startswith('kube-'),
                        name.startswith('kubernetes-'),
                        name.startswith('openshift-'),)):
                return False

            # determine if the namespace was created by a user
            if 'annotations' not in namespace['metadata']:
                return False
            return 'openshift.io/requester' in namespace['metadata']['annotations']

        checks = (
            (
                'hostsubnet',
                lambda x: x['metadata']['name'] != x['host'],
                u'hostsubnets where metadata.name != host',
            ),
            (
                'netnamespace',
                lambda x: x['metadata']['name'] != x['netname'],
                u'netnamespaces where metadata.name != netname',
            ),
            (
                'namespace',
                _is_invalid_namespace,
                u'namespaces that use reserved names and were not created by infrastructure components',
            ),
        )

        for resource, invalid_filter, invalid_msg in checks:
            success, rval, invalid = objectvalidator.get_invalid(resource, invalid_filter)
            if not success:
                return {'failed': True, 'msg': 'Failed to GET {}.'.format(resource), 'state': 'list', 'results': rval}
            if invalid:
                failed = True
                all_invalid[invalid_msg] = invalid

        if failed:
            return {
                'failed': True,
                'msg': (
                    "All objects are not valid.  If you are a supported customer please contact "
                    "Red Hat Support providing the complete output above. If you are not a customer "
                    "please contact users@lists.openshift.redhat.com for assistance."
                    ),
                'state': 'list',
                'results': all_invalid
                }

        return {'msg': 'All objects are valid.'}
