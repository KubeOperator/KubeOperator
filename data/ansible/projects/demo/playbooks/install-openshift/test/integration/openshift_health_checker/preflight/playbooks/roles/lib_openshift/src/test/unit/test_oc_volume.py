'''
 Unit tests for oc volume
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
# pylint: disable=import-error
# place class in our python path
module_path = os.path.join('/'.join(os.path.realpath(__file__).split('/')[:-4]), 'library')  # noqa: E501
sys.path.insert(0, module_path)
from oc_volume import OCVolume, locate_oc_binary  # noqa: E402


class OCVolumeTest(unittest.TestCase):
    '''
     Test class for OCVolume
    '''
    params = {'name': 'oso-rhel7-zagg-web',
              'kubeconfig': '/etc/origin/master/admin.kubeconfig',
              'namespace': 'test',
              'labels': None,
              'state': 'present',
              'kind': 'dc',
              'mount_path': None,
              'secret_name': None,
              'mount_type': 'pvc',
              'claim_name': 'testclaim',
              'claim_size': '1G',
              'configmap_name': None,
              'vol_name': 'test-volume',
              'debug': False}

    @mock.patch('oc_volume.Utils.create_tmpfile_copy')
    @mock.patch('oc_volume.OCVolume._run')
    def test_create_pvc(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a label list '''
        params = copy.deepcopy(OCVolumeTest.params)

        dc = '''{
                "kind": "DeploymentConfig",
                "apiVersion": "v1",
                "metadata": {
                    "name": "oso-rhel7-zagg-web",
                    "namespace": "new-monitoring",
                    "selfLink": "/oapi/v1/namespaces/new-monitoring/deploymentconfigs/oso-rhel7-zagg-web",
                    "uid": "f56e9dd2-7c13-11e6-b046-0e8844de0587",
                    "resourceVersion": "137095771",
                    "generation": 4,
                    "creationTimestamp": "2016-09-16T13:46:24Z",
                    "labels": {
                        "app": "oso-rhel7-ops-base",
                        "name": "oso-rhel7-zagg-web"
                    },
                    "annotations": {
                        "openshift.io/generated-by": "OpenShiftNewApp"
                    }
                },
                "spec": {
                    "strategy": {
                        "type": "Rolling",
                        "rollingParams": {
                            "updatePeriodSeconds": 1,
                            "intervalSeconds": 1,
                            "timeoutSeconds": 600,
                            "maxUnavailable": "25%",
                            "maxSurge": "25%"
                        },
                        "resources": {}
                    },
                    "triggers": [
                        {
                            "type": "ConfigChange"
                        },
                        {
                            "type": "ImageChange",
                            "imageChangeParams": {
                                "automatic": true,
                                "containerNames": [
                                    "oso-rhel7-zagg-web"
                                ],
                                "from": {
                                    "kind": "ImageStreamTag",
                                    "namespace": "new-monitoring",
                                    "name": "oso-rhel7-zagg-web:latest"
                                },
                                "lastTriggeredImage": "notused"
                            }
                        }
                    ],
                    "replicas": 10,
                    "test": false,
                    "selector": {
                        "deploymentconfig": "oso-rhel7-zagg-web"
                    },
                    "template": {
                        "metadata": {
                            "creationTimestamp": null,
                            "labels": {
                                "app": "oso-rhel7-ops-base",
                                "deploymentconfig": "oso-rhel7-zagg-web"
                            },
                            "annotations": {
                                "openshift.io/generated-by": "OpenShiftNewApp"
                            }
                        },
                        "spec": {
                            "volumes": [
                                {
                                    "name": "monitoring-secrets",
                                    "secret": {
                                        "secretName": "monitoring-secrets"
                                    }
                                }
                            ],
                            "containers": [
                                {
                                    "name": "oso-rhel7-zagg-web",
                                    "image": "notused",
                                    "resources": {},
                                    "volumeMounts": [
                                        {
                                            "name": "monitoring-secrets",
                                            "mountPath": "/secrets"
                                        }
                                    ],
                                    "terminationMessagePath": "/dev/termination-log",
                                    "imagePullPolicy": "Always",
                                    "securityContext": {
                                        "capabilities": {},
                                        "privileged": false
                                    }
                                }
                            ],
                            "restartPolicy": "Always",
                            "terminationGracePeriodSeconds": 30,
                            "dnsPolicy": "ClusterFirst",
                            "securityContext": {}
                        }
                    }
                }
            }'''

        post_dc = '''{
                "kind": "DeploymentConfig",
                "apiVersion": "v1",
                "metadata": {
                    "name": "oso-rhel7-zagg-web",
                    "namespace": "new-monitoring",
                    "selfLink": "/oapi/v1/namespaces/new-monitoring/deploymentconfigs/oso-rhel7-zagg-web",
                    "uid": "f56e9dd2-7c13-11e6-b046-0e8844de0587",
                    "resourceVersion": "137095771",
                    "generation": 4,
                    "creationTimestamp": "2016-09-16T13:46:24Z",
                    "labels": {
                        "app": "oso-rhel7-ops-base",
                        "name": "oso-rhel7-zagg-web"
                    },
                    "annotations": {
                        "openshift.io/generated-by": "OpenShiftNewApp"
                    }
                },
                "spec": {
                    "strategy": {
                        "type": "Rolling",
                        "rollingParams": {
                            "updatePeriodSeconds": 1,
                            "intervalSeconds": 1,
                            "timeoutSeconds": 600,
                            "maxUnavailable": "25%",
                            "maxSurge": "25%"
                        },
                        "resources": {}
                    },
                    "triggers": [
                        {
                            "type": "ConfigChange"
                        },
                        {
                            "type": "ImageChange",
                            "imageChangeParams": {
                                "automatic": true,
                                "containerNames": [
                                    "oso-rhel7-zagg-web"
                                ],
                                "from": {
                                    "kind": "ImageStreamTag",
                                    "namespace": "new-monitoring",
                                    "name": "oso-rhel7-zagg-web:latest"
                                },
                                "lastTriggeredImage": "notused"
                            }
                        }
                    ],
                    "replicas": 10,
                    "test": false,
                    "selector": {
                        "deploymentconfig": "oso-rhel7-zagg-web"
                    },
                    "template": {
                        "metadata": {
                            "creationTimestamp": null,
                            "labels": {
                                "app": "oso-rhel7-ops-base",
                                "deploymentconfig": "oso-rhel7-zagg-web"
                            },
                            "annotations": {
                                "openshift.io/generated-by": "OpenShiftNewApp"
                            }
                        },
                        "spec": {
                            "volumes": [
                                {
                                    "name": "monitoring-secrets",
                                    "secret": {
                                        "secretName": "monitoring-secrets"
                                    }
                                },
                                {
                                    "name": "test-volume",
                                    "persistentVolumeClaim": {
                                        "claimName": "testclass",
                                        "claimSize": "1G"
                                    }
                                }
                            ],
                            "containers": [
                                {
                                    "name": "oso-rhel7-zagg-web",
                                    "image": "notused",
                                    "resources": {},
                                    "volumeMounts": [
                                        {
                                            "name": "monitoring-secrets",
                                            "mountPath": "/secrets"
                                        },
                                        {
                                            "name": "test-volume",
                                            "mountPath": "/data"
                                        }
                                    ],
                                    "terminationMessagePath": "/dev/termination-log",
                                    "imagePullPolicy": "Always",
                                    "securityContext": {
                                        "capabilities": {},
                                        "privileged": false
                                    }
                                }
                            ],
                            "restartPolicy": "Always",
                            "terminationGracePeriodSeconds": 30,
                            "dnsPolicy": "ClusterFirst",
                            "securityContext": {}
                        }
                    }
                }
            }'''

        mock_cmd.side_effect = [
            (0, dc, ''),
            (0, dc, ''),
            (0, '', ''),
            (0, post_dc, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCVolume.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertTrue(results['results']['results'][-1]['name'] == 'test-volume')

    @mock.patch('oc_volume.Utils.create_tmpfile_copy')
    @mock.patch('oc_volume.OCVolume._run')
    def test_create_configmap(self, mock_cmd, mock_tmpfile_copy):
        ''' Testing a label list '''
        params = copy.deepcopy(OCVolumeTest.params)
        params.update({'mount_path': '/configmap',
                       'mount_type': 'configmap',
                       'configmap_name': 'configtest',
                       'vol_name': 'configvol'})

        dc = '''{
                "kind": "DeploymentConfig",
                "apiVersion": "v1",
                "metadata": {
                    "name": "oso-rhel7-zagg-web",
                    "namespace": "new-monitoring",
                    "selfLink": "/oapi/v1/namespaces/new-monitoring/deploymentconfigs/oso-rhel7-zagg-web",
                    "uid": "f56e9dd2-7c13-11e6-b046-0e8844de0587",
                    "resourceVersion": "137095771",
                    "generation": 4,
                    "creationTimestamp": "2016-09-16T13:46:24Z",
                    "labels": {
                        "app": "oso-rhel7-ops-base",
                        "name": "oso-rhel7-zagg-web"
                    },
                    "annotations": {
                        "openshift.io/generated-by": "OpenShiftNewApp"
                    }
                },
                "spec": {
                    "strategy": {
                        "type": "Rolling",
                        "rollingParams": {
                            "updatePeriodSeconds": 1,
                            "intervalSeconds": 1,
                            "timeoutSeconds": 600,
                            "maxUnavailable": "25%",
                            "maxSurge": "25%"
                        },
                        "resources": {}
                    },
                    "triggers": [
                        {
                            "type": "ConfigChange"
                        },
                        {
                            "type": "ImageChange",
                            "imageChangeParams": {
                                "automatic": true,
                                "containerNames": [
                                    "oso-rhel7-zagg-web"
                                ],
                                "from": {
                                    "kind": "ImageStreamTag",
                                    "namespace": "new-monitoring",
                                    "name": "oso-rhel7-zagg-web:latest"
                                },
                                "lastTriggeredImage": "notused"
                            }
                        }
                    ],
                    "replicas": 10,
                    "test": false,
                    "selector": {
                        "deploymentconfig": "oso-rhel7-zagg-web"
                    },
                    "template": {
                        "metadata": {
                            "creationTimestamp": null,
                            "labels": {
                                "app": "oso-rhel7-ops-base",
                                "deploymentconfig": "oso-rhel7-zagg-web"
                            },
                            "annotations": {
                                "openshift.io/generated-by": "OpenShiftNewApp"
                            }
                        },
                        "spec": {
                            "volumes": [
                                {
                                    "name": "monitoring-secrets",
                                    "secret": {
                                        "secretName": "monitoring-secrets"
                                    }
                                }
                            ],
                            "containers": [
                                {
                                    "name": "oso-rhel7-zagg-web",
                                    "image": "notused",
                                    "resources": {},
                                    "volumeMounts": [
                                        {
                                            "name": "monitoring-secrets",
                                            "mountPath": "/secrets"
                                        }
                                    ],
                                    "terminationMessagePath": "/dev/termination-log",
                                    "imagePullPolicy": "Always",
                                    "securityContext": {
                                        "capabilities": {},
                                        "privileged": false
                                    }
                                }
                            ],
                            "restartPolicy": "Always",
                            "terminationGracePeriodSeconds": 30,
                            "dnsPolicy": "ClusterFirst",
                            "securityContext": {}
                        }
                    }
                }
            }'''

        post_dc = '''{
                "kind": "DeploymentConfig",
                "apiVersion": "v1",
                "metadata": {
                    "name": "oso-rhel7-zagg-web",
                    "namespace": "new-monitoring",
                    "selfLink": "/oapi/v1/namespaces/new-monitoring/deploymentconfigs/oso-rhel7-zagg-web",
                    "uid": "f56e9dd2-7c13-11e6-b046-0e8844de0587",
                    "resourceVersion": "137095771",
                    "generation": 4,
                    "creationTimestamp": "2016-09-16T13:46:24Z",
                    "labels": {
                        "app": "oso-rhel7-ops-base",
                        "name": "oso-rhel7-zagg-web"
                    },
                    "annotations": {
                        "openshift.io/generated-by": "OpenShiftNewApp"
                    }
                },
                "spec": {
                    "strategy": {
                        "type": "Rolling",
                        "rollingParams": {
                            "updatePeriodSeconds": 1,
                            "intervalSeconds": 1,
                            "timeoutSeconds": 600,
                            "maxUnavailable": "25%",
                            "maxSurge": "25%"
                        },
                        "resources": {}
                    },
                    "triggers": [
                        {
                            "type": "ConfigChange"
                        },
                        {
                            "type": "ImageChange",
                            "imageChangeParams": {
                                "automatic": true,
                                "containerNames": [
                                    "oso-rhel7-zagg-web"
                                ],
                                "from": {
                                    "kind": "ImageStreamTag",
                                    "namespace": "new-monitoring",
                                    "name": "oso-rhel7-zagg-web:latest"
                                },
                                "lastTriggeredImage": "notused"
                            }
                        }
                    ],
                    "replicas": 10,
                    "test": false,
                    "selector": {
                        "deploymentconfig": "oso-rhel7-zagg-web"
                    },
                    "template": {
                        "metadata": {
                            "creationTimestamp": null,
                            "labels": {
                                "app": "oso-rhel7-ops-base",
                                "deploymentconfig": "oso-rhel7-zagg-web"
                            },
                            "annotations": {
                                "openshift.io/generated-by": "OpenShiftNewApp"
                            }
                        },
                        "spec": {
                            "volumes": [
                                {
                                    "name": "monitoring-secrets",
                                    "secret": {
                                        "secretName": "monitoring-secrets"
                                    }
                                },
                                {
                                    "name": "configvol",
                                    "configMap": {
                                        "name": "configtest"
                                    }
                                }
                            ],
                            "containers": [
                                {
                                    "name": "oso-rhel7-zagg-web",
                                    "image": "notused",
                                    "resources": {},
                                    "volumeMounts": [
                                        {
                                            "name": "monitoring-secrets",
                                            "mountPath": "/secrets"
                                        },
                                        {
                                            "name": "configvol",
                                            "mountPath": "/configmap"
                                        }
                                    ],
                                    "terminationMessagePath": "/dev/termination-log",
                                    "imagePullPolicy": "Always",
                                    "securityContext": {
                                        "capabilities": {},
                                        "privileged": false
                                    }
                                }
                            ],
                            "restartPolicy": "Always",
                            "terminationGracePeriodSeconds": 30,
                            "dnsPolicy": "ClusterFirst",
                            "securityContext": {}
                        }
                    }
                }
            }'''

        mock_cmd.side_effect = [
            (0, dc, ''),
            (0, dc, ''),
            (0, '', ''),
            (0, post_dc, ''),
        ]

        mock_tmpfile_copy.side_effect = [
            '/tmp/mocked_kubeconfig',
        ]

        results = OCVolume.run_ansible(params, False)

        self.assertTrue(results['changed'])
        self.assertTrue(results['results']['results'][-1]['name'] == 'configvol')

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
