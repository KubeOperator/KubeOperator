import pytest
import json

from openshift_checks.logging.elasticsearch import Elasticsearch, OpenShiftCheckExceptionList


task_vars_config_base = dict(openshift=dict(common=dict(config_base='/etc/origin')))


def canned_elasticsearch(task_vars=None, exec_oc=None):
    """Create an Elasticsearch check object with stubbed exec_oc method"""
    check = Elasticsearch(None, task_vars or {})
    if exec_oc:
        check.exec_oc = exec_oc
    return check


def assert_error_in_list(expect_err, errorlist):
    assert any(err.name == expect_err for err in errorlist), "{} in {}".format(str(expect_err), str(errorlist))


def pods_by_name(pods):
    return {pod['metadata']['name']: pod for pod in pods}


plain_es_pod = {
    "metadata": {
        "labels": {"component": "es", "deploymentconfig": "logging-es"},
        "name": "logging-es",
    },
    "spec": {},
    "status": {
        "conditions": [{"status": "True", "type": "Ready"}],
        "containerStatuses": [{"ready": True}],
        "podIP": "10.10.10.10",
    },
    "_test_master_name_str": "name logging-es",
}

split_es_pod = {
    "metadata": {
        "labels": {"component": "es", "deploymentconfig": "logging-es-2"},
        "name": "logging-es-2",
    },
    "spec": {},
    "status": {
        "conditions": [{"status": "True", "type": "Ready"}],
        "containerStatuses": [{"ready": True}],
        "podIP": "10.10.10.10",
    },
    "_test_master_name_str": "name logging-es-2",
}

unready_es_pod = {
    "metadata": {
        "labels": {"component": "es", "deploymentconfig": "logging-es-3"},
        "name": "logging-es-3",
    },
    "spec": {},
    "status": {
        "conditions": [{"status": "False", "type": "Ready"}],
        "containerStatuses": [{"ready": False}],
        "podIP": "10.10.10.10",
    },
    "_test_master_name_str": "BAD_NAME_RESPONSE",
}


def test_check_elasticsearch():
    with pytest.raises(OpenShiftCheckExceptionList) as excinfo:
        canned_elasticsearch().check_elasticsearch([])
    assert_error_in_list('NoRunningPods', excinfo.value)

    # canned oc responses to match so all the checks pass
    def exec_oc(cmd, args, **_):
        if '_cat/master' in cmd:
            return 'name logging-es'
        elif '/_nodes' in cmd:
            return json.dumps(es_node_list)
        elif '_cluster/health' in cmd:
            return '{"status": "green"}'
        elif ' df ' in cmd:
            return 'IUse% Use%\n 3%  4%\n'
        else:
            raise Exception(cmd)

    check = canned_elasticsearch({}, exec_oc)
    check.get_pods_for_component = lambda *_: [plain_es_pod]
    assert {} == check.run()


def test_check_running_es_pods():
    pods, errors = Elasticsearch().running_elasticsearch_pods([plain_es_pod, unready_es_pod])
    assert plain_es_pod in pods
    assert_error_in_list('PodNotRunning', errors)


def test_check_elasticsearch_masters():
    pods = [plain_es_pod]
    check = canned_elasticsearch(task_vars_config_base, lambda *args, **_: plain_es_pod['_test_master_name_str'])
    assert not check.check_elasticsearch_masters(pods_by_name(pods))


@pytest.mark.parametrize('pods, expect_error', [
    (
        [],
        'NoMasterFound',
    ),
    (
        [unready_es_pod],
        'NoMasterName',
    ),
    (
        [plain_es_pod, split_es_pod],
        'SplitBrainMasters',
    ),
])
def test_check_elasticsearch_masters_error(pods, expect_error):
    test_pods = list(pods)
    check = canned_elasticsearch(task_vars_config_base, lambda *args, **_: test_pods.pop(0)['_test_master_name_str'])
    assert_error_in_list(expect_error, check.check_elasticsearch_masters(pods_by_name(pods)))


es_node_list = {
    'nodes': {
        'random-es-name': {
            'host': 'logging-es',
        }}}


def test_check_elasticsearch_node_list():
    check = canned_elasticsearch(task_vars_config_base, lambda *args, **_: json.dumps(es_node_list))
    assert not check.check_elasticsearch_node_list(pods_by_name([plain_es_pod]))


@pytest.mark.parametrize('pods, node_list, expect_error', [
    (
        [],
        {},
        'MissingComponentPods',
    ),
    (
        [plain_es_pod],
        {},  # empty list of nodes triggers KeyError
        'MissingNodeList',
    ),
    (
        [split_es_pod],
        es_node_list,
        'EsPodNodeMismatch',
    ),
])
def test_check_elasticsearch_node_list_errors(pods, node_list, expect_error):
    check = canned_elasticsearch(task_vars_config_base, lambda cmd, args, **_: json.dumps(node_list))
    assert_error_in_list(expect_error, check.check_elasticsearch_node_list(pods_by_name(pods)))


def test_check_elasticsearch_cluster_health():
    test_health_data = [{"status": "green"}]
    check = canned_elasticsearch(exec_oc=lambda *args, **_: json.dumps(test_health_data.pop(0)))
    assert not check.check_es_cluster_health(pods_by_name([plain_es_pod]))


@pytest.mark.parametrize('pods, health_data, expect_error', [
    (
        [plain_es_pod],
        [{"no-status": "should bomb"}],
        'BadEsResponse',
    ),
    (
        [plain_es_pod, split_es_pod],
        [{"status": "green"}, {"status": "red"}],
        'EsClusterHealthRed',
    ),
])
def test_check_elasticsearch_cluster_health_errors(pods, health_data, expect_error):
    test_health_data = list(health_data)
    check = canned_elasticsearch(exec_oc=lambda *args, **_: json.dumps(test_health_data.pop(0)))
    assert_error_in_list(expect_error, check.check_es_cluster_health(pods_by_name(pods)))


def test_check_elasticsearch_diskspace():
    check = canned_elasticsearch(exec_oc=lambda *args, **_: 'IUse% Use%\n 3%  4%\n')
    assert not check.check_elasticsearch_diskspace(pods_by_name([plain_es_pod]))


@pytest.mark.parametrize('disk_data, expect_error', [
    (
        'df: /elasticsearch/persistent: No such file or directory\n',
        'BadDfResponse',
    ),
    (
        'IUse% Use%\n 95%  40%\n',
        'InodeUsageTooHigh',
    ),
    (
        'IUse% Use%\n 3%  94%\n',
        'DiskUsageTooHigh',
    ),
])
def test_check_elasticsearch_diskspace_errors(disk_data, expect_error):
    check = canned_elasticsearch(exec_oc=lambda *args, **_: disk_data)
    assert_error_in_list(expect_error, check.check_elasticsearch_diskspace(pods_by_name([plain_es_pod])))
