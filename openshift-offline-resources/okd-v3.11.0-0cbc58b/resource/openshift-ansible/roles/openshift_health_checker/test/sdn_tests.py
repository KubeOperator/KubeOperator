import pytest
from openshift_checks.sdn import SDNCheck
from openshift_checks import OpenShiftCheckException


def fake_execute_module(*args):
    raise AssertionError('this function should not be called')


def test_check_nodes_missing_node_name():
    task_vars = dict(
        group_names=['oo_nodes_to_config'],
    )

    check = SDNCheck(fake_execute_module, task_vars)
    check.run()

    assert 1 == len(check.failures)
    assert 'Could not determine node name' in str(check.failures[0])


def test_check_master():
    nodes = [
        {
            'apiVersion': 'v1',
            'kind': 'Node',
            'metadata': {
                'annotations': {'kubernetes.io/hostname': 'node1'},
                'name': 'ip-172-31-50-1.ec2.internal'
            },
            'status': {
                'addresses': [
                    {'address': '172.31.50.1', 'type': 'InternalIP'},
                    {'address': '52.0.0.1', 'type': 'ExternalIP'},
                    {
                        'address': 'ip-172-31-50-1.ec2.internal',
                        'type': 'Hostname'
                    }
                ]
            }
        },
        {
            'apiVersion': 'v1',
            'kind': 'Node',
            'metadata': {'name': 'ip-172-31-50-2.ec2.internal'},
            'status': {
                'addresses': [
                    {'address': '172.31.50.2', 'type': 'InternalIP'},
                    {'address': '52.0.0.2', 'type': 'ExternalIP'},
                    {
                        'address': 'ip-172-31-50-2.ec2.internal',
                        'type': 'Hostname'
                    }
                ]
            }
        }
    ]

    task_vars = dict(
        group_names=['oo_masters_to_config'],
        resources=dict(results=[
            dict(item='nodes', results=dict(results=[dict(items=nodes)])),
            dict(item='pods', results=dict(results=[dict(items={})])),
            dict(item='services', results=dict(results=[dict(items={})]))
        ])
    )

    node_addresses = {
        node['metadata']['name']: {
            address['type']: address['address']
            for address
            in node['status']['addresses']
        }
        for node in nodes
    }
    expected_hostnames = [addresses['Hostname']
                          for addresses in node_addresses.values()]
    uri_hostnames = []
    resolve_address_hostnames = []

    def execute_module(module_name, args, *_):
        if module_name == 'uri':
            for hostname in expected_hostnames:
                if hostname in args['url']:
                    uri_hostnames.append(hostname)
                    return {}
            raise ValueError('unexpected url: %s' % args['url'])
        raise ValueError('not expecting module %s' % module_name)

    def resolve_address(address):
        for hostname in expected_hostnames:
            if address == hostname:
                resolve_address_hostnames.append(hostname)
                return node_addresses[hostname]['InternalIP']
        raise ValueError('unexpected address: %s' % hostname)

    check = SDNCheck(execute_module, task_vars)
    check.resolve_address = resolve_address
    check.run()

    assert 0 == len(check.failures)
    assert set(expected_hostnames) == set(uri_hostnames), 'should try to connect to the kubelet'
    assert set(expected_hostnames) == set(resolve_address_hostnames), 'should try to resolve the node\'s address'


def test_check_nodes():
    nodes = [
        {
            'apiVersion': 'v1',
            'kind': 'Node',
            'metadata': {
                'annotations': {'kubernetes.io/hostname': 'node1'},
                'name': 'ip-172-31-50-1.ec2.internal'
            },
            'status': {
                'addresses': [
                    {'address': '172.31.50.1', 'type': 'InternalIP'},
                    {'address': '52.0.0.1', 'type': 'ExternalIP'},
                    {
                        'address': 'ip-172-31-50-1.ec2.internal',
                        'type': 'Hostname'
                    }
                ]
            }
        },
        {
            'apiVersion': 'v1',
            'kind': 'Node',
            'metadata': {'name': 'ip-172-31-50-2.ec2.internal'},
            'status': {
                'addresses': [
                    {'address': '172.31.50.2', 'type': 'InternalIP'},
                    {'address': '52.0.0.2', 'type': 'ExternalIP'},
                    {
                        'address': 'ip-172-31-50-2.ec2.internal',
                        'type': 'Hostname'
                    }
                ]
            }
        }
    ]
    hostsubnets = [
        {
            'metadata': {
                'name': 'ip-172-31-50-1.ec2.internal'
            },
            'subnet': '10.128.0.1/23'
        },
        {
            'metadata': {
                'name': 'ip-172-31-50-2.ec2.internal'
            },
            'subnet': '10.129.0.1/23'
        }
    ]

    task_vars = dict(
        group_names=['oo_nodes_to_config'],
        resources=dict(results=[
            dict(item='nodes', results=dict(results=[dict(items=nodes)])),
            dict(item='hostsubnets', results=dict(results=[dict(items=hostsubnets)]))
        ]),
        openshift=dict(node=dict(nodename='foo'))
    )

    def execute_module(module_name, args, *_):
        if module_name == 'command':
            return dict(stdout='bogus_container_id')
        raise ValueError('not expecting module %s' % module_name)

    SDNCheck(execute_module, task_vars).run()


def test_resolve_address():
    def execute_module(module_name, args, *_):
        if module_name != 'command':
            raise ValueError('not expecting module %s' % module_name)

        command_args = args['_raw_params'].split()
        if command_args[0] != '/bin/getent':
            raise ValueError('not expecting command: %s' % args.raw_params)

        # The expected command_args is ['/bin/getent', 'ahostsv4', 'foo'].
        if command_args[2] == 'foo':
            return {
                'rc': 0,
                'stdout': '''1.2.3.4         STREAM bar
1.2.3.4         DGRAM
1.2.3.4         RAW
'''
            }

        return {'rc': 2}

    check = SDNCheck(execute_module, None)
    assert check.resolve_address('foo') == '1.2.3.4'
    with pytest.raises(OpenShiftCheckException):
        check.resolve_address('baz')


def test_no_nodes():
    task_vars = dict(
        group_names=['oo_masters_to_config'],
        resources=dict(results=[
            dict(item='nodes', results=dict(results=[dict(items={})])),
            dict(item='pods', results=dict(results=[dict(items={})])),
            dict(item='services', results=dict(results=[dict(items={})]))
        ])
    )

    check = SDNCheck(fake_execute_module, task_vars)
    check.run()
    assert 1 == len(check.failures)
    assert 'No nodes' in str(check.failures[0])


@pytest.mark.parametrize('group_names,expected', [
    (['oo_masters_to_config'], True),
    (['oo_nodes_to_config'], True),
    (['oo_masters_to_config', 'oo_nodes_to_config'], True),
    (['oo_masters_to_config', 'oo_etcd_to_config'], True),
    ([], False),
    (['oo_etcd_to_config'], False),
    (['lb'], False),
    (['nfs'], False),
])
def test_sdn_skip_when_not_master_nor_node(group_names, expected):
    task_vars = dict(
        group_names=group_names,
        openshift_is_atomic=True,
    )
    assert SDNCheck(None, task_vars).is_active() == expected
