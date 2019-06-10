'''
 Unit tests for oc configmap
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
from oc_configmap import OCConfigMap, locate_oc_binary  # noqa: E402


class OCConfigMapTest(unittest.TestCase):
    '''
     Test class for OCConfigMap
    '''
    params = {'kubeconfig': '/etc/origin/master/admin.kubeconfig',
              'state': 'present',
              'debug': False,
              'name': 'configmap',
              'from_file': {},
              'from_literal': {},
              'namespace': 'test'}

    @mock.patch('oc_configmap.Utils._write')
    @mock.patch('oc_configmap.Utils.create_tmpfile_copy')
    @mock.patch('oc_configmap.OCConfigMap._run')
    def test_create_configmap(self, mock_run, mock_tmpfile_copy, mock_write):
        ''' Testing a configmap create '''
        # TODO
        return
        params = copy.deepcopy(OCConfigMapTest.params)
        params['from_file'] = {'test': '/root/file'}
        params['from_literal'] = {'foo': 'bar'}

        configmap = '''{
                "apiVersion": "v1",
                "data": {
                    "foo": "bar",
                    "test": "this is a file\\n"
                },
                "kind": "ConfigMap",
                "metadata": {
                    "creationTimestamp": "2017-03-20T20:24:35Z",
                    "name": "configmap",
                    "namespace": "test"
                }
            }'''

        mock_run.side_effect = [
            (1, '', 'Error from server (NotFound): configmaps "configmap" not found'),
            (0, '', ''),
            (0, configmap, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCConfigMap.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'configmap')

    @mock.patch('oc_configmap.Utils._write')
    @mock.patch('oc_configmap.Utils.create_tmpfile_copy')
    @mock.patch('oc_configmap.OCConfigMap._run')
    def test_update_configmap(self, mock_run, mock_tmpfile_copy, mock_write):
        ''' Testing a configmap create '''
        params = copy.deepcopy(OCConfigMapTest.params)
        params['from_file'] = {'test': '/root/file'}
        params['from_literal'] = {'foo': 'bar', 'deployment_type': 'openshift-enterprise'}

        configmap = '''{
                "apiVersion": "v1",
                "data": {
                    "foo": "bar",
                    "test": "this is a file\\n"
                },
                "kind": "ConfigMap",
                "metadata": {
                    "creationTimestamp": "2017-03-20T20:24:35Z",
                    "name": "configmap",
                    "namespace": "test"

                }
            }'''

        mod_configmap = '''{
                "apiVersion": "v1",
                "data": {
                    "foo": "bar",
                    "deployment_type": "openshift-enterprise",
                    "test": "this is a file\\n"
                },
                "kind": "ConfigMap",
                "metadata": {
                    "creationTimestamp": "2017-03-20T20:24:35Z",
                    "name": "configmap",
                    "namespace": "test"

                }
            }'''

        mock_run.side_effect = [
            (0, configmap, ''),
            (0, mod_configmap, ''),
            (0, configmap, ''),
            (0, '', ''),
            (0, mod_configmap, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCConfigMap.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'configmap')
        self.assertEqual(results['results']['results'][0]['data']['deployment_type'], 'openshift-enterprise')

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
