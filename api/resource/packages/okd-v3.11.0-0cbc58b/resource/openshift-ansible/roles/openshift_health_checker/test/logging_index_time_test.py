import json

import pytest

from openshift_checks.logging.logging_index_time import LoggingIndexTime, OpenShiftCheckException


SAMPLE_UUID = "unique-test-uuid"


def canned_loggingindextime(exec_oc=None):
    """Create a check object with a canned exec_oc method"""
    check = LoggingIndexTime()  # fails if a module is actually invoked
    if exec_oc:
        check.exec_oc = exec_oc
    return check


plain_running_elasticsearch_pod = {
    "metadata": {
        "labels": {"component": "es", "deploymentconfig": "logging-es-data-master"},
        "name": "logging-es-data-master-1",
    },
    "status": {
        "containerStatuses": [{"ready": True}, {"ready": True}],
        "phase": "Running",
    }
}
plain_running_kibana_pod = {
    "metadata": {
        "labels": {"component": "kibana", "deploymentconfig": "logging-kibana"},
        "name": "logging-kibana-1",
    },
    "status": {
        "containerStatuses": [{"ready": True}, {"ready": True}],
        "phase": "Running",
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
        "phase": "pending",
    }
}


@pytest.mark.parametrize('pods, expect_pods', [
    (
        [not_running_kibana_pod],
        [],
    ),
    (
        [plain_running_kibana_pod],
        [plain_running_kibana_pod],
    ),
    (
        [],
        [],
    )
])
def test_check_running_pods(pods, expect_pods):
    check = canned_loggingindextime()
    pods = check.running_pods(pods)
    assert pods == expect_pods


def test_bad_config_param():
    with pytest.raises(OpenShiftCheckException) as error:
        LoggingIndexTime(task_vars=dict(openshift_check_logging_index_timeout_seconds="foo")).run()
    assert 'InvalidTimeout' == error.value.name


def test_no_running_pods():
    check = LoggingIndexTime()
    check.get_pods_for_component = lambda *_: [not_running_kibana_pod]
    with pytest.raises(OpenShiftCheckException) as error:
        check.run()
    assert 'kibanaNoRunningPods' == error.value.name


def test_with_running_pods():
    check = LoggingIndexTime()
    check.get_pods_for_component = lambda *_: [plain_running_kibana_pod, plain_running_elasticsearch_pod]
    check.curl_kibana_with_uuid = lambda *_: SAMPLE_UUID
    check.wait_until_cmd_or_err = lambda *_: None
    assert not check.run().get("failed")


@pytest.mark.parametrize('name, json_response, uuid, timeout', [
    (
        'valid count in response',
        {
            "count": 1,
        },
        SAMPLE_UUID,
        0.001,
    ),
], ids=lambda argval: argval[0])
def test_wait_until_cmd_or_err_succeeds(name, json_response, uuid, timeout):
    check = canned_loggingindextime(lambda *args, **_: json.dumps(json_response))
    check.wait_until_cmd_or_err(plain_running_elasticsearch_pod, uuid, timeout)


@pytest.mark.parametrize('name, json_response, timeout, expect_error', [
    (
        'invalid json response',
        {
            "invalid_field": 1,
        },
        0.001,
        'esInvalidResponse',
    ),
    (
        'empty response',
        {},
        0.001,
        'esInvalidResponse',
    ),
    (
        'valid response but invalid match count',
        {
            "count": 0,
        },
        0.005,
        'NoMatchFound',
    )
], ids=lambda argval: argval[0])
def test_wait_until_cmd_or_err(name, json_response, timeout, expect_error):
    check = canned_loggingindextime(lambda *args, **_: json.dumps(json_response))
    with pytest.raises(OpenShiftCheckException) as error:
        check.wait_until_cmd_or_err(plain_running_elasticsearch_pod, SAMPLE_UUID, timeout)

    assert expect_error == error.value.name


def test_curl_kibana_with_uuid():
    check = canned_loggingindextime(lambda *args, **_: json.dumps({"statusCode": 404}))
    check.generate_uuid = lambda: SAMPLE_UUID
    assert SAMPLE_UUID == check.curl_kibana_with_uuid(plain_running_kibana_pod)


@pytest.mark.parametrize('name, json_response, expect_error', [
    (
        'invalid json response',
        {
            "invalid_field": "invalid",
        },
        'kibanaInvalidResponse',
    ),
    (
        'wrong error code in response',
        {
            "statusCode": 500,
        },
        'kibanaInvalidReturnCode',
    ),
], ids=lambda argval: argval[0])
def test_failed_curl_kibana_with_uuid(name, json_response, expect_error):
    check = canned_loggingindextime(lambda *args, **_: json.dumps(json_response))
    check.generate_uuid = lambda: SAMPLE_UUID

    with pytest.raises(OpenShiftCheckException) as error:
        check.curl_kibana_with_uuid(plain_running_kibana_pod)

    assert expect_error == error.value.name
