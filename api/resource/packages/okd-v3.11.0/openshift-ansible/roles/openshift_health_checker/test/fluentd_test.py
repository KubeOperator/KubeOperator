import pytest
import json

from openshift_checks.logging.fluentd import Fluentd, OpenShiftCheckExceptionList, OpenShiftCheckException


def assert_error_in_list(expect_err, errorlist):
    assert any(err.name == expect_err for err in errorlist), "{} in {}".format(str(expect_err), str(errorlist))


fluentd_pod_node1 = {
    "metadata": {
        "labels": {"component": "fluentd", "deploymentconfig": "logging-fluentd"},
        "name": "logging-fluentd-1",
    },
    "spec": {"host": "node1", "nodeName": "node1"},
    "status": {
        "containerStatuses": [{"ready": True}],
        "conditions": [{"status": "True", "type": "Ready"}],
    }
}
fluentd_pod_node2_down = {
    "metadata": {
        "labels": {"component": "fluentd", "deploymentconfig": "logging-fluentd"},
        "name": "logging-fluentd-2",
    },
    "spec": {"host": "node2", "nodeName": "node2"},
    "status": {
        "containerStatuses": [{"ready": False}],
        "conditions": [{"status": "False", "type": "Ready"}],
    }
}
fluentd_node1 = {
    "metadata": {
        "labels": {"logging-infra-fluentd": "true", "kubernetes.io/hostname": "node1"},
        "name": "node1",
    },
    "status": {"addresses": [{"type": "InternalIP", "address": "10.10.1.1"}]},
}
fluentd_node2 = {
    "metadata": {
        "labels": {"logging-infra-fluentd": "true", "kubernetes.io/hostname": "hostname"},
        "name": "node2",
    },
    "status": {"addresses": [{"type": "InternalIP", "address": "10.10.1.2"}]},
}
fluentd_node3_unlabeled = {
    "metadata": {
        "labels": {"kubernetes.io/hostname": "hostname"},
        "name": "node3",
    },
    "status": {"addresses": [{"type": "InternalIP", "address": "10.10.1.3"}]},
}


def test_get_fluentd_pods():
    check = Fluentd()
    check.exec_oc = lambda *_: json.dumps(dict(items=[fluentd_node1]))
    check.get_pods_for_component = lambda *_: [fluentd_pod_node1]
    assert not check.run()


@pytest.mark.parametrize('pods, nodes, expect_error', [
    (
        [],
        [],
        'NoNodesDefined',
    ),
    (
        [],
        [fluentd_node3_unlabeled],
        'NoNodesLabeled',
    ),
    (
        [],
        [fluentd_node1, fluentd_node3_unlabeled],
        'NodesUnlabeled',
    ),
    (
        [],
        [fluentd_node2],
        'MissingFluentdPod',
    ),
    (
        [fluentd_pod_node1, fluentd_pod_node1],
        [fluentd_node1],
        'TooManyFluentdPods',
    ),
    (
        [fluentd_pod_node2_down],
        [fluentd_node2],
        'FluentdNotRunning',
    ),
])
def test_get_fluentd_pods_errors(pods, nodes, expect_error):
    check = Fluentd()
    check.exec_oc = lambda *_: json.dumps(dict(items=nodes))

    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.check_fluentd(pods)
    if isinstance(excinfo.value, OpenShiftCheckExceptionList):
        assert_error_in_list(expect_error, excinfo.value)
    else:
        assert expect_error == excinfo.value.name


def test_bad_oc_node_list():
    check = Fluentd()
    check.exec_oc = lambda *_: "this isn't even json"
    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.get_nodes_by_name()
    assert 'BadOcNodeList' == excinfo.value.name
