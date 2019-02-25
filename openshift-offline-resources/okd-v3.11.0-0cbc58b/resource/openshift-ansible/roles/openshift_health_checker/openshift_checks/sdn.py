"""
Check that the SDN is routing traffic properly.
"""

import datetime
import os
import textwrap
import time

import ipaddress
from ansible.module_utils import six

from openshift_checks import OpenShiftCheck, OpenShiftCheckException


class SDNCheck(OpenShiftCheck):
    """A check to run relevant diagnostics on the SDN."""

    name = 'sdn'
    tags = ['health']

    def is_active(self):
        """Skip hosts that are not masters or nodes."""
        group_names = self.get_var('group_names', default=[])
        master_or_node = 'oo_masters_to_config' in group_names or \
                         'oo_nodes_to_config' in group_names
        return super(SDNCheck, self).is_active() and master_or_node

    def run(self):
        if self.want_full_results:
            # Gather diagnostic information and perform diagnostics on a master
            # or node host.
            try:
                self.save_journal()
                self.save_command_output('nmcli-dev',
                                         ['/bin/nmcli', '--nocheck', '-f', 'all',
                                          'dev', 'show'])
                self.save_command_output('nmcli-con',
                                         ['/bin/nmcli', '--nocheck', '-f', 'all',
                                          'con', 'show'])
                self.save_command_output(
                    'ifcfg',
                    'head -1000 /etc/sysconfig/network-scripts/ifcfg-*')
                self.save_command_output('addresses',
                                         ['/sbin/ip', 'addr', 'show'])
                self.save_command_output('routes',
                                         ['/sbin/ip', 'route', 'show'])
                self.save_command_output('arp',
                                         ['/sbin/ip', '-s', 'neighbor', 'show'])
                self.save_command_output('iptables', ['/sbin/iptables-save'])
                self.register_file('hosts', None, '/etc/hosts')
                self.register_file('resolv.conf', None, '/etc/resolv.conf')
                self.save_command_output('modules', ['/sbin/lsmod'])
                self.save_command_output('sysctl', ['/sbin/sysctl', '-a'])
                if self.get_var('openshift_use_crio', default=False):
                    self.save_command_output('crio-version',
                                             ['/bin/crictl', 'version'])
                if not self.get_var('openshift_use_crio_only', default=False):
                    self.save_command_output('docker-version',
                                             ['/bin/docker', 'version'])
                oc_executable = self.get_var('openshift_client_binary',
                                             default='/bin/oc')
                oc_executable = self.template_var(oc_executable)
                # The oc executable is not installed on containerized nodes, so
                # use "2>&1" to capture any error output (such as "command not
                # found"), and use "|| :" to ignore the exit code.
                self.save_command_output('oc-version',
                                         oc_executable + ' version 2>&1 || :')
                self.register_file('os-version', None,
                                   '/etc/system-release-cpe')
            except OpenShiftCheckException as exc:
                self.register_failure(exc)

        group_names = self.get_var('group_names', default=[])
        if 'oo_masters_to_config' in group_names:
            self.check_master()
        if 'oo_nodes_to_config' in group_names:
            self.check_node()
        return {}

    def save_journal(self):
        """Save the last 5 minutes of the journal."""
        out = self.read_command_output(['/bin/journalctl', '-n', '1', '-q'])
        (since, until) = SDNCheck.compute_log_interval_from(out)
        self.register_file('journal',
                           self.read_command_output(['/bin/journalctl',
                                                     '"--since=%s"' % since,
                                                     '"--until=%s"' % until]))

    @staticmethod
    def compute_log_interval_from(log_line):
        """Compute and return a 2-tuple of timestamps (ts1, ts2) where ts1
        represents the date 5 minutes prior to the timestamp of the provided log
        message and ts2 represents the date of that timestamp.  The log line is
        assumed to be from today."""
        try:
            log_ts = log_line.strip().split()[2]
            ts2_time = datetime.datetime.strptime(log_ts, '%H:%M:%S').time()
            now = datetime.datetime.now()
            ts2_date = datetime.datetime.combine(now, ts2_time)
            ts1_date = ts2_date - datetime.timedelta(minutes=5)
            time_fmt = '%Y-%m-%d %H:%M:%S'
            # pylint may infer that ts1_date is NotImplemented or a timedelta
            # object and complain about using the timetuple method from the
            # datetime class because the subtraction operation above would
            # return a timedelta if the RHS were a date or NotImplemented if the
            # RHS were neither a datetime nor a timedelta.  However,
            # datetime.datetime.combine cannot return a non-datetime value, and
            # datetime.timedelta(minutes=5) is an explicit timedelta value,
            # so the subtraction really can only return a datetime value.
            # pylint: disable=no-member
            ts1 = time.strftime(time_fmt, ts1_date.timetuple())
            ts2 = time.strftime(time_fmt, ts2_date.timetuple())
        except (ValueError, IndexError):
            ts1 = '-5m'
            ts2 = 'now'

        return (ts1, ts2)

    def save_command_output(self, path, command):
        """Execute the provided command using the command module
        and save its output to the specified file.

        If the command is a string, use a shell.  Otherwise, assume the command
        is a list, join it with spaces, and execute it without shell.
        """
        self.register_file(path, self.read_command_output(command))

    def read_command_output(self, command, utf8=True):
        """Execute the provided command using the command module
        and return its output.

        If the command is a string, use a shell.  Otherwise, assume the command
        is a list, join it with spaces, and execute it without shell.
        """
        uses_shell = False
        if isinstance(command, six.string_types):
            uses_shell = True
        else:
            command = ' '.join(command)

        command_args = dict(_raw_params=command, _uses_shell=uses_shell)
        # Use self._execute_module instead of self.execute_module because
        # the latter sets self.changed.
        result = self._execute_module('command', command_args)
        if result.get('rc', 0) != 0 or result.get('failed'):
            raise OpenShiftCheckException(
                'RemoteCommandFailure',
                'Failed to execute command on remote host: %s' % command)

        if utf8:
            return result['stdout'].encode('utf-8')
        return result['stdout']

    def check_master(self):
        """Gather diagnostic information on a master and ensure it can connect
        to kubelets."""
        if self.want_full_results:
            conf_base_path = self.get_var('openshift.common.config_base')
            master_conf_path = os.path.join(conf_base_path, 'master',
                                            'master-config.yaml')
            self.register_file('master-config.yaml', None, master_conf_path)

            self.save_component_container_logs('controllers', 'controllers')
            self.save_component_container_logs('api', 'api')

        nodes = self.get_resource('nodes')

        if self.want_full_results:
            self.register_file('nodes.json', nodes)
            self.register_file('pods.json', self.get_resource('pods'))
            self.register_file('services.json', self.get_resource('services'))
            self.register_file('endpoints.json', self.get_resource('endpoints'))
            self.register_file('routes.json', self.get_resource('routes'))
            self.register_file('clusternetworks.json',
                               self.get_resource('clusternetworks'))
            self.register_file('hostsubnets.json',
                               self.get_resource('hostsubnets'))
            self.register_file('netnamespaces.json',
                               self.get_resource('netnamespaces'))

        if not nodes:
            self.register_failure(
                'No nodes appear to be defined according to the API.'
            )

        for node in nodes:
            self.check_node_kubelet(node)

    def save_component_container_logs(self, component, container):
        """Save the first and last 2000 lines of logs for the specified
        component and container."""
        awk_script = textwrap.dedent('''\
            BEGIN {
                n = 2000
            }
            NR <= n {
                print
            }
            NR > n {
                buf[(NR - 1)%n + 1] = $0
            }
            END {
                if (NR <= n)
                    exit

                if (NR > 2*n)
                    print "..."

                for (i = NR >= 2*n ? 0 : n - NR%n; i < n; ++i)
                    print buf[(NR + i)%n + 1]
            }''')
        out = self.read_command_output(' '.join(['/usr/local/bin/master-logs',
                                                 component, container, '2>&1',
                                                 '|', '/bin/awk',
                                                 "'%s'" % awk_script]))
        self.register_file('master-logs_%s_%s' % (component, container), out)

    def get_resource(self, kind):
        """Return a list of all resources of the specified kind."""
        for resource in self.task_vars['resources']['results']:
            if resource['item'] == kind:
                return resource['results']['results'][0]['items']

        raise OpenShiftCheckException('CouldNotListResource',
                                      'Could not list resource %s' % kind)

    def check_node(self):
        """Gather diagnostic information on a node and perform connectivity
        checks on pods and services."""
        node_name = self.get_var('openshift', 'node', 'nodename', default=None)
        if not node_name:
            self.register_failure('Could not determine node name.')
            return

        # The "openvswitch" container uses the host netnamespace, but the host
        # file system may not have the ovs-appctl and ovs-ofctl binaries, which
        # we use for some diagnostics.  Thus we run these binaries inside the
        # container, and to that end, we need to determine its container id.
        exec_in_ovs_container = self.get_container_exec_command('openvswitch',
                                                                'openshift-sdn')

        if self.want_full_results:
            try:
                service_prefix = self.get_var('openshift_service_type')
                if self._templar is not None:
                    service_prefix = self._templar.template(service_prefix)
                self.save_service_logs('%s-node' % service_prefix)

                if self.get_var('openshift_use_crio', default=False):
                    self.save_command_output('crio-unit-file',
                                             ['/bin/systemctl',
                                              'cat', 'crio.service'])
                    self.save_command_output('crio-ps', ['/bin/crictl', 'ps'])

                if not self.get_var('openshift_use_crio_only', default=False):
                    self.save_command_output('docker-unit-file',
                                             ['/bin/systemctl',
                                              'cat', 'docker.service'])
                    self.save_command_output('docker-ps', ['/bin/docker', 'ps'])

                self.save_command_output('flows', exec_in_ovs_container +
                                         ['/bin/ovs-ofctl', '-O', 'OpenFlow13',
                                          'dump-flows', 'br0'])
                self.save_command_output('ovs-show', exec_in_ovs_container +
                                         ['/bin/ovs-ofctl', '-O', 'OpenFlow13',
                                          'show', 'br0'])

                self.save_command_output('tc-qdisc',
                                         ['/sbin/tc', 'qdisc', 'show'])
                self.save_command_output('tc-class',
                                         ['/sbin/tc', 'class', 'show'])
                self.save_command_output('tc-filter',
                                         ['/sbin/tc', 'filter', 'show'])
            except OpenShiftCheckException as exc:
                self.register_failure(exc)

        subnets = {hostsubnet['metadata']['name']: hostsubnet['subnet']
                   for hostsubnet in self.get_resource('hostsubnets')}

        subnet = subnets.get(node_name, None)
        if subnet is None:
            self.register_failure('Node %s has no hostsubnet.' % node_name)
            return
        subnet = six.text_type(subnet)
        address = ipaddress.ip_network(subnet)[1]

        for remote_node in self.get_resource('nodes'):
            remote_node_name = remote_node['metadata']['name']
            if remote_node_name == node_name:
                continue

            remote_subnet = subnets.get(remote_node_name, None)
            if remote_subnet is None:
                continue
            remote_subnet = six.text_type(remote_subnet)
            remote_address = ipaddress.ip_network(remote_subnet)[1]

            self.save_command_output(
                'trace_node_%s_to_node_%s' % (node_name, remote_node_name),
                exec_in_ovs_container +
                ['/bin/ovs-appctl', 'ofproto/trace', 'br0',
                 'in_port=2,reg0=0,ip,nw_src=%s,nw_dst=%s' %
                 (address, remote_address)])

            try:
                self.save_command_output('ping_node_%s_to_node_%s' %
                                         (node_name, remote_node_name),
                                         ['/bin/ping', '-c', '1', '-W', '2',
                                          str(remote_address)])
            except OpenShiftCheckException as exc:
                self.register_failure('Node %s cannot ping node %s.' %
                                      (node_name, remote_node_name))

    def get_container_exec_command(self, container_name, namespace):
        """Return an array comprising a command and arguments that can be used
        to execute commands inside the specified container running in a pod in
        the specified namespace."""
        if self.get_var('openshift_use_crio', default=False):
            container_id = self.read_command_output([
                '/bin/crictl', 'ps', '-l', '-a', '-q',
                '--label=io.kubernetes.container.name=%s' % container_name,
                '--label=io.kubernetes.pod.namespace=%s' % namespace
            ])
            command = ['/bin/crictl', 'exec', container_id]
        else:
            container_id = self.read_command_output([
                '/bin/docker', 'ps', '-l', '-a', '-q',
                '--filter=label=io.kubernetes.container.name=%s'
                % container_name,
                '--filter=label=io.kubernetes.pod.namespace=%s' % namespace
            ])
            command = ['/bin/docker', 'exec', container_id]

        return command

    def save_service_logs(self, service_name):
        """Save the first 5 minutes of logs after the specified service started
        and the last 5 minutes of logs for that service."""
        time_fmt = '%Y-%m-%d %H:%M:%S'

        out = self.read_command_output(['systemctl', 'show', service_name,
                                        '-p', 'ExecMainStartTimestamp'])
        start_timestamp = out.strip().split('=', 1)[1]
        if len(start_timestamp) == 0:
            self.register_failure('%s is not started.' % service_name)
            return

        # The timestamp should be in the format "%a %Y-%m-%d %H:%M:%S %Z".
        # However, Python cannot reliably parse timezone names
        # (see <https://bugs.python.org/issue22377>), so we must drop the
        # timezone name before parsing the timestamp.
        start_timestamp = ' '.join(start_timestamp.split()[0:3])

        since_date = datetime.datetime.strptime(start_timestamp,
                                                '%a %Y-%m-%d %H:%M:%S')
        until_date = since_date + datetime.timedelta(minutes=5)
        since = since_date.strftime(time_fmt)
        until = until_date.strftime(time_fmt)
        start_logs = self.read_command_output(['/bin/journalctl',
                                               '"--since=%s"' % since,
                                               '"--until=%s"' % until])

        out = self.read_command_output(['/bin/journalctl', '-u', service_name,
                                        '-n', '1', '-q'])
        (since, until) = SDNCheck.compute_log_interval_from(out)
        last_logs = self.read_command_output(['/bin/journalctl',
                                              '-u', service_name,
                                              '"--since=%s"' % since,
                                              '"--until=%s"' % until])

        self.register_file(service_name, (start_logs + '\n...\n' + last_logs))

    def check_node_kubelet(self, node):
        """Check that the host can find the address of the given node, resolve
        that address, and connect to the node's kubelet."""
        name = node['metadata']['name']

        preferred_addr = SDNCheck.get_node_preferred_address(node)
        if not preferred_addr:
            self.register_failure('Node %s: no preferred address' % name)
            return

        internal_addr = None
        for address in node.get('status', {}).get('addresses', []):
            if address.get('type') == 'InternalIP':
                internal_addr = address.get('address')
                break

        if not internal_addr:
            self.register_failure('Node %s: no IP address in OpenShift' % name)
        else:
            try:
                resolved_addr = self.resolve_address(preferred_addr)
            except OpenShiftCheckException as exc:
                self.register_failure(exc)
            else:
                if resolved_addr != internal_addr:
                    self.register_failure(
                        ('Node %s: the IP address in OpenShift (%s)' +
                         ' does not match DNS/hosts (%s)') %
                        (name, internal_addr, resolved_addr))

        url = 'http://%s:%d' % (preferred_addr, 10250)
        result = self.execute_module('uri', dict(url=url))
        if result.get('rc', 0) != 0 or result.get('failed'):
            self.register_failure(
                'Kubelet on node %s is not responding: %s' %
                (name, result.get('msg', 'unknown error')))

    @staticmethod
    def get_node_preferred_address(node):
        """Return a host name or address for the given node, or None.

        The host name or address is selected from the node's status.addresses
        field in accordance with the preference order used by the OpenShift
        master."""
        preferred_address_types = ['Hostname', 'InternalIP', 'ExternalIP']
        for address_type in preferred_address_types:
            for address in node.get('status', {}).get('addresses', []):
                if address.get('type') == address_type:
                    return address.get('address')

            if address_type == 'Hostname':
                hostname = node.get('metadata', {}) \
                               .get('labels', {}) \
                               .get('kubernetes.io/hostname', "")
                if len(hostname) > 0:
                    return hostname

        return None

    def resolve_address(self, addr):
        """Look up the given IPv4 address using getent."""
        command = ' '.join(['/bin/getent', 'ahostsv4', addr])
        try:
            out = self.read_command_output(command, False)
        except OpenShiftCheckException as exc:
            raise OpenShiftCheckException(
                'NameResolutionError',
                'Cannot resolve node %s: %s' % (addr, exc))

        for line in out.splitlines():
            record = line.split()
            if record[1] == 'STREAM':
                return record[0]

        return None
