"""
Check for ensuring logs from pods can be queried in a reasonable amount of time.
"""

import json
import time

from uuid import uuid4

from openshift_checks import OpenShiftCheckException
from openshift_checks.logging.logging import LoggingCheck


ES_CMD_TIMEOUT_SECONDS = 30


class LoggingIndexTime(LoggingCheck):
    """Check that pod logs are aggregated and indexed in ElasticSearch within a reasonable amount of time."""
    name = "logging_index_time"
    tags = ["health", "logging"]

    def run(self):
        """Add log entry by making unique request to Kibana. Check for unique entry in the ElasticSearch pod logs."""
        try:
            log_index_timeout = int(
                self.get_var("openshift_check_logging_index_timeout_seconds", default=ES_CMD_TIMEOUT_SECONDS)
            )
        except ValueError:
            raise OpenShiftCheckException(
                'InvalidTimeout',
                'Invalid value provided for "openshift_check_logging_index_timeout_seconds". '
                'Value must be an integer representing an amount in seconds.'
            )

        running_component_pods = dict()

        # get all component pods
        for component, name in (['kibana', 'Kibana'], ['es', 'Elasticsearch']):
            pods = self.get_pods_for_component(component)
            running_pods = self.running_pods(pods)

            if not running_pods:
                raise OpenShiftCheckException(
                    component + 'NoRunningPods',
                    'No {} pods in the "Running" state were found.'
                    'At least one pod is required in order to perform this check.'.format(name)
                )

            running_component_pods[component] = running_pods

        uuid = self.curl_kibana_with_uuid(running_component_pods["kibana"][0])
        self.wait_until_cmd_or_err(running_component_pods["es"][0], uuid, log_index_timeout)
        return {}

    def wait_until_cmd_or_err(self, es_pod, uuid, timeout_secs):
        """Retry an Elasticsearch query every second until query success, or a defined
        length of time has passed."""
        deadline = time.time() + timeout_secs
        interval = 1
        while not self.query_es_from_es(es_pod, uuid):
            if time.time() + interval > deadline:
                raise OpenShiftCheckException(
                    "NoMatchFound",
                    "expecting match in Elasticsearch for message with uuid {}, "
                    "but no matches were found after {}s.".format(uuid, timeout_secs)
                )
            time.sleep(interval)

    def curl_kibana_with_uuid(self, kibana_pod):
        """curl Kibana with a unique uuid."""
        uuid = self.generate_uuid()
        pod_name = kibana_pod["metadata"]["name"]
        exec_cmd = "exec {pod_name} -c kibana -- curl --max-time 30 -s http://localhost:5601/{uuid}"
        exec_cmd = exec_cmd.format(pod_name=pod_name, uuid=uuid)

        error_str = self.exec_oc(exec_cmd, [])

        try:
            error_code = json.loads(error_str)["statusCode"]
        except (KeyError, ValueError):
            raise OpenShiftCheckException(
                'kibanaInvalidResponse',
                'invalid response returned from Kibana request:\n'
                'Command: {}\nResponse: {}'.format(exec_cmd, error_str)
            )

        if error_code != 404:
            raise OpenShiftCheckException(
                'kibanaInvalidReturnCode',
                'invalid error code returned from Kibana request.\n'
                'Expecting error code "404", but got "{}" instead.'.format(error_code)
            )

        return uuid

    def query_es_from_es(self, es_pod, uuid):
        """curl the Elasticsearch pod and look for a unique uuid in its logs."""
        pod_name = es_pod["metadata"]["name"]
        exec_cmd = (
            "exec {pod_name} -- curl --max-time 30 -s -f "
            "--cacert /etc/elasticsearch/secret/admin-ca "
            "--cert /etc/elasticsearch/secret/admin-cert "
            "--key /etc/elasticsearch/secret/admin-key "
            "https://logging-es:9200/project.{namespace}*/_count?q=message:{uuid}"
        )
        exec_cmd = exec_cmd.format(pod_name=pod_name, namespace=self.logging_namespace(), uuid=uuid)
        result = self.exec_oc(exec_cmd, [], save_as_name="query_for_uuid.json")

        try:
            count = json.loads(result)["count"]
        except (KeyError, ValueError):
            raise OpenShiftCheckException(
                'esInvalidResponse',
                'Invalid response from Elasticsearch query:\n'
                '  {}\n'
                'Response was:\n{}'.format(exec_cmd, result)
            )

        return count

    @staticmethod
    def running_pods(pods):
        """Filter pods that are running."""
        return [pod for pod in pods if pod['status']['phase'] == 'Running']

    @staticmethod
    def generate_uuid():
        """Wrap uuid generator. Allows for testing with expected values."""
        return str(uuid4())
