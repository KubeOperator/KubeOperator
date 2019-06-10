import pytest

from openshift_checks.package_availability import PackageAvailability


@pytest.mark.parametrize('pkg_mgr,openshift_is_atomic,is_active', [
    ('yum', False, True),
    ('yum', True, False),
    ('dnf', True, False),
    ('dnf', False, False),
])
def test_is_active(pkg_mgr, openshift_is_atomic, is_active):
    task_vars = dict(
        ansible_pkg_mgr=pkg_mgr,
        openshift_is_atomic=openshift_is_atomic,
    )
    assert PackageAvailability(None, task_vars).is_active() == is_active


@pytest.mark.parametrize('task_vars,must_have_packages,must_not_have_packages', [
    (
        dict(openshift_service_type='origin'),
        set(),
        set(['openshift-hyperkube', 'openshift-node']),
    ),
    (
        dict(
            openshift_service_type='origin',
            group_names=['oo_masters_to_config'],
        ),
        set(['origin-hyperkube']),
        set(['origin-node']),
    ),
    (
        dict(
            openshift_service_type='atomic-openshift',
            group_names=['oo_nodes_to_config'],
        ),
        set(['atomic-openshift-node']),
        set(['atomic-openshift-hyperkube']),
    ),
    (
        dict(
            openshift_service_type='atomic-openshift',
            group_names=['oo_masters_to_config', 'oo_nodes_to_config'],
        ),
        set(['atomic-openshift-hyperkube', 'atomic-openshift-node']),
        set(),
    ),
])
def test_package_availability(task_vars, must_have_packages, must_not_have_packages):
    return_value = {}

    def execute_module(module_name=None, module_args=None, *_):
        assert module_name == 'check_yum_update'
        assert 'packages' in module_args
        assert set(module_args['packages']).issuperset(must_have_packages)
        assert not set(module_args['packages']).intersection(must_not_have_packages)
        return {'foo': return_value}

    result = PackageAvailability(execute_module, task_vars).run()
    assert result['foo'] is return_value
