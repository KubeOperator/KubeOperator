"""Check for an aggregated logging Fluentd deployment"""

import json


from openshift_checks import OpenShiftCheckException, OpenShiftCheckExceptionList
from openshift_checks.logging.logging import LoggingCheck


class Fluentd(LoggingCheck):
    """Check for an aggregated logging Fluentd deployment"""

    name = "fluentd"
    tags = ["health", "logging"]

    def run(self):
        """Check the Fluentd deployment and raise an error if any problems are found."""

        fluentd_pods = self.get_pods_for_component("fluentd")
        self.check_fluentd(fluentd_pods)
        return {}

    def check_fluentd(self, pods):
        """Verify fluentd is running everywhere. Raises OpenShiftCheckExceptionList if error(s) found."""

        node_selector = self.get_var(
            'openshift_logging_fluentd_nodeselector',
            default='logging-infra-fluentd=true'
        )

        nodes_by_name = self.get_nodes_by_name()
        fluentd_nodes = self.filter_fluentd_labeled_nodes(nodes_by_name, node_selector)

        errors = []
        errors += self.check_node_labeling(nodes_by_name, fluentd_nodes, node_selector)
        errors += self.check_nodes_have_fluentd(pods, fluentd_nodes)
        errors += self.check_fluentd_pods_running(pods)

        # Make sure there are no extra fluentd pods
        if len(pods) > len(fluentd_nodes):
            errors.append(OpenShiftCheckException(
                'TooManyFluentdPods',
                'There are more Fluentd pods running than nodes labeled.\n'
                'This may not cause problems with logging but it likely indicates something wrong.'
            ))

        if errors:
            raise OpenShiftCheckExceptionList(errors)

    def get_nodes_by_name(self):
        """Retrieve all the node definitions. Returns: dict(name: node)"""
        nodes_json = self.exec_oc("get nodes -o json", [])
        try:
            nodes = json.loads(nodes_json)
        except ValueError:  # no valid json - should not happen
            raise OpenShiftCheckException(
                "BadOcNodeList",
                "Could not obtain a list of nodes to validate fluentd.\n"
                "Output from oc get:\n" + nodes_json
            )
        if not nodes or not nodes.get('items'):  # also should not happen
            raise OpenShiftCheckException(
                "NoNodesDefined",
                "No nodes appear to be defined according to the API."
            )
        return {
            node['metadata']['name']: node
            for node in nodes['items']
        }

    @staticmethod
    def filter_fluentd_labeled_nodes(nodes_by_name, node_selector):
        """Filter to all nodes with fluentd label. Returns dict(name: node)"""
        label, value = node_selector.split('=', 1)
        fluentd_nodes = {
            name: node for name, node in nodes_by_name.items()
            if node['metadata']['labels'].get(label) == value
        }
        if not fluentd_nodes:
            raise OpenShiftCheckException(
                'NoNodesLabeled',
                'There are no nodes with the fluentd label {label}.\n'
                'This means no logs will be aggregated from the nodes.'.format(label=node_selector)
            )
        return fluentd_nodes

    def check_node_labeling(self, nodes_by_name, fluentd_nodes, node_selector):
        """Note if nodes are not labeled as expected. Returns: error list"""
        intended_nodes = self.get_var('openshift_logging_fluentd_hosts', default=['--all'])
        if not intended_nodes or '--all' in intended_nodes:
            intended_nodes = nodes_by_name.keys()
        nodes_missing_labels = set(intended_nodes) - set(fluentd_nodes.keys())
        if nodes_missing_labels:
            return [OpenShiftCheckException(
                'NodesUnlabeled',
                'The following nodes are supposed to be labeled with {label} but are not:\n'
                '  {nodes}\n'
                'Fluentd will not aggregate logs from these nodes.'.format(
                    label=node_selector, nodes=', '.join(nodes_missing_labels)
                ))]

        return []

    @staticmethod
    def check_nodes_have_fluentd(pods, fluentd_nodes):
        """Make sure fluentd is on all the labeled nodes. Returns: error list"""
        unmatched_nodes = fluentd_nodes.copy()
        node_names_by_label = {
            node['metadata']['labels']['kubernetes.io/hostname']: name
            for name, node in fluentd_nodes.items()
        }
        node_names_by_internal_ip = {
            address['address']: name
            for name, node in fluentd_nodes.items()
            for address in node['status']['addresses']
            if address['type'] == "InternalIP"
        }
        for pod in pods:
            for name in [
                    pod['spec']['nodeName'],
                    node_names_by_internal_ip.get(pod['spec']['nodeName']),
                    node_names_by_label.get(pod.get('spec', {}).get('host')),
            ]:
                unmatched_nodes.pop(name, None)
        if unmatched_nodes:
            return [OpenShiftCheckException(
                'MissingFluentdPod',
                'The following nodes are supposed to have a Fluentd pod but do not:\n'
                '  {nodes}\n'
                'These nodes will not have their logs aggregated.'.format(
                    nodes='\n  '.join(unmatched_nodes.keys())
                ))]

        return []

    def check_fluentd_pods_running(self, pods):
        """Make sure all fluentd pods are running. Returns: error string"""
        not_running = super(Fluentd, self).not_running_pods(pods)
        if not_running:
            return [OpenShiftCheckException(
                'FluentdNotRunning',
                'The following Fluentd pods are supposed to be running but are not:\n'
                '  {pods}\n'
                'These pods will not aggregate logs from their nodes.'.format(
                    pods='\n'.join(
                        "  {name} ({host})".format(
                            name=pod['metadata']['name'],
                            host=pod['spec'].get('host', 'None')
                        )
                        for pod in not_running
                    )
                ))]

        return []
