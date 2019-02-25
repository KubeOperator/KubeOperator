"""Check for an aggregated logging Curator deployment"""

from openshift_checks.logging.logging import OpenShiftCheckException, LoggingCheck


class Curator(LoggingCheck):
    """Check for an aggregated logging Curator deployment"""

    name = "curator"
    tags = ["health", "logging"]

    def run(self):
        """Check various things and gather errors. Returns: result as hash"""

        curator_pods = self.get_pods_for_component("curator")
        self.check_curator(curator_pods)
        # TODO(lmeyer): run it all again for the ops cluster

        return {}

    def check_curator(self, pods):
        """Check to see if curator is up and working. Returns: error string"""
        if not pods:
            raise OpenShiftCheckException(
                "MissingComponentPods",
                "There are no Curator pods for the logging stack,\n"
                "so nothing will prune Elasticsearch indexes.\n"
                "Is Curator correctly deployed?"
            )

        not_running = self.not_running_pods(pods)
        if len(not_running) == len(pods):
            raise OpenShiftCheckException(
                "CuratorNotRunning",
                "The Curator pod is not currently in a running state,\n"
                "so Elasticsearch indexes may increase without bound."
            )
        if len(pods) - len(not_running) > 1:
            raise OpenShiftCheckException(
                "TooManyCurators",
                "There is more than one Curator pod running. This should not normally happen.\n"
                "Although this doesn't cause any problems, you may want to investigate."
            )
