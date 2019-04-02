'''
 Unit tests for oc image
'''
import os
import sys
import unittest
import mock
import six

# Removing invalid variable names for tests so that I can
# keep them brief
# pylint: disable=invalid-name,no-name-in-module
# Disable import-error b/c our libraries aren't loaded in jenkins
# pylint: disable=import-error
# place class in our python path
module_path = os.path.join('/'.join(os.path.realpath(__file__).split('/')[:-4]), 'library')  # noqa: E501
sys.path.insert(0, module_path)
from oc_image import OCImage, locate_oc_binary  # noqa: E402


class OCImageTest(unittest.TestCase):
    '''
     Test class for OCImage
    '''

    @mock.patch('oc_image.Utils.create_tmpfile_copy')
    @mock.patch('oc_image.OCImage._run')
    def test_state_list(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a label list '''
        params = {'registry_url': 'registry.ops.openshift.com',
                  'image_name': 'oso-rhel7-zagg-web',
                  'image_tag': 'int',
                  'namespace': 'default',
                  'state': 'list',
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'debug': False}

        istream = '''{
            "kind": "ImageStream",
            "apiVersion": "v1",
            "metadata": {
                "name": "oso-rhel7-zagg-web",
                "namespace": "default",
                "selfLink": "/oapi/v1/namespaces/default/imagestreams/oso-rhel7-zagg-web",
                "uid": "6ca2b199-dcdb-11e6-8ffd-0a5f8e3e32be",
                "resourceVersion": "8135944",
                "generation": 1,
                "creationTimestamp": "2017-01-17T17:36:05Z",
                "annotations": {
                    "openshift.io/image.dockerRepositoryCheck": "2017-01-17T17:36:05Z"
                }
            },
            "spec": {
                "tags": [
                    {
                        "name": "int",
                        "annotations": null,
                        "from": {
                            "kind": "DockerImage",
                            "name": "registry.ops.openshift.com/ops/oso-rhel7-zagg-web:int"
                        },
                        "generation": 1,
                        "importPolicy": {}
                    }
                ]
            },
            "status": {
                "dockerImageRepository": "172.30.183.164:5000/default/oso-rhel7-zagg-web",
                "tags": [
                    {
                        "tag": "int",
                        "items": [
                            {
                                "created": "2017-01-17T17:36:05Z",
                                "dockerImageReference": "registry.ops.openshift.com/ops/oso-rhel7-zagg-web@sha256:645bab780cf18a9b764d64b02ca65c39d13cb16f19badd0a49a1668629759392",
                                "image": "sha256:645bab780cf18a9b764d64b02ca65c39d13cb16f19badd0a49a1668629759392",
                                "generation": 1
                            }
                        ]
                    }
                ]
            }
        }
        '''

        mock_cmd.side_effect = [
            (0, istream, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCImage.run_ansible(params, False)

        self.assertFalse(results['changed'])
        self.assertEquals(results['results']['results'][0]['metadata']['name'], 'oso-rhel7-zagg-web')

    @mock.patch('oc_image.Utils.create_tmpfile_copy')
    @mock.patch('oc_image.OCImage._run')
    def test_state_present(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a image present '''
        params = {'registry_url': 'registry.ops.openshift.com',
                  'image_name': 'oso-rhel7-zagg-web',
                  'image_tag': 'int',
                  'namespace': 'default',
                  'state': 'present',
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'debug': False}

        istream = '''{
            "kind": "ImageStream",
            "apiVersion": "v1",
            "metadata": {
                "name": "oso-rhel7-zagg-web",
                "namespace": "default",
                "selfLink": "/oapi/v1/namespaces/default/imagestreams/oso-rhel7-zagg-web",
                "uid": "6ca2b199-dcdb-11e6-8ffd-0a5f8e3e32be",
                "resourceVersion": "8135944",
                "generation": 1,
                "creationTimestamp": "2017-01-17T17:36:05Z",
                "annotations": {
                    "openshift.io/image.dockerRepositoryCheck": "2017-01-17T17:36:05Z"
                }
            },
            "spec": {
                "tags": [
                    {
                        "name": "int",
                        "annotations": null,
                        "from": {
                            "kind": "DockerImage",
                            "name": "registry.ops.openshift.com/ops/oso-rhel7-zagg-web:int"
                        },
                        "generation": 1,
                        "importPolicy": {}
                    }
                ]
            },
            "status": {
                "dockerImageRepository": "172.30.183.164:5000/default/oso-rhel7-zagg-web",
                "tags": [
                    {
                        "tag": "int",
                        "items": [
                            {
                                "created": "2017-01-17T17:36:05Z",
                                "dockerImageReference": "registry.ops.openshift.com/ops/oso-rhel7-zagg-web@sha256:645bab780cf18a9b764d64b02ca65c39d13cb16f19badd0a49a1668629759392",
                                "image": "sha256:645bab780cf18a9b764d64b02ca65c39d13cb16f19badd0a49a1668629759392",
                                "generation": 1
                            }
                        ]
                    }
                ]
            }
        }
        '''

        mock_cmd.side_effect = [
            (1, '', 'Error from server: imagestreams "oso-rhel7-zagg-web" not found'),
            (0, '', ''),
            (0, istream, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCImage.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertTrue(results['results']['results'][0]['metadata']['name'] == 'oso-rhel7-zagg-web')

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
