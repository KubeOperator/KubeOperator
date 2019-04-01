import pytest
import json

# pylint can't find the package when its installed in virtualenv
from ansible.module_utils.six.moves.urllib import request  # pylint: disable=import-error
# pylint: disable=import-error
from ansible.module_utils.six.moves.urllib.error import HTTPError, URLError

from openshift_checks.logging.kibana import Kibana, OpenShiftCheckException


plain_kibana_pod = {
    "metadata": {
        "labels": {"component": "kibana", "deploymentconfig": "logging-kibana"},
        "name": "logging-kibana-1",
    },
    "status": {
        "containerStatuses": [{"ready": True}, {"ready": True}],
        "conditions": [{"status": "True", "type": "Ready"}],
    }
}
not_running_kibana_pod = {
    "metadata": {
        "labels": {"component": "kibana", "deploymentconfig": "logging-kibana"},
        "name": "logging-kibana-2",
    },
    "status": {
        "containerStatuses": [{"ready": True}, {"ready": False}],
        "conditions": [{"status": "True", "type": "Ready"}],
    }
}


def test_check_kibana():
    # should run without exception:
    Kibana().check_kibana([plain_kibana_pod])


@pytest.mark.parametrize('pods, expect_error', [
    (
        [],
        "MissingComponentPods",
    ),
    (
        [not_running_kibana_pod],
        "NoRunningPods",
    ),
    (
        [plain_kibana_pod, not_running_kibana_pod],
        "PodNotRunning",
    ),
])
def test_check_kibana_error(pods, expect_error):
    with pytest.raises(OpenShiftCheckException) as excinfo:
        Kibana().check_kibana(pods)
    assert expect_error == excinfo.value.name


@pytest.mark.parametrize('comment, route, expect_error', [
    (
        "No route returned",
        None,
        "no_route_exists",
    ),

    (
        "broken route response",
        {"status": {}},
        "get_route_failed",
    ),
    (
        "route with no ingress",
        {
            "metadata": {
                "labels": {"component": "kibana", "deploymentconfig": "logging-kibana"},
                "name": "logging-kibana",
            },
            "status": {
                "ingress": [],
            },
            "spec": {
                "host": "hostname",
            }
        },
        "route_not_accepted",
    ),

    (
        "route with no host",
        {
            "metadata": {
                "labels": {"component": "kibana", "deploymentconfig": "logging-kibana"},
                "name": "logging-kibana",
            },
            "status": {
                "ingress": [{
                    "status": True,
                }],
            },
            "spec": {},
        },
        "route_missing_host",
    ),
])
def test_get_kibana_url_error(comment, route, expect_error):
    check = Kibana()
    check.exec_oc = lambda *_: json.dumps(route) if route else ""

    with pytest.raises(OpenShiftCheckException) as excinfo:
        check._get_kibana_url()
    assert excinfo.value.name == expect_error


@pytest.mark.parametrize('comment, route, expect_url', [
    (
        "test route that looks fine",
        {
            "metadata": {
                "labels": {"component": "kibana", "deploymentconfig": "logging-kibana"},
                "name": "logging-kibana",
            },
            "status": {
                "ingress": [{
                    "status": True,
                }],
            },
            "spec": {
                "host": "hostname",
            },
        },
        "https://hostname/",
    ),
])
def test_get_kibana_url(comment, route, expect_url):
    check = Kibana()
    check.exec_oc = lambda *_: json.dumps(route)
    assert expect_url == check._get_kibana_url()


@pytest.mark.parametrize('exec_result, expect', [
    (
        'urlopen error [Errno 111] Connection refused',
        'FailedToConnectInternal',
    ),
    (
        'urlopen error [Errno -2] Name or service not known',
        'FailedToResolveInternal',
    ),
    (
        'Status code was not [302]: HTTP Error 500: Server error',
        'WrongReturnCodeInternal',
    ),
    (
        'bork bork bork',
        'MiscRouteErrorInternal',
    ),
])
def test_verify_url_internal_failure(exec_result, expect):
    check = Kibana(execute_module=lambda *_: dict(failed=True, msg=exec_result))
    check._get_kibana_url = lambda: 'url'

    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.check_kibana_route()
    assert expect == excinfo.value.name


@pytest.mark.parametrize('lib_result, expect', [
    (
        HTTPError('url', 500, 'it broke', hdrs=None, fp=None),
        'MiscRouteError',
    ),
    (
        URLError('urlopen error [Errno 111] Connection refused'),
        'FailedToConnect',
    ),
    (
        URLError('urlopen error [Errno -2] Name or service not known'),
        'FailedToResolve',
    ),
    (
        302,
        'WrongReturnCode',
    ),
    (
        200,
        None,
    ),
])
def test_verify_url_external_failure(lib_result, expect, monkeypatch):

    class _http_return:

        def __init__(self, code):
            self.code = code

        def getcode(self):
            return self.code

    def urlopen(url, context):
        if type(lib_result) is int:
            return _http_return(lib_result)
        raise lib_result
    monkeypatch.setattr(request, 'urlopen', urlopen)

    check = Kibana()
    check._get_kibana_url = lambda: 'url'
    check._verify_url_internal = lambda url: None

    if not expect:
        check.check_kibana_route()
        return

    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.check_kibana_route()
    assert expect == excinfo.value.name


def test_verify_url_external_skip():
    check = Kibana(lambda *_: {}, dict(openshift_check_efk_kibana_external="false"))
    check._get_kibana_url = lambda: 'url'
    check.check_kibana_route()


# this is kind of silly but it adds coverage for the run() method...
def test_run():
    pods = ["foo"]
    ran = dict(check_kibana=False, check_route=False)

    def check_kibana(pod_list):
        ran["check_kibana"] = True
        assert pod_list == pods

    def check_kibana_route():
        ran["check_route"] = True

    check = Kibana()
    check.get_pods_for_component = lambda *_: pods
    check.check_kibana = check_kibana
    check.check_kibana_route = check_kibana_route

    check.run()
    assert ran["check_kibana"] and ran["check_route"]
