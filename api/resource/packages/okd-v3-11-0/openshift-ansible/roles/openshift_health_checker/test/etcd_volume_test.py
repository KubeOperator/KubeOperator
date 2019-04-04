import pytest

from openshift_checks.etcd_volume import EtcdVolume
from openshift_checks import OpenShiftCheckException


@pytest.mark.parametrize('ansible_mounts,extra_words', [
    ([], ['none']),  # empty ansible_mounts
    ([{'mount': '/mnt'}], ['/mnt']),  # missing relevant mount paths
])
def test_cannot_determine_available_disk(ansible_mounts, extra_words):
    task_vars = dict(
        ansible_mounts=ansible_mounts,
    )

    with pytest.raises(OpenShiftCheckException) as excinfo:
        EtcdVolume(fake_execute_module, task_vars).run()

    for word in ['Unable to determine mount point'] + extra_words:
        assert word in str(excinfo.value)


@pytest.mark.parametrize('size_limit,ansible_mounts', [
    (
        # if no size limit is specified, expect max usage
        # limit to default to 90% of size_total
        None,
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 80 * 10**9
        }],
    ),
    (
        1,
        [{
            'mount': '/',
            'size_available': 30 * 10**9,
            'size_total': 30 * 10**9,
        }],
    ),
    (
        20000000000,
        [{
            'mount': '/',
            'size_available': 20 * 10**9,
            'size_total': 40 * 10**9,
        }],
    ),
    (
        5000000000,
        [{
            # not enough space on / ...
            'mount': '/',
            'size_available': 0,
            'size_total': 0,
        }, {
            # not enough space on /var/lib ...
            'mount': '/var/lib',
            'size_available': 2 * 10**9,
            'size_total': 21 * 10**9,
        }, {
            # ... but enough on /var/lib/etcd
            'mount': '/var/lib/etcd',
            'size_available': 36 * 10**9,
            'size_total': 40 * 10**9
        }],
    )
])
def test_succeeds_with_recommended_disk_space(size_limit, ansible_mounts):
    task_vars = dict(
        etcd_device_usage_threshold_percent=size_limit,
        ansible_mounts=ansible_mounts,
    )

    if task_vars["etcd_device_usage_threshold_percent"] is None:
        task_vars.pop("etcd_device_usage_threshold_percent")

    result = EtcdVolume(fake_execute_module, task_vars).run()

    assert not result.get('failed', False)


@pytest.mark.parametrize('size_limit_percent,ansible_mounts,extra_words', [
    (
        # if no size limit is specified, expect max usage
        # limit to default to 90% of size_total
        None,
        [{
            'mount': '/',
            'size_available': 1 * 10**9,
            'size_total': 100 * 10**9,
        }],
        ['99.0%'],
    ),
    (
        70.0,
        [{
            'mount': '/',
            'size_available': 1 * 10**6,
            'size_total': 5 * 10**9,
        }],
        ['100.0%'],
    ),
    (
        40.0,
        [{
            'mount': '/',
            'size_available': 2 * 10**9,
            'size_total': 6 * 10**9,
        }],
        ['66.7%'],
    ),
    (
        None,
        [{
            # enough space on /var ...
            'mount': '/var',
            'size_available': 20 * 10**9,
            'size_total': 20 * 10**9,
        }, {
            # .. but not enough on /var/lib
            'mount': '/var/lib',
            'size_available': 1 * 10**9,
            'size_total': 20 * 10**9,
        }],
        ['95.0%'],
    ),
])
def test_fails_with_insufficient_disk_space(size_limit_percent, ansible_mounts, extra_words):
    task_vars = dict(
        etcd_device_usage_threshold_percent=size_limit_percent,
        ansible_mounts=ansible_mounts,
    )

    if task_vars["etcd_device_usage_threshold_percent"] is None:
        task_vars.pop("etcd_device_usage_threshold_percent")

    result = EtcdVolume(fake_execute_module, task_vars).run()

    assert result['failed']
    for word in extra_words:
        assert word in result['msg']


def fake_execute_module(*args):
    raise AssertionError('this function should not be called')
