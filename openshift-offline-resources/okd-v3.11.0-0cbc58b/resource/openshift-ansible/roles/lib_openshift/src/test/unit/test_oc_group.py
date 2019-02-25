'''
 Unit tests for oc group
'''

import copy
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
from oc_group import OCGroup, locate_oc_binary  # noqa: E402


class OCGroupTest(unittest.TestCase):
    '''
     Test class for OCGroup
    '''
    params = {'kubeconfig': '/etc/origin/master/admin.kubeconfig',
              'state': 'present',
              'debug': False,
              'name': 'acme',
              'namespace': 'test'}

    @mock.patch('oc_group.Utils.create_tmpfile_copy')
    @mock.patch('oc_group.OCGroup._run')
    def test_create_group(self, mock_run, mock_tmpfile_copy):
        ''' Testing a group create '''
        params = copy.deepcopy(OCGroupTest.params)

        group = '''{
            "kind": "Group",
            "apiVersion": "v1",
            "metadata": {
                "name": "acme"
            },
            "users": []
        }'''

        mock_run.side_effect = [
            (1, '', 'Error from server: groups.user.openshift.io "acme" not found'),
            (1, '', 'Error from server: groups.user.openshift.io "acme" not found'),
            (0, '', ''),
            (0, group, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCGroup.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'acme')

    @mock.patch('oc_group.Utils.create_tmpfile_copy')
    @mock.patch('oc_group.OCGroup._run')
    def test_failed_get_group(self, mock_run, mock_tmpfile_copy):
        ''' Testing a group create '''
        params = copy.deepcopy(OCGroupTest.params)
        params['state'] = 'list'
        params['name'] = 'noexist'

        mock_run.side_effect = [
            (1, '', 'Error from server: groups.user.openshift.io "acme" not found'),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCGroup.run_ansible(params, False)

        self.assertTrue(results['failed'])

    @mock.patch('oc_group.Utils.create_tmpfile_copy')
    @mock.patch('oc_group.OCGroup._run')
    def test_delete_group(self, mock_run, mock_tmpfile_copy):
        ''' Testing a group create '''
        params = copy.deepcopy(OCGroupTest.params)
        params['state'] = 'absent'

        group = '''{
            "kind": "Group",
            "apiVersion": "v1",
            "metadata": {
                "name": "acme"
            },
            "users": [
              "user1"
            ]
        }'''

        mock_run.side_effect = [
            (0, group, ''),
            (0, '', ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCGroup.run_ansible(params, False)

        self.assertTrue(results['changed'])

    @mock.patch('oc_group.Utils.create_tmpfile_copy')
    @mock.patch('oc_group.OCGroup._run')
    def test_get_group(self, mock_run, mock_tmpfile_copy):
        ''' Testing a group create '''
        params = copy.deepcopy(OCGroupTest.params)
        params['state'] = 'list'

        group = '''{
            "kind": "Group",
            "apiVersion": "v1",
            "metadata": {
                "name": "acme"
            },
            "users": [
              "user1"
            ]
        }'''

        mock_run.side_effect = [
            (0, group, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCGroup.run_ansible(params, False)

        self.assertFalse(results['changed'])
        self.assertEqual(results['results'][0]['metadata']['name'], 'acme')
        self.assertEqual(results['results'][0]['users'][0], 'user1')

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
