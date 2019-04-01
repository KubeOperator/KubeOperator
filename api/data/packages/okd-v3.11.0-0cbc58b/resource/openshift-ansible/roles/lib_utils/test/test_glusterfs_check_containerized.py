import os
import sys

import pytest

try:
    # python3, mock is built in.
    from unittest.mock import patch
except ImportError:
    # In python2, mock is installed via pip.
    from mock import patch

MODULE_PATH = os.path.realpath(os.path.join(__file__, os.pardir, os.pardir, 'library'))
sys.path.insert(1, MODULE_PATH)

import glusterfs_check_containerized  # noqa


NODE_LIST_STD_OUT_1 = ("""
NAME                       STATUS    ROLES                  AGE       VERSION
fedora1.openshift.io   Ready     compute,infra,master   1d        v1.11.0+d4cacc0
fedora2.openshift.io   Ready     infra                  1d        v1.11.0+d4cacc0
fedora3.openshift.io   Ready     infra                  1d        v1.11.0+d4cacc0
""")

NODE_LIST_STD_OUT_2 = ("""
NAME                       STATUS    ROLES                  AGE       VERSION
fedora1.openshift.io   Ready     compute,infra,master   1d        v1.11.0+d4cacc0
fedora2.openshift.io   NotReady     infra                  1d        v1.11.0+d4cacc0
fedora3.openshift.io   Ready     infra                  1d        v1.11.0+d4cacc0
""")

NODE_LIST_STD_OUT_3 = ("""
NAME                       STATUS    ROLES                  AGE       VERSION
fedora1.openshift.io   Ready     compute,infra,master   1d        v1.11.0+d4cacc0
fedora2.openshift.io   NotReady     infra                  1d        v1.11.0+d4cacc0
fedora3.openshift.io   Invalid     infra                  1d        v1.11.0+d4cacc0
""")

POD_SELECT_STD_OUT = ("""NAME                                          READY     STATUS    RESTARTS   AGE       IP                NODE
glusterblock-storage-provisioner-dc-1-ks5zt   1/1       Running   0          1d        10.130.0.5        fedora3.openshift.io
glusterfs-storage-fzdn2                       1/1       Running   0          1d        192.168.124.175   fedora1.openshift.io
glusterfs-storage-mp9nk                       1/1       Running   4          1d        192.168.124.233   fedora2.openshift.io
glusterfs-storage-t9c6d                       1/1       Running   0          1d        192.168.124.50    fedora3.openshift.io
heketi-storage-1-rgj8b                        1/1       Running   0          1d        10.130.0.4        fedora3.openshift.io""")

# Need to ensure we have extra empty lines in this output;
# thus the quotes are one line above and below the text.
VOLUME_LIST_STDOUT = ("""
heketidbstorage
volume1
""")

VOLUME_HEAL_INFO_GOOD = ("""
Brick 192.168.124.233:/var/lib/heketi/mounts/vg_936ddf24061d55788f50496757d2f3b2/brick_9df1b6229025ea45521ab1b370d24a06/brick
Status: Connected
Number of entries: 0

Brick 192.168.124.175:/var/lib/heketi/mounts/vg_95975e77a6dc7a8e45586eac556b0f24/brick_172b6be6704a3d9f706535038f7f2e52/brick
Status: Connected
Number of entries: 0

Brick 192.168.124.50:/var/lib/heketi/mounts/vg_6523756fe1becfefd3224d3082373344/brick_359e4cf44cd1b82674f7d931cb5c481e/brick
Status: Connected
Number of entries: 0
""")

VOLUME_HEAL_INFO_BAD = ("""
Brick 192.168.124.233:/var/lib/heketi/mounts/vg_936ddf24061d55788f50496757d2f3b2/brick_9df1b6229025ea45521ab1b370d24a06/brick
Status: Connected
Number of entries: 0

Brick 192.168.124.175:/var/lib/heketi/mounts/vg_95975e77a6dc7a8e45586eac556b0f24/brick_172b6be6704a3d9f706535038f7f2e52/brick
Status: Connected
Number of entries: 0

Brick 192.168.124.50:/var/lib/heketi/mounts/vg_6523756fe1becfefd3224d3082373344/brick_359e4cf44cd1b82674f7d931cb5c481e/brick
Status: Connected
Number of entries: -
""")


class DummyModule(object):
    def exit_json(*args, **kwargs):
        return 0

    def fail_json(*args, **kwargs):
        raise Exception(kwargs['msg'])


def test_get_valid_nodes():
    with patch('glusterfs_check_containerized.call_or_fail') as call_mock:
        module = DummyModule()
        oc_exec = []
        exclude_node = "fedora1.openshift.io"

        call_mock.return_value = NODE_LIST_STD_OUT_1
        valid_nodes = glusterfs_check_containerized.get_valid_nodes(module, oc_exec, exclude_node)
        assert valid_nodes == ['fedora2.openshift.io', 'fedora3.openshift.io']

        call_mock.return_value = NODE_LIST_STD_OUT_2
        valid_nodes = glusterfs_check_containerized.get_valid_nodes(module, oc_exec, exclude_node)
        assert valid_nodes == ['fedora3.openshift.io']

        call_mock.return_value = NODE_LIST_STD_OUT_3
        with pytest.raises(Exception) as err:
            valid_nodes = glusterfs_check_containerized.get_valid_nodes(module, oc_exec, exclude_node)
        assert 'Exception: Unable to find suitable node in get nodes output' in str(err)


def test_select_pod():
    with patch('glusterfs_check_containerized.call_or_fail') as call_mock:
        module = DummyModule()
        oc_exec = []
        cluster_name = "storage"
        valid_nodes = ["fedora2.openshift.io", "fedora3.openshift.io"]
        call_mock.return_value = POD_SELECT_STD_OUT
        # Should select first valid podname in call_or_fail output.
        pod_name = glusterfs_check_containerized.select_pod(module, oc_exec, cluster_name, valid_nodes)
        assert pod_name == 'glusterfs-storage-mp9nk'
        with pytest.raises(Exception) as err:
            pod_name = glusterfs_check_containerized.select_pod(module, oc_exec, "does not exist", valid_nodes)
        assert 'Exception: Unable to find suitable pod in get pods output' in str(err)


def test_get_volume_list():
    with patch('glusterfs_check_containerized.call_or_fail') as call_mock:
        module = DummyModule()
        oc_exec = []
        pod_name = ''
        call_mock.return_value = VOLUME_LIST_STDOUT
        volume_list = glusterfs_check_containerized.get_volume_list(module, oc_exec, pod_name)
        assert volume_list == ['heketidbstorage', 'volume1']


def test_check_volume_health_info():
    with patch('glusterfs_check_containerized.call_or_fail') as call_mock:
        module = DummyModule()
        oc_exec = []
        pod_name = ''
        volume = 'somevolume'
        call_mock.return_value = VOLUME_HEAL_INFO_GOOD
        # this should just complete quietly.
        glusterfs_check_containerized.check_volume_health_info(module, oc_exec, pod_name, volume)

        call_mock.return_value = VOLUME_HEAL_INFO_BAD
        expected_error = 'volume {} is not ready'.format(volume)
        with pytest.raises(Exception) as err:
            glusterfs_check_containerized.check_volume_health_info(module, oc_exec, pod_name, volume)
        assert expected_error in str(err)


if __name__ == '__main__':
    test_get_valid_nodes()
    test_select_pod()
    test_get_volume_list()
    test_check_volume_health_info()
