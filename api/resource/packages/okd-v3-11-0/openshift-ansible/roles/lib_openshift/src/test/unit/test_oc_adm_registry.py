#!/usr/bin/env python
'''
 Unit tests for oc adm registry
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
from oc_adm_registry import Registry, locate_oc_binary  # noqa: E402


# pylint: disable=too-many-public-methods
class RegistryTest(unittest.TestCase):
    '''
     Test class for Registry
    '''
    dry_run = '''{
        "kind": "List",
        "apiVersion": "v1",
        "metadata": {},
        "items": [
            {
                "kind": "ServiceAccount",
                "apiVersion": "v1",
                "metadata": {
                    "name": "registry",
                    "creationTimestamp": null
                }
            },
            {
                "kind": "ClusterRoleBinding",
                "apiVersion": "v1",
                "metadata": {
                    "name": "registry-registry-role",
                    "creationTimestamp": null
                },
                "userNames": [
                    "system:serviceaccount:default:registry"
                ],
                "groupNames": null,
                "subjects": [
                    {
                        "kind": "ServiceAccount",
                        "namespace": "default",
                        "name": "registry"
                    }
                ],
                "roleRef": {
                    "kind": "ClusterRole",
                    "name": "system:registry"
                }
            },
            {
                "kind": "DeploymentConfig",
                "apiVersion": "v1",
                "metadata": {
                    "name": "docker-registry",
                    "creationTimestamp": null,
                    "labels": {
                        "docker-registry": "default"
                    }
                },
                "spec": {
                    "strategy": {
                        "resources": {}
                    },
                    "triggers": [
                        {
                            "type": "ConfigChange"
                        }
                    ],
                    "replicas": 1,
                    "test": false,
                    "selector": {
                        "docker-registry": "default"
                    },
                    "template": {
                        "metadata": {
                            "creationTimestamp": null,
                            "labels": {
                                "docker-registry": "default"
                            }
                        },
                        "spec": {
                            "volumes": [
                                {
                                    "name": "registry-storage",
                                    "emptyDir": {}
                                }
                            ],
                            "containers": [
                                {
                                    "name": "registry",
                                    "image": "registry.redhat.io/openshift3/ose-docker-registry:v3.5.0.39",
                                    "ports": [
                                        {
                                            "containerPort": 5000
                                        }
                                    ],
                                    "env": [
                                        {
                                            "name": "REGISTRY_HTTP_ADDR",
                                            "value": ":5000"
                                        },
                                        {
                                            "name": "REGISTRY_HTTP_NET",
                                            "value": "tcp"
                                        },
                                        {
                                            "name": "REGISTRY_HTTP_SECRET",
                                            "value": "WQjSGeUu5KFZRTwGeIXgwIjyraNDLmdJblsFbtzZdF8="
                                        },
                                        {
                                            "name": "REGISTRY_MIDDLEWARE_REPOSITORY_OPENSHIFT_ENFORCEQUOTA",
                                            "value": "false"
                                        }
                                    ],
                                    "resources": {
                                        "requests": {
                                            "cpu": "100m",
                                            "memory": "256Mi"
                                        }
                                    },
                                    "volumeMounts": [
                                        {
                                            "name": "registry-storage",
                                            "mountPath": "/registry"
                                        }
                                    ],
                                    "livenessProbe": {
                                        "httpGet": {
                                            "path": "/healthz",
                                            "port": 5000
                                        },
                                        "initialDelaySeconds": 10,
                                        "timeoutSeconds": 5
                                    },
                                    "readinessProbe": {
                                        "httpGet": {
                                            "path": "/healthz",
                                            "port": 5000
                                        },
                                        "timeoutSeconds": 5
                                    },
                                    "securityContext": {
                                        "privileged": false
                                    }
                                }
                            ],
                            "nodeSelector": {
                                "type": "infra"
                            },
                            "serviceAccountName": "registry",
                            "serviceAccount": "registry"
                        }
                    }
                },
                "status": {
                    "latestVersion": 0,
                    "observedGeneration": 0,
                    "replicas": 0,
                    "updatedReplicas": 0,
                    "availableReplicas": 0,
                    "unavailableReplicas": 0
                }
            },
            {
                "kind": "Service",
                "apiVersion": "v1",
                "metadata": {
                    "name": "docker-registry",
                    "creationTimestamp": null,
                    "labels": {
                        "docker-registry": "default"
                    }
                },
                "spec": {
                    "ports": [
                        {
                            "name": "5000-tcp",
                            "port": 5000,
                            "targetPort": 5000
                        }
                    ],
                    "selector": {
                        "docker-registry": "default"
                    },
                    "clusterIP": "172.30.119.110",
                    "sessionAffinity": "ClientIP"
                },
                "status": {
                    "loadBalancer": {}
                }
            }
        ]}'''

    @mock.patch('oc_adm_registry.locate_oc_binary')
    @mock.patch('oc_adm_registry.Utils._write')
    @mock.patch('oc_adm_registry.Utils.create_tmpfile_copy')
    @mock.patch('oc_adm_registry.Registry._run')
    def test_state_present(self, mock_cmd, mock_tmpfile_copy, mock_write, mock_oc_binary):
        ''' Testing state present '''
        params = {'state': 'present',
                  'debug': False,
                  'namespace': 'default',
                  'name': 'docker-registry',
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'images': None,
                  'latest_images': None,
                  'labels': {"docker-registry": "default", "another-label": "val"},
                  'ports': ['5000'],
                  'replicas': 1,
                  'selector': 'type=infra',
                  'service_account': 'registry',
                  'mount_host': None,
                  'volume_mounts': None,
                  'env_vars': {},
                  'enforce_quota': False,
                  'force': False,
                  'daemonset': False,
                  'tls_key': None,
                  'tls_certificate': None,
                  'edits': []}

        mock_cmd.side_effect = [
            (1, '', 'Error from server (NotFound): deploymentconfigs "docker-registry" not found'),
            (1, '', 'Error from server (NotFound): service "docker-registry" not found'),
            (0, RegistryTest.dry_run, ''),
            (0, '', ''),
            (0, '', ''),
        ]

        mock_tmpfile_copy.return_value = '/tmp/mocked_kubeconfig'

        mock_oc_binary.return_value = 'oc'

        results = Registry.run_ansible(params, False)

        self.assertTrue(results['changed'])
        for result in results['results']['results']:
            self.assertEqual(result['returncode'], 0)

        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'dc', 'docker-registry', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'svc', 'docker-registry', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'adm', 'registry',
                       "--labels=another-label=val,docker-registry=default",
                       '--ports=5000', '--replicas=1', '--selector=type=infra',
                       '--service-account=registry', '--dry-run=True', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'create', '-f', mock.ANY, '-n', 'default'], None),
            mock.call(['oc', 'create', '-f', mock.ANY, '-n', 'default'], None), ])

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
