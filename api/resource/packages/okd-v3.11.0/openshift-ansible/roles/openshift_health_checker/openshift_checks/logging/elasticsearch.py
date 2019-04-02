"""Check for an aggregated logging Elasticsearch deployment"""

import json
import re

from openshift_checks import OpenShiftCheckException, OpenShiftCheckExceptionList
from openshift_checks.logging.logging import LoggingCheck


class Elasticsearch(LoggingCheck):
    """Check for an aggregated logging Elasticsearch deployment"""

    name = "elasticsearch"
    tags = ["health", "logging"]

    def run(self):
        """Check various things and gather errors. Returns: result as hash"""

        es_pods = self.get_pods_for_component("es")
        self.check_elasticsearch(es_pods)
        # TODO(lmeyer): run it all again for the ops cluster

        return {}

    def check_elasticsearch(self, es_pods):
        """Perform checks for Elasticsearch. Raises OpenShiftCheckExceptionList on any errors."""
        running_pods, errors = self.running_elasticsearch_pods(es_pods)
        pods_by_name = {
            pod['metadata']['name']: pod for pod in running_pods
            # Filter out pods that are not members of a DC
            if pod['metadata'].get('labels', {}).get('deploymentconfig')
        }
        if not pods_by_name:
            # nothing running, cannot run the rest of the check
            errors.append(OpenShiftCheckException(
                'NoRunningPods',
                'No logging Elasticsearch pods were found running, so no logs are being aggregated.'
            ))
            raise OpenShiftCheckExceptionList(errors)

        errors += self.check_elasticsearch_masters(pods_by_name)
        errors += self.check_elasticsearch_node_list(pods_by_name)
        errors += self.check_es_cluster_health(pods_by_name)
        errors += self.check_elasticsearch_diskspace(pods_by_name)
        if errors:
            raise OpenShiftCheckExceptionList(errors)

    def running_elasticsearch_pods(self, es_pods):
        """Returns: list of running pods, list of errors about non-running pods"""
        not_running = self.not_running_pods(es_pods)
        running_pods = [pod for pod in es_pods if pod not in not_running]
        if not_running:
            return running_pods, [OpenShiftCheckException(
                'PodNotRunning',
                'The following Elasticsearch pods are defined but not running:\n'
                '{pods}'.format(pods=''.join(
                    "  {} ({})\n".format(pod['metadata']['name'], pod['spec'].get('host', 'None'))
                    for pod in not_running
                ))
            )]
        return running_pods, []

    @staticmethod
    def _build_es_curl_cmd(pod_name, url):
        base = "exec {name} -- curl -s --cert {base}cert --key {base}key --cacert {base}ca -XGET '{url}'"
        return base.format(base="/etc/elasticsearch/secret/admin-", name=pod_name, url=url)

    def check_elasticsearch_masters(self, pods_by_name):
        """Check that Elasticsearch masters are sane. Returns: list of errors"""
        es_master_names = set()
        errors = []
        for pod_name in pods_by_name.keys():
            # Compare what each ES node reports as master and compare for split brain
            get_master_cmd = self._build_es_curl_cmd(pod_name, "https://localhost:9200/_cat/master")
            master_name_str = self.exec_oc(get_master_cmd, [], save_as_name="get_master_names.json")
            master_names = (master_name_str or '').split(' ')
            if len(master_names) > 1:
                es_master_names.add(master_names[1])
            else:
                errors.append(OpenShiftCheckException(
                    'NoMasterName',
                    'Elasticsearch {pod} gave unexpected response when asked master name:\n'
                    '  {response}'.format(pod=pod_name, response=master_name_str)
                ))

        if not es_master_names:
            errors.append(OpenShiftCheckException(
                'NoMasterFound',
                'No logging Elasticsearch masters were found.'
            ))
            return errors

        if len(es_master_names) > 1:
            errors.append(OpenShiftCheckException(
                'SplitBrainMasters',
                'Found multiple Elasticsearch masters according to the pods:\n'
                '{master_list}\n'
                'This implies that the masters have "split brain" and are not correctly\n'
                'replicating data for the logging cluster. Log loss is likely to occur.'
                .format(master_list='\n'.join('  ' + master for master in es_master_names))
            ))

        return errors

    def check_elasticsearch_node_list(self, pods_by_name):
        """Check that reported ES masters are accounted for by pods. Returns: list of errors"""

        if not pods_by_name:
            return [OpenShiftCheckException(
                'MissingComponentPods',
                'No logging Elasticsearch pods were found.'
            )]

        # get ES cluster nodes
        node_cmd = self._build_es_curl_cmd(list(pods_by_name.keys())[0], 'https://localhost:9200/_nodes')
        cluster_node_data = self.exec_oc(node_cmd, [], save_as_name="get_es_nodes.json")
        try:
            cluster_nodes = json.loads(cluster_node_data)['nodes']
        except (ValueError, KeyError):
            return [OpenShiftCheckException(
                'MissingNodeList',
                'Failed to query Elasticsearch for the list of ES nodes. The output was:\n' +
                cluster_node_data
            )]

        # Try to match all ES-reported node hosts to known pods.
        errors = []
        for node in cluster_nodes.values():
            # Note that with 1.4/3.4 the pod IP may be used as the master name
            if not any(node['host'] in (pod_name, pod['status'].get('podIP'))
                       for pod_name, pod in pods_by_name.items()):
                errors.append(OpenShiftCheckException(
                    'EsPodNodeMismatch',
                    'The Elasticsearch cluster reports a member node "{node}"\n'
                    'that does not correspond to any known ES pod.'.format(node=node['host'])
                ))

        return errors

    def check_es_cluster_health(self, pods_by_name):
        """Exec into the elasticsearch pods and check the cluster health. Returns: list of errors"""
        errors = []
        for pod_name in pods_by_name.keys():
            cluster_health_cmd = self._build_es_curl_cmd(pod_name, 'https://localhost:9200/_cluster/health?pretty=true')
            cluster_health_data = self.exec_oc(cluster_health_cmd, [], save_as_name='get_es_health.json')
            try:
                health_res = json.loads(cluster_health_data)
                if not health_res or not health_res.get('status'):
                    raise ValueError()
            except ValueError:
                errors.append(OpenShiftCheckException(
                    'BadEsResponse',
                    'Could not retrieve cluster health status from logging ES pod "{pod}".\n'
                    'Response was:\n{output}'.format(pod=pod_name, output=cluster_health_data)
                ))
                continue

            if health_res['status'] not in ['green', 'yellow']:
                errors.append(OpenShiftCheckException(
                    'EsClusterHealthRed',
                    'Elasticsearch cluster health status is RED according to pod "{}"'.format(pod_name)
                ))

        return errors

    def check_elasticsearch_diskspace(self, pods_by_name):
        """
        Exec into an ES pod and query the diskspace on the persistent volume.
        Returns: list of errors
        """
        errors = []
        for pod_name in pods_by_name.keys():
            df_cmd = '-c elasticsearch exec {} -- df --output=ipcent,pcent /elasticsearch/persistent'.format(pod_name)
            disk_output = self.exec_oc(df_cmd, [], save_as_name='get_pv_diskspace.json')
            lines = disk_output.splitlines()
            # expecting one header looking like 'IUse% Use%' and one body line
            body_re = r'\s*(\d+)%?\s+(\d+)%?\s*$'
            if len(lines) != 2 or len(lines[0].split()) != 2 or not re.match(body_re, lines[1]):
                errors.append(OpenShiftCheckException(
                    'BadDfResponse',
                    'Could not retrieve storage usage from logging ES pod "{pod}".\n'
                    'Response to `df` command was:\n{output}'.format(pod=pod_name, output=disk_output)
                ))
                continue
            inode_pct, disk_pct = re.match(body_re, lines[1]).groups()

            inode_pct_thresh = self.get_var('openshift_check_efk_es_inode_pct', default='90')
            if int(inode_pct) >= int(inode_pct_thresh):
                errors.append(OpenShiftCheckException(
                    'InodeUsageTooHigh',
                    'Inode percent usage on the storage volume for logging ES pod "{pod}"\n'
                    '  is {pct}, greater than threshold {limit}.\n'
                    '  Note: threshold can be specified in inventory with {param}'.format(
                        pod=pod_name,
                        pct=str(inode_pct),
                        limit=str(inode_pct_thresh),
                        param='openshift_check_efk_es_inode_pct',
                    )))
            disk_pct_thresh = self.get_var('openshift_check_efk_es_storage_pct', default='80')
            if int(disk_pct) >= int(disk_pct_thresh):
                errors.append(OpenShiftCheckException(
                    'DiskUsageTooHigh',
                    'Disk percent usage on the storage volume for logging ES pod "{pod}"\n'
                    '  is {pct}, greater than threshold {limit}.\n'
                    '  Note: threshold can be specified in inventory with {param}'.format(
                        pod=pod_name,
                        pct=str(disk_pct),
                        limit=str(disk_pct_thresh),
                        param='openshift_check_efk_es_storage_pct',
                    )))

        return errors
