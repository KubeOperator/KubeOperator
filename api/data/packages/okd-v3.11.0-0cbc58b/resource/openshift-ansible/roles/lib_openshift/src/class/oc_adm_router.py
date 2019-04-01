# pylint: skip-file
# flake8: noqa


class RouterException(Exception):
    ''' Router exception'''
    pass


class RouterConfig(OpenShiftCLIConfig):
    ''' RouterConfig is a DTO for the router.  '''
    def __init__(self, rname, namespace, kubeconfig, router_options):
        super(RouterConfig, self).__init__(rname, namespace, kubeconfig, router_options)


class Router(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''
    def __init__(self,
                 router_config,
                 verbose=False):
        ''' Constructor for OpenshiftOC

           a router consists of 3 or more parts
           - dc/router
           - svc/router
           - sa/router
           - secret/router-certs
           - clusterrolebinding/router-router-role
        '''
        super(Router, self).__init__(router_config.namespace, router_config.kubeconfig, verbose)
        self.config = router_config
        self.verbose = verbose
        self.router_parts = [{'kind': 'dc', 'name': self.config.name},
                             {'kind': 'svc', 'name': self.config.name},
                             {'kind': 'sa', 'name': self.config.config_options['service_account']['value']},
                             {'kind': 'secret', 'name': self.config.name + '-certs'},
                             {'kind': 'clusterrolebinding', 'name': 'router-' + self.config.name + '-role'},
                            ]

        self.__prepared_router = None
        self.dconfig = None
        self.svc = None
        self._secret = None
        self._serviceaccount = None
        self._rolebinding = None

    @property
    def prepared_router(self):
        ''' property for the prepared router'''
        if self.__prepared_router is None:
            results = self._prepare_router()
            if not results or 'returncode' in results and results['returncode'] != 0:
                if 'stderr' in results:
                    raise RouterException('Could not perform router preparation: %s' % results['stderr'])

                raise RouterException('Could not perform router preparation.')
            self.__prepared_router = results

        return self.__prepared_router

    @prepared_router.setter
    def prepared_router(self, obj):
        '''setter for the prepared_router'''
        self.__prepared_router = obj

    @property
    def deploymentconfig(self):
        ''' property deploymentconfig'''
        return self.dconfig

    @deploymentconfig.setter
    def deploymentconfig(self, config):
        ''' setter for property deploymentconfig '''
        self.dconfig = config

    @property
    def service(self):
        ''' property for service '''
        return self.svc

    @service.setter
    def service(self, config):
        ''' setter for property service '''
        self.svc = config

    @property
    def secret(self):
        ''' property secret '''
        return self._secret

    @secret.setter
    def secret(self, config):
        ''' setter for property secret '''
        self._secret = config

    @property
    def serviceaccount(self):
        ''' property for serviceaccount '''
        return self._serviceaccount

    @serviceaccount.setter
    def serviceaccount(self, config):
        ''' setter for property serviceaccount '''
        self._serviceaccount = config

    @property
    def rolebinding(self):
        ''' property rolebinding '''
        return self._rolebinding

    @rolebinding.setter
    def rolebinding(self, config):
        ''' setter for property rolebinding '''
        self._rolebinding = config

    def get_object_by_kind(self, kind):
        '''return the current object kind by name'''
        if re.match("^(dc|deploymentconfig)$", kind, flags=re.IGNORECASE):
            return self.deploymentconfig
        elif re.match("^(svc|service)$", kind, flags=re.IGNORECASE):
            return self.service
        elif re.match("^(sa|serviceaccount)$", kind, flags=re.IGNORECASE):
            return self.serviceaccount
        elif re.match("secret", kind, flags=re.IGNORECASE):
            return self.secret
        elif re.match("clusterrolebinding", kind, flags=re.IGNORECASE):
            return self.rolebinding

        return None

    def get(self):
        ''' return the self.router_parts '''
        self.service = None
        self.deploymentconfig = None
        self.serviceaccount = None
        self.secret = None
        self.rolebinding = None
        for part in self.router_parts:
            result = self._get(part['kind'], name=part['name'])
            if result['returncode'] == 0 and part['kind'] == 'dc':
                self.deploymentconfig = DeploymentConfig(result['results'][0])
            elif result['returncode'] == 0 and part['kind'] == 'svc':
                self.service = Service(content=result['results'][0])
            elif result['returncode'] == 0 and part['kind'] == 'sa':
                self.serviceaccount = ServiceAccount(content=result['results'][0])
            elif result['returncode'] == 0 and part['kind'] == 'secret':
                self.secret = Secret(content=result['results'][0])
            elif result['returncode'] == 0 and part['kind'] == 'clusterrolebinding':
                self.rolebinding = RoleBinding(content=result['results'][0])

        return {'deploymentconfig': self.deploymentconfig,
                'service': self.service,
                'serviceaccount': self.serviceaccount,
                'secret': self.secret,
                'clusterrolebinding': self.rolebinding,
               }

    def exists(self):
        '''return a whether svc or dc exists '''
        if self.deploymentconfig and self.service and self.secret and self.serviceaccount:
            return True

        return False

    def delete(self):
        '''return all pods '''
        parts = []
        for part in self.router_parts:
            parts.append(self._delete(part['kind'], part['name']))

        rval = 0
        for part in parts:
            if part['returncode'] != 0 and not 'already exist' in part['stderr']:
                rval = part['returncode']

        return {'returncode': rval, 'results': parts}

    def add_modifications(self, deploymentconfig):
        '''modify the deployment config'''
        # We want modifications in the form of edits coming in from the module.
        # Let's apply these here

        # If extended validation is enabled, set the corresponding environment
        # variable.
        if self.config.config_options['extended_validation']['value']:
            if not deploymentconfig.exists_env_key('EXTENDED_VALIDATION'):
                deploymentconfig.add_env_value('EXTENDED_VALIDATION', "true")
            else:
                deploymentconfig.update_env_var('EXTENDED_VALIDATION', "true")

        # Apply any edits.
        edit_results = []
        for edit in self.config.config_options['edits'].get('value', []):
            if edit['action'] == 'put':
                edit_results.append(deploymentconfig.put(edit['key'],
                                                         edit['value']))
            if edit['action'] == 'update':
                edit_results.append(deploymentconfig.update(edit['key'],
                                                            edit['value'],
                                                            edit.get('index', None),
                                                            edit.get('curr_value', None)))
            if edit['action'] == 'append':
                edit_results.append(deploymentconfig.append(edit['key'],
                                                            edit['value']))

        if edit_results and not any([res[0] for res in edit_results]):
            return None

        return deploymentconfig

    # pylint: disable=too-many-branches
    def _prepare_router(self):
        '''prepare router for instantiation'''
        # if cacert, key, and cert were passed, combine them into a pem file
        if (self.config.config_options['cacert_file']['value'] and
                self.config.config_options['cert_file']['value'] and
                self.config.config_options['key_file']['value']):

            router_pem = '/tmp/router.pem'
            with open(router_pem, 'w') as rfd:
                rfd.write(open(self.config.config_options['cert_file']['value']).read())
                rfd.write(open(self.config.config_options['key_file']['value']).read())
                if self.config.config_options['cacert_file']['value'] and \
                   os.path.exists(self.config.config_options['cacert_file']['value']):
                    rfd.write(open(self.config.config_options['cacert_file']['value']).read())

            atexit.register(Utils.cleanup, [router_pem])

            self.config.config_options['default_cert']['value'] = router_pem

        elif self.config.config_options['default_cert']['value'] is None:
            # No certificate was passed to us.  do not pass one to oc adm router
            self.config.config_options['default_cert']['include'] = False

        options = self.config.to_option_list(ascommalist='labels')

        cmd = ['router', self.config.name]
        cmd.extend(options)
        cmd.extend(['--dry-run=True', '-o', 'json'])

        results = self.openshift_cmd(cmd, oadm=True, output=True, output_type='json')

        # pylint: disable=maybe-no-member
        if results['returncode'] != 0 or 'items' not in results['results']:
            return results

        oc_objects = {'DeploymentConfig': {'obj': None, 'path': None, 'update': False},
                      'Secret': {'obj': None, 'path': None, 'update': False},
                      'ServiceAccount': {'obj': None, 'path': None, 'update': False},
                      'ClusterRoleBinding': {'obj': None, 'path': None, 'update': False},
                      'Service': {'obj': None, 'path': None, 'update': False},
                     }
        # pylint: disable=invalid-sequence-index
        for res in results['results']['items']:
            if res['kind'] == 'DeploymentConfig':
                oc_objects['DeploymentConfig']['obj'] = DeploymentConfig(res)
            elif res['kind'] == 'Service':
                oc_objects['Service']['obj'] = Service(res)
            elif res['kind'] == 'ServiceAccount':
                oc_objects['ServiceAccount']['obj'] = ServiceAccount(res)
            elif res['kind'] == 'Secret':
                oc_objects['Secret']['obj'] = Secret(res)
            elif res['kind'] == 'ClusterRoleBinding':
                oc_objects['ClusterRoleBinding']['obj'] = RoleBinding(res)

        # Currently only deploymentconfig needs updating
        # Verify we got a deploymentconfig
        if not oc_objects['DeploymentConfig']['obj']:
            return results

        # add modifications added
        oc_objects['DeploymentConfig']['obj'] = self.add_modifications(oc_objects['DeploymentConfig']['obj'])

        for oc_type, oc_data in oc_objects.items():
            if oc_data['obj'] is not None:
                oc_data['path'] = Utils.create_tmp_file_from_contents(oc_type, oc_data['obj'].yaml_dict)

        return oc_objects

    def create(self):
        '''Create a router

           This includes the different parts:
           - deploymentconfig
           - service
           - serviceaccount
           - secrets
           - clusterrolebinding
        '''
        results = []
        self.needs_update()

        # pylint: disable=maybe-no-member
        for kind, oc_data in self.prepared_router.items():
            if oc_data['obj'] is not None:
                time.sleep(1)
                if self.get_object_by_kind(kind) is None:
                    results.append(self._create(oc_data['path']))

                elif oc_data['update']:
                    results.append(self._replace(oc_data['path']))


        rval = 0
        for result in results:
            if result['returncode'] != 0 and not 'already exist' in result['stderr']:
                rval = result['returncode']

        return {'returncode': rval, 'results': results}

    def update(self):
        '''run update for the router.  This performs a replace'''
        results = []

        # pylint: disable=maybe-no-member
        for _, oc_data in self.prepared_router.items():
            if oc_data['update']:
                results.append(self._replace(oc_data['path']))

        rval = 0
        for result in results:
            if result['returncode'] != 0:
                rval = result['returncode']

        return {'returncode': rval, 'results': results}

    # pylint: disable=too-many-return-statements,too-many-branches
    def needs_update(self):
        ''' check to see if we need to update '''
        # ServiceAccount:
        #   Need to determine changes from the pregenerated ones from the original
        #   Since these are auto generated, we can skip
        skip = ['secrets', 'imagePullSecrets']
        if self.serviceaccount is None or \
                not Utils.check_def_equal(self.prepared_router['ServiceAccount']['obj'].yaml_dict,
                                          self.serviceaccount.yaml_dict,
                                          skip_keys=skip,
                                          debug=self.verbose):
            self.prepared_router['ServiceAccount']['update'] = True

        # Secret:
        #   See if one was generated from our dry-run and verify it if needed
        if self.prepared_router['Secret']['obj']:
            if not self.secret:
                self.prepared_router['Secret']['update'] = True

            if self.secret is None or \
                    not Utils.check_def_equal(self.prepared_router['Secret']['obj'].yaml_dict,
                                              self.secret.yaml_dict,
                                              skip_keys=skip,
                                              debug=self.verbose):
                self.prepared_router['Secret']['update'] = True

        # Service:
        #   Fix the ports to have protocol=TCP
        for port in self.prepared_router['Service']['obj'].get('spec.ports'):
            port['protocol'] = 'TCP'

        skip = ['portalIP', 'clusterIP', 'sessionAffinity', 'type']
        if self.service is None or \
                not Utils.check_def_equal(self.prepared_router['Service']['obj'].yaml_dict,
                                          self.service.yaml_dict,
                                          skip_keys=skip,
                                          debug=self.verbose):
            self.prepared_router['Service']['update'] = True

        # DeploymentConfig:
        #   Router needs some exceptions.
        #   We do not want to check the autogenerated password for stats admin
        if self.deploymentconfig is not None:
            if not self.config.config_options['stats_password']['value']:
                for idx, env_var in enumerate(self.prepared_router['DeploymentConfig']['obj'].get(\
                            'spec.template.spec.containers[0].env') or []):
                    if env_var['name'] == 'STATS_PASSWORD':
                        env_var['value'] = \
                          self.deploymentconfig.get('spec.template.spec.containers[0].env[%s].value' % idx)
                        break

            # dry-run doesn't add the protocol to the ports section.  We will manually do that.
            for idx, port in enumerate(self.prepared_router['DeploymentConfig']['obj'].get(\
                            'spec.template.spec.containers[0].ports') or []):
                if not 'protocol' in port:
                    port['protocol'] = 'TCP'

        # These are different when generating
        skip = ['dnsPolicy',
                'terminationGracePeriodSeconds',
                'restartPolicy', 'timeoutSeconds',
                'livenessProbe', 'readinessProbe',
                'terminationMessagePath', 'hostPort',
                'defaultMode',
               ]

        if self.deploymentconfig is None or \
                not Utils.check_def_equal(self.prepared_router['DeploymentConfig']['obj'].yaml_dict,
                                          self.deploymentconfig.yaml_dict,
                                          skip_keys=skip,
                                          debug=self.verbose):
            self.prepared_router['DeploymentConfig']['update'] = True

        # Check if any of the parts need updating, if so, return True
        # else, no need to update
        # pylint: disable=no-member
        return any([self.prepared_router[oc_type]['update'] for oc_type in self.prepared_router.keys()])

    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_adm_router module'''

        rconfig = RouterConfig(params['name'],
                               params['namespace'],
                               params['kubeconfig'],
                               {'default_cert': {'value': params['default_cert'], 'include': True},
                                'cert_file': {'value': params['cert_file'], 'include': False},
                                'key_file': {'value': params['key_file'], 'include': False},
                                'images': {'value': params['images'], 'include': True},
                                'latest_images': {'value': params['latest_images'], 'include': True},
                                'labels': {'value': params['labels'], 'include': True},
                                'ports': {'value': ','.join(params['ports']), 'include': True},
                                'replicas': {'value': params['replicas'], 'include': True},
                                'selector': {'value': params['selector'], 'include': True},
                                'service_account': {'value': params['service_account'], 'include': True},
                                'router_type': {'value': params['router_type'], 'include': False},
                                'host_network': {'value': params['host_network'], 'include': True},
                                'extended_validation': {'value': params['extended_validation'], 'include': False},
                                'external_host': {'value': params['external_host'], 'include': True},
                                'external_host_vserver': {'value': params['external_host_vserver'],
                                                          'include': True},
                                'external_host_insecure': {'value': params['external_host_insecure'],
                                                           'include': True},
                                'external_host_partition_path': {'value': params['external_host_partition_path'],
                                                                 'include': True},
                                'external_host_username': {'value': params['external_host_username'],
                                                           'include': True},
                                'external_host_password': {'value': params['external_host_password'],
                                                           'include': True},
                                'external_host_private_key': {'value': params['external_host_private_key'],
                                                              'include': True},
                                'stats_user': {'value': params['stats_user'], 'include': True},
                                'stats_password': {'value': params['stats_password'], 'include': True},
                                'stats_port': {'value': params['stats_port'], 'include': True},
                                # extra
                                'cacert_file': {'value': params['cacert_file'], 'include': False},
                                # edits
                                'edits': {'value': params['edits'], 'include': False},
                               })


        state = params['state']

        ocrouter = Router(rconfig, verbose=params['debug'])

        api_rval = ocrouter.get()

        ########
        # get
        ########
        if state == 'list':
            return {'changed': False, 'results': api_rval, 'state': state}

        ########
        # Delete
        ########
        if state == 'absent':
            if not ocrouter.exists():
                return {'changed': False, 'state': state}

            if check_mode:
                return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a delete.'}

            # In case of delete we return a list of each object
            # that represents a router and its result in a list
            # pylint: disable=redefined-variable-type
            api_rval = ocrouter.delete()

            return {'changed': True, 'results': api_rval, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not ocrouter.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}

                api_rval = ocrouter.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if not ocrouter.needs_update():
                return {'changed': False, 'state': state}

            if check_mode:
                return {'changed': False, 'msg': 'CHECK_MODE: Would have performed an update.'}

            api_rval = ocrouter.update()

            if api_rval['returncode'] != 0:
                return {'failed': True, 'msg': api_rval}

            return {'changed': True, 'results': api_rval, 'state': state}
