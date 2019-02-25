"""
Util functions for performing checks on an Elasticsearch, Fluentd, and Kibana stack
"""

import json
import os

from openshift_checks import OpenShiftCheck, OpenShiftCheckException


class MissingComponentPods(OpenShiftCheckException):
    """Raised when a component has no pods in the namespace."""
    pass


class CouldNotUseOc(OpenShiftCheckException):
    """Raised when ocutil has a failure running oc."""
    pass


class LoggingCheck(OpenShiftCheck):
    """Base class for OpenShift aggregated logging component checks"""

    # FIXME: this should not be listed as a check, since it is not meant to be
    # run by itself.

    name = "logging"

    def is_active(self):
        logging_deployed = self.get_var("openshift_logging_install_logging", convert=bool, default=False)
        return logging_deployed and super(LoggingCheck, self).is_active() and self.is_first_master()

    def run(self):
        return {}

    def get_pods_for_component(self, logging_component):
        """Get all pods for a given component. Returns: list of pods."""
        pod_output = self.exec_oc(
            "get pods -l component={} -o json".format(logging_component),
            [],
        )
        try:
            pods = json.loads(pod_output)  # raises ValueError if deserialize fails
            if not pods or not pods.get('items'):  # also a broken response, treat the same
                raise ValueError()
        except ValueError:
            # successful run but non-parsing data generally means there were no pods to be found
            raise MissingComponentPods(
                'There are no "{}" component pods in the "{}" namespace.\n'
                'Is logging deployed?'.format(logging_component, self.logging_namespace())
            )

        return pods['items']

    @staticmethod
    def not_running_pods(pods):
        """Returns: list of pods not in a ready and running state"""
        return [
            pod for pod in pods
            if not pod.get("status", {}).get("containerStatuses") or any(
                container['ready'] is False
                for container in pod['status']['containerStatuses']
            ) or not any(
                condition['type'] == 'Ready' and condition['status'] == 'True'
                for condition in pod['status'].get('conditions', [])
            )
        ]

    def logging_namespace(self):
        """Returns the namespace in which logging is configured to deploy."""
        return self.get_var("openshift_logging_namespace", default="openshift-logging")

    def exec_oc(self, cmd_str="", extra_args=None, save_as_name=None):
        """
        Execute an 'oc' command in the remote host.
        Returns: output of command and namespace,
        or raises CouldNotUseOc on error
        """
        config_base = self.get_var("openshift", "common", "config_base")
        args = {
            "namespace": self.logging_namespace(),
            "config_file": os.path.join(config_base, "master", "admin.kubeconfig"),
            "cmd": cmd_str,
            "extra_args": list(extra_args) if extra_args else [],
        }

        result = self.execute_module("ocutil", args, save_as_name=save_as_name)
        if result.get("failed"):
            if result['result'] == '[Errno 2] No such file or directory':
                raise CouldNotUseOc(
                    "This host is supposed to be a master but does not have the `oc` command where expected.\n"
                    "Has an installation been run on this host yet?"
                )

            raise CouldNotUseOc(
                'Unexpected error using `oc` to validate the logging stack components.\n'
                'Error executing `oc {cmd}`:\n'
                '{error}'.format(cmd=args['cmd'], error=result['result'])
            )

        return result.get("result", "")
