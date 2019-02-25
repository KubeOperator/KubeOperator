#!/usr/bin/env python2
'''
 Unit tests for oc user
'''
# To run
# ./oc_user.py
#
# ..
# ----------------------------------------------------------------------
# Ran 2 tests in 0.003s
#
# OK

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
from oc_user import OCUser  # noqa: E402


class OCUserTest(unittest.TestCase):
    '''
     Test class for OCUser
    '''

    def setUp(self):
        ''' setup method will create a file and set to known configuration '''
        pass

    @mock.patch('oc_user.Utils.create_tmpfile_copy')
    @mock.patch('oc_user.OCUser._run')
    def test_state_list(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a user list '''
        params = {'username': 'testuser@email.com',
                  'state': 'list',
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'full_name': None,
                  'groups': [],
                  'debug': False}

        user = '''{
               "kind": "User",
               "apiVersion": "v1",
               "metadata": {
                   "name": "testuser@email.com",
                   "selfLink": "/oapi/v1/users/testuser@email.com",
                   "uid": "02fee6c9-f20d-11e6-b83b-12e1a7285e80",
                   "resourceVersion": "38566887",
                   "creationTimestamp": "2017-02-13T16:53:58Z"
               },
               "fullName": "Test User",
               "identities": null,
               "groups": null
           }'''

        mock_cmd.side_effect = [
            (0, user, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCUser.run_ansible(params, False)

        self.assertFalse(results['changed'])
        self.assertTrue(results['results'][0]['metadata']['name'] == "testuser@email.com")

    @mock.patch('oc_user.Utils.create_tmpfile_copy')
    @mock.patch('oc_user.OCUser._run')
    def test_state_present(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a user list '''
        params = {'username': 'testuser@email.com',
                  'state': 'present',
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'full_name': 'Test User',
                  'groups': [],
                  'debug': False}

        created_user = '''{
                          "kind": "User",
                          "apiVersion": "v1",
                          "metadata": {
                              "name": "testuser@email.com",
                              "selfLink": "/oapi/v1/users/testuser@email.com",
                              "uid": "8d508039-f224-11e6-b83b-12e1a7285e80",
                              "resourceVersion": "38646241",
                              "creationTimestamp": "2017-02-13T19:42:28Z"
                          },
                          "fullName": "Test User",
                          "identities": null,
                          "groups": null
                      }'''

        mock_cmd.side_effect = [
            (1, '', 'Error from server: users "testuser@email.com" not found'),  # get
            (1, '', 'Error from server: users "testuser@email.com" not found'),  # get
            (0, 'user "testuser@email.com" created', ''),  # create
            (0, created_user, ''),  # get
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCUser.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertTrue(results['results']['results'][0]['metadata']['name'] ==
                        "testuser@email.com")

    def tearDown(self):
        '''TearDown method'''
        pass


if __name__ == "__main__":
    unittest.main()
