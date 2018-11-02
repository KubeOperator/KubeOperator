"""
Module for performing checks on a Fluentd logging deployment configuration
"""

from openshift_checks import OpenShiftCheckException
from openshift_checks.logging.logging import LoggingCheck


class FluentdConfig(LoggingCheck):
    """Module that checks logging configuration of an integrated logging Fluentd deployment"""
    name = "fluentd_config"
    tags = ["health"]

    def is_active(self):
        logging_deployed = self.get_var("openshift_hosted_logging_deploy", default=False)

        try:
            version = self.get_major_minor_version(self.get_var("openshift_image_tag"))
        except ValueError:
            # if failed to parse OpenShift version, perform check anyway (if logging enabled)
            return logging_deployed

        return logging_deployed and version < (3, 6)

    def run(self):
        """Check that Fluentd has running pods, and that its logging config matches Docker's logging config."""
        config_error = self.check_logging_config()
        if config_error:
            msg = ("The following Fluentd logging configuration problem was found:"
                   "\n{}".format(config_error))
            return {"failed": True, "msg": msg}

        return {}

    def check_logging_config(self):
        """Ensure that the configured Docker logging driver matches fluentd settings.
        This means that, at least for now, if the following condition is met:

            openshift_logging_fluentd_use_journal == True

        then the value of the configured Docker logging driver should be "journald".
        Otherwise, the value of the Docker logging driver should be "json-file".
        Returns an error string if the above condition is not met, or None otherwise."""
        use_journald = self.get_var("openshift_logging_fluentd_use_journal", default=True)

        # if check is running on a master, retrieve all running pods
        # and check any pod's container for the env var "USE_JOURNAL"
        group_names = self.get_var("group_names")
        if "oo_masters_to_config" in group_names:
            use_journald = self.check_fluentd_env_var()

        docker_info = self.execute_module("docker_info", {})
        try:
            logging_driver = docker_info["info"]["LoggingDriver"]
        except KeyError:
            return "Unable to determine Docker logging driver."

        logging_driver = docker_info["info"]["LoggingDriver"]
        recommended_logging_driver = "journald"
        error = None

        # If fluentd is set to use journald but Docker is not, recommend setting the `--log-driver`
        # option as an inventory file variable, or adding the log driver value as part of the
        # Docker configuration in /etc/docker/daemon.json. There is no global --log-driver flag that
        # can be passed to the Docker binary; the only other recommendation that can be made, would be
        # to pass the `--log-driver` flag to the "run" sub-command of the `docker` binary when running
        # individual containers.
        if use_journald and logging_driver != "journald":
            error = ('Your Fluentd configuration is set to aggregate Docker container logs from "journald".\n'
                     'This differs from your Docker configuration, which has been set to use "{driver}" '
                     'as the default method of storing logs.\n'
                     'This discrepancy in configuration will prevent Fluentd from receiving any logs'
                     'from your Docker containers.').format(driver=logging_driver)
        elif not use_journald and logging_driver != "json-file":
            recommended_logging_driver = "json-file"
            error = ('Your Fluentd configuration is set to aggregate Docker container logs from '
                     'individual json log files per container.\n '
                     'This differs from your Docker configuration, which has been set to use '
                     '"{driver}" as the default method of storing logs.\n'
                     'This discrepancy in configuration will prevent Fluentd from receiving any logs'
                     'from your Docker containers.').format(driver=logging_driver)

        if error:
            error += ('\nTo resolve this issue, add the following variable to your Ansible inventory file:\n\n'
                      '  openshift_docker_options="--log-driver={driver}"\n\n'
                      'Alternatively, you can add the following option to your Docker configuration, located in'
                      '"/etc/docker/daemon.json":\n\n'
                      '{{ "log-driver": "{driver}" }}\n\n'
                      'See https://docs.docker.com/engine/admin/logging/json-file '
                      'for more information.').format(driver=recommended_logging_driver)

        return error

    def check_fluentd_env_var(self):
        """Read and return the value of the 'USE_JOURNAL' environment variable on a fluentd pod."""
        running_pods = self.running_fluentd_pods()

        try:
            pod_containers = running_pods[0]["spec"]["containers"]
        except KeyError:
            return "Unable to detect running containers on selected Fluentd pod."

        if not pod_containers:
            msg = ('There are no running containers on selected Fluentd pod "{}".\n'
                   'Unable to calculate expected logging driver.').format(running_pods[0]["metadata"].get("name", ""))
            raise OpenShiftCheckException(msg)

        pod_env = pod_containers[0].get("env")
        if not pod_env:
            msg = ('There are no environment variables set on the Fluentd container "{}".\n'
                   'Unable to calculate expected logging driver.').format(pod_containers[0].get("name"))
            raise OpenShiftCheckException(msg)

        for env in pod_env:
            if env["name"] == "USE_JOURNAL":
                return env.get("value", "false") != "false"

        return False

    def running_fluentd_pods(self):
        """Return a list of running fluentd pods."""
        fluentd_pods = self.get_pods_for_component("fluentd")

        running_fluentd_pods = [pod for pod in fluentd_pods if pod['status']['phase'] == 'Running']
        if not running_fluentd_pods:
            raise OpenShiftCheckException(
                'No Fluentd pods were found to be in the "Running" state. '
                'At least one Fluentd pod is required in order to perform this check.'
            )

        return running_fluentd_pods
