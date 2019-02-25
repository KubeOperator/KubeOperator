import pytest
import json

from openshift_checks.logging.logging import LoggingCheck, MissingComponentPods, CouldNotUseOc

task_vars_config_base = dict(openshift=dict(common=dict(config_base='/etc/origin')))


def canned_loggingcheck(exec_oc=None, execute_module=None):
    """Create a LoggingCheck object with canned exec_oc method"""
    check = LoggingCheck(execute_module)
    if exec_oc:
        check.exec_oc = exec_oc
    return check


def assert_error(error, expect_error):
    if expect_error:
        assert error
        assert expect_error in error
    else:
        assert not error


plain_es_pod = {
    "metadata": {
        "labels": {"component": "es", "deploymentconfig": "logging-es"},
        "name": "logging-es",
    },
    "status": {
        "conditions": [{"status": "True", "type": "Ready"}],
        "containerStatuses": [{"ready": True}],
        "podIP": "10.10.10.10",
    },
    "_test_master_name_str": "name logging-es",
}

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

plain_kibana_pod_no_containerstatus = {
    "metadata": {
        "labels": {"component": "kibana", "deploymentconfig": "logging-kibana"},
        "name": "logging-kibana-1",
    },
    "status": {
        "conditions": [{"status": "True", "type": "Ready"}],
    }
}

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

plain_curator_pod = {
    "metadata": {
        "labels": {"component": "curator", "deploymentconfig": "logging-curator"},
        "name": "logging-curator-1",
    },
    "status": {
        "containerStatuses": [{"ready": True}],
        "conditions": [{"status": "True", "type": "Ready"}],
        "podIP": "10.10.10.10",
    }
}


@pytest.mark.parametrize('problem, expect', [
    ("[Errno 2] No such file or directory", "supposed to be a master"),
    ("Permission denied", "Unexpected error using `oc`"),
])
def test_oc_failure(problem, expect):
    def execute_module(module_name, *_):
        if module_name == "ocutil":
            return dict(failed=True, result=problem)
        return dict(changed=False)

    check = LoggingCheck(execute_module, task_vars_config_base)

    with pytest.raises(CouldNotUseOc) as excinfo:
        check.exec_oc('get foo', [])
    assert expect in str(excinfo)


groups_with_first_master = dict(oo_first_master=['this-host'])
groups_not_a_master = dict(oo_first_master=['other-host'], oo_masters=['other-host'])


@pytest.mark.parametrize('groups, logging_deployed, is_active', [
    (groups_with_first_master, True, True),
    (groups_with_first_master, False, False),
    (groups_not_a_master, True, False),
    (groups_not_a_master, True, False),
])
def test_is_active(groups, logging_deployed, is_active):
    task_vars = dict(
        ansible_host='this-host',
        groups=groups,
        openshift_logging_install_logging=logging_deployed,
    )

    assert LoggingCheck(None, task_vars).is_active() == is_active


@pytest.mark.parametrize('pod_output, expect_pods', [
    (
        json.dumps({'items': [plain_es_pod]}),
        [plain_es_pod],
    ),
])
def test_get_pods_for_component(pod_output, expect_pods):
    check = canned_loggingcheck(lambda *_: pod_output)
    pods = check.get_pods_for_component("es")
    assert pods == expect_pods


@pytest.mark.parametrize('exec_oc_output, expect_error', [
    (
        'No resources found.',
        MissingComponentPods,
    ),
    (
        '{"items": null}',
        MissingComponentPods,
    ),
])
def test_get_pods_for_component_fail(exec_oc_output, expect_error):
    check = canned_loggingcheck(lambda *_: exec_oc_output)
    with pytest.raises(expect_error):
        check.get_pods_for_component("es")


@pytest.mark.parametrize('name, pods, expected_pods', [
    (
        'test single pod found, scheduled, but no containerStatuses field',
        [plain_kibana_pod_no_containerstatus],
        [plain_kibana_pod_no_containerstatus],
    ),
    (
        'set of pods has at least one pod with containerStatuses (scheduled); should still fail',
        [plain_kibana_pod_no_containerstatus, plain_kibana_pod],
        [plain_kibana_pod_no_containerstatus],
    ),

], ids=lambda argvals: argvals[0])
def test_get_not_running_pods_no_container_status(name, pods, expected_pods):
    check = canned_loggingcheck(lambda *_: '')
    result = check.not_running_pods(pods)

    assert result == expected_pods
