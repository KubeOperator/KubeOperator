# flake8: noqa
# pylint: skip-file


# pylint: disable=too-many-instance-attributes
class OCVersion(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    # pylint allows 5
    # pylint: disable=too-many-arguments
    def __init__(self,
                 config,
                 debug):
        ''' Constructor for OCVersion '''
        super(OCVersion, self).__init__(None, config)
        self.debug = debug

    def get(self):
        '''get and return version information '''

        results = {}

        version_results = self._version()

        if version_results['returncode'] == 0:
            filtered_vers = Utils.filter_versions(version_results['results'])
            custom_vers = Utils.add_custom_versions(filtered_vers)

            results['returncode'] = version_results['returncode']
            results.update(filtered_vers)
            results.update(custom_vers)

            return results

        raise OpenShiftCLIError('Problem detecting openshift version.')

    @staticmethod
    def run_ansible(params):
        '''run the oc_version module'''
        oc_version = OCVersion(params['kubeconfig'], params['debug'])

        if params['state'] == 'list':

            #pylint: disable=protected-access
            result = oc_version.get()
            return {'state': params['state'],
                    'results': result,
                    'changed': False}
