import pytest

from openshift_checks.logging.fluentd_config import FluentdConfig, OpenShiftCheckException


def canned_fluentd_pod(containers):
    return {
        "metadata": {
            "labels": {"component": "fluentd", "deploymentconfig": "logging-fluentd"},
            "name": "logging-fluentd-1",
        },
        "spec": {
            "host": "node1",
            "nodeName": "node1",
            "containers": containers,
        },
        "status": {
            "phase": "Running",
            "containerStatuses": [{"ready": True}],
            "conditions": [{"status": "True", "type": "Ready"}],
        }
    }


fluentd_pod = {
    "metadata": {
        "labels": {"component": "fluentd", "deploymentconfig": "logging-fluentd"},
        "name": "logging-fluentd-1",
    },
    "spec": {
        "host": "node1",
        "nodeName": "node1",
        "containers": [
            {
                "name": "container1",
                "env": [
                    {
                        "name": "USE_JOURNAL",
                        "value": "true",
                    }
                ],
            }
        ],
    },
    "status": {
        "phase": "Running",
        "containerStatuses": [{"ready": True}],
        "conditions": [{"status": "True", "type": "Ready"}],
    }
}

not_running_fluentd_pod = {
    "metadata": {
        "labels": {"component": "fluentd", "deploymentconfig": "logging-fluentd"},
        "name": "logging-fluentd-2",
    },
    "status": {
        "phase": "Unknown",
        "containerStatuses": [{"ready": True}, {"ready": False}],
        "conditions": [{"status": "True", "type": "Ready"}],
    }
}


@pytest.mark.parametrize('name, use_journald, logging_driver, extra_words', [
    (
        'test success with use_journald=false, and docker config set to use "json-file"',
        False,
        "json-file",
        [],
    ),
], ids=lambda argvals: argvals[0])
def test_check_logging_config_non_master(name, use_journald, logging_driver, extra_words):
    def execute_module(module_name, args):
        if module_name == "docker_info":
            return {
                "info": {
                    "LoggingDriver": logging_driver,
                }
            }

        return {}

    task_vars = dict(
        group_names=["oo_nodes_to_config", "oo_etcd_to_config"],
        openshift_logging_fluentd_use_journal=use_journald,
        openshift=dict(
            common=dict(config_base=""),
        ),
    )

    check = FluentdConfig(execute_module, task_vars)
    check.execute_module = execute_module
    error = check.check_logging_config()

    assert error is None


@pytest.mark.parametrize('name, use_journald, logging_driver, words', [
    (
        'test failure with use_journald=false, but docker config set to use "journald"',
        False,
        "journald",
        ['json log files', 'has been set to use "journald"'],
    ),
    (
        'test failure with use_journald=false, but docker config set to use an "unsupported" driver',
        False,
        "unsupported",
        ["json log files", 'has been set to use "unsupported"'],
    ),
    (
        'test failure with use_journald=true, but docker config set to use "json-file"',
        True,
        "json-file",
        ['logs from "journald"', 'has been set to use "json-file"'],
    ),
], ids=lambda argvals: argvals[0])
def test_check_logging_config_non_master_failed(name, use_journald, logging_driver, words):
    def execute_module(module_name, args):
        if module_name == "docker_info":
            return {
                "info": {
                    "LoggingDriver": logging_driver,
                }
            }

        return {}

    task_vars = dict(
        group_names=["oo_nodes_to_config", "oo_etcd_to_config"],
        openshift_logging_fluentd_use_journal=use_journald,
        openshift=dict(
            common=dict(config_base=""),
        ),
    )

    check = FluentdConfig(execute_module, task_vars)
    check.execute_module = execute_module
    error = check.check_logging_config()

    assert error is not None
    for word in words:
        assert word in error


@pytest.mark.parametrize('name, pods, logging_driver, extra_words', [
    # use_journald returns false (not using journald), but check succeeds
    # since docker is set to use json-file
    (
        'test success with use_journald=false, and docker config set to use default driver "json-file"',
        [canned_fluentd_pod(
            [
                {
                    "name": "container1",
                    "env": [{
                        "name": "USE_JOURNAL",
                        "value": "false",
                    }],
                },
            ]
        )],
        "json-file",
        [],
    ),
    (
        'test success with USE_JOURNAL env var missing and docker config set to use default driver "json-file"',
        [canned_fluentd_pod(
            [
                {
                    "name": "container1",
                    "env": [{
                        "name": "RANDOM",
                        "value": "value",
                    }],
                },
            ]
        )],
        "json-file",
        [],
    ),
], ids=lambda argvals: argvals[0])
def test_check_logging_config_master(name, pods, logging_driver, extra_words):
    def execute_module(module_name, args):
        if module_name == "docker_info":
            return {
                "info": {
                    "LoggingDriver": logging_driver,
                }
            }

        return {}

    task_vars = dict(
        group_names=["oo_masters_to_config"],
        openshift=dict(
            common=dict(config_base=""),
        ),
    )

    check = FluentdConfig(execute_module, task_vars)
    check.execute_module = execute_module
    check.get_pods_for_component = lambda _: pods
    error = check.check_logging_config()

    assert error is None


@pytest.mark.parametrize('name, pods, logging_driver, words', [
    (
        'test failure with use_journald=false, but docker config set to use "journald"',
        [canned_fluentd_pod(
            [
                {
                    "name": "container1",
                    "env": [{
                        "name": "USE_JOURNAL",
                        "value": "false",
                    }],
                },
            ]
        )],
        "journald",
        ['json log files', 'has been set to use "journald"'],
    ),
    (
        'test failure with use_journald=true, but docker config set to use "json-file"',
        [fluentd_pod],
        "json-file",
        ['logs from "journald"', 'has been set to use "json-file"'],
    ),
    (
        'test failure with use_journald=false, but docker set to use an "unsupported" driver',
        [canned_fluentd_pod(
            [
                {
                    "name": "container1",
                    "env": [{
                        "name": "USE_JOURNAL",
                        "value": "false",
                    }],
                },
            ]
        )],
        "unsupported",
        ["json log files", 'has been set to use "unsupported"'],
    ),
    (
        'test failure with USE_JOURNAL env var missing and docker config set to use "journald"',
        [canned_fluentd_pod(
            [
                {
                    "name": "container1",
                    "env": [{
                        "name": "RANDOM",
                        "value": "value",
                    }],
                },
            ]
        )],
        "journald",
        ["configuration is set to", "json log files"],
    ),
], ids=lambda argvals: argvals[0])
def test_check_logging_config_master_failed(name, pods, logging_driver, words):
    def execute_module(module_name, args):
        if module_name == "docker_info":
            return {
                "info": {
                    "LoggingDriver": logging_driver,
                }
            }

        return {}

    task_vars = dict(
        group_names=["oo_masters_to_config"],
        openshift=dict(
            common=dict(config_base=""),
        ),
    )

    check = FluentdConfig(execute_module, task_vars)
    check.execute_module = execute_module
    check.get_pods_for_component = lambda _: pods
    error = check.check_logging_config()

    assert error is not None
    for word in words:
        assert word in error


@pytest.mark.parametrize('name, pods, response, logging_driver, extra_words', [
    (
        'test OpenShiftCheckException with no running containers',
        [canned_fluentd_pod([])],
        {
            "failed": True,
            "result": "unexpected",
        },
        "json-file",
        ['no running containers'],
    ),
    (
        'test OpenShiftCheckException one container and no env vars set',
        [canned_fluentd_pod(
            [
                {
                    "name": "container1",
                    "env": [],
                },
            ]
        )],
        {
            "failed": True,
            "result": "unexpected",
        },
        "json-file",
        ['no environment variables'],
    ),
], ids=lambda argvals: argvals[0])
def test_check_logging_config_master_fails_on_unscheduled_deployment(name, pods, response, logging_driver, extra_words):
    def execute_module(module_name, args):
        if module_name == "docker_info":
            return {
                "info": {
                    "LoggingDriver": logging_driver,
                }
            }

        return {}

    task_vars = dict(
        group_names=["oo_masters_to_config"],
        openshift=dict(
            common=dict(config_base=""),
        ),
    )

    check = FluentdConfig(execute_module, task_vars)
    check.get_pods_for_component = lambda _: pods

    with pytest.raises(OpenShiftCheckException) as error:
        check.check_logging_config()

    assert error is not None
    for word in extra_words:
        assert word in str(error)
