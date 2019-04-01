import pytest

from openshift_checks import OpenShiftCheckException
from openshift_checks.docker_storage import DockerStorage


@pytest.mark.parametrize('openshift_is_atomic, group_names, is_active', [
    (False, ["oo_masters_to_config", "oo_etcd_to_config"], False),
    (False, ["oo_masters_to_config", "oo_nodes_to_config"], True),
    (True, ["oo_etcd_to_config"], False),
])
def test_is_active(openshift_is_atomic, group_names, is_active):
    task_vars = dict(
        openshift_is_atomic=openshift_is_atomic,
        group_names=group_names,
    )
    assert DockerStorage(None, task_vars).is_active() == is_active


def non_atomic_task_vars():
    return {"openshift_is_atomic": False}


@pytest.mark.parametrize('docker_info, failed, expect_msg', [
    (
        dict(failed=True, msg="Error connecting: Error while fetching server API version"),
        True,
        ["Is docker running on this host?"],
    ),
    (
        dict(msg="I have no info"),
        True,
        ["missing info"],
    ),
    (
        dict(info={
            "Driver": "devicemapper",
            "DriverStatus": [("Pool Name", "docker-docker--pool")],
        }),
        False,
        [],
    ),
    (
        dict(info={
            "Driver": "devicemapper",
            "DriverStatus": [("Data loop file", "true")],
        }),
        True,
        ["loopback devices with the Docker devicemapper storage driver"],
    ),
    (
        dict(info={
            "Driver": "overlay2",
            "DriverStatus": [("Backing Filesystem", "xfs")],
        }),
        False,
        [],
    ),
    (
        dict(info={
            "Driver": "overlay",
            "DriverStatus": [("Backing Filesystem", "btrfs")],
        }),
        True,
        ["storage is type 'btrfs'", "only supported with\n'xfs'"],
    ),
    (
        dict(info={
            "Driver": "overlay2",
            "DriverStatus": [("Backing Filesystem", "xfs")],
            "OperatingSystem": "Red Hat Enterprise Linux Server release 7.2 (Maipo)",
            "KernelVersion": "3.10.0-327.22.2.el7.x86_64",
        }),
        True,
        ["Docker reports kernel version 3.10.0-327"],
    ),
    (
        dict(info={
            "Driver": "overlay",
            "DriverStatus": [("Backing Filesystem", "xfs")],
            "OperatingSystem": "CentOS",
            "KernelVersion": "3.10.0-514",
        }),
        False,
        [],
    ),
    (
        dict(info={
            "Driver": "unsupported",
        }),
        True,
        ["unsupported Docker storage driver"],
    ),
])
def test_check_storage_driver(docker_info, failed, expect_msg):
    def execute_module(module_name, *_):
        if module_name == "yum":
            return {}
        if module_name != "docker_info":
            raise ValueError("not expecting module " + module_name)
        return docker_info

    check = DockerStorage(execute_module, non_atomic_task_vars())
    check.check_dm_usage = lambda status: dict()  # stub out for this test
    check.check_overlay_usage = lambda info: dict()  # stub out for this test
    result = check.run()

    if failed:
        assert result["failed"]
    else:
        assert not result.get("failed", False)

    for word in expect_msg:
        assert word in result["msg"]


enough_space = {
    "Pool Name": "docker--vg-docker--pool",
    "Data Space Used": "19.92 MB",
    "Data Space Total": "8.535 GB",
    "Metadata Space Used": "40.96 kB",
    "Metadata Space Total": "25.17 MB",
}

not_enough_space = {
    "Pool Name": "docker--vg-docker--pool",
    "Data Space Used": "10 GB",
    "Data Space Total": "10 GB",
    "Metadata Space Used": "42 kB",
    "Metadata Space Total": "43 kB",
}


@pytest.mark.parametrize('task_vars, driver_status, vg_free, success, expect_msg', [
    (
        {"max_thinpool_data_usage_percent": "not a float"},
        enough_space,
        "12g",
        False,
        ["is not a percentage"],
    ),
    (
        {},
        {},  # empty values from driver status
        "bogus",  # also does not parse as bytes
        False,
        ["Could not interpret", "as bytes"],
    ),
    (
        {},
        enough_space,
        "12.00g",
        True,
        [],
    ),
    (
        {},
        not_enough_space,
        "0.00",
        False,
        ["data usage", "metadata usage", "higher than threshold"],
    ),
])
def test_dm_usage(task_vars, driver_status, vg_free, success, expect_msg):
    check = DockerStorage(None, task_vars)
    check.get_vg_free = lambda pool: vg_free
    result = check.check_dm_usage(driver_status)
    result_success = not result.get("failed")

    assert result_success is success
    for msg in expect_msg:
        assert msg in result["msg"]


@pytest.mark.parametrize('pool, command_returns, raises, returns', [
    (
        "foo-bar",
        {  # vgs missing
            "msg": "[Errno 2] No such file or directory",
            "failed": True,
            "cmd": "/sbin/vgs",
            "rc": 2,
        },
        "Failed to run /sbin/vgs",
        None,
    ),
    (
        "foo",  # no hyphen in name - should not happen
        {},
        "name does not have the expected format",
        None,
    ),
    (
        "foo-bar",
        dict(stdout="  4.00g\n"),
        None,
        "4.00g",
    ),
    (
        "foo-bar",
        dict(stdout="\n"),  # no matching VG
        "vgs did not find this VG",
        None,
    )
])
def test_vg_free(pool, command_returns, raises, returns):
    def execute_module(module_name, *_):
        if module_name != "command":
            raise ValueError("not expecting module " + module_name)
        return command_returns

    check = DockerStorage(execute_module)
    if raises:
        with pytest.raises(OpenShiftCheckException) as err:
            check.get_vg_free(pool)
        assert raises in str(err.value)
    else:
        ret = check.get_vg_free(pool)
        assert ret == returns


@pytest.mark.parametrize('string, expect_bytes', [
    ("12", 12.0),
    ("12 k", 12.0 * 1024),
    ("42.42 MB", 42.42 * 1024**2),
    ("12g", 12.0 * 1024**3),
])
def test_convert_to_bytes(string, expect_bytes):
    got = DockerStorage.convert_to_bytes(string)
    assert got == expect_bytes


@pytest.mark.parametrize('string', [
    "bork",
    "42 Qs",
])
def test_convert_to_bytes_error(string):
    with pytest.raises(ValueError) as err:
        DockerStorage.convert_to_bytes(string)
    assert "Cannot convert" in str(err.value)
    assert string in str(err.value)


ansible_mounts_enough = [{
    'mount': '/var/lib/docker',
    'size_available': 50 * 10**9,
    'size_total': 50 * 10**9,
}]
ansible_mounts_not_enough = [{
    'mount': '/var/lib/docker',
    'size_available': 0,
    'size_total': 50 * 10**9,
}]
ansible_mounts_missing_fields = [dict(mount='/var/lib/docker')]
ansible_mounts_zero_size = [{
    'mount': '/var/lib/docker',
    'size_available': 0,
    'size_total': 0,
}]


@pytest.mark.parametrize('ansible_mounts, threshold, expect_fail, expect_msg', [
    (
        ansible_mounts_enough,
        None,
        False,
        [],
    ),
    (
        ansible_mounts_not_enough,
        None,
        True,
        ["usage percentage", "higher than threshold"],
    ),
    (
        ansible_mounts_not_enough,
        "bogus percent",
        True,
        ["is not a percentage"],
    ),
    (
        ansible_mounts_missing_fields,
        None,
        True,
        ["Ansible bug"],
    ),
    (
        ansible_mounts_zero_size,
        None,
        True,
        ["Ansible bug"],
    ),
])
def test_overlay_usage(ansible_mounts, threshold, expect_fail, expect_msg):
    task_vars = non_atomic_task_vars()
    task_vars["ansible_mounts"] = ansible_mounts
    if threshold is not None:
        task_vars["max_overlay_usage_percent"] = threshold
    check = DockerStorage(None, task_vars)
    docker_info = dict(DockerRootDir="/var/lib/docker", Driver="overlay")
    result = check.check_overlay_usage(docker_info)

    assert expect_fail == bool(result.get("failed"))
    for msg in expect_msg:
        assert msg in result["msg"]
