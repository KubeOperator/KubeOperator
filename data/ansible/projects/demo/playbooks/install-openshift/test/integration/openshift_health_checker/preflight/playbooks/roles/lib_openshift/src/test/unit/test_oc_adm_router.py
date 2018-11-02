#!/usr/bin/env python
'''
 Unit tests for oc adm router
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
from oc_adm_router import Router, locate_oc_binary  # noqa: E402


# pylint: disable=too-many-public-methods
class RouterTest(unittest.TestCase):
    '''
     Test class for Router
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
                "name": "router",
                "creationTimestamp": null
            }
        },
        {
            "kind": "ClusterRoleBinding",
            "apiVersion": "v1",
            "metadata": {
                "name": "router-router-role",
                "creationTimestamp": null
            },
            "userNames": [
                "system:serviceaccount:default:router"
            ],
            "groupNames": null,
            "subjects": [
                {
                    "kind": "ServiceAccount",
                    "namespace": "default",
                    "name": "router"
                }
            ],
            "roleRef": {
                "kind": "ClusterRole",
                "name": "system:router"
            }
        },
        {
            "kind": "DeploymentConfig",
            "apiVersion": "v1",
            "metadata": {
                "name": "router",
                "creationTimestamp": null,
                "labels": {
                    "router": "router"
                }
            },
            "spec": {
                "strategy": {
                    "type": "Rolling",
                    "rollingParams": {
                        "maxUnavailable": "25%",
                        "maxSurge": 0
                    },
                    "resources": {}
                },
                "triggers": [
                    {
                        "type": "ConfigChange"
                    }
                ],
                "replicas": 2,
                "test": false,
                "selector": {
                    "router": "router"
                },
                "template": {
                    "metadata": {
                        "creationTimestamp": null,
                        "labels": {
                            "router": "router"
                        }
                    },
                    "spec": {
                        "volumes": [
                            {
                                "name": "server-certificate",
                                "secret": {
                                    "secretName": "router-certs"
                                }
                            }
                        ],
                        "containers": [
                            {
                                "name": "router",
                                "image": "registry.access.redhat.com/openshift3/ose-haproxy-router:v3.5.0.39",
                                "ports": [
                                    {
                                        "containerPort": 80
                                    },
                                    {
                                        "containerPort": 443
                                    },
                                    {
                                        "name": "stats",
                                        "containerPort": 1936,
                                        "protocol": "TCP"
                                    }
                                ],
                                "env": [
                                    {
                                        "name": "DEFAULT_CERTIFICATE_DIR",
                                        "value": "/etc/pki/tls/private"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_HOSTNAME"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_HTTPS_VSERVER"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_HTTP_VSERVER"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_INSECURE",
                                        "value": "false"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_INTERNAL_ADDRESS"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_PARTITION_PATH"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_PASSWORD"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_PRIVKEY",
                                        "value": "/etc/secret-volume/router.pem"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_USERNAME"
                                    },
                                    {
                                        "name": "ROUTER_EXTERNAL_HOST_VXLAN_GW_CIDR"
                                    },
                                    {
                                        "name": "ROUTER_SERVICE_HTTPS_PORT",
                                        "value": "443"
                                    },
                                    {
                                        "name": "ROUTER_SERVICE_HTTP_PORT",
                                        "value": "80"
                                    },
                                    {
                                        "name": "ROUTER_SERVICE_NAME",
                                        "value": "router"
                                    },
                                    {
                                        "name": "ROUTER_SERVICE_NAMESPACE",
                                        "value": "default"
                                    },
                                    {
                                        "name": "ROUTER_SUBDOMAIN"
                                    },
                                    {
                                        "name": "STATS_PASSWORD",
                                        "value": "eSfUICQyyr"
                                    },
                                    {
                                        "name": "STATS_PORT",
                                        "value": "1936"
                                    },
                                    {
                                        "name": "STATS_USERNAME",
                                        "value": "admin"
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
                                        "name": "server-certificate",
                                        "readOnly": true,
                                        "mountPath": "/etc/pki/tls/private"
                                    }
                                ],
                                "livenessProbe": {
                                    "httpGet": {
                                        "path": "/healthz",
                                        "port": 1936,
                                        "host": "localhost"
                                    },
                                    "initialDelaySeconds": 10
                                },
                                "readinessProbe": {
                                    "httpGet": {
                                        "path": "/healthz",
                                        "port": 1936,
                                        "host": "localhost"
                                    },
                                    "initialDelaySeconds": 10
                                },
                                "imagePullPolicy": "IfNotPresent"
                            }
                        ],
                        "nodeSelector": {
                            "type": "infra"
                        },
                        "serviceAccountName": "router",
                        "serviceAccount": "router",
                        "hostNetwork": true,
                        "securityContext": {}
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
                "name": "router",
                "creationTimestamp": null,
                "labels": {
                    "router": "router"
                },
                "annotations": {
                    "service.alpha.openshift.io/serving-cert-secret-name": "router-certs"
                }
            },
            "spec": {
                "ports": [
                    {
                        "name": "80-tcp",
                        "port": 80,
                        "targetPort": 80
                    },
                    {
                        "name": "443-tcp",
                        "port": 443,
                        "targetPort": 443
                    },
                    {
                        "name": "1936-tcp",
                        "protocol": "TCP",
                        "port": 1936,
                        "targetPort": 1936
                    }
                ],
                "selector": {
                    "router": "router"
                }
            },
            "status": {
                "loadBalancer": {}
            }
        }
    ]
}'''

    @mock.patch('oc_adm_router.locate_oc_binary')
    @mock.patch('oc_adm_router.Utils._write')
    @mock.patch('oc_adm_router.Utils.create_tmpfile_copy')
    @mock.patch('oc_adm_router.Router._run')
    def test_state_present(self, mock_cmd, mock_tmpfile_copy, mock_write, mock_oc_binary):
        ''' Testing a create '''
        params = {'state': 'present',
                  'debug': False,
                  'namespace': 'default',
                  'name': 'router',
                  'default_cert': None,
                  'cert_file': None,
                  'key_file': None,
                  'cacert_file': None,
                  'labels': {"router": "router", "another-label": "val"},
                  'ports': ['80:80', '443:443'],
                  'images': None,
                  'latest_images': None,
                  'clusterip': None,
                  'portalip': None,
                  'session_affinity': None,
                  'service_type': None,
                  'kubeconfig': '/etc/origin/master/admin.kubeconfig',
                  'replicas': 2,
                  'selector': 'type=infra',
                  'service_account': 'router',
                  'router_type': None,
                  'host_network': None,
                  'external_host': None,
                  'external_host_vserver': None,
                  'external_host_insecure': False,
                  'external_host_partition_path': None,
                  'external_host_username': None,
                  'external_host_password': None,
                  'external_host_private_key': None,
                  'stats_user': None,
                  'stats_password': None,
                  'stats_port': 1936,
                  'edits': []}

        mock_cmd.side_effect = [
            (1, '', 'Error from server (NotFound): deploymentconfigs "router" not found'),
            (1, '', 'Error from server (NotFound): service "router" not found'),
            (1, '', 'Error from server (NotFound): serviceaccount "router" not found'),
            (1, '', 'Error from server (NotFound): secret "router-certs" not found'),
            (1, '', 'Error from server (NotFound): clsuterrolebinding "router-router-role" not found'),
            (0, RouterTest.dry_run, ''),
            (0, '', ''),
            (0, '', ''),
            (0, '', ''),
            (0, '', ''),
            (0, '', ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        mock_oc_binary.side_effect = [
            'oc',
        ]

        results = Router.run_ansible(params, False)

        self.assertTrue(results['changed'])
        for result in results['results']['results']:
            self.assertEqual(result['returncode'], 0)

        mock_cmd.assert_has_calls([
            mock.call(['oc', 'get', 'dc', 'router', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'svc', 'router', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'sa', 'router', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'secret', 'router-certs', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'get', 'clusterrolebinding', 'router-router-role', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'adm', 'router', 'router', '--external-host-insecure=False',
                       "--labels=another-label=val,router=router",
                       '--ports=80:80,443:443', '--replicas=2', '--selector=type=infra', '--service-account=router',
                       '--stats-port=1936', '--dry-run=True', '-o', 'json', '-n', 'default'], None),
            mock.call(['oc', 'create', '-f', mock.ANY, '-n', 'default'], None),
            mock.call(['oc', 'create', '-f', mock.ANY, '-n', 'default'], None),
            mock.call(['oc', 'create', '-f', mock.ANY, '-n', 'default'], None),
            mock.call(['oc', 'create', '-f', mock.ANY, '-n', 'default'], None)])

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
