#!/usr/bin/python
# -*- coding: utf-8 -*-
# pylint: disable=line-too-long,invalid-name

"""For details on this module see DOCUMENTATION (below)"""

import base64
import datetime
import io
import os
import subprocess
import yaml
import dateutil.parser

# pylint import-error disabled because pylint cannot find the package
# when installed in a virtualenv
from ansible.module_utils.six.moves import configparser  # pylint: disable=import-error
from ansible.module_utils.basic import AnsibleModule

try:
    # You can comment this import out and include a 'pass' in this
    # block if you're manually testing this module on a NON-ATOMIC
    # HOST (or any host that just doesn't have PyOpenSSL
    # available). That will force the `load_and_handle_cert` function
    # to use the Fake OpenSSL classes.
    import OpenSSL.crypto
    HAS_OPENSSL = True
except ImportError:
    # Some platforms (such as RHEL Atomic) may not have the Python
    # OpenSSL library installed. In this case we will use a manual
    # work-around to parse each certificate.
    #
    # Check for 'OpenSSL.crypto' in `sys.modules` later.
    HAS_OPENSSL = False

DOCUMENTATION = '''
---
module: openshift_cert_expiry
short_description: Check OpenShift Container Platform (OCP) and Kube certificate expirations on a cluster
description:
  - The M(openshift_cert_expiry) module has two basic functions: to flag certificates which will expire in a set window of time from now, and to notify you about certificates which have already expired.
  - When the module finishes, a summary of the examination is returned. Each certificate in the summary has a C(health) key with a value of one of the following:
  - C(ok) - not expired, and outside of the expiration C(warning_days) window.
  - C(warning) - not expired, but will expire between now and the C(warning_days) window.
  - C(expired) - an expired certificate.
  - Certificate flagging follow this logic:
  - If the expiration date is before now then the certificate is classified as C(expired).
  - The certificates time to live (expiration date - now) is calculated, if that time window is less than C(warning_days) the certificate is classified as C(warning).
  - All other conditions are classified as C(ok).
  - The following keys are ALSO present in the certificate summary:
  - C(cert_cn) - The common name of the certificate (additional CNs present in SAN extensions are omitted)
  - C(days_remaining) - The number of days until the certificate expires.
  - C(expiry) - The date the certificate expires on.
  - C(path) - The full path to the certificate on the examined host.
version_added: "1.0"
options:
  config_base:
    description:
      - Base path to OCP system settings.
    required: false
    default: /etc/origin
  warning_days:
    description:
      - Flag certificates which will expire in C(warning_days) days from now.
    required: false
    default: 30
  show_all:
    description:
      - Enable this option to show analysis of ALL certificates examined by this module.
      - By default only certificates which have expired, or will expire within the C(warning_days) window will be reported.
    required: false
    default: false

author: "Tim Bielawa (@tbielawa) <tbielawa@redhat.com>"
'''

EXAMPLES = '''
# Default invocation, only notify about expired certificates or certificates which will expire within 30 days from now
- openshift_cert_expiry:

# Expand the warning window to show certificates expiring within a year from now
- openshift_cert_expiry: warning_days=365

# Show expired, soon to expire (now + 30 days), and all other certificates examined
- openshift_cert_expiry: show_all=true
'''


class FakeOpenSSLCertificate(object):
    """This provides a rough mock of what you get from
`OpenSSL.crypto.load_certificate()`. This is a work-around for
platforms missing the Python OpenSSL library.
    """
    def __init__(self, cert_string):
        """`cert_string` is a certificate in the form you get from running a
.crt through 'openssl x509 -in CERT.cert -text'"""
        self.cert_string = cert_string
        self.serial = None
        self.subject = None
        self.extensions = []
        self.not_after = None
        self._parse_cert()

    def _parse_cert(self):
        """Manually parse the certificate line by line"""
        self.extensions = []

        PARSING_ALT_NAMES = False
        PARSING_HEX_SERIAL = False
        for line in self.cert_string.split('\n'):
            l = line.strip()
            if PARSING_ALT_NAMES:
                # We're parsing a 'Subject Alternative Name' line
                self.extensions.append(
                    FakeOpenSSLCertificateSANExtension(l))

                PARSING_ALT_NAMES = False
                continue

            if PARSING_HEX_SERIAL:
                # Hex serials arrive colon-delimited
                serial_raw = l.replace(':', '')
                # Convert to decimal
                self.serial = int('0x' + serial_raw, base=16)
                PARSING_HEX_SERIAL = False
                continue

            # parse out the bits that we can
            if l.startswith('Serial Number:'):
                # Decimal format:
                #   Serial Number: 11 (0xb)
                #   => 11
                # Hex Format (large serials):
                #   Serial Number:
                #       0a:de:eb:24:04:75:ab:56:39:14:e9:5a:22:e2:85:bf
                #   => 14449739080294792594019643629255165375
                if l.endswith(':'):
                    PARSING_HEX_SERIAL = True
                    continue
                self.serial = int(l.split()[-2])

            elif l.startswith('Not After :'):
                # Not After : Feb  7 18:19:35 2019 GMT
                # => strptime(str, '%b %d %H:%M:%S %Y %Z')
                # => strftime('%Y%m%d%H%M%SZ')
                # => 20190207181935Z
                not_after_raw = l.partition(' : ')[-1]
                # Last item: ('Not After', ' : ', 'Feb  7 18:19:35 2019 GMT')
                not_after_parsed = dateutil.parser.parse(not_after_raw)
                self.not_after = not_after_parsed.strftime('%Y%m%d%H%M%SZ')

            elif l.startswith('X509v3 Subject Alternative Name:'):
                PARSING_ALT_NAMES = True
                continue

            elif l.startswith('Subject:'):
                # O = system:nodes, CN = system:node:m01.example.com
                self.subject = FakeOpenSSLCertificateSubjects(l.partition(': ')[-1])

    def get_serial_number(self):
        """Return the serial number of the cert"""
        return self.serial

    def get_subject(self):
        """Subjects must implement get_components() and return dicts or
tuples. An 'openssl x509 -in CERT.cert -text' with 'Subject':

    Subject: Subject: O=system:nodes, CN=system:node:m01.example.com

might return: [('O=system', 'nodes'), ('CN=system', 'node:m01.example.com')]
        """
        return self.subject

    def get_extension(self, i):
        """Extensions must implement get_short_name() and return the string
'subjectAltName'"""
        return self.extensions[i]

    def get_extension_count(self):
        """ get_extension_count """
        return len(self.extensions)

    def get_notAfter(self):
        """Returns a date stamp as a string in the form
'20180922170439Z'. strptime the result with format param:
'%Y%m%d%H%M%SZ'."""
        return self.not_after


class FakeOpenSSLCertificateSANExtension(object):  # pylint: disable=too-few-public-methods
    """Mocks what happens when `get_extension` is called on a certificate
object"""

    def __init__(self, san_string):
        """With `san_string` as you get from:

    $ openssl x509 -in certificate.crt -text
        """
        self.san_string = san_string
        self.short_name = 'subjectAltName'

    def get_short_name(self):
        """Return the 'type' of this extension. It's always the same though
because we only care about subjectAltName's"""
        return self.short_name

    def __str__(self):
        """Return this extension and the value as a simple string"""
        return self.san_string


# pylint: disable=too-few-public-methods
class FakeOpenSSLCertificateSubjects(object):
    """Mocks what happens when `get_subject` is called on a certificate
object"""

    def __init__(self, subject_string):
        """With `subject_string` as you get from:

    $ openssl x509 -in certificate.crt -text
        """
        self.subjects = []
        for s in subject_string.split(', '):
            name, _, value = s.partition(' = ')
            self.subjects.append((name, value))

    def get_components(self):
        """Returns a list of tuples"""
        return self.subjects


######################################################################
def filter_paths(path_list):
    """`path_list` - A list of file paths to check. Only files which exist
will be returned
    """
    return [p for p in path_list if os.path.exists(os.path.realpath(p))]


# pylint: disable=too-many-locals,too-many-branches
#
# TODO: Break this function down into smaller chunks
def load_and_handle_cert(cert_string, now, base64decode=False, ans_module=None):
    """Load a certificate, split off the good parts, and return some
useful data

Params:

- `cert_string` (string) - a certificate loaded into a string object
- `now` (datetime) - a datetime object of the time to calculate the certificate 'time_remaining' against
- `base64decode` (bool) - run base64.b64decode() on the input
- `ans_module` (AnsibleModule) - The AnsibleModule object for this module (so we can raise errors)

Returns:
A tuple of the form:
    (cert_subject, cert_expiry_date, time_remaining, cert_serial_number)
    """
    if base64decode:
        _cert_string = base64.b64decode(cert_string).decode('utf-8')
    else:
        _cert_string = cert_string

    # Disable this. We 'redefine' the type because we are working
    # around a missing library on the target host.
    #
    # pylint: disable=redefined-variable-type
    if HAS_OPENSSL:
        # No work-around required
        cert_loaded = OpenSSL.crypto.load_certificate(
            OpenSSL.crypto.FILETYPE_PEM, _cert_string)
    else:
        # Missing library, work-around required. Run the 'openssl'
        # command on it to decode it
        cmd = 'openssl x509 -text'
        try:
            openssl_proc = subprocess.Popen(cmd.split(),
                                            stdout=subprocess.PIPE,
                                            stdin=subprocess.PIPE)
        except OSError:
            ans_module.fail_json(msg="Error: The 'OpenSSL' python library and CLI command were not found on the target host. Unable to parse any certificates. This host will not be included in generated reports.")
        else:
            openssl_decoded = openssl_proc.communicate(_cert_string.encode('utf-8'))[0].decode('utf-8')
            cert_loaded = FakeOpenSSLCertificate(openssl_decoded)

    ######################################################################
    # Read all possible names from the cert
    cert_subjects = []
    for name, value in cert_loaded.get_subject().get_components():
        if isinstance(name, bytes) or isinstance(value, bytes):
            name = name.decode('utf-8')
            value = value.decode('utf-8')
        cert_subjects.append('{}:{}'.format(name, value))

    # To read SANs from a cert we must read the subjectAltName
    # extension from the X509 Object. What makes this more difficult
    # is that pyOpenSSL does not give extensions as an iterable
    san = None
    for i in range(cert_loaded.get_extension_count()):
        ext = cert_loaded.get_extension(i)
        if ext.get_short_name() == 'subjectAltName':
            san = ext

    if san is not None:
        # The X509Extension object for subjectAltName prints as a
        # string with the alt names separated by a comma and a
        # space. Split the string by ', ' and then add our new names
        # to the list of existing names
        cert_subjects.extend(str(san).split(', '))

    cert_subject = ', '.join(cert_subjects)
    ######################################################################

    # Grab the expiration date
    not_after = cert_loaded.get_notAfter()
    # example get_notAfter() => 20180922170439Z
    if isinstance(not_after, bytes):
        not_after = not_after.decode('utf-8')

    cert_expiry_date = datetime.datetime.strptime(
        not_after,
        '%Y%m%d%H%M%SZ')

    time_remaining = cert_expiry_date - now

    return (cert_subject, cert_expiry_date, time_remaining, cert_loaded.get_serial_number())


def classify_cert(cert_meta, now, time_remaining, expire_window, cert_list):
    """Given metadata about a certificate under examination, classify it
    into one of three categories, 'ok', 'warning', and 'expired'.

Params:

- `cert_meta` dict - A dict with certificate metadata. Required fields
  include: 'cert_cn', 'path', 'expiry', 'days_remaining', 'health'.
- `now` (datetime) - a datetime object of the time to calculate the certificate 'time_remaining' against
- `time_remaining` (datetime.timedelta) - a timedelta for how long until the cert expires
- `expire_window` (datetime.timedelta) - a timedelta for how long the warning window is
- `cert_list` list - A list to shove the classified cert into

Return:
- `cert_list` - The updated list of classified certificates
    """
    expiry_str = str(cert_meta['expiry'])
    # Categorization
    if cert_meta['expiry'] < now:
        # This already expired, must NOTIFY
        cert_meta['health'] = 'expired'
    elif time_remaining < expire_window:
        # WARN about this upcoming expirations
        cert_meta['health'] = 'warning'
    else:
        # Not expired or about to expire
        cert_meta['health'] = 'ok'

    cert_meta['expiry'] = expiry_str
    cert_meta['serial_hex'] = hex(int(cert_meta['serial']))
    cert_list.append(cert_meta)
    return cert_list


def tabulate_summary(certificates, kubeconfigs, etcd_certs, router_certs, registry_certs):
    """Calculate the summary text for when the module finishes
running. This includes counts of each classification and what have
you.

Params:

- `certificates` (list of dicts) - Processed `expire_check_result`
  dicts with filled in `health` keys for system certificates.
- `kubeconfigs` - as above for kubeconfigs
- `etcd_certs` - as above for etcd certs

Return:

- `summary_results` (dict) - Counts of each cert type classification
  and total items examined.
    """
    items = certificates + kubeconfigs + etcd_certs + router_certs + registry_certs

    summary_results = {
        'system_certificates': len(certificates),
        'kubeconfig_certificates': len(kubeconfigs),
        'etcd_certificates': len(etcd_certs),
        'router_certs': len(router_certs),
        'registry_certs': len(registry_certs),
        'total': len(items),
        'ok': 0,
        'warning': 0,
        'expired': 0
    }

    summary_results['expired'] = len([c for c in items if c['health'] == 'expired'])
    summary_results['warning'] = len([c for c in items if c['health'] == 'warning'])
    summary_results['ok'] = len([c for c in items if c['health'] == 'ok'])

    return summary_results


######################################################################
# This is our module MAIN function after all, so there's bound to be a
# lot of code bundled up into one block
#
# Reason: These checks are disabled because the issue was introduced
# during a period where the pylint checks weren't enabled for this file
# Status: temporarily disabled pending future refactoring
# pylint: disable=too-many-locals,too-many-statements,too-many-branches
def main():
    """This module examines certificates (in various forms) which compose
an OpenShift Container Platform cluster
    """

    module = AnsibleModule(
        argument_spec=dict(
            config_base=dict(
                required=False,
                default="/etc/origin",
                type='str'),
            warning_days=dict(
                required=False,
                default=30,
                type='int'),
            show_all=dict(
                required=False,
                default=False,
                type='bool')
        ),
        supports_check_mode=True,
    )

    # Basic scaffolding for OpenShift specific certs
    openshift_base_config_path = os.path.realpath(module.params['config_base'])
    openshift_master_config_path = os.path.join(openshift_base_config_path,
                                                "master", "master-config.yaml")
    openshift_node_config_path = os.path.join(openshift_base_config_path,
                                              "node", "node-config.yaml")
    openshift_node_bootstrap_config_path = os.path.join(openshift_base_config_path,
                                                        "node", "bootstrap-node-config.yaml")
    openshift_cert_check_paths = [
        openshift_master_config_path,
        openshift_node_config_path,
        openshift_node_bootstrap_config_path,
    ]

    # Paths for Kubeconfigs. Additional kubeconfigs are conditionally
    # checked later in the code
    master_kube_configs = ['admin', 'openshift-master',
                           'openshift-node', 'openshift-router',
                           'openshift-registry']

    kubeconfig_paths = []
    for m_kube_config in master_kube_configs:
        kubeconfig_paths.append(
            os.path.join(openshift_base_config_path, "master", m_kube_config + ".kubeconfig")
        )

    # Validate some paths we have the ability to do ahead of time
    openshift_cert_check_paths = filter_paths(openshift_cert_check_paths)
    kubeconfig_paths = filter_paths(kubeconfig_paths)

    # etcd, where do you hide your certs? Used when parsing etcd.conf
    etcd_cert_params = [
        "ETCD_TRUSTED_CA_FILE",
        "ETCD_CERT_FILE",
        "ETCD_PEER_TRUSTED_CA_FILE",
        "ETCD_PEER_CERT_FILE",
    ]

    # Expiry checking stuff
    now = datetime.datetime.now()
    # todo, catch exception for invalid input and return a fail_json
    warning_days = int(module.params['warning_days'])
    expire_window = datetime.timedelta(days=warning_days)

    # Module stuff
    #
    # The results of our cert checking to return from the task call
    check_results = {}
    check_results['meta'] = {}
    check_results['meta']['warning_days'] = warning_days
    check_results['meta']['checked_at_time'] = str(now)
    check_results['meta']['warn_before_date'] = str(now + expire_window)
    check_results['meta']['show_all'] = str(module.params['show_all'])
    # All the analyzed certs accumulate here
    ocp_certs = []

    ######################################################################
    # Sure, why not? Let's enable check mode.
    if module.check_mode:
        check_results['ocp_certs'] = []
        module.exit_json(
            check_results=check_results,
            msg="Checked 0 total certificates. Expired/Warning/OK: 0/0/0. Warning window: %s days" % module.params['warning_days'],
            rc=0,
            changed=False
        )

    ######################################################################
    # Check for OpenShift Container Platform specific certs
    ######################################################################
    for os_cert in filter_paths(openshift_cert_check_paths):
        # Open up that config file and locate the cert and CA
        with io.open(os_cert, 'r', encoding='utf-8') as fp:
            cert_meta = {}
            cfg = yaml.load(fp)
            # cert files are specified in parsed `fp` as relative to the path
            # of the original config file. 'master-config.yaml' with certFile
            # = 'foo.crt' implies that 'foo.crt' is in the same
            # directory. certFile = '../foo.crt' is in the parent directory.
            cfg_path = os.path.dirname(fp.name)

            servingInfoFile = cfg.get('servingInfo', {}).get('certFile')
            if servingInfoFile:
                cert_meta['certFile'] = os.path.join(cfg_path, servingInfoFile)

            servingInfoCA = cfg.get('servingInfo', {}).get('clientCA')
            if servingInfoCA:
                cert_meta['clientCA'] = os.path.join(cfg_path, servingInfoCA)

            serviceSigner = cfg.get('controllerConfig', {}).get('serviceServingCert', {}).get('signer', {}).get('certFile')
            if serviceSigner:
                cert_meta['serviceSigner'] = os.path.join(cfg_path, serviceSigner)

            etcdClientCA = cfg.get('etcdClientInfo', {}).get('ca')
            if etcdClientCA:
                cert_meta['etcdClientCA'] = os.path.join(cfg_path, etcdClientCA)

            etcdClientCert = cfg.get('etcdClientInfo', {}).get('certFile')
            if etcdClientCert:
                cert_meta['etcdClientCert'] = os.path.join(cfg_path, etcdClientCert)

            kubeletCert = cfg.get('kubeletClientInfo', {}).get('certFile')
            if kubeletCert:
                cert_meta['kubeletCert'] = os.path.join(cfg_path, kubeletCert)

            proxyClient = cfg.get('kubernetesMasterConfig', {}).get('proxyClientInfo', {}).get('certFile')
            if proxyClient:
                cert_meta['proxyClient'] = os.path.join(cfg_path, proxyClient)

        ######################################################################
        # Load the certificate and the CA, parse their expiration dates into
        # datetime objects so we can manipulate them later
        for v in cert_meta.values():
            with io.open(v, 'r', encoding='utf-8') as fp:
                cert = fp.read()
                (cert_subject,
                 cert_expiry_date,
                 time_remaining,
                 cert_serial) = load_and_handle_cert(cert, now, ans_module=module)

                expire_check_result = {
                    'cert_cn': cert_subject,
                    'path': fp.name,
                    'expiry': cert_expiry_date,
                    'days_remaining': time_remaining.days,
                    'health': None,
                    'serial': cert_serial
                }

                classify_cert(expire_check_result, now, time_remaining, expire_window, ocp_certs)

    ######################################################################
    # /Check for OpenShift Container Platform specific certs
    ######################################################################

    ######################################################################
    # Check service Kubeconfigs
    ######################################################################
    kubeconfigs = []

    # There may be additional kubeconfigs to check, but their naming
    # is less predictable than the ones we've already assembled.

    for node_config in [openshift_node_config_path, openshift_node_bootstrap_config_path]:
        try:
            # Try to read the standard 'node-config.yaml' file to check if
            # this host is a node.
            with io.open(node_config, 'r', encoding='utf-8') as fp:
                cfg = yaml.load(fp)

            # OK, the config file exists, therefore this is a
            # node. Nodes have their own kubeconfig files to
            # communicate with the master API. Let's read the relative
            # path to that file from the node config.
            node_masterKubeConfig = cfg['masterKubeConfig']
            # As before, the path to the 'masterKubeConfig' file is
            # relative to `fp`
            cfg_path = os.path.dirname(fp.name)
            node_kubeconfig = os.path.join(cfg_path, node_masterKubeConfig)

            with io.open(node_kubeconfig, 'r', encoding='utf8') as fp:
                # Read in the nodes kubeconfig file and grab the good stuff
                cfg = yaml.load(fp)

            c = cfg['users'][0]['user'].get('client-certificate-data')
            if not c:
                # This is not a node
                raise IOError
            (cert_subject,
             cert_expiry_date,
             time_remaining,
             cert_serial) = load_and_handle_cert(c, now, base64decode=True, ans_module=module)

            expire_check_result = {
                'cert_cn': cert_subject,
                'path': fp.name,
                'expiry': cert_expiry_date,
                'days_remaining': time_remaining.days,
                'health': None,
                'serial': cert_serial
            }

            classify_cert(expire_check_result, now, time_remaining, expire_window, kubeconfigs)
        except IOError:
            # This is not a node
            pass

    for kube in filter_paths(kubeconfig_paths):
        with io.open(kube, 'r', encoding='utf-8') as fp:
            # TODO: Maybe consider catching exceptions here?
            cfg = yaml.load(fp)

        # Per conversation, "the kubeconfigs you care about:
        # admin, router, registry should all be single
        # value". Following that advice we only grab the data for
        # the user at index 0 in the 'users' list. There should
        # not be more than one user.
        c = cfg['users'][0]['user']['client-certificate-data']
        (cert_subject,
         cert_expiry_date,
         time_remaining,
         cert_serial) = load_and_handle_cert(c, now, base64decode=True, ans_module=module)

        expire_check_result = {
            'cert_cn': cert_subject,
            'path': fp.name,
            'expiry': cert_expiry_date,
            'days_remaining': time_remaining.days,
            'health': None,
            'serial': cert_serial
        }

        classify_cert(expire_check_result, now, time_remaining, expire_window, kubeconfigs)

    ######################################################################
    # /Check service Kubeconfigs
    ######################################################################

    ######################################################################
    # Check etcd certs
    #
    # Two things to check: 'external' etcd, and embedded etcd.
    ######################################################################
    # FIRST: The 'external' etcd
    #
    # Some values may be duplicated, make this a set for now so we
    # unique them all
    etcd_certs_to_check = set([])
    etcd_certs = []
    etcd_cert_params.append('dne')
    try:
        with io.open('/etc/etcd/etcd.conf', 'r', encoding='utf-8') as fp:
            # Add dummy header section.
            config = io.StringIO()
            config.write(u'[ETCD]\n')
            config.write(fp.read().replace('%', '%%'))
            config.seek(0, os.SEEK_SET)

            etcd_config = configparser.ConfigParser()
            etcd_config.readfp(config)

        for param in etcd_cert_params:
            try:
                etcd_certs_to_check.add(etcd_config.get('ETCD', param))
            except configparser.NoOptionError:
                # That parameter does not exist, oh well...
                pass
    except IOError:
        # No etcd to see here, move along
        pass

    for etcd_cert in filter_paths(etcd_certs_to_check):
        with io.open(etcd_cert, 'r', encoding='utf-8') as fp:
            c = fp.read()
            (cert_subject,
             cert_expiry_date,
             time_remaining,
             cert_serial) = load_and_handle_cert(c, now, ans_module=module)

            expire_check_result = {
                'cert_cn': cert_subject,
                'path': fp.name,
                'expiry': cert_expiry_date,
                'days_remaining': time_remaining.days,
                'health': None,
                'serial': cert_serial
            }

            classify_cert(expire_check_result, now, time_remaining, expire_window, etcd_certs)

    ######################################################################
    # /Check etcd certs
    ######################################################################

    ######################################################################
    # Check router/registry certs
    #
    # These are saved as secrets in etcd. That means that we can not
    # simply read a file to grab the data. Instead we're going to
    # subprocess out to the 'oc get' command. On non-masters this
    # command will fail, that is expected so we catch that exception.
    ######################################################################
    router_certs = []
    registry_certs = []

    ######################################################################
    # First the router certs
    try:
        router_secrets_raw = subprocess.Popen('oc get -n default secret router-certs -o yaml'.split(),
                                              stdout=subprocess.PIPE)
        router_ds = yaml.load(router_secrets_raw.communicate()[0])
        router_c = router_ds['data']['tls.crt']
        router_path = router_ds['metadata']['selfLink']
    except TypeError:
        # YAML couldn't load the result, this is not a master
        pass
    except OSError:
        # The OC command doesn't exist here. Move along.
        pass
    else:
        (cert_subject,
         cert_expiry_date,
         time_remaining,
         cert_serial) = load_and_handle_cert(router_c, now, base64decode=True, ans_module=module)

        expire_check_result = {
            'cert_cn': cert_subject,
            'path': router_path,
            'expiry': cert_expiry_date,
            'days_remaining': time_remaining.days,
            'health': None,
            'serial': cert_serial
        }

        classify_cert(expire_check_result, now, time_remaining, expire_window, router_certs)

    ######################################################################
    # Now for registry
    try:
        registry_secrets_raw = subprocess.Popen('oc get -n default secret registry-certificates -o yaml'.split(),
                                                stdout=subprocess.PIPE)
        registry_ds = yaml.load(registry_secrets_raw.communicate()[0])
        registry_c = registry_ds['data']['registry.crt']
        registry_path = registry_ds['metadata']['selfLink']
    except TypeError:
        # YAML couldn't load the result, this is not a master
        pass
    except OSError:
        # The OC command doesn't exist here. Move along.
        pass
    else:
        (cert_subject,
         cert_expiry_date,
         time_remaining,
         cert_serial) = load_and_handle_cert(registry_c, now, base64decode=True, ans_module=module)

        expire_check_result = {
            'cert_cn': cert_subject,
            'path': registry_path,
            'expiry': cert_expiry_date,
            'days_remaining': time_remaining.days,
            'health': None,
            'serial': cert_serial
        }

        classify_cert(expire_check_result, now, time_remaining, expire_window, registry_certs)

    ######################################################################
    # /Check router/registry certs
    ######################################################################

    res = tabulate_summary(ocp_certs, kubeconfigs, etcd_certs, router_certs, registry_certs)
    warn_certs = bool(res['expired'] + res['warning'])
    msg = "Checked {count} total certificates. Expired/Warning/OK: {exp}/{warn}/{ok}. Warning window: {window} days".format(
        count=res['total'],
        exp=res['expired'],
        warn=res['warning'],
        ok=res['ok'],
        window=int(module.params['warning_days']),
    )

    # By default we only return detailed information about expired or
    # warning certificates. If show_all is true then we will print all
    # the certificates examined.
    if not module.params['show_all']:
        check_results['ocp_certs'] = [crt for crt in ocp_certs if crt['health'] in ['expired', 'warning']]
        check_results['kubeconfigs'] = [crt for crt in kubeconfigs if crt['health'] in ['expired', 'warning']]
        check_results['etcd'] = [crt for crt in etcd_certs if crt['health'] in ['expired', 'warning']]
        check_results['registry'] = [crt for crt in registry_certs if crt['health'] in ['expired', 'warning']]
        check_results['router'] = [crt for crt in router_certs if crt['health'] in ['expired', 'warning']]
    else:
        check_results['ocp_certs'] = ocp_certs
        check_results['kubeconfigs'] = kubeconfigs
        check_results['etcd'] = etcd_certs
        check_results['registry'] = registry_certs
        check_results['router'] = router_certs

    # Sort the final results to report in order of ascending safety
    # time. That is to say, the certificates which will expire sooner
    # will be at the front of the list and certificates which will
    # expire later are at the end. Router and registry certs should be
    # limited to just 1 result, so don't bother sorting those.
    def cert_key(item):
        ''' return the days_remaining key '''
        return item['days_remaining']

    check_results['ocp_certs'] = sorted(check_results['ocp_certs'], key=cert_key)
    check_results['kubeconfigs'] = sorted(check_results['kubeconfigs'], key=cert_key)
    check_results['etcd'] = sorted(check_results['etcd'], key=cert_key)

    # This module will never change anything, but we might want to
    # change the return code parameter if there is some catastrophic
    # error we noticed earlier
    module.exit_json(
        check_results=check_results,
        warn_certs=warn_certs,
        summary=res,
        msg=msg,
        rc=0,
        changed=False
    )


if __name__ == '__main__':
    main()
