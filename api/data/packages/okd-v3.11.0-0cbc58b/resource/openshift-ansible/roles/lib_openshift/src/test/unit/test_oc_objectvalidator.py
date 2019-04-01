'''
 Unit tests for oc_objectvalidator
'''

import os
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
from oc_objectvalidator import OCObjectValidator  # noqa: E402


class OCObjectValidatorTest(unittest.TestCase):
    '''
     Test class for OCObjectValidator
    '''

    maxDiff = None

    @mock.patch('oc_objectvalidator.locate_oc_binary')
    @mock.patch('oc_objectvalidator.Utils.create_tmpfile_copy')
    @mock.patch('oc_objectvalidator.OCObjectValidator._run')
    def test_no_data(self, mock_cmd, mock_tmpfile_copy, mock_oc_binary):
        ''' Testing when both all objects are empty '''

        # Arrange

        # run_ansible input parameters
        params = {
            'kubeconfig': '/etc/origin/master/admin.kubeconfig',
        }

        empty = '''{
    "apiVersion": "v1",
    "items": [],
    "kind": "List",
    "metadata": {},
    "resourceVersion": "",
    "selfLink": ""
}'''

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            # First call to mock
            (0, empty, ''),

            # Second call to mock
            (0, empty, ''),

            # Third call to mock
            (0, empty, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        mock_oc_binary.side_effect = [
            'oc',
        ]

        # Act
        results = OCObjectValidator.run_ansible(params)

        # Assert
        self.assertNotIn('failed', results)
        self.assertEqual(results['msg'], 'All objects are valid.')

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'hostsubnet', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'netnamespace', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'namespace', '-o', 'json', '-n', 'default'], None),
        ])

    @mock.patch('oc_objectvalidator.locate_oc_binary')
    @mock.patch('oc_objectvalidator.Utils.create_tmpfile_copy')
    @mock.patch('oc_objectvalidator.OCObjectValidator._run')
    def test_error_code(self, mock_cmd, mock_tmpfile_copy, mock_oc_binary):
        ''' Testing when we fail to get objects '''

        # Arrange

        # run_ansible input parameters
        params = {
            'kubeconfig': '/etc/origin/master/admin.kubeconfig',
        }

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            # First call to mock
            (1, '', 'Error.'),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        mock_oc_binary.side_effect = [
            'oc'
        ]

        error_results = {
            'returncode': 1,
            'stderr': 'Error.',
            'stdout': '',
            'cmd': 'oc get hostsubnet -o json -n default',
            'results': [{}]
        }

        # Act
        results = OCObjectValidator.run_ansible(params)

        # Assert
        self.assertTrue(results['failed'])
        self.assertEqual(results['msg'], 'Failed to GET hostsubnet.')
        self.assertEqual(results['state'], 'list')
        self.assertEqual(results['results'], error_results)

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'hostsubnet', '-o', 'json', '-n', 'default'], None),
        ])

    @mock.patch('oc_objectvalidator.locate_oc_binary')
    @mock.patch('oc_objectvalidator.Utils.create_tmpfile_copy')
    @mock.patch('oc_objectvalidator.OCObjectValidator._run')
    def test_valid_both(self, mock_cmd, mock_tmpfile_copy, mock_oc_binary):
        ''' Testing when both all objects are valid '''

        # Arrange

        # run_ansible input parameters
        params = {
            'kubeconfig': '/etc/origin/master/admin.kubeconfig',
        }

        valid_hostsubnet = '''{
    "apiVersion": "v1",
    "items": [
        {
            "apiVersion": "v1",
            "host": "bar0",
            "hostIP": "1.1.1.1",
            "kind": "HostSubnet",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:47:09Z",
                "name": "bar0",
                "namespace": "",
                "resourceVersion": "986",
                "selfLink": "/oapi/v1/hostsubnetsbar0",
                "uid": "528dbb41-f478-11e6-aae0-507b9dac97ff"
            },
            "subnet": "1.1.0.0/24"
        },
        {
            "apiVersion": "v1",
            "host": "bar1",
            "hostIP": "1.1.1.1",
            "kind": "HostSubnet",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:47:18Z",
                "name": "bar1",
                "namespace": "",
                "resourceVersion": "988",
                "selfLink": "/oapi/v1/hostsubnetsbar1",
                "uid": "57710d84-f478-11e6-aae0-507b9dac97ff"
            },
            "subnet": "1.1.0.0/24"
        },
        {
            "apiVersion": "v1",
            "host": "bar2",
            "hostIP": "1.1.1.1",
            "kind": "HostSubnet",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:47:26Z",
                "name": "bar2",
                "namespace": "",
                "resourceVersion": "991",
                "selfLink": "/oapi/v1/hostsubnetsbar2",
                "uid": "5c59a28c-f478-11e6-aae0-507b9dac97ff"
            },
            "subnet": "1.1.0.0/24"
        }
    ],
    "kind": "List",
    "metadata": {},
    "resourceVersion": "",
    "selfLink": ""
    }'''

        valid_netnamespace = '''{
    "apiVersion": "v1",
    "items": [
        {
            "apiVersion": "v1",
            "kind": "NetNamespace",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:45:16Z",
                "name": "foo0",
                "namespace": "",
                "resourceVersion": "959",
                "selfLink": "/oapi/v1/netnamespacesfoo0",
                "uid": "0f1c85b2-f478-11e6-aae0-507b9dac97ff"
            },
            "netid": 100,
            "netname": "foo0"
        },
        {
            "apiVersion": "v1",
            "kind": "NetNamespace",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:45:26Z",
                "name": "foo1",
                "namespace": "",
                "resourceVersion": "962",
                "selfLink": "/oapi/v1/netnamespacesfoo1",
                "uid": "14effa0d-f478-11e6-aae0-507b9dac97ff"
            },
            "netid": 100,
            "netname": "foo1"
        },
        {
            "apiVersion": "v1",
            "kind": "NetNamespace",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:45:36Z",
                "name": "foo2",
                "namespace": "",
                "resourceVersion": "965",
                "selfLink": "/oapi/v1/netnamespacesfoo2",
                "uid": "1aabdf84-f478-11e6-aae0-507b9dac97ff"
            },
            "netid": 100,
            "netname": "foo2"
        }
    ],
    "kind": "List",
    "metadata": {},
    "resourceVersion": "",
    "selfLink": ""
    }'''

        valid_namespace = '''{
    "apiVersion": "v1",
    "items": [
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/sa.scc.mcs": "s0:c1,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000000000/10000",
                    "openshift.io/sa.scc.uid-range": "1000000000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:49Z",
                "name": "default",
                "namespace": "",
                "resourceVersion": "165",
                "selfLink": "/api/v1/namespacesdefault",
                "uid": "23c0c6aa-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/sa.scc.mcs": "s0:c3,c2",
                    "openshift.io/sa.scc.supplemental-groups": "1000010000/10000",
                    "openshift.io/sa.scc.uid-range": "1000010000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:49Z",
                "name": "kube-system",
                "namespace": "",
                "resourceVersion": "533",
                "selfLink": "/api/v1/namespaceskube-system",
                "uid": "23c21758-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/description": "",
                    "openshift.io/display-name": "",
                    "openshift.io/requester": "developer",
                    "openshift.io/sa.scc.mcs": "s0:c9,c4",
                    "openshift.io/sa.scc.supplemental-groups": "1000080000/10000",
                    "openshift.io/sa.scc.uid-range": "1000080000/10000"
                },
                "creationTimestamp": "2017-03-02T02:17:16Z",
                "name": "myproject",
                "namespace": "",
                "resourceVersion": "2898",
                "selfLink": "/api/v1/namespacesmyproject",
                "uid": "5ae3764d-feee-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "openshift.io/origin",
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/sa.scc.mcs": "s0:c6,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000030000/10000",
                    "openshift.io/sa.scc.uid-range": "1000030000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:51Z",
                "name": "openshift",
                "namespace": "",
                "resourceVersion": "171",
                "selfLink": "/api/v1/namespacesopenshift",
                "uid": "24f7b34d-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/sa.scc.mcs": "s0:c5,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000020000/10000",
                    "openshift.io/sa.scc.uid-range": "1000020000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:51Z",
                "name": "openshift-infra",
                "namespace": "",
                "resourceVersion": "169",
                "selfLink": "/api/v1/namespacesopenshift-infra",
                "uid": "24a2ed75-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/description": "",
                    "openshift.io/display-name": "",
                    "openshift.io/requester": "developer1",
                    "openshift.io/sa.scc.mcs": "s0:c10,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000090000/10000",
                    "openshift.io/sa.scc.uid-range": "1000090000/10000"
                },
                "creationTimestamp": "2017-03-02T02:17:56Z",
                "name": "yourproject",
                "namespace": "",
                "resourceVersion": "2955",
                "selfLink": "/api/v1/namespacesyourproject",
                "uid": "72df7fb9-feee-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "openshift.io/origin",
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        }
    ],
    "kind": "List",
    "metadata": {},
    "resourceVersion": "",
    "selfLink": ""
}'''

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            # First call to mock
            (0, valid_hostsubnet, ''),

            # Second call to mock
            (0, valid_netnamespace, ''),

            # Third call to mock
            (0, valid_namespace, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        mock_oc_binary.side_effect = [
            'oc'
        ]

        # Act
        results = OCObjectValidator.run_ansible(params)

        # Assert
        self.assertNotIn('failed', results)
        self.assertEqual(results['msg'], 'All objects are valid.')

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'hostsubnet', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'netnamespace', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'namespace', '-o', 'json', '-n', 'default'], None),
        ])

    @mock.patch('oc_objectvalidator.locate_oc_binary')
    @mock.patch('oc_objectvalidator.Utils.create_tmpfile_copy')
    @mock.patch('oc_objectvalidator.OCObjectValidator._run')
    def test_invalid_both(self, mock_cmd, mock_tmpfile_copy, mock_oc_binary):
        ''' Testing when all objects are invalid '''

        # Arrange

        # run_ansible input parameters
        params = {
            'kubeconfig': '/etc/origin/master/admin.kubeconfig',
        }

        invalid_hostsubnet = '''{
    "apiVersion": "v1",
    "items": [
        {
            "apiVersion": "v1",
            "host": "bar0",
            "hostIP": "1.1.1.1",
            "kind": "HostSubnet",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:47:09Z",
                "name": "bar0",
                "namespace": "",
                "resourceVersion": "986",
                "selfLink": "/oapi/v1/hostsubnetsbar0",
                "uid": "528dbb41-f478-11e6-aae0-507b9dac97ff"
            },
            "subnet": "1.1.0.0/24"
        },
        {
            "apiVersion": "v1",
            "host": "bar1",
            "hostIP": "1.1.1.1",
            "kind": "HostSubnet",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:47:18Z",
                "name": "bar1",
                "namespace": "",
                "resourceVersion": "988",
                "selfLink": "/oapi/v1/hostsubnetsbar1",
                "uid": "57710d84-f478-11e6-aae0-507b9dac97ff"
            },
            "subnet": "1.1.0.0/24"
        },
        {
            "apiVersion": "v1",
            "host": "bar2",
            "hostIP": "1.1.1.1",
            "kind": "HostSubnet",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:47:26Z",
                "name": "bar2",
                "namespace": "",
                "resourceVersion": "991",
                "selfLink": "/oapi/v1/hostsubnetsbar2",
                "uid": "5c59a28c-f478-11e6-aae0-507b9dac97ff"
            },
            "subnet": "1.1.0.0/24"
        },
        {
            "apiVersion": "v1",
            "host": "baz1",
            "hostIP": "1.1.1.1",
            "kind": "HostSubnet",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:47:49Z",
                "name": "baz0",
                "namespace": "",
                "resourceVersion": "996",
                "selfLink": "/oapi/v1/hostsubnetsbaz0",
                "uid": "69f75f87-f478-11e6-aae0-507b9dac97ff"
            },
            "subnet": "1.1.0.0/24"
        }
    ],
    "kind": "List",
    "metadata": {},
    "resourceVersion": "",
    "selfLink": ""
}'''

        invalid_netnamespace = '''{
    "apiVersion": "v1",
    "items": [
        {
            "apiVersion": "v1",
            "kind": "NetNamespace",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:45:52Z",
                "name": "bar0",
                "namespace": "",
                "resourceVersion": "969",
                "selfLink": "/oapi/v1/netnamespacesbar0",
                "uid": "245d416e-f478-11e6-aae0-507b9dac97ff"
            },
            "netid": 100,
            "netname": "bar1"
        },
        {
            "apiVersion": "v1",
            "kind": "NetNamespace",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:45:16Z",
                "name": "foo0",
                "namespace": "",
                "resourceVersion": "959",
                "selfLink": "/oapi/v1/netnamespacesfoo0",
                "uid": "0f1c85b2-f478-11e6-aae0-507b9dac97ff"
            },
            "netid": 100,
            "netname": "foo0"
        },
        {
            "apiVersion": "v1",
            "kind": "NetNamespace",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:45:26Z",
                "name": "foo1",
                "namespace": "",
                "resourceVersion": "962",
                "selfLink": "/oapi/v1/netnamespacesfoo1",
                "uid": "14effa0d-f478-11e6-aae0-507b9dac97ff"
            },
            "netid": 100,
            "netname": "foo1"
        },
        {
            "apiVersion": "v1",
            "kind": "NetNamespace",
            "metadata": {
                "creationTimestamp": "2017-02-16T18:45:36Z",
                "name": "foo2",
                "namespace": "",
                "resourceVersion": "965",
                "selfLink": "/oapi/v1/netnamespacesfoo2",
                "uid": "1aabdf84-f478-11e6-aae0-507b9dac97ff"
            },
            "netid": 100,
            "netname": "foo2"
        }
    ],
    "kind": "List",
    "metadata": {},
    "resourceVersion": "",
    "selfLink": ""
}'''

        invalid_namespace = '''{
    "apiVersion": "v1",
    "items": [
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/sa.scc.mcs": "s0:c1,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000000000/10000",
                    "openshift.io/sa.scc.uid-range": "1000000000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:49Z",
                "name": "default",
                "namespace": "",
                "resourceVersion": "165",
                "selfLink": "/api/v1/namespacesdefault",
                "uid": "23c0c6aa-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/requester": "",
                    "openshift.io/sa.scc.mcs": "s0:c3,c2",
                    "openshift.io/sa.scc.supplemental-groups": "1000010000/10000",
                    "openshift.io/sa.scc.uid-range": "1000010000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:49Z",
                "name": "kube-system",
                "namespace": "",
                "resourceVersion": "3052",
                "selfLink": "/api/v1/namespaceskube-system",
                "uid": "23c21758-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/description": "",
                    "openshift.io/display-name": "",
                    "openshift.io/requester": "developer",
                    "openshift.io/sa.scc.mcs": "s0:c9,c4",
                    "openshift.io/sa.scc.supplemental-groups": "1000080000/10000",
                    "openshift.io/sa.scc.uid-range": "1000080000/10000"
                },
                "creationTimestamp": "2017-03-02T02:17:16Z",
                "name": "myproject",
                "namespace": "",
                "resourceVersion": "2898",
                "selfLink": "/api/v1/namespacesmyproject",
                "uid": "5ae3764d-feee-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "openshift.io/origin",
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/requester": "",
                    "openshift.io/sa.scc.mcs": "s0:c6,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000030000/10000",
                    "openshift.io/sa.scc.uid-range": "1000030000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:51Z",
                "name": "openshift",
                "namespace": "",
                "resourceVersion": "3057",
                "selfLink": "/api/v1/namespacesopenshift",
                "uid": "24f7b34d-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/description": "",
                    "openshift.io/display-name": "",
                    "openshift.io/requester": "system:admin",
                    "openshift.io/sa.scc.mcs": "s0:c10,c5",
                    "openshift.io/sa.scc.supplemental-groups": "1000100000/10000",
                    "openshift.io/sa.scc.uid-range": "1000100000/10000"
                },
                "creationTimestamp": "2017-03-02T02:21:15Z",
                "name": "openshift-fancy",
                "namespace": "",
                "resourceVersion": "3072",
                "selfLink": "/api/v1/namespacesopenshift-fancy",
                "uid": "e958063c-feee-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "openshift.io/origin",
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/sa.scc.mcs": "s0:c5,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000020000/10000",
                    "openshift.io/sa.scc.uid-range": "1000020000/10000"
                },
                "creationTimestamp": "2017-03-02T00:49:51Z",
                "name": "openshift-infra",
                "namespace": "",
                "resourceVersion": "169",
                "selfLink": "/api/v1/namespacesopenshift-infra",
                "uid": "24a2ed75-fee2-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "kubernetes",
                    "openshift.io/origin"
                ]
            },
            "status": {
                "phase": "Active"
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Namespace",
            "metadata": {
                "annotations": {
                    "openshift.io/description": "",
                    "openshift.io/display-name": "",
                    "openshift.io/requester": "developer1",
                    "openshift.io/sa.scc.mcs": "s0:c10,c0",
                    "openshift.io/sa.scc.supplemental-groups": "1000090000/10000",
                    "openshift.io/sa.scc.uid-range": "1000090000/10000"
                },
                "creationTimestamp": "2017-03-02T02:17:56Z",
                "name": "yourproject",
                "namespace": "",
                "resourceVersion": "2955",
                "selfLink": "/api/v1/namespacesyourproject",
                "uid": "72df7fb9-feee-11e6-b45a-507b9dac97ff"
            },
            "spec": {
                "finalizers": [
                    "openshift.io/origin",
                    "kubernetes"
                ]
            },
            "status": {
                "phase": "Active"
            }
        }
    ],
    "kind": "List",
    "metadata": {},
    "resourceVersion": "",
    "selfLink": ""
}'''

        invalid_results = {
            'hostsubnets where metadata.name != host': [{
                'apiVersion': 'v1',
                'host': 'baz1',
                'hostIP': '1.1.1.1',
                'kind': 'HostSubnet',
                'metadata': {
                    'creationTimestamp': '2017-02-16T18:47:49Z',
                    'name': 'baz0',
                    'namespace': '',
                    'resourceVersion': '996',
                    'selfLink': '/oapi/v1/hostsubnetsbaz0',
                    'uid': '69f75f87-f478-11e6-aae0-507b9dac97ff'
                },
                'subnet': '1.1.0.0/24'
            }],
            'netnamespaces where metadata.name != netname': [{
                'apiVersion': 'v1',
                'kind': 'NetNamespace',
                'metadata': {
                    'creationTimestamp': '2017-02-16T18:45:52Z',
                    'name': 'bar0',
                    'namespace': '',
                    'resourceVersion': '969',
                    'selfLink': '/oapi/v1/netnamespacesbar0',
                    'uid': '245d416e-f478-11e6-aae0-507b9dac97ff'
                },
                'netid': 100,
                'netname': 'bar1'
            }],
            'namespaces that use reserved names and were not created by infrastructure components': [{
                'apiVersion': 'v1',
                'kind': 'Namespace',
                'metadata': {'annotations': {'openshift.io/requester': '',
                                             'openshift.io/sa.scc.mcs': 's0:c3,c2',
                                             'openshift.io/sa.scc.supplemental-groups': '1000010000/10000',
                                             'openshift.io/sa.scc.uid-range': '1000010000/10000'},
                             'creationTimestamp': '2017-03-02T00:49:49Z',
                             'name': 'kube-system',
                             'namespace': '',
                             'resourceVersion': '3052',
                             'selfLink': '/api/v1/namespaceskube-system',
                             'uid': '23c21758-fee2-11e6-b45a-507b9dac97ff'},
                'spec': {'finalizers': ['kubernetes', 'openshift.io/origin']},
                'status': {'phase': 'Active'}},
                {'apiVersion': 'v1',
                 'kind': 'Namespace',
                 'metadata': {'annotations': {'openshift.io/requester': '',
                                              'openshift.io/sa.scc.mcs': 's0:c6,c0',
                                              'openshift.io/sa.scc.supplemental-groups': '1000030000/10000',
                                              'openshift.io/sa.scc.uid-range': '1000030000/10000'},
                              'creationTimestamp': '2017-03-02T00:49:51Z',
                              'name': 'openshift',
                              'namespace': '',
                              'resourceVersion': '3057',
                              'selfLink': '/api/v1/namespacesopenshift',
                              'uid': '24f7b34d-fee2-11e6-b45a-507b9dac97ff'},
                 'spec': {'finalizers': ['kubernetes', 'openshift.io/origin']},
                 'status': {'phase': 'Active'}},
                {'apiVersion': 'v1',
                 'kind': 'Namespace',
                 'metadata': {'annotations': {'openshift.io/description': '',
                                              'openshift.io/display-name': '',
                                              'openshift.io/requester': 'system:admin',
                                              'openshift.io/sa.scc.mcs': 's0:c10,c5',
                                              'openshift.io/sa.scc.supplemental-groups': '1000100000/10000',
                                              'openshift.io/sa.scc.uid-range': '1000100000/10000'},
                              'creationTimestamp': '2017-03-02T02:21:15Z',
                              'name': 'openshift-fancy',
                              'namespace': '',
                              'resourceVersion': '3072',
                              'selfLink': '/api/v1/namespacesopenshift-fancy',
                              'uid': 'e958063c-feee-11e6-b45a-507b9dac97ff'},
                 'spec': {'finalizers': ['openshift.io/origin', 'kubernetes']},
                 'status': {'phase': 'Active'}
                 }],
        }

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            # First call to mock
            (0, invalid_hostsubnet, ''),

            # Second call to mock
            (0, invalid_netnamespace, ''),

            # Third call to mock
            (0, invalid_namespace, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        mock_oc_binary.side_effect = [
            'oc'
        ]

        # Act
        results = OCObjectValidator.run_ansible(params)

        # Assert
        self.assertTrue(results['failed'])
        self.assertIn('All objects are not valid.', results['msg'])
        self.assertEqual(results['state'], 'list')
        self.assertEqual(results['results'], invalid_results)

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'hostsubnet', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'netnamespace', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'namespace', '-o', 'json', '-n', 'default'], None),
        ])
