# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCService(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    kind = 'service'

    # pylint allows 5
    # pylint: disable=too-many-arguments
    def __init__(self,
                 sname,
                 namespace,
                 labels,
                 annotations,
                 selector,
                 cluster_ip,
                 portal_ip,
                 ports,
                 session_affinity,
                 service_type,
                 external_ips,
                 kubeconfig='/etc/origin/master/admin.kubeconfig',
                 verbose=False):
        ''' Constructor for OCVolume '''
        super(OCService, self).__init__(namespace, kubeconfig, verbose)
        self.namespace = namespace
        self.config = ServiceConfig(sname, namespace, ports, annotations, selector, labels,
                                    cluster_ip, portal_ip, session_affinity, service_type,
                                    external_ips)
        self.user_svc = Service(content=self.config.data)
        self.svc = None

    @property
    def service(self):
        ''' property function service'''
        if not self.svc:
            self.get()
        return self.svc

    @service.setter
    def service(self, data):
        ''' setter function for service var '''
        self.svc = data

    def exists(self):
        ''' return whether a service exists '''
        if self.service:
            return True

        return False

    def get(self):
        '''return service information '''
        result = self._get(self.kind, self.config.name)
        if result['returncode'] == 0:
            self.service = Service(content=result['results'][0])
            result['clusterip'] = self.service.get('spec.clusterIP')
        elif 'services \"%s\" not found' % self.config.name  in result['stderr']:
            result['clusterip'] = ''
            result['returncode'] = 0
        elif 'namespaces \"%s\" not found' % self.config.namespace  in result['stderr']:
            result['clusterip'] = ''
            result['returncode'] = 0

        return result

    def delete(self):
        '''delete the service'''
        return self._delete(self.kind, self.config.name)

    def create(self):
        '''create a service '''
        return self._create_from_content(self.config.name, self.user_svc.yaml_dict)

    def update(self):
        '''create a service '''
        # Need to copy over the portalIP and the serviceIP settings

        self.user_svc.add_cluster_ip(self.service.get('spec.clusterIP'))
        self.user_svc.add_portal_ip(self.service.get('spec.portalIP'))
        return self._replace_content(self.kind, self.config.name, self.user_svc.yaml_dict)

    def needs_update(self):
        ''' verify an update is needed '''
        skip = ['clusterIP', 'portalIP']
        return not Utils.check_def_equal(self.user_svc.yaml_dict, self.service.yaml_dict, skip_keys=skip, debug=True)

    # pylint: disable=too-many-return-statements,too-many-branches
    @staticmethod
    def run_ansible(params, check_mode):
        '''Run the oc_service module'''
        oc_svc = OCService(params['name'],
                           params['namespace'],
                           params['labels'],
                           params['annotations'],
                           params['selector'],
                           params['clusterip'],
                           params['portalip'],
                           params['ports'],
                           params['session_affinity'],
                           params['service_type'],
                           params['external_ips'],
                           params['kubeconfig'],
                           params['debug'])

        state = params['state']

        api_rval = oc_svc.get()

        if api_rval['returncode'] != 0:
            return {'failed': True, 'msg': api_rval}

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval, 'state': state}

        ########
        # Delete
        ########
        if state == 'absent':
            if oc_svc.exists():

                if check_mode:
                    return {'changed': True,
                            'msg': 'CHECK_MODE: Would have performed a delete.'}  # noqa: E501

                api_rval = oc_svc.delete()

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not oc_svc.exists():

                if check_mode:
                    return {'changed': True,
                            'msg': 'CHECK_MODE: Would have performed a create.'}  # noqa: E501

                # Create it here
                api_rval = oc_svc.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_svc.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if oc_svc.needs_update():
                api_rval = oc_svc.update()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_svc.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True, 'msg': 'UNKNOWN state passed. [%s]' % state}
