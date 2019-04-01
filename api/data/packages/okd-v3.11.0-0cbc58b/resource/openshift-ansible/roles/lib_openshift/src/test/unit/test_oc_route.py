'''
 Unit tests for oc route
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
# pylint: disable=import-error,wrong-import-position
# place class in our python path
module_path = os.path.join('/'.join(os.path.realpath(__file__).split('/')[:-4]), 'library')  # noqa: E501
sys.path.insert(0, module_path)
from oc_route import OCRoute, locate_oc_binary  # noqa: E402


class OCRouteTest(unittest.TestCase):
    '''
     Test class for OCRoute
    '''

    @mock.patch('oc_route.locate_oc_binary')
    @mock.patch('oc_route.Utils.create_tmpfile_copy')
    @mock.patch('oc_route.OCRoute._run')
    def test_list_route(self, mock_cmd, mock_tmpfile_copy, mock_oc_binary):
        ''' Testing getting a route '''

        # Arrange

        # run_ansible input parameters
        params = {
            'kubeconfig': '/etc/origin/master/admin.kubeconfig',
            'state': 'list',
            'debug': False,
            'name': 'test',
            'namespace': 'default',
            'labels': {'route': 'route'},
            'tls_termination': 'passthrough',
            'dest_cacert_path': None,
            'cacert_path': None,
            'cert_path': None,
            'key_path': None,
            'dest_cacert_content': None,
            'cacert_content': None,
            'cert_content': None,
            'key_content': None,
            'service_name': 'testservice',
            'host': 'test.openshift.com',
            'wildcard_policy': None,
            'weight': None,
            'port': None
        }

        route_result = '''{
            "kind": "Route",
            "apiVersion": "v1",
            "metadata": {
                "name": "test",
                "namespace": "default",
                "selfLink": "/oapi/v1/namespaces/default/routes/test",
                "uid": "1b127c67-ecd9-11e6-96eb-0e0d9bdacd26",
                "resourceVersion": "439182",
                "creationTimestamp": "2017-02-07T01:59:48Z",
                "labels": {
                    "route": "route"
                }
            },
            "spec": {
                "host": "test.example",
                "to": {
                    "kind": "Service",
                    "name": "test",
                    "weight": 100
                },
                "port": {
                    "targetPort": 8443
                },
                "tls": {
                    "termination": "passthrough"
                },
                "wildcardPolicy": "None"
            },
            "status": {
                "ingress": [
                    {
                        "host": "test.example",
                        "routerName": "router",
                        "conditions": [
                            {
                                "type": "Admitted",
                                "status": "True",
                                "lastTransitionTime": "2017-02-07T01:59:48Z"
                            }
                        ],
                        "wildcardPolicy": "None"
                    }
                ]
            }
        }'''

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            # First call to mock
            (0, route_result, ''),
        ]

        mock_oc_binary.side_effect = [
            'oc'
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mock.kubeconfig',
        ]

        # Act
        results = OCRoute.run_ansible(params, False)

        # Assert
        self.assertFalse(results['changed'])
        self.assertEqual(results['state'], 'list')
        self.assertEqual(results['results'][0]['metadata']['name'], 'test')

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'route', 'test', '-o', 'json', '-n', 'default'], None),
        ])

    @mock.patch('oc_route.locate_oc_binary')
    @mock.patch('oc_route.Utils.create_tmpfile_copy')
    @mock.patch('oc_route.Yedit._write')
    @mock.patch('oc_route.OCRoute._run')
    def test_create_route(self, mock_cmd, mock_write, mock_tmpfile_copy, mock_oc_binary):
        ''' Testing getting a route '''
        # Arrange

        # run_ansible input parameters
        params = {
            'kubeconfig': '/etc/origin/master/admin.kubeconfig',
            'state': 'present',
            'debug': False,
            'name': 'test',
            'namespace': 'default',
            'labels': {'route': 'route'},
            'tls_termination': 'edge',
            'dest_cacert_path': None,
            'cacert_path': None,
            'cert_path': None,
            'key_path': None,
            'dest_cacert_content': None,
            'cacert_content': 'testing',
            'cert_content': 'testing',
            'key_content': 'testing',
            'service_name': 'testservice',
            'host': 'test.openshift.com',
            'wildcard_policy': None,
            'weight': None,
            'port': None
        }

        route_result = '''{
                "apiVersion": "v1",
                "kind": "Route",
                "metadata": {
                    "creationTimestamp": "2017-02-07T20:55:10Z",
                    "name": "test",
                    "namespace": "default",
                    "resourceVersion": "517745",
                    "selfLink": "/oapi/v1/namespaces/default/routes/test",
                    "uid": "b6f25898-ed77-11e6-9755-0e737db1e63a",
                    "labels": {"route": "route"}
                },
                "spec": {
                    "host": "test.openshift.com",
                    "tls": {
                        "caCertificate": "testing",
                        "certificate": "testing",
                        "key": "testing",
                        "termination": "edge"
                    },
                    "to": {
                        "kind": "Service",
                        "name": "testservice",
                        "weight": 100
                    },
                    "wildcardPolicy": "None"
                },
                "status": {
                    "ingress": [
                        {
                            "conditions": [
                                {
                                    "lastTransitionTime": "2017-02-07T20:55:10Z",
                                    "status": "True",
                                    "type": "Admitted"
                                }
                            ],
                            "host": "test.openshift.com",
                            "routerName": "router",
                            "wildcardPolicy": "None"
                        }
                    ]
                }
            }'''

        test_route = '''\
kind: Route
spec:
  tls:
    caCertificate: testing
    termination: edge
    certificate: testing
    key: testing
  to:
    kind: Service
    name: testservice
    weight: 100
  host: test.openshift.com
  wildcardPolicy: None
apiVersion: v1
metadata:
  namespace: default
  name: test
'''

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            # First call to mock
            (1, '', 'Error from server: routes "test" not found'),
            (1, '', 'Error from server: routes "test" not found'),
            (0, 'route "test" created', ''),
            (0, route_result, ''),
        ]

        mock_oc_binary.side_effect = [
            'oc'
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mock.kubeconfig',
        ]

        mock_write.assert_has_calls = [
            # First call to mock
            mock.call('/tmp/test', test_route)
        ]

        # Act
        results = OCRoute.run_ansible(params, False)

        # Assert
        self.assertTrue(results['changed'])
        self.assertEqual(results['state'], 'present')
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'test')
        self.assertEqual(results['results']['results'][0]['metadata']['labels']['route'], 'route')

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'route', 'test', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'create', '-f', mock.ANY, '-n', 'default'], None),
            mock.call(['oc', 'get', 'route', 'test', '-o', 'json', '-n', 'default'], None),
        ])

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
