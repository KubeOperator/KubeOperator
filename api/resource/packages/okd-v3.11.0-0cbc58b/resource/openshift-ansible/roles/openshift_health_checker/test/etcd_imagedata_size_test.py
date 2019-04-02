import pytest

from collections import namedtuple
from openshift_checks.etcd_imagedata_size import EtcdImageDataSize
from openshift_checks import OpenShiftCheckException
from etcdkeysize import check_etcd_key_size


def fake_etcd_client(root):
    fake_nodes = dict()
    fake_etcd_node(root, fake_nodes)

    clientclass = namedtuple("client", ["read"])
    return clientclass(lambda key, recursive: fake_etcd_result(fake_nodes[key]))


def fake_etcd_result(fake_node):
    resultclass = namedtuple("result", ["leaves"])
    if not fake_node.dir:
        return resultclass([fake_node])

    return resultclass(fake_node.leaves)


def fake_etcd_node(node, visited):
    min_req_fields = ["dir", "key"]
    fields = list(node)
    leaves = []

    if node["dir"] and node.get("leaves"):
        for leaf in node["leaves"]:
            leaves.append(fake_etcd_node(leaf, visited))

    if len(set(min_req_fields) - set(fields)) > 0:
        raise ValueError("fake etcd nodes require at least {} fields.".format(min_req_fields))

    if node.get("leaves"):
        node["leaves"] = leaves

    nodeclass = namedtuple("node", fields)
    nodeinst = nodeclass(**node)
    visited[nodeinst.key] = nodeinst

    return nodeinst


@pytest.mark.parametrize('ansible_mounts,extra_words', [
    ([], ['none']),  # empty ansible_mounts
    ([{'mount': '/mnt'}], ['/mnt']),  # missing relevant mount paths
])
def test_cannot_determine_available_mountpath(ansible_mounts, extra_words):
    task_vars = dict(
        ansible_mounts=ansible_mounts,
    )
    check = EtcdImageDataSize(fake_execute_module, task_vars)

    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.run()

    for word in ['Unable to determine mount point'] + extra_words:
        assert word in str(excinfo.value)


@pytest.mark.parametrize('ansible_mounts,tree,size_limit,should_fail,extra_words', [
    (
        # test that default image size limit evals to 1/2 * (total size in use)
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 80 * 10**9,
        }],
        {"dir": False, "key": "/", "value": "1234"},
        None,
        False,
        [],
    ),
    (
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 48 * 10**9,
        }],
        {"dir": False, "key": "/", "value": "1234"},
        None,
        False,
        [],
    ),
    (
        # set max size limit for image data to be below total node value
        # total node value is defined as the sum of the value field
        # from every node
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 48 * 10**9,
        }],
        {"dir": False, "key": "/", "value": "12345678"},
        7,
        True,
        ["exceeds the maximum recommended limit", "0.00 GB"],
    ),
    (
        [{
            'mount': '/',
            'size_available': 48 * 10**9 - 1,
            'size_total': 48 * 10**9,
        }],
        {"dir": False, "key": "/", "value": "1234"},
        None,
        True,
        ["exceeds the maximum recommended limit", "0.00 GB"],
    )
])
def test_check_etcd_key_size_calculates_correct_limit(ansible_mounts, tree, size_limit, should_fail, extra_words):
    def execute_module(module_name, module_args, *_):
        if module_name != "etcdkeysize":
            return {
                "changed": False,
            }

        client = fake_etcd_client(tree)
        s, limit_exceeded = check_etcd_key_size(client, tree["key"], module_args["size_limit_bytes"])

        return {"size_limit_exceeded": limit_exceeded}

    task_vars = dict(
        etcd_max_image_data_size_bytes=size_limit,
        ansible_mounts=ansible_mounts,
        openshift=dict(
            common=dict(config_base="/var/lib/origin")
        ),
        openshift_master_etcd_hosts=["localhost"]
    )
    if size_limit is None:
        task_vars.pop("etcd_max_image_data_size_bytes")

    check = EtcdImageDataSize(execute_module, task_vars).run()

    if should_fail:
        assert check["failed"]

        for word in extra_words:
            assert word in check["msg"]
    else:
        assert not check.get("failed", False)


@pytest.mark.parametrize('ansible_mounts,tree,root_path,expected_size,extra_words', [
    (
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 80 * 10**9,
        }],
        # test recursive size check on tree with height > 1
        {
            "dir": True,
            "key": "/",
            "leaves": [
                {"dir": False, "key": "/foo1", "value": "1234"},
                {"dir": False, "key": "/foo2", "value": "1234"},
                {"dir": False, "key": "/foo3", "value": "1234"},
                {"dir": False, "key": "/foo4", "value": "1234"},
                {
                    "dir": True,
                    "key": "/foo5",
                    "leaves": [
                        {"dir": False, "key": "/foo/bar1", "value": "56789"},
                        {"dir": False, "key": "/foo/bar2", "value": "56789"},
                        {"dir": False, "key": "/foo/bar3", "value": "56789"},
                        {
                            "dir": True,
                            "key": "/foo/bar4",
                            "leaves": [
                                {"dir": False, "key": "/foo/bar/baz1", "value": "123"},
                                {"dir": False, "key": "/foo/bar/baz2", "value": "123"},
                            ]
                        },
                    ]
                },
            ]
        },
        "/",
        37,
        [],
    ),
    (
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 80 * 10**9,
        }],
        # test correct sub-tree size calculation
        {
            "dir": True,
            "key": "/",
            "leaves": [
                {"dir": False, "key": "/foo1", "value": "1234"},
                {"dir": False, "key": "/foo2", "value": "1234"},
                {"dir": False, "key": "/foo3", "value": "1234"},
                {"dir": False, "key": "/foo4", "value": "1234"},
                {
                    "dir": True,
                    "key": "/foo5",
                    "leaves": [
                        {"dir": False, "key": "/foo/bar1", "value": "56789"},
                        {"dir": False, "key": "/foo/bar2", "value": "56789"},
                        {"dir": False, "key": "/foo/bar3", "value": "56789"},
                        {
                            "dir": True,
                            "key": "/foo/bar4",
                            "leaves": [
                                {"dir": False, "key": "/foo/bar/baz1", "value": "123"},
                                {"dir": False, "key": "/foo/bar/baz2", "value": "123"},
                            ]
                        },
                    ]
                },
            ]
        },
        "/foo5",
        21,
        [],
    ),
    (
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 80 * 10**9,
        }],
        # test that a non-existing key is handled correctly
        {
            "dir": False,
            "key": "/",
            "value": "1234",
        },
        "/missing",
        0,
        [],
    ),
    (
        [{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 80 * 10**9,
        }],
        # test etcd cycle handling
        {
            "dir": True,
            "key": "/",
            "leaves": [
                {"dir": False, "key": "/foo1", "value": "1234"},
                {"dir": False, "key": "/foo2", "value": "1234"},
                {"dir": False, "key": "/foo3", "value": "1234"},
                {"dir": False, "key": "/foo4", "value": "1234"},
                {
                    "dir": True,
                    "key": "/",
                    "leaves": [
                        {"dir": False, "key": "/foo1", "value": "1"},
                    ],
                },
            ]
        },
        "/",
        16,
        [],
    ),
])
def test_etcd_key_size_check_calculates_correct_size(ansible_mounts, tree, root_path, expected_size, extra_words):
    def execute_module(module_name, module_args, *_):
        if module_name != "etcdkeysize":
            return {
                "changed": False,
            }

        client = fake_etcd_client(tree)
        size, limit_exceeded = check_etcd_key_size(client, root_path, module_args["size_limit_bytes"])

        assert size == expected_size
        return {
            "size_limit_exceeded": limit_exceeded,
        }

    task_vars = dict(
        ansible_mounts=ansible_mounts,
        openshift=dict(
            common=dict(config_base="/var/lib/origin")
        ),
        openshift_master_etcd_hosts=["localhost"]
    )

    check = EtcdImageDataSize(execute_module, task_vars).run()
    assert not check.get("failed", False)


def test_etcdkeysize_module_failure():
    def execute_module(module_name, *_):
        if module_name != "etcdkeysize":
            return {
                "changed": False,
            }

        return {
            "rc": 1,
            "module_stderr": "failure",
        }

    task_vars = dict(
        ansible_mounts=[{
            'mount': '/',
            'size_available': 40 * 10**9,
            'size_total': 80 * 10**9,
        }],
        openshift=dict(
            common=dict(config_base="/var/lib/origin")
        ),
        openshift_master_etcd_hosts=["localhost"]
    )

    check = EtcdImageDataSize(execute_module, task_vars).run()

    assert check["failed"]
    for word in "Failed to retrieve stats":
        assert word in check["msg"]


def fake_execute_module(*args):
    raise AssertionError('this function should not be called')
