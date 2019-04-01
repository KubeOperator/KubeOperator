"""
Module for performing checks on a Kibana logging deployment
"""

import json
import ssl

# pylint can't find the package when its installed in virtualenv
# pylint: disable=import-error,no-name-in-module
from ansible.module_utils.six.moves.urllib import request
# pylint: disable=import-error,no-name-in-module
from ansible.module_utils.six.moves.urllib.error import HTTPError, URLError

from openshift_checks.logging.logging import LoggingCheck, OpenShiftCheckException


class Kibana(LoggingCheck):
    """Module that checks an integrated logging Kibana deployment"""

    name = "kibana"
    tags = ["health", "logging"]

    def run(self):
        """Check various things and gather errors. Returns: result as hash"""

        kibana_pods = self.get_pods_for_component("kibana")
        self.check_kibana(kibana_pods)
        self.check_kibana_route()
        # TODO(lmeyer): run it all again for the ops cluster

        return {}

    def _verify_url_internal(self, url):
        """
        Try to reach a URL from the host.
        Returns: success (bool), reason (for failure)
        """
        args = dict(
            url=url,
            follow_redirects='none',
            validate_certs='no',  # likely to be signed with internal CA
            # TODO(lmeyer): give users option to validate certs
            status_code=302,
        )
        result = self.execute_module('uri', args)
        if result.get('failed'):
            return result['msg']
        return None

    @staticmethod
    def _verify_url_external(url):
        """
        Try to reach a URL from ansible control host.
        Raise an OpenShiftCheckException if anything goes wrong.
        """
        # This actually checks from the ansible control host, which may or may not
        # really be "external" to the cluster.

        # Disable SSL cert validation to work around internally signed certs
        ctx = ssl.create_default_context()
        ctx.check_hostname = False  # or setting CERT_NONE is refused
        ctx.verify_mode = ssl.CERT_NONE

        # Verify that the url is returning a valid response
        try:
            # We only care if the url connects and responds
            return_code = request.urlopen(url, context=ctx).getcode()
        except HTTPError as httperr:
            return httperr.reason
        except URLError as urlerr:
            return str(urlerr)

        # there appears to be no way to prevent urlopen from following redirects
        if return_code != 200:
            return 'Expected success (200) but got return code {}'.format(int(return_code))

        return None

    def check_kibana(self, pods):
        """Check to see if Kibana is up and working. Raises OpenShiftCheckException if not."""

        if not pods:
            raise OpenShiftCheckException(
                "MissingComponentPods",
                "There are no Kibana pods deployed, so no access to the logging UI."
            )

        not_running = self.not_running_pods(pods)
        if len(not_running) == len(pods):
            raise OpenShiftCheckException(
                "NoRunningPods",
                "No Kibana pod is in a running state, so there is no access to the logging UI."
            )
        elif not_running:
            raise OpenShiftCheckException(
                "PodNotRunning",
                "The following Kibana pods are not currently in a running state:\n"
                "  {pods}\n"
                "However at least one is, so service may not be impacted.".format(
                    pods="\n  ".join(pod['metadata']['name'] for pod in not_running)
                )
            )

    def _get_kibana_url(self):
        """
        Get kibana route or report error.
        Returns: url
        """

        # Get logging url
        get_route = self.exec_oc("get route logging-kibana -o json", [])
        if not get_route:
            raise OpenShiftCheckException(
                'no_route_exists',
                'No route is defined for Kibana in the logging namespace,\n'
                'so the logging stack is not accessible. Is logging deployed?\n'
                'Did something remove the logging-kibana route?'
            )

        try:
            route = json.loads(get_route)
            # check that the route has been accepted by a router
            ingress = route["status"]["ingress"]
        except (ValueError, KeyError):
            raise OpenShiftCheckException(
                'get_route_failed',
                '"oc get route" returned an unexpected response:\n' + get_route
            )

        # ingress can be null if there is no router, or empty if not routed
        if not ingress or not ingress[0]:
            raise OpenShiftCheckException(
                'route_not_accepted',
                'The logging-kibana route is not being routed by any router.\n'
                'Is the router deployed and working?'
            )

        host = route.get("spec", {}).get("host")
        if not host:
            raise OpenShiftCheckException(
                'route_missing_host',
                'The logging-kibana route has no hostname defined,\n'
                'which should never happen. Did something alter its definition?'
            )

        return 'https://{}/'.format(host)

    def check_kibana_route(self):
        """
        Check to see if kibana route is up and working.
        Raises exception if not.
        """

        kibana_url = self._get_kibana_url()

        # first, check that kibana is reachable from the master.
        error = self._verify_url_internal(kibana_url)
        if error:
            if 'urlopen error [Errno 111] Connection refused' in error:
                raise OpenShiftCheckException(
                    'FailedToConnectInternal',
                    'Failed to connect from this master to Kibana URL {url}\n'
                    'Is kibana running, and is at least one router routing to it?'.format(url=kibana_url)
                )
            elif 'urlopen error [Errno -2] Name or service not known' in error:
                raise OpenShiftCheckException(
                    'FailedToResolveInternal',
                    'Failed to connect from this master to Kibana URL {url}\n'
                    'because the hostname does not resolve.\n'
                    'Is DNS configured for the Kibana hostname?'.format(url=kibana_url)
                )
            elif 'Status code was not' in error:
                raise OpenShiftCheckException(
                    'WrongReturnCodeInternal',
                    'A request from this master to the Kibana URL {url}\n'
                    'did not return the correct status code (302).\n'
                    'This could mean that Kibana is malfunctioning, the hostname is\n'
                    'resolving incorrectly, or other network issues. The output was:\n'
                    '  {error}'.format(url=kibana_url, error=error)
                )
            raise OpenShiftCheckException(
                'MiscRouteErrorInternal',
                'Error validating the logging Kibana route internally:\n' + error
            )

        # in production we would like the kibana route to work from outside the
        # cluster too; but that may not be the case, so allow disabling just this part.
        if self.get_var("openshift_check_efk_kibana_external", default="True").lower() != "true":
            return
        error = self._verify_url_external(kibana_url)

        if not error:
            return

        error_fmt = (
            'Error validating the logging Kibana route:\n{error}\n'
            'To disable external Kibana route validation, set the variable:\n'
            '  openshift_check_efk_kibana_external=False'
        )
        if 'urlopen error [Errno 111] Connection refused' in error:
            msg = (
                'Failed to connect from the Ansible control host to Kibana URL {url}\n'
                'Is the router for the Kibana hostname exposed externally?'
            ).format(url=kibana_url)
            raise OpenShiftCheckException('FailedToConnect', error_fmt.format(error=msg))
        elif 'urlopen error [Errno -2] Name or service not known' in error:
            msg = (
                'Failed to resolve the Kibana hostname in {url}\n'
                'from the Ansible control host.\n'
                'Is DNS configured to resolve this Kibana hostname externally?'
            ).format(url=kibana_url)
            raise OpenShiftCheckException('FailedToResolve', error_fmt.format(error=msg))
        elif 'Expected success (200)' in error:
            msg = (
                'A request to Kibana at {url}\n'
                'returned the wrong error code:\n'
                '  {error}\n'
                'This could mean that Kibana is malfunctioning, the hostname is\n'
                'resolving incorrectly, or other network issues.'
            ).format(url=kibana_url, error=error)
            raise OpenShiftCheckException('WrongReturnCode', error_fmt.format(error=msg))
        raise OpenShiftCheckException(
            'MiscRouteError',
            'Error validating the logging Kibana route externally:\n' + error
        )
