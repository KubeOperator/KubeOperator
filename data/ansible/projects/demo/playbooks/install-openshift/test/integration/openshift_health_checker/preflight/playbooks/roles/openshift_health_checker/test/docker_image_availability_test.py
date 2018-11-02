from openshift_checks.docker_image_availability import DockerImageAvailability

try:
    # python3, mock is built in.
    from unittest.mock import patch
except ImportError:
    # In python2, mock is installed via pip.
    from mock import patch


def test_is_available_skopeo_image():
    result = {'rc': 0}
    # test unauth secure and insecure
    openshift_docker_insecure_registries = ['insecure.redhat.io']
    task_vars = {'openshift_docker_insecure_registries': openshift_docker_insecure_registries}
    dia = DockerImageAvailability(task_vars=task_vars)
    with patch.object(DockerImageAvailability, 'execute_module_with_retries') as m1:
        m1.return_value = result
        assert dia.is_available_skopeo_image('registry.redhat.io/openshift3/ose-pod') is True
        m1.assert_called_with('command', {'_uses_shell': True, '_raw_params': ' timeout 10 skopeo inspect --tls-verify=true  docker://registry.redhat.io/openshift3/ose-pod'})
        assert dia.is_available_skopeo_image('insecure.redhat.io/openshift3/ose-pod') is True
        m1.assert_called_with('command', {'_uses_shell': True, '_raw_params': ' timeout 10 skopeo inspect --tls-verify=false  docker://insecure.redhat.io/openshift3/ose-pod'})

    # test auth
    task_vars = {'oreg_auth_user': 'test_user', 'oreg_auth_password': 'test_pass'}
    dia = DockerImageAvailability(task_vars=task_vars)
    with patch.object(DockerImageAvailability, 'execute_module_with_retries') as m1:
        m1.return_value = result
        assert dia.is_available_skopeo_image('registry.redhat.io/openshift3/ose-pod') is True
        m1.assert_called_with('command', {'_uses_shell': True, '_raw_params': ' timeout 10 skopeo inspect --tls-verify=true --creds=test_user:test_pass docker://registry.redhat.io/openshift3/ose-pod'})


def test_available_images():
    images = ['image1', 'image2']
    dia = DockerImageAvailability(task_vars={})

    with patch('openshift_checks.docker_image_availability.DockerImageAvailability.is_available_skopeo_image') as call_mock:
        call_mock.return_value = True

        images_available = dia.available_images(images)
        assert images_available == images
