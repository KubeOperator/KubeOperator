#!/usr/bin/python
# pylint: disable=too-many-lines
# -*- coding: utf-8 -*-
# Reason: Disable pylint too-many-lines because we don't want to split up this file.
# Status: Permanently disabled to keep this module as self-contained as possible.

"""Ansible module for retrieving and setting openshift related facts"""

# pylint: disable=no-name-in-module, import-error, wrong-import-order
import copy
import errno
import json
import re
import os
import yaml
import struct
import socket
import ipaddress
from distutils.util import strtobool
from ansible.module_utils.six import text_type
from ansible.module_utils.six import string_types

# ignore pylint errors related to the module_utils import
# pylint: disable=redefined-builtin, unused-wildcard-import, wildcard-import
# import module snippets
from ansible.module_utils.basic import *  # noqa: F403
from ansible.module_utils.facts import *  # noqa: F403
from ansible.module_utils.urls import *  # noqa: F403
from ansible.module_utils.six import iteritems, itervalues
from ansible.module_utils.six.moves.urllib.parse import urlparse, urlunparse
from ansible.module_utils._text import to_native


DOCUMENTATION = '''
---
module: openshift_facts
short_description: Cluster Facts
author: Jason DeTiberus
requirements: [ ]
'''
EXAMPLES = '''
'''


def migrate_admission_plugin_facts(facts):
    """ Apply migrations for admission plugin facts """
    if 'master' in facts:
        if 'kube_admission_plugin_config' in facts['master']:
            if 'admission_plugin_config' not in facts['master']:
                facts['master']['admission_plugin_config'] = dict()
            # Merge existing kube_admission_plugin_config with admission_plugin_config.
            facts['master']['admission_plugin_config'] = merge_facts(facts['master']['admission_plugin_config'],
                                                                     facts['master']['kube_admission_plugin_config'],
                                                                     additive_facts_to_overwrite=[])
            # Remove kube_admission_plugin_config fact
            facts['master'].pop('kube_admission_plugin_config', None)
    return facts


def migrate_local_facts(facts):
    """ Apply migrations of local facts """
    migrated_facts = copy.deepcopy(facts)
    migrated_facts = migrate_admission_plugin_facts(migrated_facts)
    return migrated_facts


def first_ip(network):
    """ Return the first IPv4 address in network

        Args:
            network (str): network in CIDR format
        Returns:
            str: first IPv4 address
    """
    atoi = lambda addr: struct.unpack("!I", socket.inet_aton(addr))[0]  # noqa: E731
    itoa = lambda addr: socket.inet_ntoa(struct.pack("!I", addr))  # noqa: E731

    (address, netmask) = network.split('/')
    netmask_i = (0xffffffff << (32 - atoi(netmask))) & 0xffffffff
    return itoa((atoi(address) & netmask_i) + 1)


def hostname_valid(hostname):
    """ Test if specified hostname should be considered valid

        Args:
            hostname (str): hostname to test
        Returns:
            bool: True if valid, otherwise False
    """
    if (not hostname or
            hostname.startswith('localhost') or
            hostname.endswith('localdomain') or
            # OpenShift will not allow a node with more than 63 chars in name.
            len(hostname) > 63):
        return False

    return True


def choose_hostname(hostnames=None, fallback=''):
    """ Choose a hostname from the provided hostnames

        Given a list of hostnames and a fallback value, choose a hostname to
        use. This function will prefer fqdns if they exist (excluding any that
        begin with localhost or end with localdomain) over ip addresses.

        Args:
            hostnames (list): list of hostnames
            fallback (str): default value to set if hostnames does not contain
                            a valid hostname
        Returns:
            str: chosen hostname
    """
    hostname = fallback
    if hostnames is None:
        return hostname

    ip_regex = r'\A\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\Z'
    ips = [i for i in hostnames if i is not None and isinstance(i, string_types) and re.match(ip_regex, i)]
    hosts = [i for i in hostnames if i is not None and i != '' and i not in ips]

    for host_list in (hosts, ips):
        for host in host_list:
            if hostname_valid(host):
                return host

    return hostname


def query_metadata(metadata_url, headers=None, expect_json=False):
    """ Return metadata from the provided metadata_url

        Args:
            metadata_url (str): metadata url
            headers (dict): headers to set for metadata request
            expect_json (bool): does the metadata_url return json
        Returns:
            dict or list: metadata request result
    """
    result, info = fetch_url(module, metadata_url, headers=headers)  # noqa: F405
    if info['status'] != 200:
        raise OpenShiftFactsMetadataUnavailableError("Metadata unavailable")
    if expect_json:
        return module.from_json(to_native(result.read()))  # noqa: F405
    else:
        return [to_native(line.strip()) for line in result.readlines()]


def walk_metadata(metadata_url, headers=None, expect_json=False):
    """ Walk the metadata tree and return a dictionary of the entire tree

        Args:
            metadata_url (str): metadata url
            headers (dict): headers to set for metadata request
            expect_json (bool): does the metadata_url return json
        Returns:
            dict: the result of walking the metadata tree
    """
    metadata = dict()

    for line in query_metadata(metadata_url, headers, expect_json):
        if line.endswith('/') and not line == 'public-keys/':
            key = line[:-1]
            metadata[key] = walk_metadata(metadata_url + line,
                                          headers, expect_json)
        else:
            results = query_metadata(metadata_url + line, headers,
                                     expect_json)
            if len(results) == 1:
                # disable pylint maybe-no-member because overloaded use of
                # the module name causes pylint to not detect that results
                # is an array or hash
                # pylint: disable=maybe-no-member
                metadata[line] = results.pop()
            else:
                metadata[line] = results
    return metadata


def get_provider_metadata(metadata_url, supports_recursive=False,
                          headers=None, expect_json=False):
    """ Retrieve the provider metadata

        Args:
            metadata_url (str): metadata url
            supports_recursive (bool): does the provider metadata api support
                                       recursion
            headers (dict): headers to set for metadata request
            expect_json (bool): does the metadata_url return json
        Returns:
            dict: the provider metadata
    """
    try:
        if supports_recursive:
            metadata = query_metadata(metadata_url, headers,
                                      expect_json)
        else:
            metadata = walk_metadata(metadata_url, headers,
                                     expect_json)
    except OpenShiftFactsMetadataUnavailableError:
        metadata = None
    return metadata


def normalize_gce_facts(metadata, facts):
    """ Normalize gce facts

        Args:
            metadata (dict): provider metadata
            facts (dict): facts to update
        Returns:
            dict: the result of adding the normalized metadata to the provided
                  facts dict
    """
    for interface in metadata['instance']['networkInterfaces']:
        int_info = dict(ips=[interface['ip']], network_type='gce')
        int_info['public_ips'] = [ac['externalIp'] for ac
                                  in interface['accessConfigs']]
        int_info['public_ips'].extend(interface['forwardedIps'])
        _, _, network_id = interface['network'].rpartition('/')
        int_info['network_id'] = network_id
        facts['network']['interfaces'].append(int_info)
    _, _, zone = metadata['instance']['zone'].rpartition('/')
    facts['zone'] = zone

    # GCE currently only supports a single interface
    facts['network']['ip'] = facts['network']['interfaces'][0]['ips'][0]
    pub_ip = facts['network']['interfaces'][0]['public_ips'][0]
    facts['network']['public_ip'] = pub_ip
    # Split instance hostname from GCE metadata to use the short instance name
    facts['network']['hostname'] = metadata['instance']['hostname'].split('.')[0]

    # TODO: attempt to resolve public_hostname
    facts['network']['public_hostname'] = facts['network']['public_ip']

    return facts


def normalize_aws_facts(metadata, facts):
    """ Normalize aws facts

        Args:
            metadata (dict): provider metadata
            facts (dict): facts to update
        Returns:
            dict: the result of adding the normalized metadata to the provided
                  facts dict
    """
    for interface in sorted(
            metadata['network']['interfaces']['macs'].values(),
            key=lambda x: x['device-number']
    ):
        int_info = dict()
        var_map = {'ips': 'local-ipv4s', 'public_ips': 'public-ipv4s'}
        for ips_var, int_var in iteritems(var_map):
            ips = interface.get(int_var)
            if isinstance(ips, string_types):
                int_info[ips_var] = [ips]
            else:
                int_info[ips_var] = ips
        if 'vpc-id' in interface:
            int_info['network_type'] = 'vpc'
        else:
            int_info['network_type'] = 'classic'
        if int_info['network_type'] == 'vpc':
            int_info['network_id'] = interface['subnet-id']
        else:
            int_info['network_id'] = None
        facts['network']['interfaces'].append(int_info)
    facts['zone'] = metadata['placement']['availability-zone']

    # TODO: actually attempt to determine default local and public ips
    # by using the ansible default ip fact and the ipv4-associations
    # from the ec2 metadata
    facts['network']['ip'] = metadata.get('local-ipv4')
    facts['network']['public_ip'] = metadata.get('public-ipv4')

    # TODO: verify that local hostname makes sense and is resolvable
    facts['network']['hostname'] = metadata.get('local-hostname')

    # TODO: verify that public hostname makes sense and is resolvable
    facts['network']['public_hostname'] = metadata.get('public-hostname')

    return facts


def normalize_openstack_facts(metadata, facts):
    """ Normalize openstack facts

        Args:
            metadata (dict): provider metadata
            facts (dict): facts to update
        Returns:
            dict: the result of adding the normalized metadata to the provided
                  facts dict
    """
    # openstack ec2 compat api does not support network interfaces and
    # the version tested on did not include the info in the openstack
    # metadata api, should be updated if neutron exposes this.

    facts['zone'] = metadata['availability_zone']
    local_ipv4 = metadata['ec2_compat']['local-ipv4'].split(',')[0]
    facts['network']['ip'] = local_ipv4
    facts['network']['public_ip'] = metadata['ec2_compat']['public-ipv4']

    for f_var, h_var, ip_var in [('hostname', 'hostname', 'local-ipv4'),
                                 ('public_hostname', 'public-hostname', 'public-ipv4')]:
        try:
            if socket.gethostbyname(metadata['ec2_compat'][h_var]) == metadata['ec2_compat'][ip_var]:
                facts['network'][f_var] = metadata['ec2_compat'][h_var]
            else:
                facts['network'][f_var] = metadata['ec2_compat'][ip_var]
        except socket.gaierror:
            facts['network'][f_var] = metadata['ec2_compat'][ip_var]

    return facts


def normalize_provider_facts(provider, metadata):
    """ Normalize provider facts

        Args:
            provider (str): host provider
            metadata (dict): provider metadata
        Returns:
            dict: the normalized provider facts
    """
    if provider is None or metadata is None:
        return {}

    # TODO: test for ipv6_enabled where possible (gce, aws do not support)
    # and configure ipv6 facts if available

    # TODO: add support for setting user_data if available

    facts = dict(name=provider, metadata=metadata,
                 network=dict(interfaces=[], ipv6_enabled=False))

    if provider == 'gce':
        facts = normalize_gce_facts(metadata, facts)
    elif provider == 'aws':
        facts = normalize_aws_facts(metadata, facts)
    elif provider == 'openstack':
        facts = normalize_openstack_facts(metadata, facts)
    return facts


def set_url_facts_if_unset(facts):
    """ Set url facts if not already present in facts dict

        Args:
            facts (dict): existing facts
        Returns:
            dict: the facts dict updated with the generated url facts if they
                  were not already present
    """
    if 'master' in facts:
        hostname = facts['common']['hostname']
        cluster_hostname = facts['master'].get('cluster_hostname')
        cluster_public_hostname = facts['master'].get('cluster_public_hostname')
        public_hostname = facts['common']['public_hostname']
        api_hostname = cluster_hostname if cluster_hostname else hostname
        api_public_hostname = cluster_public_hostname if cluster_public_hostname else public_hostname
        console_path = facts['master']['console_path']

        use_ssl = dict(
            api=facts['master']['api_use_ssl'],
            public_api=facts['master']['api_use_ssl'],
            loopback_api=facts['master']['api_use_ssl'],
            console=facts['master']['console_use_ssl'],
            public_console=facts['master']['console_use_ssl'],
        )

        ports = dict(
            api=facts['master']['api_port'],
            public_api=facts['master']['api_port'],
            loopback_api=facts['master']['api_port'],
            console=facts['master']['console_port'],
            public_console=facts['master']['console_port'],
        )

        prefix_hosts = [('api', api_hostname),
                        ('public_api', api_public_hostname),
                        ('loopback_api', hostname)]

        for prefix, host in prefix_hosts:
            facts['master'].setdefault(prefix + '_url', format_url(use_ssl[prefix],
                                                                   host,
                                                                   ports[prefix]))

        r_lhn = "{0}:{1}".format(hostname, ports['api']).replace('.', '-')
        r_lhu = "system:openshift-master/{0}:{1}".format(api_hostname, ports['api']).replace('.', '-')
        facts['master'].setdefault('loopback_cluster_name', r_lhn)
        facts['master'].setdefault('loopback_context_name', "default/{0}/system:openshift-master".format(r_lhn))
        facts['master'].setdefault('loopback_user', r_lhu)

        prefix_hosts = [('console', api_hostname), ('public_console', api_public_hostname)]
        for prefix, host in prefix_hosts:
            facts['master'].setdefault(prefix + '_url', format_url(use_ssl[prefix],
                                                                   host,
                                                                   ports[prefix],
                                                                   console_path))

    return facts


def set_aggregate_facts(facts):
    """ Set aggregate facts

        Args:
            facts (dict): existing facts
        Returns:
            dict: the facts dict updated with aggregated facts
    """
    all_hostnames = set()
    internal_hostnames = set()
    kube_svc_ip = first_ip(facts['common']['portal_net'])
    if 'common' in facts:
        all_hostnames.add(facts['common']['hostname'])
        all_hostnames.add(facts['common']['public_hostname'])
        all_hostnames.add(facts['common']['ip'])
        all_hostnames.add(facts['common']['public_ip'])
        facts['common']['kube_svc_ip'] = kube_svc_ip

        internal_hostnames.add(facts['common']['hostname'])
        internal_hostnames.add(facts['common']['ip'])

        cluster_domain = facts['common']['dns_domain']

        if 'master' in facts:
            if 'cluster_hostname' in facts['master']:
                all_hostnames.add(facts['master']['cluster_hostname'])
            if 'cluster_public_hostname' in facts['master']:
                all_hostnames.add(facts['master']['cluster_public_hostname'])
            svc_names = ['openshift', 'openshift.default', 'openshift.default.svc',
                         'openshift.default.svc.' + cluster_domain, 'kubernetes', 'kubernetes.default',
                         'kubernetes.default.svc', 'kubernetes.default.svc.' + cluster_domain]
            all_hostnames.update(svc_names)
            internal_hostnames.update(svc_names)
            all_hostnames.add(kube_svc_ip)
            internal_hostnames.add(kube_svc_ip)

        facts['common']['all_hostnames'] = list(all_hostnames)
        facts['common']['internal_hostnames'] = list(internal_hostnames)

    return facts


def set_sdn_facts_if_unset(facts, system_facts):
    """ Set sdn facts if not already present in facts dict

        Args:
            facts (dict): existing facts
            system_facts (dict): ansible_facts
        Returns:
            dict: the facts dict updated with the generated sdn facts if they
                  were not already present
    """

    if 'node' in facts and 'sdn_mtu' not in facts['node']:
        node_ip = facts['common']['ip']

        # default MTU if interface MTU cannot be detected
        facts['node']['sdn_mtu'] = '1450'

        for val in itervalues(system_facts):
            if isinstance(val, dict) and 'mtu' in val:
                mtu = val['mtu']

                if 'ipv4' in val and val['ipv4'].get('address') == node_ip:
                    facts['node']['sdn_mtu'] = str(mtu - 50)

    return facts


def set_nodename(facts):
    """ set nodename """
    if 'node' in facts and 'common' in facts:
        if 'cloudprovider' in facts and facts['cloudprovider']['kind'] == 'gce':
            facts['node']['nodename'] = facts['provider']['metadata']['instance']['hostname'].split('.')[0]

        # TODO: The openstack cloudprovider nodename setting was too opinionaed.
        #       It needs to be generalized before it can be enabled again.
        # elif 'cloudprovider' in facts and facts['cloudprovider']['kind'] == 'openstack':
        #     facts['node']['nodename'] = facts['provider']['metadata']['hostname'].replace('.novalocal', '')
        else:
            facts['node']['nodename'] = facts['common']['raw_hostname'].lower()
    return facts


def make_allowed_registries(registry_list):
    """ turns a list of wildcard registries to allowedRegistriesForImport json setting """
    return {
        "allowedRegistriesForImport": [
            {'domainName': reg} if isinstance(reg, str) else reg for reg in registry_list
        ]
    }


def set_allowed_registries(facts):
    """ override allowedRegistriesForImport in imagePolicyConfig """
    if 'master' in facts:
        image_policy = {}
        overriden = False
        if facts['master'].get('image_policy_config', None):
            image_policy = facts['master']['image_policy_config']
            overriden = True

        overrides = facts['master'].get('image_policy_allowed_registries_for_import', None)
        if overrides:
            image_policy = merge_facts(image_policy, make_allowed_registries(overrides), None)
            overriden = True

        if overriden:
            facts['master']['image_policy_config'] = image_policy
    return facts


def format_url(use_ssl, hostname, port, path=''):
    """ Format url based on ssl flag, hostname, port and path

        Args:
            use_ssl (bool): is ssl enabled
            hostname (str): hostname
            port (str): port
            path (str): url path
        Returns:
            str: The generated url string
    """
    scheme = 'https' if use_ssl else 'http'
    netloc = hostname
    if (use_ssl and port != '443') or (not use_ssl and port != '80'):
        netloc += ":%s" % port
    try:
        url = urlparse.urlunparse((scheme, netloc, path, '', '', ''))
    except AttributeError:
        # pylint: disable=undefined-variable
        url = urlunparse((scheme, netloc, path, '', '', ''))
    return url


def get_current_config(facts):
    """ Get current openshift config

        Args:
            facts (dict): existing facts
        Returns:
            dict: the facts dict updated with the current openshift config
    """
    current_config = dict()
    roles = [role for role in facts if role not in ['common', 'provider']]
    for role in roles:
        if 'roles' in current_config:
            current_config['roles'].append(role)
        else:
            current_config['roles'] = [role]

        # TODO: parse the /etc/sysconfig/openshift-{master,node} config to
        # determine the location of files.
        # TODO: I suspect this isn't working right now, but it doesn't prevent
        # anything from working properly as far as I can tell, perhaps because
        # we override the kubeconfig path everywhere we use it?
        # Query kubeconfig settings
        kubeconfig_dir = '/var/lib/origin/openshift.local.certificates'
        if role == 'node':
            kubeconfig_dir = os.path.join(
                kubeconfig_dir, "node-%s" % facts['common']['hostname']
            )

        kubeconfig_path = os.path.join(kubeconfig_dir, '.kubeconfig')
        if os.path.isfile('/usr/bin/openshift') and os.path.isfile(kubeconfig_path):
            try:
                _, output, _ = module.run_command(  # noqa: F405
                    ["/usr/bin/openshift", "ex", "config", "view", "-o",
                     "json", "--kubeconfig=%s" % kubeconfig_path],
                    check_rc=False
                )
                config = json.loads(output)

                cad = 'certificate-authority-data'
                try:
                    for cluster in config['clusters']:
                        config['clusters'][cluster][cad] = 'masked'
                except KeyError:
                    pass
                try:
                    for user in config['users']:
                        config['users'][user][cad] = 'masked'
                        config['users'][user]['client-key-data'] = 'masked'
                except KeyError:
                    pass

                current_config['kubeconfig'] = config

            # override pylint broad-except warning, since we do not want
            # to bubble up any exceptions if oc config view
            # fails
            # pylint: disable=broad-except
            except Exception:
                pass

    return current_config


def build_controller_args(facts):
    """ Build master controller_args """
    cloud_cfg_path = os.path.join(facts['common']['config_base'],
                                  'cloudprovider')
    if 'master' in facts:
        controller_args = {}

        if 'cloudprovider' in facts:
            if 'kind' in facts['cloudprovider']:
                if facts['cloudprovider']['kind'] == 'aws':
                    controller_args['cloud-provider'] = ['aws']
                    controller_args['cloud-config'] = [cloud_cfg_path + '/aws.conf']
                    controller_args['disable-attach-detach-reconcile-sync'] = ['true']
                if facts['cloudprovider']['kind'] == 'azure':
                    controller_args['cloud-provider'] = ['azure']
                    controller_args['cloud-config'] = [cloud_cfg_path + '/azure.conf']
                if facts['cloudprovider']['kind'] == 'openstack':
                    controller_args['cloud-provider'] = ['openstack']
                    controller_args['cloud-config'] = [cloud_cfg_path + '/openstack.conf']
                if facts['cloudprovider']['kind'] == 'gce':
                    controller_args['cloud-provider'] = ['gce']
                    controller_args['cloud-config'] = [cloud_cfg_path + '/gce.conf']
        if controller_args != {}:
            facts = merge_facts({'master': {'controller_args': controller_args}}, facts, [])
    return facts


def build_api_server_args(facts):
    """ Build master api_server_args """
    cloud_cfg_path = os.path.join(facts['common']['config_base'],
                                  'cloudprovider')
    if 'master' in facts:
        api_server_args = {}
        if 'cloudprovider' in facts:
            if 'kind' in facts['cloudprovider']:
                if facts['cloudprovider']['kind'] == 'aws':
                    api_server_args['cloud-provider'] = ['aws']
                    api_server_args['cloud-config'] = [cloud_cfg_path + '/aws.conf']
                if facts['cloudprovider']['kind'] == 'azure':
                    api_server_args['cloud-provider'] = ['azure']
                    api_server_args['cloud-config'] = [cloud_cfg_path + '/azure.conf']
                if facts['cloudprovider']['kind'] == 'openstack':
                    api_server_args['cloud-provider'] = ['openstack']
                    api_server_args['cloud-config'] = [cloud_cfg_path + '/openstack.conf']
                if facts['cloudprovider']['kind'] == 'gce':
                    api_server_args['cloud-provider'] = ['gce']
                    api_server_args['cloud-config'] = [cloud_cfg_path + '/gce.conf']
        if api_server_args != {}:
            facts = merge_facts({'master': {'api_server_args': api_server_args}}, facts, [])
    return facts


def apply_provider_facts(facts, provider_facts):
    """ Apply provider facts to supplied facts dict

        Args:
            facts (dict): facts dict to update
            provider_facts (dict): provider facts to apply
            roles: host roles
        Returns:
            dict: the merged facts
    """
    if not provider_facts:
        return facts

    common_vars = [('hostname', 'ip'), ('public_hostname', 'public_ip')]
    for h_var, ip_var in common_vars:
        ip_value = provider_facts['network'].get(ip_var)
        if ip_value:
            facts['common'][ip_var] = ip_value

        facts['common'][h_var] = choose_hostname(
            [provider_facts['network'].get(h_var)],
            facts['common'][h_var]
        )

    facts['provider'] = provider_facts
    return facts


# Disabling pylint too many branches. This function needs refactored
# but is a very core part of openshift_facts.
# pylint: disable=too-many-branches, too-many-nested-blocks
def merge_facts(orig, new, additive_facts_to_overwrite):
    """ Recursively merge facts dicts

        Args:
            orig (dict): existing facts
            new (dict): facts to update
            additive_facts_to_overwrite (list): additive facts to overwrite in jinja
                                                '.' notation ex: ['master.named_certificates']

        Returns:
            dict: the merged facts
    """
    additive_facts = ['named_certificates']

    # Facts we do not ever want to merge. These originate in inventory variables
    # and contain JSON dicts. We don't ever want to trigger a merge
    # here, just completely overwrite with the new if they are present there.
    inventory_json_facts = ['admission_plugin_config',
                            'kube_admission_plugin_config',
                            'image_policy_config',
                            "builddefaults",
                            "buildoverrides"]

    facts = dict()
    for key, value in iteritems(orig):
        # Key exists in both old and new facts.
        if key in new:
            if key in inventory_json_facts:
                # Watchout for JSON facts that sometimes load as strings.
                # (can happen if the JSON contains a boolean)
                if isinstance(new[key], string_types):
                    facts[key] = yaml.safe_load(new[key])
                else:
                    facts[key] = copy.deepcopy(new[key])
            # Continue to recurse if old and new fact is a dictionary.
            elif isinstance(value, dict) and isinstance(new[key], dict):
                # Collect the subset of additive facts to overwrite if
                # key matches. These will be passed to the subsequent
                # merge_facts call.
                relevant_additive_facts = []
                for item in additive_facts_to_overwrite:
                    if '.' in item and item.startswith(key + '.'):
                        relevant_additive_facts.append(item)

                facts[key] = merge_facts(value, new[key], relevant_additive_facts)
            # Key matches an additive fact and we are not overwriting
            # it so we will append the new value to the existing value.
            elif key in additive_facts and key not in [x.split('.')[-1] for x in additive_facts_to_overwrite]:
                if isinstance(value, list) and isinstance(new[key], list):
                    new_fact = []
                    for item in copy.deepcopy(value) + copy.deepcopy(new[key]):
                        if item not in new_fact:
                            new_fact.append(item)
                    facts[key] = new_fact
            # No other condition has been met. Overwrite the old fact
            # with the new value.
            else:
                facts[key] = copy.deepcopy(new[key])
        # Key isn't in new so add it to facts to keep it.
        else:
            facts[key] = copy.deepcopy(value)
    new_keys = set(new.keys()) - set(orig.keys())
    for key in new_keys:
        # Watchout for JSON facts that sometimes load as strings.
        # (can happen if the JSON contains a boolean)
        if key in inventory_json_facts and isinstance(new[key], string_types):
            facts[key] = yaml.safe_load(new[key])
        else:
            facts[key] = copy.deepcopy(new[key])
    return facts


def save_local_facts(filename, facts):
    """ Save local facts

        Args:
            filename (str): local facts file
            facts (dict): facts to set
    """
    try:
        fact_dir = os.path.dirname(filename)
        try:
            os.makedirs(fact_dir)  # try to make the directory
        except OSError as exception:
            if exception.errno != errno.EEXIST:  # but it is okay if it is already there
                raise  # pass any other exceptions up the chain
        with open(filename, 'w') as fact_file:
            fact_file.write(module.jsonify(facts))  # noqa: F405
        os.chmod(filename, 0o600)
    except (IOError, OSError) as ex:
        raise OpenShiftFactsFileWriteError(
            "Could not create fact file: %s, error: %s" % (filename, ex)
        )


def get_local_facts_from_file(filename):
    """ Retrieve local facts from fact file

        Args:
            filename (str): local facts file
        Returns:
            dict: the retrieved facts
    """
    try:
        with open(filename, 'r') as facts_file:
            local_facts = json.load(facts_file)
    except (ValueError, IOError):
        local_facts = {}

    return local_facts


def sort_unique(alist):
    """ Sorts and de-dupes a list

        Args:
            list: a list
        Returns:
            list: a sorted de-duped list
    """
    return sorted(list(set(alist)))


def safe_get_bool(fact):
    """ Get a boolean fact safely.

        Args:
            facts: fact to convert
        Returns:
            bool: given fact as a bool
    """
    return bool(strtobool(str(fact)))


def set_proxy_facts(facts):
    """ Set global proxy facts

        Args:
            facts(dict): existing facts
        Returns:
            facts(dict): Updated facts with missing values
    """
    if 'common' in facts:
        common = facts['common']
        if 'http_proxy' in common or 'https_proxy' in common or 'no_proxy' in common:
            if 'no_proxy' in common and isinstance(common['no_proxy'], string_types):
                common['no_proxy'] = common['no_proxy'].split(",")
            elif 'no_proxy' not in common:
                common['no_proxy'] = []

            # See https://bugzilla.redhat.com/show_bug.cgi?id=1466783
            # masters behind a proxy need to connect to etcd via IP
            if 'no_proxy_etcd_host_ips' in common:
                if isinstance(common['no_proxy_etcd_host_ips'], string_types):
                    common['no_proxy'].extend(common['no_proxy_etcd_host_ips'].split(','))

            # Master IPs should be added to no proxy lists to make liveness probes to pass
            if 'no_proxy_master_ips' in common:
                if isinstance(common['no_proxy_master_ips'], string_types):
                    common['no_proxy'].extend(common['no_proxy_master_ips'].split(','))

            if 'generate_no_proxy_hosts' in common and safe_get_bool(common['generate_no_proxy_hosts']):
                if 'no_proxy_internal_hostnames' in common:
                    common['no_proxy'].extend(common['no_proxy_internal_hostnames'].split(','))
            # TODO: This is Azure specific and should be scoped out to only Azure installs
            common['no_proxy'].append('169.254.169.254')
            # We always add local dns domain and ourselves no matter what
            kube_svc_ip = str(ipaddress.ip_network(text_type(common['portal_net']))[1])
            common['no_proxy'].append(kube_svc_ip)
            common['no_proxy'].append('.' + common['dns_domain'])
            common['no_proxy'].append('.svc')
            common['no_proxy'].append(common['hostname'])
            if 'master' in facts:
                if 'cluster_hostname' in facts['master']:
                    common['no_proxy'].append(facts['master']['cluster_hostname'])
            common['no_proxy'] = ','.join(sort_unique(common['no_proxy']))
        facts['common'] = common

    return facts


def set_builddefaults_facts(facts):
    """ Set build defaults including setting proxy values from http_proxy, https_proxy,
        no_proxy to the more specific builddefaults and builddefaults_git vars.
           1. http_proxy, https_proxy, no_proxy
           2. builddefaults_*
           3. builddefaults_git_*

        Args:
            facts(dict): existing facts
        Returns:
            facts(dict): Updated facts with missing values
    """

    if 'builddefaults' in facts:
        builddefaults = facts['builddefaults']
        common = facts['common']
        # Copy values from common to builddefaults
        if 'http_proxy' not in builddefaults and 'http_proxy' in common:
            builddefaults['http_proxy'] = common['http_proxy']
        if 'https_proxy' not in builddefaults and 'https_proxy' in common:
            builddefaults['https_proxy'] = common['https_proxy']
        if 'no_proxy' not in builddefaults and 'no_proxy' in common:
            builddefaults['no_proxy'] = common['no_proxy']

        # Create git specific facts from generic values, if git specific values are
        # not defined.
        if 'git_http_proxy' not in builddefaults and 'http_proxy' in builddefaults:
            builddefaults['git_http_proxy'] = builddefaults['http_proxy']
        if 'git_https_proxy' not in builddefaults and 'https_proxy' in builddefaults:
            builddefaults['git_https_proxy'] = builddefaults['https_proxy']
        if 'git_no_proxy' not in builddefaults and 'no_proxy' in builddefaults:
            builddefaults['git_no_proxy'] = builddefaults['no_proxy']
        # If we're actually defining a builddefaults config then create admission_plugin_config
        # then merge builddefaults[config] structure into admission_plugin_config

        # 'config' is the 'openshift_builddefaults_json' inventory variable
        if 'config' in builddefaults:
            if 'admission_plugin_config' not in facts['master']:
                # Scaffold out the full expected datastructure
                facts['master']['admission_plugin_config'] = {'BuildDefaults': {'configuration': {'env': {}}}}
            facts['master']['admission_plugin_config'].update(builddefaults['config'])
            if 'env' in facts['master']['admission_plugin_config']['BuildDefaults']['configuration']:
                delete_empty_keys(facts['master']['admission_plugin_config']['BuildDefaults']['configuration']['env'])

    return facts


def delete_empty_keys(keylist):
    """ Delete dictionary elements from keylist where "value" is empty.

        Args:
          keylist(list): A list of builddefault configuration envs.

        Returns:
          none

        Example:
          keylist = [{'name': 'HTTP_PROXY', 'value': 'http://file.rdu.redhat.com:3128'},
                     {'name': 'HTTPS_PROXY', 'value': 'http://file.rdu.redhat.com:3128'},
                     {'name': 'NO_PROXY', 'value': ''}]

          After calling delete_empty_keys the provided list is modified to become:

                    [{'name': 'HTTP_PROXY', 'value': 'http://file.rdu.redhat.com:3128'},
                     {'name': 'HTTPS_PROXY', 'value': 'http://file.rdu.redhat.com:3128'}]
    """
    count = 0
    for i in range(0, len(keylist)):
        if len(keylist[i - count]['value']) == 0:
            del keylist[i - count]
            count += 1


def set_buildoverrides_facts(facts):
    """ Set build overrides

        Args:
            facts(dict): existing facts
        Returns:
            facts(dict): Updated facts with missing values
    """
    if 'buildoverrides' in facts:
        buildoverrides = facts['buildoverrides']
        # If we're actually defining a buildoverrides config then create admission_plugin_config
        # then merge buildoverrides[config] structure into admission_plugin_config
        if 'config' in buildoverrides:
            if 'admission_plugin_config' not in facts['master']:
                facts['master']['admission_plugin_config'] = dict()
            facts['master']['admission_plugin_config'].update(buildoverrides['config'])

    return facts


def pop_obsolete_local_facts(local_facts):
    """Remove unused keys from local_facts"""
    keys_to_remove = {
        'master': ('etcd_port', 'etcd_use_ssl', 'etcd_hosts')
    }
    for role in keys_to_remove:
        if role in local_facts:
            for key in keys_to_remove[role]:
                local_facts[role].pop(key, None)


class OpenShiftFactsUnsupportedRoleError(Exception):
    """Origin Facts Unsupported Role Error"""
    pass


class OpenShiftFactsFileWriteError(Exception):
    """Origin Facts File Write Error"""
    pass


class OpenShiftFactsMetadataUnavailableError(Exception):
    """Origin Facts Metadata Unavailable Error"""
    pass


class OpenShiftFacts(object):
    """ Origin Facts

        Attributes:
            facts (dict): facts for the host

        Args:
            module (AnsibleModule): an AnsibleModule object
            role (str): role for setting local facts
            filename (str): local facts file to use
            local_facts (dict): local facts to set
            additive_facts_to_overwrite (list): additive facts to overwrite in jinja
                                                '.' notation ex: ['master.named_certificates']

        Raises:
            OpenShiftFactsUnsupportedRoleError:
    """
    known_roles = ['builddefaults',
                   'buildoverrides',
                   'cloudprovider',
                   'common',
                   'etcd',
                   'master',
                   'node']

    # Disabling too-many-arguments, this should be cleaned up as a TODO item.
    # pylint: disable=too-many-arguments,no-value-for-parameter
    def __init__(self, role, filename, local_facts,
                 additive_facts_to_overwrite=None):
        self.changed = False
        self.filename = filename
        if role not in self.known_roles:
            raise OpenShiftFactsUnsupportedRoleError(
                "Role %s is not supported by this module" % role
            )
        self.role = role

        # Collect system facts and preface each fact with 'ansible_'.
        try:
            # pylint: disable=too-many-function-args,invalid-name
            self.system_facts = ansible_facts(module, ['hardware', 'network', 'virtual', 'facter'])  # noqa: F405
            additional_facts = {}
            for (k, v) in self.system_facts.items():
                additional_facts["ansible_%s" % k.replace('-', '_')] = v
            self.system_facts.update(additional_facts)
        except UnboundLocalError:
            # ansible-2.2,2.3
            self.system_facts = get_all_facts(module)['ansible_facts']  # noqa: F405

        self.facts = self.generate_facts(local_facts,
                                         additive_facts_to_overwrite)

    def generate_facts(self,
                       local_facts,
                       additive_facts_to_overwrite):
        """ Generate facts

            Args:
                local_facts (dict): local_facts for overriding generated defaults
                additive_facts_to_overwrite (list): additive facts to overwrite in jinja
                                                    '.' notation ex: ['master.named_certificates']
            Returns:
                dict: The generated facts
        """

        local_facts = self.init_local_facts(local_facts,
                                            additive_facts_to_overwrite)
        roles = local_facts.keys()

        defaults = self.get_defaults(roles)
        provider_facts = {}
        if 'common' in local_facts and 'cloudprovider' in local_facts['common']:
            provider_facts = self.init_provider_facts()
        facts = apply_provider_facts(defaults, provider_facts)
        facts = merge_facts(facts,
                            local_facts,
                            additive_facts_to_overwrite)
        facts['current_config'] = get_current_config(facts)
        facts = set_url_facts_if_unset(facts)
        facts = set_sdn_facts_if_unset(facts, self.system_facts)
        facts = build_controller_args(facts)
        facts = build_api_server_args(facts)
        facts = set_aggregate_facts(facts)
        facts = set_proxy_facts(facts)
        facts = set_builddefaults_facts(facts)
        facts = set_buildoverrides_facts(facts)
        facts = set_nodename(facts)
        facts = set_allowed_registries(facts)
        return dict(openshift=facts)

    def get_defaults(self, roles):
        """ Get default fact values

            Args:
                roles (list): list of roles for this host

            Returns:
                dict: The generated default facts
        """
        defaults = {}
        ip_addr = self.system_facts['ansible_default_ipv4']['address']
        exit_code, output, _ = module.run_command(['hostname', '-f'])  # noqa: F405
        hostname_f = output.strip() if exit_code == 0 else ''
        hostname_values = [hostname_f, self.system_facts['ansible_nodename'],
                           self.system_facts['ansible_fqdn']]
        hostname = choose_hostname(hostname_values, ip_addr).lower()
        exit_code, output, _ = module.run_command(['hostname'])  # noqa: F405
        raw_hostname = output.strip() if exit_code == 0 else hostname

        defaults['common'] = dict(ip=ip_addr,
                                  public_ip=ip_addr,
                                  raw_hostname=raw_hostname,
                                  hostname=hostname,
                                  public_hostname=hostname,
                                  portal_net='172.30.0.0/16',
                                  dns_domain='cluster.local',
                                  config_base='/etc/origin')

        if 'master' in roles:
            defaults['master'] = dict(api_use_ssl=True, api_port='8443',
                                      controllers_port='8444',
                                      console_use_ssl=True,
                                      console_path='/console',
                                      console_port='8443',
                                      portal_net='172.30.0.0/16',
                                      bind_addr='0.0.0.0',
                                      session_max_seconds=3600,
                                      session_name='ssn')

        if 'cloudprovider' in roles:
            defaults['cloudprovider'] = dict(kind=None)

        return defaults

    def guess_host_provider(self):
        """ Guess the host provider

            Returns:
                dict: The generated default facts for the detected provider
        """
        # TODO: cloud provider facts should probably be submitted upstream
        product_name = self.system_facts['ansible_product_name']
        product_version = self.system_facts['ansible_product_version']
        virt_type = self.system_facts['ansible_virtualization_type']
        virt_role = self.system_facts['ansible_virtualization_role']
        bios_vendor = self.system_facts['ansible_system_vendor']
        provider = None
        metadata = None

        if bios_vendor == 'Google':
            provider = 'gce'
            metadata_url = ('http://metadata.google.internal/'
                            'computeMetadata/v1/?recursive=true')
            headers = {'Metadata-Flavor': 'Google'}
            metadata = get_provider_metadata(metadata_url, True, headers,
                                             True)

            # Filter sshKeys and serviceAccounts from gce metadata
            if metadata:
                metadata['project']['attributes'].pop('sshKeys', None)
                metadata['instance'].pop('serviceAccounts', None)
        elif bios_vendor == 'Amazon EC2':
            # Adds support for Amazon EC2 C5 instance types
            provider = 'aws'
            metadata_url = 'http://169.254.169.254/latest/meta-data/'
            metadata = get_provider_metadata(metadata_url)
        elif virt_type == 'xen' and virt_role == 'guest' and re.match(r'.*\.amazon$', product_version):
            provider = 'aws'
            metadata_url = 'http://169.254.169.254/latest/meta-data/'
            metadata = get_provider_metadata(metadata_url)
        elif re.search(r'OpenStack', product_name):
            provider = 'openstack'
            metadata_url = ('http://169.254.169.254/openstack/latest/'
                            'meta_data.json')
            metadata = get_provider_metadata(metadata_url, True, None,
                                             True)

            if metadata:
                ec2_compat_url = 'http://169.254.169.254/latest/meta-data/'
                metadata['ec2_compat'] = get_provider_metadata(
                    ec2_compat_url
                )

                # disable pylint maybe-no-member because overloaded use of
                # the module name causes pylint to not detect that results
                # is an array or hash
                # pylint: disable=maybe-no-member
                # Filter public_keys  and random_seed from openstack metadata
                metadata.pop('public_keys', None)
                metadata.pop('random_seed', None)

                if not metadata['ec2_compat']:
                    metadata = None

        return dict(name=provider, metadata=metadata)

    def init_provider_facts(self):
        """ Initialize the provider facts

            Returns:
                dict: The normalized provider facts
        """
        provider_info = self.guess_host_provider()
        provider_facts = normalize_provider_facts(
            provider_info.get('name'),
            provider_info.get('metadata')
        )
        return provider_facts

    # Disabling too-many-branches and too-many-locals.
    # This should be cleaned up as a TODO item.
    # pylint: disable=too-many-branches, too-many-locals
    def init_local_facts(self, facts=None,
                         additive_facts_to_overwrite=None):
        """ Initialize the local facts

            Args:
                facts (dict): local facts to set
                additive_facts_to_overwrite (list): additive facts to overwrite in jinja
                                                    '.' notation ex: ['master.named_certificates']
            Returns:
                dict: The result of merging the provided facts with existing
                      local facts
        """
        changed = False

        facts_to_set = dict()

        if facts is not None:
            facts_to_set[self.role] = facts

        local_facts = get_local_facts_from_file(self.filename)

        migrated_facts = migrate_local_facts(local_facts)

        new_local_facts = merge_facts(migrated_facts,
                                      facts_to_set,
                                      additive_facts_to_overwrite)

        new_local_facts = self.remove_empty_facts(new_local_facts)
        pop_obsolete_local_facts(new_local_facts)

        if new_local_facts != local_facts:
            changed = True
            if not module.check_mode:  # noqa: F405
                save_local_facts(self.filename, new_local_facts)

        self.changed = changed
        return new_local_facts

    def remove_empty_facts(self, facts=None):
        """ Remove empty facts

            Args:
                facts (dict): facts to clean
        """
        facts_to_remove = []
        for fact, value in iteritems(facts):
            if isinstance(facts[fact], dict):
                facts[fact] = self.remove_empty_facts(facts[fact])
            else:
                if value == "" or value == [""] or value is None:
                    facts_to_remove.append(fact)
        for fact in facts_to_remove:
            del facts[fact]
        return facts


def main():
    """ main """
    # disabling pylint errors for global-variable-undefined and invalid-name
    # for 'global module' usage, since it is required to use ansible_facts
    # pylint: disable=global-variable-undefined, invalid-name
    global module
    module = AnsibleModule(  # noqa: F405
        argument_spec=dict(
            role=dict(default='common', required=False,
                      choices=OpenShiftFacts.known_roles),
            local_facts=dict(default=None, type='dict', required=False),
            additive_facts_to_overwrite=dict(default=[], type='list', required=False),
        ),
        supports_check_mode=True,
        add_file_common_args=True,
    )

    module.params['gather_subset'] = ['hardware', 'network', 'virtual', 'facter']  # noqa: F405
    module.params['gather_timeout'] = 10  # noqa: F405
    module.params['filter'] = '*'  # noqa: F405

    role = module.params['role']  # noqa: F405
    local_facts = module.params['local_facts']  # noqa: F405
    additive_facts_to_overwrite = module.params['additive_facts_to_overwrite']  # noqa: F405

    fact_file = '/etc/ansible/facts.d/openshift.fact'

    openshift_facts = OpenShiftFacts(role,
                                     fact_file,
                                     local_facts,
                                     additive_facts_to_overwrite)

    file_params = module.params.copy()  # noqa: F405
    file_params['path'] = fact_file
    file_args = module.load_file_common_arguments(file_params)  # noqa: F405
    changed = module.set_fs_attributes_if_different(file_args,  # noqa: F405
                                                    openshift_facts.changed)

    return module.exit_json(changed=changed,  # noqa: F405
                            ansible_facts=openshift_facts.facts)


if __name__ == '__main__':
    main()
