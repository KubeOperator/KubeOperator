'''
 Unit tests for oc clusterrole
'''

import copy
import os
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
from oc_clusterrole import OCClusterRole  # noqa: E402


class OCClusterRoleTest(unittest.TestCase):
    '''
     Test class for OCClusterRole
    '''

    # run_ansible input parameters
    params = {
        'state': 'present',
        'name': 'operations',
        'rules': [
            {'apiGroups': [''],
             'attributeRestrictions': None,
             'verbs': ['create', 'delete', 'deletecollection',
                       'get', 'list', 'patch', 'update', 'watch'],
             'resources': ['persistentvolumes']}
        ],
        'kubeconfig': '/etc/origin/master/admin.kubeconfig',
        'debug': False,
    }

    @mock.patch('oc_clusterrole.locate_oc_binary')
    @mock.patch('oc_clusterrole.Utils.create_tmpfile_copy')
    @mock.patch('oc_clusterrole.Utils._write')
    @mock.patch('oc_clusterrole.OCClusterRole._run')
    def test_adding_a_clusterrole(self, mock_cmd, mock_write, mock_tmpfile_copy, mock_loc_binary):
        ''' Testing adding a project '''

        params = copy.deepcopy(OCClusterRoleTest.params)

        clusterrole = '''{
            "apiVersion": "v1",
            "kind": "ClusterRole",
            "metadata": {
                "creationTimestamp": "2017-03-27T14:19:09Z",
                "name": "operations",
                "resourceVersion": "23",
                "selfLink": "/oapi/v1/clusterrolesoperations",
                "uid": "57d358fe-12f8-11e7-874a-0ec502977670"
            },
            "rules": [
                {
                    "apiGroups": [
                        ""
                    ],
                    "attributeRestrictions": null,
                    "resources": [
                        "persistentvolumes"
                    ],
                    "verbs": [
                        "create",
                        "delete",
                        "deletecollection",
                        "get",
                        "list",
                        "patch",
                        "update",
                        "watch"
                    ]
                }
            ]
        }'''

        # Return values of our mocked function call. These get returned once per call.
        mock_cmd.side_effect = [
            (1, '', 'Error from server: clusterrole "operations" not found'),
            (1, '', 'Error from server: namespaces "operations" not found'),
            (0, '', ''),  # created
            (0, clusterrole, ''),  # fetch it
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        mock_loc_binary.side_effect = [
            'oc',
        ]

        # Act
        results = OCClusterRole.run_ansible(params, False)

        # Assert
        self.assertTrue(results['changed'])
        self.assertEqual(results['results']['returncode'], 0)
        self.assertEqual(results['results']['results']['metadata']['name'], 'operations')
        self.assertEqual(results['state'], 'present')

        # Making sure our mock was called as we expected
        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'clusterrole', 'operations', '-o', 'json'], None),
            mock.call(['oc', 'get', 'clusterrole', 'operations', '-o', 'json'], None),
            mock.call(['oc', 'create', '-f', mock.ANY], None),
            mock.call(['oc', 'get', 'clusterrole', 'operations', '-o', 'json'], None),
        ])
