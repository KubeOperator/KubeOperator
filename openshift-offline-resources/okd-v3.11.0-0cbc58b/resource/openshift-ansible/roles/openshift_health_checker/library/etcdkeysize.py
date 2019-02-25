#!/usr/bin/python
"""Ansible module that recursively determines if the size of a key in an etcd cluster exceeds a given limit."""

from ansible.module_utils.basic import AnsibleModule


try:
    import etcd

    IMPORT_EXCEPTION_MSG = None
except ImportError as err:
    IMPORT_EXCEPTION_MSG = str(err)

    from collections import namedtuple
    EtcdMock = namedtuple("etcd", ["EtcdKeyNotFound"])
    etcd = EtcdMock(KeyError)


# pylint: disable=too-many-arguments
def check_etcd_key_size(client, key, size_limit, total_size=0, depth=0, depth_limit=1000, visited=None):
    """Check size of an etcd path starting at given key. Returns tuple (string, bool)"""
    if visited is None:
        visited = set()

    if key in visited:
        return 0, False

    visited.add(key)

    try:
        result = client.read(key, recursive=False)
    except etcd.EtcdKeyNotFound:
        return 0, False

    size = 0
    limit_exceeded = False

    for node in result.leaves:
        if depth >= depth_limit:
            raise Exception("Maximum recursive stack depth ({}) exceeded.".format(depth_limit))

        if size_limit and total_size + size > size_limit:
            return size, True

        if not node.dir:
            size += len(node.value)
            continue

        key_size, limit_exceeded = check_etcd_key_size(client, node.key,
                                                       size_limit,
                                                       total_size + size,
                                                       depth + 1,
                                                       depth_limit, visited)
        size += key_size

    max_limit_exceeded = limit_exceeded or (total_size + size > size_limit)
    return size, max_limit_exceeded


def main():  # pylint: disable=missing-docstring,too-many-branches
    module = AnsibleModule(
        argument_spec=dict(
            size_limit_bytes=dict(type="int", default=0),
            paths=dict(type="list", default=["/openshift.io/images"]),
            host=dict(type="str", default="127.0.0.1"),
            port=dict(type="int", default=4001),
            protocol=dict(type="str", default="http"),
            version_prefix=dict(type="str", default=""),
            allow_redirect=dict(type="bool", default=False),
            cert=dict(type="dict", default=""),
            ca_cert=dict(type="str", default=None),
        ),
        supports_check_mode=True
    )

    module.params["cert"] = (
        module.params["cert"]["cert"],
        module.params["cert"]["key"],
    )

    size_limit = module.params.pop("size_limit_bytes")
    paths = module.params.pop("paths")

    limit_exceeded = False

    try:
        # pylint: disable=no-member
        client = etcd.Client(**module.params)
    except AttributeError as attrerr:
        msg = str(attrerr)
        if IMPORT_EXCEPTION_MSG:
            msg = IMPORT_EXCEPTION_MSG
            if "No module named etcd" in IMPORT_EXCEPTION_MSG:
                # pylint: disable=redefined-variable-type
                msg = ('Unable to import the python "etcd" dependency. '
                       'Make sure python-etcd is installed on the host.')

        module.exit_json(
            failed=True,
            changed=False,
            size_limit_exceeded=limit_exceeded,
            msg=msg,
        )

        return

    size = 0
    for path in paths:
        path_size, limit_exceeded = check_etcd_key_size(client, path, size_limit - size)
        size += path_size

        if limit_exceeded:
            break

    module.exit_json(
        changed=False,
        size_limit_exceeded=limit_exceeded,
    )


if __name__ == '__main__':
    main()
