import pytest

from openshift_checks.package_version import PackageVersion


def task_vars_for(openshift_release, deployment_type):
    service_type_dict = {'origin': 'origin',
                         'openshift-enterprise': 'atomic-openshift'}
    service_type = service_type_dict[deployment_type]
    return dict(
        ansible_pkg_mgr='yum',
        openshift_service_type=service_type,
        openshift_release=openshift_release,
        openshift_image_tag='v' + openshift_release,
        openshift_deployment_type=deployment_type,
    )


@pytest.mark.parametrize('openshift_release', [
    "111.7.0",
    "3.7",
    "3.6",
    "3.5.1.2.3",
    "3.5",
    "3.4",
    "3.3",
    "2.1.0",
])
def test_package_version(openshift_release):

    return_value = {"foo": object()}

    def execute_module(module_name=None, module_args=None, tmp=None, task_vars=None, *_):
        assert module_name == 'aos_version'
        assert "package_list" in module_args

        for pkg in module_args["package_list"]:
            if "-master" in pkg["name"] or "-node" in pkg["name"]:
                assert pkg["version"] == task_vars["openshift_release"]

        return return_value

    check = PackageVersion(execute_module, task_vars_for(openshift_release, 'origin'))
    result = check.run()
    assert result == return_value


@pytest.mark.parametrize('group_names,openshift_is_atomic,is_active', [
    (['oo_masters_to_config'], False, True),
    # ensure check is skipped on containerized installs
    (['oo_masters_to_config'], True, False),
    (['oo_nodes_to_config'], False, True),
    (['oo_masters_to_config', 'oo_nodes_to_config'], False, True),
    (['oo_masters_to_config', 'oo_etcd_to_config'], False, True),
    ([], False, False),
    (['oo_etcd_to_config'], False, False),
    (['lb'], False, False),
    (['nfs'], False, False),
])
def test_package_version_skip_when_not_master_nor_node(group_names, openshift_is_atomic, is_active):
    task_vars = dict(
        group_names=group_names,
        openshift_is_atomic=openshift_is_atomic,
    )
    assert PackageVersion(None, task_vars).is_active() == is_active
