'''
 Unit tests for oc service
'''

import os
import six
import sys
import unittest
import mock

# Removing invalid variable names for tests so that I can
# keep them brief
# pylint: disable=invalid-name,no-name-in-module
# Disable import-error b/c our libraries aren't loaded in jenkins
# pylint: disable=import-error
# place class in our python path
module_path = os.path.join('/'.join(os.path.realpath(__file__).split('/')[:-4]), 'library')  # noqa: E501
sys.path.insert(0, module_path)
from oc_service import OCService, locate_oc_binary  # noqa: E402


# pylint: disable=too-many-public-methods
class OCServiceTest(unittest.TestCase):
    '''
     Test class for OCService
    '''

    @mock.patch('oc_service.Utils.create_tmpfile_copy')
    @mock.patch('oc_service.OCService._run')
    def test_state_list(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a get '''
        params = {'name': 'router',
                  'namespace': 'default',
                  'ports': None,
                  'state': 'list',
                  'labels': None,
                  'annotations': None,
                  'clusterip': None,
                  'portalip': None,
                  'selector': None,
                  'session_affinity': None,
                  'service_type': None,
                  'external_ips': None,
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'debug': False}

        service = '''{
            "kind": "Service",
            "apiVersion": "v1",
            "metadata": {
                "name": "router",
                "namespace": "default",
                "selfLink": "/api/v1/namespaces/default/services/router",
                "uid": "fabd2440-e3d8-11e6-951c-0e3dd518cefa",
                "resourceVersion": "3206",
                "creationTimestamp": "2017-01-26T15:06:14Z",
                "labels": {
                    "router": "router"
                }
            },
            "spec": {
                "ports": [
                    {
                        "name": "80-tcp",
                        "protocol": "TCP",
                        "port": 80,
                        "targetPort": 80
                    },
                    {
                        "name": "443-tcp",
                        "protocol": "TCP",
                        "port": 443,
                        "targetPort": 443
                    },
                    {
                        "name": "1936-tcp",
                        "protocol": "TCP",
                        "port": 1936,
                        "targetPort": 1936
                    },
                    {
                        "name": "5000-tcp",
                        "protocol": "TCP",
                        "port": 5000,
                        "targetPort": 5000
                    }
                ],
                "selector": {
                    "router": "router"
                },
                "clusterIP": "172.30.129.161",
                "type": "ClusterIP",
                "sessionAffinity": "None"
            },
            "status": {
                "loadBalancer": {}
            }
        }'''
        mock_cmd.side_effect = [
            (0, service, '')
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCService.run_ansible(params, False)

        self.assertFalse(results['changed'])
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'router')

    @mock.patch('oc_service.Utils.create_tmpfile_copy')
    @mock.patch('oc_service.OCService._run')
    def test_create(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a create service '''
        params = {'name': 'router',
                  'namespace': 'default',
                  'ports': {'name': '9000-tcp',
                            'port': 9000,
                            'protocol': 'TCP',
                            'targetPOrt': 9000},
                  'state': 'present',
                  'labels': None,
                  'annotations': None,
                  'clusterip': None,
                  'portalip': None,
                  'selector': {'router': 'router'},
                  'session_affinity': 'ClientIP',
                  'service_type': 'ClusterIP',
                  'external_ips': None,
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'debug': False}

        service = '''{
            "kind": "Service",
            "apiVersion": "v1",
            "metadata": {
                "name": "router",
                "namespace": "default",
                "selfLink": "/api/v1/namespaces/default/services/router",
                "uid": "fabd2440-e3d8-11e6-951c-0e3dd518cefa",
                "resourceVersion": "3206",
                "creationTimestamp": "2017-01-26T15:06:14Z",
                "labels": {
                    "router": "router"
                }
            },
            "spec": {
                "ports": [
                    {
                        "name": "80-tcp",
                        "protocol": "TCP",
                        "port": 80,
                        "targetPort": 80
                    },
                    {
                        "name": "443-tcp",
                        "protocol": "TCP",
                        "port": 443,
                        "targetPort": 443
                    },
                    {
                        "name": "1936-tcp",
                        "protocol": "TCP",
                        "port": 1936,
                        "targetPort": 1936
                    },
                    {
                        "name": "5000-tcp",
                        "protocol": "TCP",
                        "port": 5000,
                        "targetPort": 5000
                    }
                ],
                "selector": {
                    "router": "router"
                },
                "clusterIP": "172.30.129.161",
                "type": "ClusterIP",
                "sessionAffinity": "None"
            },
            "status": {
                "loadBalancer": {}
            }
        }'''
        mock_cmd.side_effect = [
            (1, '', 'Error from server: services "router" not found'),
            (1, '', 'Error from server: services "router" not found'),
            (0, service, ''),
            (0, service, '')
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCService.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertTrue(results['results']['returncode'] == 0)
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'router')

    @unittest.skipIf(six.PY3, 'py2 test only')
    @mock.patch('os.path.exists')
    @mock.patch('os.environ.get')
    def test_binary_lookup_fallback(self, mock_env_get, mock_path_exists):
        ''' Testing binary lookup fallback '''

        mock_env_get.side_effect = lambda _v, _d: ''

        mock_path_exists.side_effect = lambda _: False

        self.assertEqual(locate_oc_binary(), 'oc')

    @unittest.skipIf(six.PY3, 'py2 test only')
    @mock.patch('os.path.exists')
    @mock.patch('os.environ.get')
    def test_binary_lookup_in_path(self, mock_env_get, mock_path_exists):
        ''' Testing binary lookup in path '''

        oc_bin = '/usr/bin/oc'

        mock_env_get.side_effect = lambda _v, _d: '/bin:/usr/bin'

        mock_path_exists.side_effect = lambda f: f == oc_bin

        self.assertEqual(locate_oc_binary(), oc_bin)

    @unittest.skipIf(six.PY3, 'py2 test only')
    @mock.patch('os.path.exists')
    @mock.patch('os.environ.get')
    def test_binary_lookup_in_usr_local(self, mock_env_get, mock_path_exists):
        ''' Testing binary lookup in /usr/local/bin '''

        oc_bin = '/usr/local/bin/oc'

        mock_env_get.side_effect = lambda _v, _d: '/bin:/usr/bin'

        mock_path_exists.side_effect = lambda f: f == oc_bin

        self.assertEqual(locate_oc_binary(), oc_bin)

    @unittest.skipIf(six.PY3, 'py2 test only')
    @mock.patch('os.path.exists')
    @mock.patch('os.environ.get')
    def test_binary_lookup_in_home(self, mock_env_get, mock_path_exists):
        ''' Testing binary lookup in ~/bin '''

        oc_bin = os.path.expanduser('~/bin/oc')

        mock_env_get.side_effect = lambda _v, _d: '/bin:/usr/bin'

        mock_path_exists.side_effect = lambda f: f == oc_bin

        self.assertEqual(locate_oc_binary(), oc_bin)

    @unittest.skipIf(six.PY2, 'py3 test only')
    @mock.patch('shutil.which')
    @mock.patch('os.environ.get')
    def test_binary_lookup_fallback_py3(self, mock_env_get, mock_shutil_which):
        ''' Testing binary lookup fallback '''

        mock_env_get.side_effect = lambda _v, _d: ''

        mock_shutil_which.side_effect = lambda _f, path=None: None

        self.assertEqual(locate_oc_binary(), 'oc')

    @unittest.skipIf(six.PY2, 'py3 test only')
    @mock.patch('shutil.which')
    @mock.patch('os.environ.get')
    def test_binary_lookup_in_path_py3(self, mock_env_get, mock_shutil_which):
        ''' Testing binary lookup in path '''

        oc_bin = '/usr/bin/oc'

        mock_env_get.side_effect = lambda _v, _d: '/bin:/usr/bin'

        mock_shutil_which.side_effect = lambda _f, path=None: oc_bin

        self.assertEqual(locate_oc_binary(), oc_bin)

    @unittest.skipIf(six.PY2, 'py3 test only')
    @mock.patch('shutil.which')
    @mock.patch('os.environ.get')
    def test_binary_lookup_in_usr_local_py3(self, mock_env_get, mock_shutil_which):
        ''' Testing binary lookup in /usr/local/bin '''

        oc_bin = '/usr/local/bin/oc'

        mock_env_get.side_effect = lambda _v, _d: '/bin:/usr/bin'

        mock_shutil_which.side_effect = lambda _f, path=None: oc_bin

        self.assertEqual(locate_oc_binary(), oc_bin)

    @unittest.skipIf(six.PY2, 'py3 test only')
    @mock.patch('shutil.which')
    @mock.patch('os.environ.get')
    def test_binary_lookup_in_home_py3(self, mock_env_get, mock_shutil_which):
        ''' Testing binary lookup in ~/bin '''

        oc_bin = os.path.expanduser('~/bin/oc')

        mock_env_get.side_effect = lambda _v, _d: '/bin:/usr/bin'

        mock_shutil_which.side_effect = lambda _f, path=None: oc_bin

        self.assertEqual(locate_oc_binary(), oc_bin)

    @mock.patch('oc_service.Utils.create_tmpfile_copy')
    @mock.patch('oc_service.OCService._run')
    def test_create_with_labels(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a create service '''
        params = {'name': 'router',
                  'namespace': 'default',
                  'ports': {'name': '9000-tcp',
                            'port': 9000,
                            'protocol': 'TCP',
                            'targetPOrt': 9000},
                  'state': 'present',
                  'labels': {'component': 'some_component', 'infra': 'true'},
                  'annotations': None,
                  'clusterip': None,
                  'portalip': None,
                  'selector': {'router': 'router'},
                  'session_affinity': 'ClientIP',
                  'service_type': 'ClusterIP',
                  'external_ips': None,
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'debug': False}

        service = '''{
            "kind": "Service",
            "apiVersion": "v1",
            "metadata": {
                "name": "router",
                "namespace": "default",
                "selfLink": "/api/v1/namespaces/default/services/router",
                "uid": "fabd2440-e3d8-11e6-951c-0e3dd518cefa",
                "resourceVersion": "3206",
                "creationTimestamp": "2017-01-26T15:06:14Z",
                "labels": {"component": "some_component", "infra": "true"}
            },
            "spec": {
                "ports": [
                    {
                        "name": "80-tcp",
                        "protocol": "TCP",
                        "port": 80,
                        "targetPort": 80
                    },
                    {
                        "name": "443-tcp",
                        "protocol": "TCP",
                        "port": 443,
                        "targetPort": 443
                    },
                    {
                        "name": "1936-tcp",
                        "protocol": "TCP",
                        "port": 1936,
                        "targetPort": 1936
                    },
                    {
                        "name": "5000-tcp",
                        "protocol": "TCP",
                        "port": 5000,
                        "targetPort": 5000
                    }
                ],
                "selector": {
                    "router": "router"
                },
                "clusterIP": "172.30.129.161",
                "type": "ClusterIP",
                "sessionAffinity": "None"
            },
            "status": {
                "loadBalancer": {}
            }
        }'''
        mock_cmd.side_effect = [
            (1, '', 'Error from server: services "router" not found'),
            (1, '', 'Error from server: services "router" not found'),
            (0, service, ''),
            (0, service, '')
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCService.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertTrue(results['results']['returncode'] == 0)
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'router')
        self.assertEqual(results['results']['results'][0]['metadata']['labels'], {"component": "some_component", "infra": "true"})

    @mock.patch('oc_service.Utils.create_tmpfile_copy')
    @mock.patch('oc_service.OCService._run')
    def test_create_with_external_ips(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a create service '''
        params = {'name': 'router',
                  'namespace': 'default',
                  'ports': {'name': '9000-tcp',
                            'port': 9000,
                            'protocol': 'TCP',
                            'targetPOrt': 9000},
                  'state': 'present',
                  'labels': {'component': 'some_component', 'infra': 'true'},
                  'annotations': None,
                  'clusterip': None,
                  'portalip': None,
                  'selector': {'router': 'router'},
                  'session_affinity': 'ClientIP',
                  'service_type': 'ClusterIP',
                  'external_ips': ['1.2.3.4', '5.6.7.8'],
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'debug': False}

        service = '''{
            "kind": "Service",
            "apiVersion": "v1",
            "metadata": {
                "name": "router",
                "namespace": "default",
                "selfLink": "/api/v1/namespaces/default/services/router",
                "uid": "fabd2440-e3d8-11e6-951c-0e3dd518cefa",
                "resourceVersion": "3206",
                "creationTimestamp": "2017-01-26T15:06:14Z",
                "labels": {"component": "some_component", "infra": "true"}
            },
            "spec": {
                "ports": [
                    {
                        "name": "80-tcp",
                        "protocol": "TCP",
                        "port": 80,
                        "targetPort": 80
                    },
                    {
                        "name": "443-tcp",
                        "protocol": "TCP",
                        "port": 443,
                        "targetPort": 443
                    },
                    {
                        "name": "1936-tcp",
                        "protocol": "TCP",
                        "port": 1936,
                        "targetPort": 1936
                    },
                    {
                        "name": "5000-tcp",
                        "protocol": "TCP",
                        "port": 5000,
                        "targetPort": 5000
                    }
                ],
                "selector": {
                    "router": "router"
                },
                "clusterIP": "172.30.129.161",
                "externalIPs": ["1.2.3.4", "5.6.7.8"],
                "type": "ClusterIP",
                "sessionAffinity": "None"
            },
            "status": {
                "loadBalancer": {}
            }
        }'''
        mock_cmd.side_effect = [
            (1, '', 'Error from server: services "router" not found'),
            (1, '', 'Error from server: services "router" not found'),
            (0, service, ''),
            (0, service, '')
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCService.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertTrue(results['results']['returncode'] == 0)
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'router')
        self.assertEqual(results['results']['results'][0]['metadata']['labels'], {"component": "some_component", "infra": "true"})
        self.assertEqual(results['results']['results'][0]['spec']['externalIPs'], ["1.2.3.4", "5.6.7.8"])
