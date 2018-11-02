import pytest

from openshift_checks.etcd_traffic import EtcdTraffic


@pytest.mark.parametrize('group_names,version,is_active', [
    (['oo_masters_to_config'], "3.5", False),
    (['oo_masters_to_config'], "3.6", False),
    (['oo_nodes_to_config'], "3.4", False),
    (['oo_etcd_to_config'], "3.4", True),
    (['oo_etcd_to_config'], "1.5", True),
    (['oo_etcd_to_config'], "3.1", False),
    (['oo_masters_to_config', 'oo_nodes_to_config'], "3.5", False),
    (['oo_masters_to_config', 'oo_etcd_to_config'], "3.5", True),
    ([], "3.4", False),
])
def test_is_active(group_names, version, is_active):
    task_vars = dict(
        group_names=group_names,
        openshift_image_tag=version,
    )
    assert EtcdTraffic(task_vars=task_vars).is_active() == is_active


@pytest.mark.parametrize('group_names,matched,failed,extra_words', [
    (["oo_masters_to_config"], True, True, ["Higher than normal", "traffic"]),
    (["oo_masters_to_config", "oo_etcd_to_config"], False, False, []),
    (["oo_etcd_to_config"], False, False, []),
])
def test_log_matches_high_traffic_msg(group_names, matched, failed, extra_words):
    def execute_module(module_name, *_):
        return {
            "matched": matched,
            "failed": failed,
        }

    task_vars = dict(
        group_names=group_names,
        openshift_is_containerized=False,
        openshift_service_type="origin"
    )

    result = EtcdTraffic(execute_module, task_vars).run()

    for word in extra_words:
        assert word in result.get("msg", "")

    assert result.get("failed", False) == failed


@pytest.mark.parametrize('openshift_is_containerized,expected_unit_value', [
    (False, "etcd"),
    (True, "etcd_container"),
])
def test_systemd_unit_matches_deployment_type(openshift_is_containerized, expected_unit_value):
    task_vars = dict(
        openshift_is_containerized=openshift_is_containerized
    )

    def execute_module(module_name, args, *_):
        assert module_name == "search_journalctl"
        matchers = args["log_matchers"]

        for matcher in matchers:
            assert matcher["unit"] == expected_unit_value

        return {"failed": False}

    EtcdTraffic(execute_module, task_vars).run()
