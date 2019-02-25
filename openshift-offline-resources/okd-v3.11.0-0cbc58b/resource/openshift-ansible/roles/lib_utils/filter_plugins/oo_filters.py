#!/usr/bin/python
# -*- coding: utf-8 -*-
# pylint: disable=too-many-lines
"""
Custom filters for use in openshift-ansible
"""
import ast
import json
import os
import pdb
import random

from base64 import b64encode
from collections import Mapping
# pylint no-name-in-module and import-error disabled here because pylint
# fails to properly detect the packages when installed in a virtualenv
from distutils.util import strtobool  # pylint:disable=no-name-in-module,import-error
from operator import itemgetter

import yaml

from ansible import errors
from ansible.parsing.yaml.dumper import AnsibleDumper

# pylint: disable=import-error,no-name-in-module
from ansible.module_utils.six import iteritems, string_types, u
# pylint: disable=import-error,no-name-in-module
from ansible.module_utils.six.moves.urllib.parse import urlparse

HAS_OPENSSL = False
try:
    import OpenSSL.crypto
    HAS_OPENSSL = True
except ImportError:
    pass


# pylint: disable=C0103

def lib_utils_oo_pdb(arg):
    """ This pops you into a pdb instance where arg is the data passed in
        from the filter.
        Ex: "{{ hostvars | lib_utils_oo_pdb }}"
    """
    pdb.set_trace()
    return arg


def get_attr(data, attribute=None):
    """ This looks up dictionary attributes of the form a.b.c and returns
        the value.

        If the key isn't present, None is returned.
        Ex: data = {'a': {'b': {'c': 5}}}
            attribute = "a.b.c"
            returns 5
    """
    if not attribute:
        raise errors.AnsibleFilterError("|failed expects attribute to be set")

    ptr = data
    for attr in attribute.split('.'):
        if attr in ptr:
            ptr = ptr[attr]
        else:
            ptr = None
            break

    return ptr


def oo_flatten(data):
    """ This filter plugin will flatten a list of lists
    """
    if not isinstance(data, list):
        raise errors.AnsibleFilterError("|failed expects to flatten a List")

    return [item for sublist in data for item in sublist]


def lib_utils_oo_collect(data_list, attribute=None, filters=None):
    """ This takes a list of dict and collects all attributes specified into a
        list. If filter is specified then we will include all items that
        match _ALL_ of filters.  If a dict entry is missing the key in a
        filter it will be excluded from the match.
        Ex: data_list = [ {'a':1, 'b':5, 'z': 'z'}, # True, return
                          {'a':2, 'z': 'z'},        # True, return
                          {'a':3, 'z': 'z'},        # True, return
                          {'a':4, 'z': 'b'},        # FAILED, obj['z'] != obj['z']
                        ]
            attribute = 'a'
            filters   = {'z': 'z'}
            returns [1, 2, 3]

        This also deals with lists of lists with dict as elements.
        Ex: data_list = [
                          [ {'a':1, 'b':5, 'z': 'z'}, # True, return
                            {'a':2, 'b':6, 'z': 'z'}  # True, return
                          ],
                          [ {'a':3, 'z': 'z'},        # True, return
                            {'a':4, 'z': 'b'}         # FAILED, obj['z'] != obj['z']
                          ],
                          {'a':5, 'z': 'z'},          # True, return
                        ]
            attribute = 'a'
            filters   = {'z': 'z'}
            returns [1, 2, 3, 5]
    """
    if not isinstance(data_list, list):
        raise errors.AnsibleFilterError("lib_utils_oo_collect expects to filter on a List")

    if not attribute:
        raise errors.AnsibleFilterError("lib_utils_oo_collect expects attribute to be set")

    data = []
    retval = []

    for item in data_list:
        if isinstance(item, list):
            retval.extend(lib_utils_oo_collect(item, attribute, filters))
        else:
            data.append(item)

    if filters is not None:
        if not isinstance(filters, dict):
            raise errors.AnsibleFilterError(
                "lib_utils_oo_collect expects filter to be a dict")
        retval.extend([get_attr(d, attribute) for d in data if (
            all([get_attr(d, key) == filters[key] for key in filters]))])
    else:
        retval.extend([get_attr(d, attribute) for d in data])

    retval = [val for val in retval if val is not None]

    return retval


def lib_utils_oo_select_keys_from_list(data, keys):
    """ This returns a list, which contains the value portions for the keys
        Ex: data = { 'a':1, 'b':2, 'c':3 }
            keys = ['a', 'c']
            returns [1, 3]
    """

    if not isinstance(data, list):
        raise errors.AnsibleFilterError("|lib_utils_oo_select_keys_from_list failed expects to filter on a list")

    if not isinstance(keys, list):
        raise errors.AnsibleFilterError("|lib_utils_oo_select_keys_from_list failed expects first param is a list")

    # Gather up the values for the list of keys passed in
    retval = [lib_utils_oo_select_keys(item, keys) for item in data]

    return oo_flatten(retval)


def lib_utils_oo_select_keys(data, keys):
    """ This returns a list, which contains the value portions for the keys
        Ex: data = { 'a':1, 'b':2, 'c':3 }
            keys = ['a', 'c']
            returns [1, 3]
    """

    if not isinstance(data, Mapping):
        raise errors.AnsibleFilterError("|lib_utils_oo_select_keys failed expects to filter on a dict or object")

    if not isinstance(keys, list):
        raise errors.AnsibleFilterError("|lib_utils_oo_select_keys failed expects first param is a list")

    # Gather up the values for the list of keys passed in
    retval = [data[key] for key in keys if key in data]

    return retval


def lib_utils_oo_prepend_strings_in_list(data, prepend):
    """ This takes a list of strings and prepends a string to each item in the
        list
        Ex: data = ['cart', 'tree']
            prepend = 'apple-'
            returns ['apple-cart', 'apple-tree']
    """
    if not isinstance(data, list):
        raise errors.AnsibleFilterError("|failed expects first param is a list")
    if not all(isinstance(x, string_types) for x in data):
        raise errors.AnsibleFilterError("|failed expects first param is a list"
                                        " of strings")
    retval = [prepend + s for s in data]
    return retval


def lib_utils_oo_dict_to_list_of_dict(data, key_title='key', value_title='value'):
    """Take a dict and arrange them as a list of dicts

       Input data:
       {'region': 'infra', 'test_k': 'test_v'}

       Return data:
       [{'key': 'region', 'value': 'infra'}, {'key': 'test_k', 'value': 'test_v'}]

       Written for use of the oc_label module
    """
    if not isinstance(data, dict):
        # pylint: disable=line-too-long
        raise errors.AnsibleFilterError("|failed expects first param is a dict. Got %s. Type: %s" % (str(data), str(type(data))))

    rval = []
    for label in data.items():
        rval.append({key_title: label[0], value_title: label[1]})

    return rval


def oo_ami_selector(data, image_name):
    """ This takes a list of amis and an image name and attempts to return
        the latest ami.
    """
    if not isinstance(data, list):
        raise errors.AnsibleFilterError("|failed expects first param is a list")

    if not data:
        return None
    else:
        if image_name is None or not image_name.endswith('_*'):
            ami = sorted(data, key=itemgetter('name'), reverse=True)[0]
            return ami['ami_id']
        else:
            ami_info = [(ami, ami['name'].split('_')[-1]) for ami in data]
            ami = sorted(ami_info, key=itemgetter(1), reverse=True)[0][0]
            return ami['ami_id']


def lib_utils_oo_split(string, separator=','):
    """ This splits the input string into a list. If the input string is
    already a list we will return it as is.
    """
    if isinstance(string, list):
        return string
    return string.split(separator)


def lib_utils_oo_dict_to_keqv_list(data):
    """Take a dict and return a list of k=v pairs

        Input data:
        {'a': 1, 'b': 2}

        Return data:
        ['a=1', 'b=2']
    """
    if not isinstance(data, dict):
        try:
            # This will attempt to convert something that looks like a string
            # representation of a dictionary (including json) into a dictionary.
            data = ast.literal_eval(data)
        except ValueError:
            msg = "|failed expects first param is a dict. Got {}. Type: {}"
            msg = msg.format(str(data), str(type(data)))
            raise errors.AnsibleFilterError(msg)
    return ['='.join(str(e) for e in x) for x in data.items()]


def lib_utils_oo_list_to_dict(lst, separator='='):
    """ This converts a list of ["k=v"] to a dictionary {k: v}.
    """
    kvs = [i.split(separator) for i in lst]
    return {k: v for k, v in kvs}


def haproxy_backend_masters(hosts, port):
    """ This takes an array of dicts and returns an array of dicts
        to be used as a backend for the haproxy role
    """
    servers = []
    for idx, host_info in enumerate(hosts):
        server = dict(name="master%s" % idx)
        server_ip = host_info['openshift']['common']['ip']
        server['address'] = "%s:%s" % (server_ip, port)
        server['opts'] = 'check'
        servers.append(server)
    return servers


# pylint: disable=too-many-branches, too-many-nested-blocks
def lib_utils_oo_parse_named_certificates(certificates, named_certs_dir, internal_hostnames):
    """ Parses names from list of certificate hashes.

        Ex: certificates = [{ "certfile": "/root/custom1.crt",
                              "keyfile": "/root/custom1.key",
                               "cafile": "/root/custom-ca1.crt" },
                            { "certfile": "custom2.crt",
                              "keyfile": "custom2.key",
                              "cafile": "custom-ca2.crt" }]

            returns [{ "certfile": "/etc/origin/master/named_certificates/custom1.crt",
                       "keyfile": "/etc/origin/master/named_certificates/custom1.key",
                       "cafile": "/etc/origin/master/named_certificates/custom-ca1.crt",
                       "names": [ "public-master-host.com",
                                  "other-master-host.com" ] },
                     { "certfile": "/etc/origin/master/named_certificates/custom2.crt",
                       "keyfile": "/etc/origin/master/named_certificates/custom2.key",
                       "cafile": "/etc/origin/master/named_certificates/custom-ca-2.crt",
                       "names": [ "some-hostname.com" ] }]
    """
    if not isinstance(named_certs_dir, string_types):
        raise errors.AnsibleFilterError("|failed expects named_certs_dir is str or unicode")

    if not isinstance(internal_hostnames, list):
        raise errors.AnsibleFilterError("|failed expects internal_hostnames is list")

    if not HAS_OPENSSL:
        raise errors.AnsibleFilterError("|missing OpenSSL python bindings")

    for certificate in certificates:
        if 'names' in certificate.keys():
            continue
        else:
            certificate['names'] = []

        if not os.path.isfile(certificate['certfile']) or not os.path.isfile(certificate['keyfile']):
            raise errors.AnsibleFilterError("|certificate and/or key does not exist '%s', '%s'" %
                                            (certificate['certfile'], certificate['keyfile']))

        try:
            st_cert = open(certificate['certfile'], 'rt').read()
            cert = OpenSSL.crypto.load_certificate(OpenSSL.crypto.FILETYPE_PEM, st_cert)
            certificate['names'].append(str(cert.get_subject().commonName.decode()))
            for i in range(cert.get_extension_count()):
                if cert.get_extension(i).get_short_name() == 'subjectAltName':
                    for name in str(cert.get_extension(i)).split(', '):
                        if 'DNS:' in name:
                            certificate['names'].append(name.replace('DNS:', ''))
        except Exception:
            raise errors.AnsibleFilterError(("|failed to parse certificate '%s', " % certificate['certfile'] +
                                             "please specify certificate names in host inventory"))

        certificate['names'] = list(set(certificate['names']))
        if 'cafile' not in certificate:
            certificate['names'] = [name for name in certificate['names'] if name not in internal_hostnames]
            if not certificate['names']:
                raise errors.AnsibleFilterError(("|failed to parse certificate '%s' or " % certificate['certfile'] +
                                                 "detected a collision with internal hostname, please specify " +
                                                 "certificate names in host inventory"))

    for certificate in certificates:
        # Update paths for configuration
        certificate['certfile'] = os.path.join(named_certs_dir, os.path.basename(certificate['certfile']))
        certificate['keyfile'] = os.path.join(named_certs_dir, os.path.basename(certificate['keyfile']))
        if 'cafile' in certificate:
            certificate['cafile'] = os.path.join(named_certs_dir, os.path.basename(certificate['cafile']))
    return certificates


def lib_utils_oo_parse_certificate_san(certificate):
    """ Parses SubjectAlternativeNames from a PEM certificate.

        Ex: certificate = '''-----BEGIN CERTIFICATE-----
                MIIEcjCCAlqgAwIBAgIBAzANBgkqhkiG9w0BAQsFADAhMR8wHQYDVQQDDBZldGNk
                LXNpZ25lckAxNTE2ODIwNTg1MB4XDTE4MDEyNDE5MDMzM1oXDTIzMDEyMzE5MDMz
                M1owHzEdMBsGA1UEAwwUbWFzdGVyMS5hYnV0Y2hlci5jb20wggEiMA0GCSqGSIb3
                DQEBAQUAA4IBDwAwggEKAoIBAQD4wBdWXNI3TF1M0b0bEIGyJPvdqKeGwF5XlxWg
                NoA1Ain/Xz0N1SW5pXW2CDo9HX+ay8DyhzR532yrBa+RO3ivNCmfnexTQinfSLWG
                mBEdiu7HO3puR/GNm74JNyXoEKlMAIRiTGq9HPoTo7tNV5MLodgYirpHrkSutOww
                DfFSrNjH/ehqxwQtrIOnTAHigdTOrKVdoYxqXblDEMONTPLI5LMvm4/BqnAVaOyb
                9RUzND6lxU/ei3FbUS5IoeASOHx0l1ifxae3OeSNAimm/RIRo9rieFNUFh45TzID
                elsdGrLB75LH/gnRVV1xxVbwPN6xW1mEwOceRMuhIArJQ2G5AgMBAAGjgbYwgbMw
                UQYDVR0jBEowSIAUXTqN88vCI6E7wONls3QJ4/63unOhJaQjMCExHzAdBgNVBAMM
                FmV0Y2Qtc2lnbmVyQDE1MTY4MjA1ODWCCQDMaopfom6OljAMBgNVHRMBAf8EAjAA
                MBMGA1UdJQQMMAoGCCsGAQUFBwMBMAsGA1UdDwQEAwIFoDAdBgNVHQ4EFgQU7l05
                OYeY3HppL6/0VJSirudj8t0wDwYDVR0RBAgwBocEwKh6ujANBgkqhkiG9w0BAQsF
                AAOCAgEAFU8sicE5EeQsUPnFEqDvoJd1cVE+8aCBqkW0++4GsVw2A/JOJ3OBJL6r
                BV3b1u8/e8xBNi8hPi42Q+LWBITZZ/COFyhwEAK94hcr7eZLCV2xfUdMJziP4Qkh
                /WRN7vXHTtJ6NP/d6A22SPbtnMSt9Y6G8y9qa5HBrqIqmkYbLzDw/SdZbDbuGhRk
                xUwg2ahXNblVoE5P6rxPONgXliA94telZ1/61iyrVaiGQb1/GUP/DRfvvR4dOCrA
                lMosW6fm37Wdi/8iYW+aDPWGS+yVK/sjSnHNjxqvrzkfGk+COa5riT9hJ7wZY0Hb
                YiJS74SZgZt/nnr5PI2zFRUiZLECqCkZnC/sz29i+irLabnq7Cif9Mv+TUcXWvry
                TdJuaaYdTSMRSUkDd/c9Ife8tOr1i1xhFzDNKNkZjTVRk1MBquSXndVCDKucdfGi
                YoWm+NDFrayw8yxK/KTHo3Db3lu1eIXTHxriodFx898b//hysHr4hs4/tsEFUTZi
                705L2ScIFLfnyaPby5GK/3sBIXtuhOFM3QV3JoYKlJB5T6wJioVoUmSLc+UxZMeE
                t9gGVQbVxtLvNHUdW7uKQ5pd76nIJqApQf8wg2Pja8oo56fRZX2XLt8nm9cswcC4
                Y1mDMvtfxglQATwMTuoKGdREuu1mbdb8QqdyQmZuMa72q+ax2kQ=
                -----END CERTIFICATE-----'''

            returns ['192.168.122.186']
    """

    if not HAS_OPENSSL:
        raise errors.AnsibleFilterError("|missing OpenSSL python bindings")

    names = []

    try:
        lcert = OpenSSL.crypto.load_certificate(OpenSSL.crypto.FILETYPE_PEM, certificate)
        for i in range(lcert.get_extension_count()):
            if lcert.get_extension(i).get_short_name() == 'subjectAltName':
                sanstr = str(lcert.get_extension(i))
                sanstr = sanstr.replace('DNS:', '')
                sanstr = sanstr.replace('IP Address:', '')
                names = sanstr.split(', ')
    except Exception:
        raise errors.AnsibleFilterError("|failed to parse certificate")

    return names


def lib_utils_oo_generate_secret(num_bytes):
    """ generate a session secret """

    if not isinstance(num_bytes, int):
        raise errors.AnsibleFilterError("|failed expects num_bytes is int")

    return b64encode(os.urandom(num_bytes)).decode('utf-8')


def lib_utils_to_padded_yaml(data, level=0, indent=2, **kw):
    """ returns a yaml snippet padded to match the indent level you specify """
    if data in [None, ""]:
        return ""

    try:
        transformed = u(yaml.dump(data, indent=indent, allow_unicode=True,
                                  default_flow_style=False,
                                  Dumper=AnsibleDumper, **kw))
        padded = "\n".join([" " * level * indent + line for line in transformed.splitlines()])
        return "\n{0}".format(padded)
    except Exception as my_e:
        raise errors.AnsibleFilterError('Failed to convert: %s' % my_e)


def lib_utils_oo_image_tag_to_rpm_version(version, include_dash=False):
    """ Convert an image tag string to an RPM version if necessary
        Empty strings and strings that are already in rpm version format
        are ignored. Also remove non semantic version components.

        Ex. v3.2.0.10 -> -3.2.0.10
            v1.2.0-rc1 -> -1.2.0
    """
    if not isinstance(version, string_types):
        raise errors.AnsibleFilterError("|failed expects a string or unicode")
    if version.startswith("v"):
        version = version[1:]
        # Strip release from requested version, we no longer support this.
        version = version.split('-')[0]

    if include_dash and version and not version.startswith("-"):
        version = "-" + version

    return version


def lib_utils_oo_hostname_from_url(url):
    """ Returns the hostname contained in a URL

        Ex: https://ose3-master.example.com/v1/api -> ose3-master.example.com
    """
    if not isinstance(url, string_types):
        raise errors.AnsibleFilterError("|failed expects a string or unicode")
    parse_result = urlparse(url)
    if parse_result.netloc != '':
        return parse_result.netloc
    else:
        # netloc wasn't parsed, assume url was missing scheme and path
        return parse_result.path


# pylint: disable=invalid-name, unused-argument
def lib_utils_oo_loadbalancer_frontends(
        api_port, servers_hostvars, use_nuage=False, nuage_rest_port=None):
    """TODO: Document me."""
    loadbalancer_frontends = [{'name': 'atomic-openshift-api',
                               'mode': 'tcp',
                               'options': ['tcplog'],
                               'binds': ["*:{0}".format(api_port)],
                               'default_backend': 'atomic-openshift-api'}]
    if bool(strtobool(str(use_nuage))) and nuage_rest_port is not None:
        loadbalancer_frontends.append({'name': 'nuage-monitor',
                                       'mode': 'tcp',
                                       'options': ['tcplog'],
                                       'binds': ["*:{0}".format(nuage_rest_port)],
                                       'default_backend': 'nuage-monitor'})
    return loadbalancer_frontends


# pylint: disable=invalid-name
def lib_utils_oo_loadbalancer_backends(
        api_port, servers_hostvars, use_nuage=False, nuage_rest_port=None):
    """TODO: Document me."""
    loadbalancer_backends = [{'name': 'atomic-openshift-api',
                              'mode': 'tcp',
                              'option': 'tcplog',
                              'balance': 'source',
                              'servers': haproxy_backend_masters(servers_hostvars, api_port)}]
    if bool(strtobool(str(use_nuage))) and nuage_rest_port is not None:
        # pylint: disable=line-too-long
        loadbalancer_backends.append({'name': 'nuage-monitor',
                                      'mode': 'tcp',
                                      'option': 'tcplog',
                                      'balance': 'source',
                                      'servers': haproxy_backend_masters(servers_hostvars, nuage_rest_port)})
    return loadbalancer_backends


def lib_utils_oo_random_word(length, source='abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'):
    """Generates a random string of given length from a set of alphanumeric characters.
       The default source uses [a-z][A-Z][0-9]
       Ex:
       - lib_utils_oo_random_word(3)                => aB9
       - lib_utils_oo_random_word(4, source='012')  => 0123
    """
    return ''.join(random.choice(source) for i in range(length))


def lib_utils_oo_selector_to_string_list(user_dict):
    """Convert a dict of selectors to a key=value list of strings

Given input of {'region': 'infra', 'zone': 'primary'} returns a list
of items as ['node-role.kubernetes.io/infra=true', 'zone=primary']
    """
    selectors = []
    for key in user_dict:
        selectors.append("{}={}".format(key, user_dict[key]))
    return selectors


def lib_utils_oo_filter_sa_secrets(sa_secrets, secret_hint='-token-'):
    """Parse the Service Account Secrets list, `sa_secrets`, (as from
oc_serviceaccount_secret:state=list) and return the name of the secret
containing the `secret_hint` string. For example, by default this will
return the name of the secret holding the SA bearer token.

Only provide the 'results' object to this filter. This filter expects
to receive a list like this:

    [
        {
            "name": "management-admin-dockercfg-p31s2"
        },
        {
            "name": "management-admin-token-bnqsh"
        }
    ]


Returns:

* `secret_name` [string] - The name of the secret matching the
  `secret_hint` parameter. By default this is the secret holding the
  SA's bearer token.

Example playbook usage:

Register a return value from oc_serviceaccount_secret with and pass
that result to this filter plugin.

    - name: Get all SA Secrets
      oc_serviceaccount_secret:
        state: list
        service_account: management-admin
        namespace: management-infra
      register: sa

    - name: Save the SA bearer token secret name
      set_fact:
        management_token: "{{ sa.results | lib_utils_oo_filter_sa_secrets }}"

    - name: Get the SA bearer token value
      oc_secret:
        state: list
        name: "{{ management_token }}"
        namespace: management-infra
        decode: true
      register: sa_secret

    - name: Print the bearer token value
      debug:
        var: sa_secret.results.decoded.token

    """
    secret_name = None

    for secret in sa_secrets:
        # each secret is a hash
        if secret['name'].find(secret_hint) == -1:
            continue
        else:
            secret_name = secret['name']
            break

    return secret_name


def lib_utils_oo_l_of_d_to_csv(input_list):
    """Map a list of dictionaries, input_list, into a csv string
    of json values.

    Example input:
    [{'var1': 'val1', 'var2': 'val2'}, {'var1': 'val3', 'var2': 'val4'}]
    Example output:
    u'{"var1": "val1", "var2": "val2"},{"var1": "val3", "var2": "val4"}'
    """
    return ','.join(json.dumps(x) for x in input_list)


def map_from_pairs(source, delim="="):
    ''' Returns a dict given the source and delim delimited '''
    if source == '':
        return dict()

    return dict(item.split(delim) for item in source.split(","))


def map_to_pairs(source, delim="="):
    ''' Returns a comma separated str given the source as a dict '''

    # Some default selectors are empty strings.
    if source == {} or source == '':
        return str()

    return ','.join(["{}{}{}".format(key, delim, value) for key, value in iteritems(source)])


def lib_utils_oo_etcd_host_urls(hosts, use_ssl=True, port='2379'):
    '''Return a list of urls for etcd hosts'''
    urls = []
    port = str(port)
    proto = "https://" if use_ssl else "http://"
    for host in hosts:
        url_string = "{}{}:{}".format(proto, host, port)
        urls.append(url_string)
    return urls


def lib_utils_mutate_htpass_provider(idps):
    '''Updates identityProviders list to mutate filename of htpasswd auth
    to hardcode filename = /etc/origin/master/htpasswd'''
    old_keys = ('filename', 'fileName', 'file_name')
    for idp in idps:
        if 'provider' in idp:
            idp_p = idp['provider']
            if idp_p['kind'] == 'HTPasswdPasswordIdentityProvider':
                for old_key in old_keys:
                    if old_key in idp_p:
                        idp_p.pop(old_key)
                idp_p['file'] = '/etc/origin/master/htpasswd'
    return idps


def lib_utils_oo_oreg_image(image_default, oreg_url):
    '''Converts default image string to utilize oreg_url, if defined.
       oreg_url should be passed in as string "None" if undefined.

       Example input:  "quay.io/coreos/etcd:v99",
                       "example.com/openshift/origin-${component}:${version}"
       Example output: "example.com/coreos/etcd:v99"'''
    # if no oreg_url is specified, we just return the original default
    if oreg_url == 'None':
        return image_default
    oreg_parts = oreg_url.rsplit('/', 2)
    if len(oreg_parts) < 2:
        raise errors.AnsibleFilterError("oreg_url malformed: {}".format(oreg_url))
    if not (len(oreg_parts) >= 3 and '.' in oreg_parts[0]):
        # oreg_url does not include host information; we'll just return etcd default
        return image_default

    image_parts = image_default.split('/')
    if len(image_parts) < 3:
        raise errors.AnsibleFilterError("default image dictionary malformed, do not adjust this value.")
    return '/'.join([oreg_parts[0], image_parts[1], image_parts[2]])


def lib_utils_oo_list_of_dict_to_dict_from_key(input_list, keyname):
    '''Converts a list of dictionaries to a dictionary with keyname: dictionary

       Example input: [{'name': 'first', 'url': 'x.com'}, {'name': 'second', 'url': 'y.com'}],
                      'name'
       Example output: {'first': {'url': 'x.com', 'name': 'first'}, 'second': {'url': 'y.com', 'name': 'second'}}'''
    output_dict = {}
    for item in input_list:
        retrieved_val = item.get(keyname)
        if keyname is not None:
            output_dict[retrieved_val] = item
    return output_dict


class FilterModule(object):
    """ Custom ansible filter mapping """

    # pylint: disable=no-self-use, too-few-public-methods
    def filters(self):
        """ returns a mapping of filters to methods """
        return {
            "lib_utils_oo_select_keys": lib_utils_oo_select_keys,
            "lib_utils_oo_select_keys_from_list": lib_utils_oo_select_keys_from_list,
            "lib_utils_oo_collect": lib_utils_oo_collect,
            "lib_utils_oo_pdb": lib_utils_oo_pdb,
            "lib_utils_oo_prepend_strings_in_list": lib_utils_oo_prepend_strings_in_list,
            "lib_utils_oo_dict_to_list_of_dict": lib_utils_oo_dict_to_list_of_dict,
            "lib_utils_oo_split": lib_utils_oo_split,
            "lib_utils_oo_dict_to_keqv_list": lib_utils_oo_dict_to_keqv_list,
            "lib_utils_oo_list_to_dict": lib_utils_oo_list_to_dict,
            "lib_utils_oo_parse_named_certificates": lib_utils_oo_parse_named_certificates,
            "lib_utils_oo_parse_certificate_san": lib_utils_oo_parse_certificate_san,
            "lib_utils_oo_generate_secret": lib_utils_oo_generate_secret,
            "lib_utils_oo_image_tag_to_rpm_version": lib_utils_oo_image_tag_to_rpm_version,
            "lib_utils_oo_hostname_from_url": lib_utils_oo_hostname_from_url,
            "lib_utils_oo_loadbalancer_frontends": lib_utils_oo_loadbalancer_frontends,
            "lib_utils_oo_loadbalancer_backends": lib_utils_oo_loadbalancer_backends,
            "lib_utils_to_padded_yaml": lib_utils_to_padded_yaml,
            "lib_utils_oo_random_word": lib_utils_oo_random_word,
            "lib_utils_oo_selector_to_string_list": lib_utils_oo_selector_to_string_list,
            "lib_utils_oo_filter_sa_secrets": lib_utils_oo_filter_sa_secrets,
            "lib_utils_oo_l_of_d_to_csv": lib_utils_oo_l_of_d_to_csv,
            "map_from_pairs": map_from_pairs,
            "map_to_pairs": map_to_pairs,
            "lib_utils_oo_etcd_host_urls": lib_utils_oo_etcd_host_urls,
            "lib_utils_mutate_htpass_provider": lib_utils_mutate_htpass_provider,
            "lib_utils_oo_oreg_image": lib_utils_oo_oreg_image,
            "lib_utils_oo_list_of_dict_to_dict_from_key": lib_utils_oo_list_of_dict_to_dict_from_key,
        }
