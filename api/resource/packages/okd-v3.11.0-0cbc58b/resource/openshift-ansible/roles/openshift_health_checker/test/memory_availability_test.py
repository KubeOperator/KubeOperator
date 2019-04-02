import pytest

from openshift_checks.memory_availability import MemoryAvailability


@pytest.mark.parametrize('group_names,is_active', [
    (['oo_masters_to_config'], True),
    (['oo_nodes_to_config'], True),
    (['oo_etcd_to_config'], True),
    (['oo_masters_to_config', 'oo_nodes_to_config'], True),
    (['oo_masters_to_config', 'oo_etcd_to_config'], True),
    ([], False),
    (['lb'], False),
    (['nfs'], False),
])
def test_is_active(group_names, is_active):
    task_vars = dict(
        group_names=group_names,
    )
    assert MemoryAvailability(None, task_vars).is_active() == is_active


@pytest.mark.parametrize('group_names,configured_min,ansible_memtotal_mb', [
    (
        ['oo_masters_to_config'],
        0,
        17200,
    ),
    (
        ['oo_nodes_to_config'],
        0,
        8200,
    ),
    (
        ['oo_nodes_to_config'],
        1,  # configure lower threshold
        2000,  # too low for recommended but not for configured
    ),
    (
        ['oo_nodes_to_config'],
        2,  # configure threshold where adjustment pushes it over
        1900,
    ),
    (
        ['oo_etcd_to_config'],
        0,
        8200,
    ),
    (
        ['oo_masters_to_config', 'oo_nodes_to_config'],
        0,
        17000,
    ),
])
def test_succeeds_with_recommended_memory(group_names, configured_min, ansible_memtotal_mb):
    task_vars = dict(
        group_names=group_names,
        openshift_check_min_host_memory_gb=configured_min,
        ansible_memtotal_mb=ansible_memtotal_mb,
    )

    result = MemoryAvailability(fake_execute_module, task_vars).run()

    assert not result.get('failed', False)


@pytest.mark.parametrize('group_names,configured_min,ansible_memtotal_mb,extra_words', [
    (
        ['oo_masters_to_config'],
        0,
        0,
        ['0.0 GiB'],
    ),
    (
        ['oo_nodes_to_config'],
        0,
        100,
        ['0.1 GiB'],
    ),
    (
        ['oo_nodes_to_config'],
        24,  # configure higher threshold
        20 * 1024,  # enough to meet recommended but not configured
        ['20.0 GiB'],
    ),
    (
        ['oo_nodes_to_config'],
        24,  # configure higher threshold
        22 * 1024,  # not enough for adjustment to push over threshold
        ['22.0 GiB'],
    ),
    (
        ['oo_etcd_to_config'],
        0,
        6 * 1024,
        ['6.0 GiB'],
    ),
    (
        ['oo_etcd_to_config', 'oo_masters_to_config'],
        0,
        9 * 1024,  # enough memory for etcd, not enough for a master
        ['9.0 GiB'],
    ),
    (
        ['oo_nodes_to_config', 'oo_masters_to_config'],
        0,
        # enough memory for a node, not enough for a master
        11 * 1024,
        ['11.0 GiB'],
    ),
])
def test_fails_with_insufficient_memory(group_names, configured_min, ansible_memtotal_mb, extra_words):
    task_vars = dict(
        group_names=group_names,
        openshift_check_min_host_memory_gb=configured_min,
        ansible_memtotal_mb=ansible_memtotal_mb,
    )

    result = MemoryAvailability(fake_execute_module, task_vars).run()

    assert result.get('failed', False)
    for word in 'below recommended'.split() + extra_words:
        assert word in result['msg']


def fake_execute_module(*args):
    raise AssertionError('this function should not be called')
