# pylint: skip-file
# flake8: noqa

class RegistryException(Exception):
    ''' Registry Exception Class '''
    pass


class RegistryConfig(OpenShiftCLIConfig):
    ''' RegistryConfig is a DTO for the registry.  '''
    def __init__(self, rname, namespace, kubeconfig, registry_options):
        super(RegistryConfig, self).__init__(rname, namespace, kubeconfig, registry_options)


class Registry(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''

    volume_mount_path = 'spec.template.spec.containers[0].volumeMounts'
    volume_path = 'spec.template.spec.volumes'
    env_path = 'spec.template.spec.containers[0].env'

    def __init__(self,
                 registry_config,
                 verbose=False):
        ''' Constructor for Registry

           a registry consists of 3 or more parts
           - dc/docker-registry
           - svc/docker-registry

           Parameters:
           :registry_config:
           :verbose:
        '''
        super(Registry, self).__init__(registry_config.namespace, registry_config.kubeconfig, verbose)
        self.version = OCVersion(registry_config.kubeconfig, verbose)
        self.svc_ip = None
        self.portal_ip = None
        self.config = registry_config
        self.verbose = verbose
        self.registry_parts = [{'kind': 'dc', 'name': self.config.name},
                               {'kind': 'svc', 'name': self.config.name},
                              ]

        self.__prepared_registry = None
        self.volume_mounts = []
        self.volumes = []
        if self.config.config_options['volume_mounts']['value']:
            for volume in self.config.config_options['volume_mounts']['value']:
                volume_info = {'secret_name': volume.get('secret_name', None),
                               'name':        volume.get('name', None),
                               'type':        volume.get('type', None),
                               'path':        volume.get('path', None),
                               'claimName':   volume.get('claim_name', None),
                               'claimSize':   volume.get('claim_size', None),
                              }

                vol, vol_mount = Volume.create_volume_structure(volume_info)
                self.volumes.append(vol)
                self.volume_mounts.append(vol_mount)

        self.dconfig = None
        self.svc = None

    @property
    def deploymentconfig(self):
        ''' deploymentconfig property '''
        return self.dconfig

    @deploymentconfig.setter
    def deploymentconfig(self, config):
        ''' setter for deploymentconfig property '''
        self.dconfig = config

    @property
    def service(self):
        ''' service property '''
        return self.svc

    @service.setter
    def service(self, config):
        ''' setter for service property '''
        self.svc = config

    @property
    def prepared_registry(self):
        ''' prepared_registry property '''
        if not self.__prepared_registry:
            results = self.prepare_registry()
            if not results or ('returncode' in results and results['returncode'] != 0):
                raise RegistryException('Could not perform registry preparation. {}'.format(results))
            self.__prepared_registry = results

        return self.__prepared_registry

    @prepared_registry.setter
    def prepared_registry(self, data):
        ''' setter method for prepared_registry attribute '''
        self.__prepared_registry = data

    def get(self):
        ''' return the self.registry_parts '''
        self.deploymentconfig = None
        self.service = None

        rval = 0
        for part in self.registry_parts:
            result = self._get(part['kind'], name=part['name'])
            if result['returncode'] == 0 and part['kind'] == 'dc':
                self.deploymentconfig = DeploymentConfig(result['results'][0])
            elif result['returncode'] == 0 and part['kind'] == 'svc':
                self.service = Service(result['results'][0])

            if result['returncode'] != 0:
                rval = result['returncode']


        return {'returncode': rval, 'deploymentconfig': self.deploymentconfig, 'service': self.service}

    def exists(self):
        '''does the object exist?'''
        if self.deploymentconfig and self.service:
            return True

        return False

    def delete(self, complete=True):
        '''return all pods '''
        parts = []
        for part in self.registry_parts:
            if not complete and part['kind'] == 'svc':
                continue
            parts.append(self._delete(part['kind'], part['name']))

        # Clean up returned results
        rval = 0
        for part in parts:
            # pylint: disable=invalid-sequence-index
            if 'returncode' in part and part['returncode'] != 0:
                rval = part['returncode']

        return {'returncode': rval, 'results': parts}

    def prepare_registry(self):
        ''' prepare a registry for instantiation '''
        options = self.config.to_option_list(ascommalist='labels')

        cmd = ['registry']
        cmd.extend(options)
        cmd.extend(['--dry-run=True', '-o', 'json'])

        results = self.openshift_cmd(cmd, oadm=True, output=True, output_type='json')
        # probably need to parse this
        # pylint thinks results is a string
        # pylint: disable=no-member
        if results['returncode'] != 0 and 'items' not in results['results']:
            raise RegistryException('Could not perform registry preparation. {}'.format(results))

        service = None
        deploymentconfig = None
        # pylint: disable=invalid-sequence-index
        for res in results['results']['items']:
            if res['kind'] == 'DeploymentConfig':
                deploymentconfig = DeploymentConfig(res)
            elif res['kind'] == 'Service':
                service = Service(res)

        # Verify we got a service and a deploymentconfig
        if not service or not deploymentconfig:
            return results

        # results will need to get parsed here and modifications added
        deploymentconfig = DeploymentConfig(self.add_modifications(deploymentconfig))

        # modify service ip
        if self.svc_ip:
            service.put('spec.clusterIP', self.svc_ip)
        if self.portal_ip:
            service.put('spec.portalIP', self.portal_ip)

        # the dry-run doesn't apply the selector correctly
        if self.service:
            service.put('spec.selector', self.service.get_selector())

        # need to create the service and the deploymentconfig
        service_file = Utils.create_tmp_file_from_contents('service', service.yaml_dict)
        deployment_file = Utils.create_tmp_file_from_contents('deploymentconfig', deploymentconfig.yaml_dict)

        return {"service": service,
                "service_file": service_file,
                "service_update": False,
                "deployment": deploymentconfig,
                "deployment_file": deployment_file,
                "deployment_update": False}

    def create(self):
        '''Create a registry'''
        results = []
        self.needs_update()
        # if the object is none, then we need to create it
        # if the object needs an update, then we should call replace
        # Handle the deploymentconfig
        if self.deploymentconfig is None:
            results.append(self._create(self.prepared_registry['deployment_file']))
        elif self.prepared_registry['deployment_update']:
            results.append(self._replace(self.prepared_registry['deployment_file']))

        # Handle the service
        if self.service is None:
            results.append(self._create(self.prepared_registry['service_file']))
        elif self.prepared_registry['service_update']:
            results.append(self._replace(self.prepared_registry['service_file']))

        # Clean up returned results
        rval = 0
        for result in results:
            # pylint: disable=invalid-sequence-index
            if 'returncode' in result and result['returncode'] != 0:
                rval = result['returncode']

        return {'returncode': rval, 'results': results}

    def update(self):
        '''run update for the registry.  This performs a replace if required'''
        # Store the current service IP
        if self.service:
            svcip = self.service.get('spec.clusterIP')
            if svcip:
                self.svc_ip = svcip
            portip = self.service.get('spec.portalIP')
            if portip:
                self.portal_ip = portip

        results = []
        if self.prepared_registry['deployment_update']:
            results.append(self._replace(self.prepared_registry['deployment_file']))
        if self.prepared_registry['service_update']:
            results.append(self._replace(self.prepared_registry['service_file']))

        # Clean up returned results
        rval = 0
        for result in results:
            if result['returncode'] != 0:
                rval = result['returncode']

        return {'returncode': rval, 'results': results}

    def add_modifications(self, deploymentconfig):
        ''' update a deployment config with changes '''
        # The environment variable for REGISTRY_HTTP_SECRET is autogenerated
        # We should set the generated deploymentconfig to the in memory version
        # the following modifications will overwrite if needed
        if self.deploymentconfig:
            result = self.deploymentconfig.get_env_var('REGISTRY_HTTP_SECRET')
            if result:
                deploymentconfig.update_env_var('REGISTRY_HTTP_SECRET', result['value'])

        # Currently we know that our deployment of a registry requires a few extra modifications
        # Modification 1
        # we need specific environment variables to be set
        for key, value in self.config.config_options['env_vars'].get('value', {}).items():
            if not deploymentconfig.exists_env_key(key):
                deploymentconfig.add_env_value(key, value)
            else:
                deploymentconfig.update_env_var(key, value)

        # Modification 2
        # we need specific volume variables to be set
        for volume in self.volumes:
            deploymentconfig.update_volume(volume)

        for vol_mount in self.volume_mounts:
            deploymentconfig.update_volume_mount(vol_mount)

        # Modification 3
        # Edits
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

        return deploymentconfig.yaml_dict

    def needs_update(self):
        ''' check to see if we need to update '''
        exclude_list = ['clusterIP', 'portalIP', 'type', 'protocol']
        if self.service is None or \
                not Utils.check_def_equal(self.prepared_registry['service'].yaml_dict,
                                          self.service.yaml_dict,
                                          exclude_list,
                                          debug=self.verbose):
            self.prepared_registry['service_update'] = True

        exclude_list = ['dnsPolicy',
                        'terminationGracePeriodSeconds',
                        'restartPolicy', 'timeoutSeconds',
                        'livenessProbe', 'readinessProbe',
                        'terminationMessagePath',
                        'securityContext',
                        'imagePullPolicy',
                        'protocol', # ports.portocol: TCP
                        'type', # strategy: {'type': 'rolling'}
                        'defaultMode', # added on secrets
                        'activeDeadlineSeconds', # added in 1.5 for timeouts
                       ]

        if self.deploymentconfig is None or \
                not Utils.check_def_equal(self.prepared_registry['deployment'].yaml_dict,
                                          self.deploymentconfig.yaml_dict,
                                          exclude_list,
                                          debug=self.verbose):
            self.prepared_registry['deployment_update'] = True

        return self.prepared_registry['deployment_update'] or self.prepared_registry['service_update'] or False

    # In the future, we would like to break out each ansible state into a function.
    # pylint: disable=too-many-branches,too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_adm_registry module'''

        registry_options = {'images': {'value': params['images'], 'include': True},
                            'latest_images': {'value': params['latest_images'], 'include': True},
                            'labels': {'value': params['labels'], 'include': True},
                            'ports': {'value': ','.join(params['ports']), 'include': True},
                            'replicas': {'value': params['replicas'], 'include': True},
                            'selector': {'value': params['selector'], 'include': True},
                            'service_account': {'value': params['service_account'], 'include': True},
                            'mount_host': {'value': params['mount_host'], 'include': True},
                            'env_vars': {'value': params['env_vars'], 'include': False},
                            'volume_mounts': {'value': params['volume_mounts'], 'include': False},
                            'edits': {'value': params['edits'], 'include': False},
                            'tls_key': {'value': params['tls_key'], 'include': True},
                            'tls_certificate': {'value': params['tls_certificate'], 'include': True},
                           }

        # Do not always pass the daemonset and enforce-quota parameters because they are not understood
        # by old versions of oc.
        # Default value is false. So, it's safe to not pass an explicit false value to oc versions which
        # understand these parameters.
        if params['daemonset']:
            registry_options['daemonset'] = {'value': params['daemonset'], 'include': True}
        if params['enforce_quota']:
            registry_options['enforce_quota'] = {'value': params['enforce_quota'], 'include': True}

        rconfig = RegistryConfig(params['name'],
                                 params['namespace'],
                                 params['kubeconfig'],
                                 registry_options)


        ocregistry = Registry(rconfig, params['debug'])

        api_rval = ocregistry.get()

        state = params['state']
        ########
        # get
        ########
        if state == 'list':

            if api_rval['returncode'] != 0:
                return {'failed': True, 'msg': api_rval}

            return {'changed': False, 'results': api_rval, 'state': state}

        ########
        # Delete
        ########
        if state == 'absent':
            if not ocregistry.exists():
                return {'changed': False, 'state': state}

            if check_mode:
                return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a delete.'}

            # Unsure as to why this is angry with the return type.
            # pylint: disable=redefined-variable-type
            api_rval = ocregistry.delete()

            if api_rval['returncode'] != 0:
                return {'failed': True, 'msg': api_rval}

            return {'changed': True, 'results': api_rval, 'state': state}

        if state == 'present':
            ########
            # Create
            ########
            if not ocregistry.exists():

                if check_mode:
                    return {'changed': True, 'msg': 'CHECK_MODE: Would have performed a create.'}

                api_rval = ocregistry.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Update
            ########
            if not params['force'] and not ocregistry.needs_update():
                return {'changed': False, 'state': state}

            if check_mode:
                return {'changed': True, 'msg': 'CHECK_MODE: Would have performed an update.'}

            api_rval = ocregistry.update()

            if api_rval['returncode'] != 0:
                return {'failed': True, 'msg': api_rval}

            return {'changed': True, 'results': api_rval, 'state': state}

        return {'failed': True, 'msg': 'Unknown state passed. %s' % state}
