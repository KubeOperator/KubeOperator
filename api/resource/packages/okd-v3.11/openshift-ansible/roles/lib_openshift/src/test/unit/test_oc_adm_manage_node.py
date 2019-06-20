'''
 Unit tests for oc_adm_manage_node
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
from oc_adm_manage_node import ManageNode, locate_oc_binary  # noqa: E402


class ManageNodeTest(unittest.TestCase):
    '''
     Test class for oc_adm_manage_node
    '''

    @mock.patch('oc_adm_manage_node.Utils.create_tmpfile_copy')
    @mock.patch('oc_adm_manage_node.ManageNode.openshift_cmd')
    def test_list_pods(self, mock_openshift_cmd, mock_tmpfile_copy):
        ''' Testing a get '''
        params = {'node': ['ip-172-31-49-140.ec2.internal'],
                  'schedulable': None,
                  'selector': None,
                  'pod_selector': None,
                  'list_pods': True,
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'evacuate': False,
                  'grace_period': False,
                  'dry_run': False,
                  'force': False}

        pod_list = '''{
    "metadata": {},
    "items": [
        {
            "metadata": {
                "name": "docker-registry-1-xuhik",
                "generateName": "docker-registry-1-",
                "namespace": "default",
                "selfLink": "/api/v1/namespaces/default/pods/docker-registry-1-xuhik",
                "uid": "ae2a25a2-e316-11e6-80eb-0ecdc51fcfc4",
                "resourceVersion": "1501",
                "creationTimestamp": "2017-01-25T15:55:23Z",
                "labels": {
                    "deployment": "docker-registry-1",
                    "deploymentconfig": "docker-registry",
                    "docker-registry": "default"
                },
                "annotations": {
                    "openshift.io/deployment-config.latest-version": "1",
                    "openshift.io/deployment-config.name": "docker-registry",
                    "openshift.io/deployment.name": "docker-registry-1",
                    "openshift.io/scc": "restricted"
                }
            },
            "spec": {}
        },
        {
            "metadata": {
                "name": "router-1-kp3m3",
                "generateName": "router-1-",
                "namespace": "default",
                "selfLink": "/api/v1/namespaces/default/pods/router-1-kp3m3",
                "uid": "9e71f4a5-e316-11e6-80eb-0ecdc51fcfc4",
                "resourceVersion": "1456",
                "creationTimestamp": "2017-01-25T15:54:56Z",
                "labels": {
                    "deployment": "router-1",
                    "deploymentconfig": "router",
                    "router": "router"
                },
                "annotations": {
                    "openshift.io/deployment-config.latest-version": "1",
                    "openshift.io/deployment-config.name": "router",
                    "openshift.io/deployment.name": "router-1",
                    "openshift.io/scc": "hostnetwork"
                }
            },
            "spec": {}
        }]
}'''

        mock_openshift_cmd.side_effect = [
            {"cmd": "/usr/bin/oadm manage-node ip-172-31-49-140.ec2.internal --list-pods",
             "results": pod_list,
             "returncode": 0}
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = ManageNode.run_ansible(params, False)

        # returned a single node
        self.assertTrue(len(results['results']['nodes']) == 1)
        # returned 2 pods
        self.assertTrue(len(results['results']['nodes']['ip-172-31-49-140.ec2.internal']) == 2)

    @mock.patch('oc_adm_manage_node.Utils.create_tmpfile_copy')
    @mock.patch('oc_adm_manage_node.ManageNode.openshift_cmd')
    def test_schedulable_false(self, mock_openshift_cmd, mock_tmpfile_copy):
        ''' Testing a get '''
        params = {'node': ['ip-172-31-49-140.ec2.internal'],
                  'schedulable': False,
                  'selector': None,
                  'pod_selector': None,
                  'list_pods': False,
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'evacuate': False,
                  'grace_period': False,
                  'dry_run': False,
                  'force': False}

        node = [{
            "apiVersion": "v1",
            "kind": "Node",
            "metadata": {
                "creationTimestamp": "2017-01-26T14:34:43Z",
                "labels": {
                    "beta.kubernetes.io/arch": "amd64",
                    "beta.kubernetes.io/instance-type": "m4.large",
                    "beta.kubernetes.io/os": "linux",
                    "failure-domain.beta.kubernetes.io/region": "us-east-1",
                    "failure-domain.beta.kubernetes.io/zone": "us-east-1c",
                    "hostname": "opstest-node-compute-0daaf",
                    "kubernetes.io/hostname": "ip-172-31-51-111.ec2.internal",
                    "ops_node": "old",
                    "region": "us-east-1",
                    "type": "compute"
                },
                "name": "ip-172-31-51-111.ec2.internal",
                "resourceVersion": "6936",
                "selfLink": "/api/v1/nodes/ip-172-31-51-111.ec2.internal",
                "uid": "93d7fdfb-e3d4-11e6-a982-0e84250fc302"
            },
            "spec": {
                "externalID": "i-06bb330e55c699b0f",
                "providerID": "aws:///us-east-1c/i-06bb330e55c699b0f",
            }}]

        mock_openshift_cmd.side_effect = [
            {"cmd": "/usr/bin/oc get node -o json ip-172-31-49-140.ec2.internal",
             "results": node,
             "returncode": 0},
            {"cmd": "/usr/bin/oadm manage-node ip-172-31-49-140.ec2.internal --schedulable=False",
             "results": "NAME                            STATUS    AGE\n" +
                        "ip-172-31-49-140.ec2.internal   Ready,SchedulingDisabled     5h\n",
             "returncode": 0}]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = ManageNode.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertEqual(results['results']['nodes'][0]['name'], 'ip-172-31-49-140.ec2.internal')
        self.assertEqual(results['results']['nodes'][0]['schedulable'], False)

    @unittest.skipIf(six.PY3, 'py2 test only')
    @mock.patch('os.path.exists')
    @mock.patch('os.environ.get')
    def test_binary_lookup_fallback(self, mock_env_get, mock_path_exists):
        ''' Testing binary lookup fallback '''

        mock_env_get.side_effect = lambda _v, _d: ''

        mock_path_exists.side_effect = lambda _: False

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
