# pylint: skip-file
# flake8: noqa

class CAServerCertConfig(OpenShiftCLIConfig):
    ''' CAServerCertConfig is a DTO for the oc adm ca command '''
    def __init__(self, kubeconfig, verbose, ca_options):
        super(CAServerCertConfig, self).__init__('ca', None, kubeconfig, ca_options)
        self.kubeconfig = kubeconfig
        self.verbose = verbose
        self._ca = ca_options


class CAServerCert(OpenShiftCLI):
    ''' Class to wrap the oc adm ca create-server-cert command line'''
    def __init__(self,
                 config,
                 verbose=False):
        ''' Constructor for oadm ca '''
        super(CAServerCert, self).__init__(None, config.kubeconfig, verbose)
        self.config = config
        self.verbose = verbose

    def get(self):
        '''get the current cert file

           If a file exists by the same name in the specified location then the cert exists
        '''
        cert = self.config.config_options['cert']['value']
        if cert and os.path.exists(cert):
            return open(cert).read()

        return None

    def create(self):
        '''run openshift oc adm ca create-server-cert cmd'''

        # Added this here as a safegaurd for stomping on the
        # cert and key files if they exist
        if self.config.config_options['backup']['value']:
            ext = time.strftime("%Y-%m-%d@%H:%M:%S", time.localtime(time.time()))
            date_str = "%s_" + "%s" % ext

            if os.path.exists(self.config.config_options['key']['value']):
                shutil.copy(self.config.config_options['key']['value'],
                            date_str % self.config.config_options['key']['value'])
            if os.path.exists(self.config.config_options['cert']['value']):
                shutil.copy(self.config.config_options['cert']['value'],
                            date_str % self.config.config_options['cert']['value'])

        options = self.config.to_option_list()

        cmd = ['ca', 'create-server-cert']
        cmd.extend(options)

        return self.openshift_cmd(cmd, oadm=True)

    def exists(self):
        ''' check whether the certificate exists and has the clusterIP '''

        cert_path = self.config.config_options['cert']['value']
        if not os.path.exists(cert_path):
            return False

        # Would prefer pyopenssl but is not installed.
        # When we verify it is, switch this code
        # Here is the code to get the subject and the SAN
        # openssl x509 -text -noout -certopt \
        #  no_header,no_version,no_serial,no_signame,no_validity,no_issuer,no_pubkey,no_sigdump,no_aux \
        #  -in /etc/origin/master/registry.crt
        # Instead of this solution we will use a regex.
        cert_names = []
        hostnames = self.config.config_options['hostnames']['value'].split(',')
        proc = subprocess.Popen(['openssl', 'x509', '-noout', '-text', '-in', cert_path],
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE)

        x509output, _ = proc.communicate()
        if proc.returncode == 0:
            regex = re.compile(r"^\s*X509v3 Subject Alternative Name:\s*?\n\s*(.*)\s*\n", re.MULTILINE)
            match = regex.search(x509output.decode())  # E501
            if not match:
                return False

            for entry in re.split(r", *", match.group(1)):
                if entry.startswith('DNS') or entry.startswith('IP Address'):
                    cert_names.append(entry.split(':')[1])
            # now that we have cert names let's compare
            cert_set = set(cert_names)
            hname_set = set(hostnames)
            if cert_set.issubset(hname_set) and hname_set.issubset(cert_set):
                return True

        return False

    @staticmethod
    def run_ansible(params, check_mode):
        '''run the oc_adm_ca_server_cert module'''

        # Filter non-strings from hostnames list (Such as boolean: False)
        params['hostnames'] = [host for host in params['hostnames'] if isinstance(host, string_types)]

        config = CAServerCertConfig(params['kubeconfig'],
                                    params['debug'],
                                    {'cert':          {'value': params['cert'], 'include': True},
                                     'hostnames':     {'value': ','.join(params['hostnames']), 'include': True},
                                     'overwrite':     {'value': True, 'include': True},
                                     'key':           {'value': params['key'], 'include': True},
                                     'signer_cert':   {'value': params['signer_cert'], 'include': True},
                                     'signer_key':    {'value': params['signer_key'], 'include': True},
                                     'signer_serial': {'value': params['signer_serial'], 'include': True},
                                     'expire_days':   {'value': params['expire_days'], 'include': True},
                                     'backup':        {'value': params['backup'], 'include': False},
                                    })

        server_cert = CAServerCert(config)

        state = params['state']

        if state == 'present':
            ########
            # Create
            ########
            if not server_cert.exists() or params['force']:

                if check_mode:
                    return {'changed': True,
                            'msg': "CHECK_MODE: Would have created the certificate.",
                            'state': state}

                api_rval = server_cert.create()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': state}

            ########
            # Exists
            ########
            api_rval = server_cert.get()
            return {'changed': False, 'results': api_rval, 'state': state}

        return {'failed': True,
                'msg': 'Unknown state passed. %s' % state}
