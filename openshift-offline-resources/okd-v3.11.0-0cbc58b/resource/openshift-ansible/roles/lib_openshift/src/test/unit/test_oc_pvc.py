'''
 Unit tests for oc pvc
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
from oc_pvc import OCPVC, locate_oc_binary  # noqa: E402


class OCPVCTest(unittest.TestCase):
    '''
     Test class for OCPVC
    '''
    params = {'kubeconfig': '/etc/origin/master/admin.kubeconfig',
              'state': 'present',
              'debug': False,
              'name': 'mypvc',
              'namespace': 'test',
              'volume_capacity': '1G',
              'selector': {'foo': 'bar', 'abc': 'a123'},
              'storage_class_name': 'mystorage',
              'access_modes': 'ReadWriteMany'}

    @mock.patch('oc_pvc.Utils.create_tmpfile_copy')
    @mock.patch('oc_pvc.OCPVC._run')
    def test_create_pvc(self, mock_run, mock_tmpfile_copy):
        ''' Testing a pvc create '''
        params = copy.deepcopy(OCPVCTest.params)

        pvc = '''{"kind": "PersistentVolumeClaim",
               "apiVersion": "v1",
               "metadata": {
                   "name": "mypvc",
                   "namespace": "test",
                   "selfLink": "/api/v1/namespaces/test/persistentvolumeclaims/mypvc",
                   "uid": "77597898-d8d8-11e6-aea5-0e3c0c633889",
                   "resourceVersion": "126510787",
                   "creationTimestamp": "2017-01-12T15:04:50Z",
                   "labels": {
                       "mypvc": "database"
                   },
                   "annotations": {
                       "pv.kubernetes.io/bind-completed": "yes",
                       "pv.kubernetes.io/bound-by-controller": "yes",
                       "v1.2-volume.experimental.kubernetes.io/provisioning-required": "volume.experimental.kubernetes.io/provisioning-completed"
                   }
               },
               "spec": {
                   "accessModes": [
                       "ReadWriteOnce"
                   ],
                   "resources": {
                       "requests": {
                           "storage": "1Gi"
                       }
                   },
                   "selector": {
                       "matchLabels": {
                           "foo": "bar",
                           "abc": "a123"
                       }
                   },
                   "storageClassName": "myStorage",
                   "volumeName": "pv-aws-ow5vl"
               },
               "status": {
                  "phase": "Bound",
                   "accessModes": [
                       "ReadWriteOnce"
                   ],
                    "capacity": {
                      "storage": "1Gi"
                    }
               }
              }'''

        mock_run.side_effect = [
            (1, '', 'Error from server: persistentvolumeclaims "mypvc" not found'),
            (1, '', 'Error from server: persistentvolumeclaims "mypvc" not found'),
            (0, '', ''),
            (0, pvc, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCPVC.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertEqual(results['results']['results'][0]['metadata']['name'], 'mypvc')
        self.assertEqual(results['results']['results'][0]['spec']['storageClassName'], 'myStorage')
        self.assertEqual(results['results']['results'][0]['spec']['selector']['matchLabels']['foo'], 'bar')

    @mock.patch('oc_pvc.Utils.create_tmpfile_copy')
    @mock.patch('oc_pvc.OCPVC._run')
    def test_update_pvc(self, mock_run, mock_tmpfile_copy):
        ''' Testing a pvc create '''
        params = copy.deepcopy(OCPVCTest.params)
        params['access_modes'] = 'ReadWriteMany'

        pvc = '''{"kind": "PersistentVolumeClaim",
               "apiVersion": "v1",
               "metadata": {
                   "name": "mypvc",
                   "namespace": "test",
                   "selfLink": "/api/v1/namespaces/test/persistentvolumeclaims/mypvc",
                   "uid": "77597898-d8d8-11e6-aea5-0e3c0c633889",
                   "resourceVersion": "126510787",
                   "creationTimestamp": "2017-01-12T15:04:50Z",
                   "labels": {
                       "mypvc": "database"
                   },
                   "annotations": {
                       "pv.kubernetes.io/bind-completed": "yes",
                       "pv.kubernetes.io/bound-by-controller": "yes",
                       "v1.2-volume.experimental.kubernetes.io/provisioning-required": "volume.experimental.kubernetes.io/provisioning-completed"
                   }
               },
               "spec": {
                   "accessModes": [
                       "ReadWriteOnce"
                   ],
                   "resources": {
                       "requests": {
                           "storage": "1Gi"
                       }
                   },
                   "volumeName": "pv-aws-ow5vl"
               },
               "status": {
                  "phase": "Bound",
                   "accessModes": [
                       "ReadWriteOnce"
                   ],
                    "capacity": {
                      "storage": "1Gi"
                    }
               }
              }'''

        mod_pvc = '''{"kind": "PersistentVolumeClaim",
               "apiVersion": "v1",
               "metadata": {
                   "name": "mypvc",
                   "namespace": "test",
                   "selfLink": "/api/v1/namespaces/test/persistentvolumeclaims/mypvc",
                   "uid": "77597898-d8d8-11e6-aea5-0e3c0c633889",
                   "resourceVersion": "126510787",
                   "creationTimestamp": "2017-01-12T15:04:50Z",
                   "labels": {
                       "mypvc": "database"
                   },
                   "annotations": {
                       "pv.kubernetes.io/bind-completed": "yes",
                       "pv.kubernetes.io/bound-by-controller": "yes",
                       "v1.2-volume.experimental.kubernetes.io/provisioning-required": "volume.experimental.kubernetes.io/provisioning-completed"
                   }
               },
               "spec": {
                   "accessModes": [
                       "ReadWriteMany"
                   ],
                   "resources": {
                       "requests": {
                           "storage": "1Gi"
                       }
                   },
                   "volumeName": "pv-aws-ow5vl"
               },
               "status": {
                  "phase": "Bound",
                   "accessModes": [
                       "ReadWriteOnce"
                   ],
                    "capacity": {
                      "storage": "1Gi"
                    }
               }
              }'''

        mock_run.side_effect = [
            (0, pvc, ''),
            (0, pvc, ''),
            (0, '', ''),
            (0, mod_pvc, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCPVC.run_ansible(params, False)

        self.assertFalse(results['changed'])
        self.assertEqual(results['results']['msg'], '##### - This volume is currently bound.  Will not update - ####')

    @mock.patch('oc_pvc.Utils.create_tmpfile_copy')
    @mock.patch('oc_pvc.OCPVC._run')
    def test_delete_pvc(self, mock_run, mock_tmpfile_copy):
        ''' Testing a pvc create '''
        params = copy.deepcopy(OCPVCTest.params)
        params['state'] = 'absent'

        pvc = '''{"kind": "PersistentVolumeClaim",
               "apiVersion": "v1",
               "metadata": {
                   "name": "mypvc",
                   "namespace": "test",
                   "selfLink": "/api/v1/namespaces/test/persistentvolumeclaims/mypvc",
                   "uid": "77597898-d8d8-11e6-aea5-0e3c0c633889",
                   "resourceVersion": "126510787",
                   "creationTimestamp": "2017-01-12T15:04:50Z",
                   "labels": {
                       "mypvc": "database"
                   },
                   "annotations": {
                       "pv.kubernetes.io/bind-completed": "yes",
                       "pv.kubernetes.io/bound-by-controller": "yes",
                       "v1.2-volume.experimental.kubernetes.io/provisioning-required": "volume.experimental.kubernetes.io/provisioning-completed"
                   }
               },
               "spec": {
                   "accessModes": [
                       "ReadWriteOnce"
                   ],
                   "resources": {
                       "requests": {
                           "storage": "1Gi"
                       }
                   },
                   "volumeName": "pv-aws-ow5vl"
               },
               "status": {
                  "phase": "Bound",
                   "accessModes": [
                       "ReadWriteOnce"
                   ],
                    "capacity": {
                      "storage": "1Gi"
                    }
               }
              }'''

        mock_run.side_effect = [
            (0, pvc, ''),
            (0, '', ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCPVC.run_ansible(params, False)

        self.assertTrue(results['changed'])

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
