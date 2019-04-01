import pytest

from openshift_checks.logging.curator import Curator, OpenShiftCheckException


plain_curator_pod = {
    "metadata": {
        "labels": {"component": "curator", "deploymentconfig": "logging-curator"},
        "name": "logging-curator-1",
    },
    "status": {
        "containerStatuses": [{"ready": True}],
        "conditions": [{"status": "True", "type": "Ready"}],
        "podIP": "10.10.10.10",
    }
}

not_running_curator_pod = {
    "metadata": {
        "labels": {"component": "curator", "deploymentconfig": "logging-curator"},
        "name": "logging-curator-2",
    },
    "status": {
        "containerStatuses": [{"ready": False}],
        "conditions": [{"status": "False", "type": "Ready"}],
        "podIP": "10.10.10.10",
    }
}


def test_get_curator_pods():
    check = Curator()
    check.get_pods_for_component = lambda *_: [plain_curator_pod]
    result = check.run()
    assert "failed" not in result or not result["failed"]


@pytest.mark.parametrize('pods, expect_error', [
    (
        [],
        'MissingComponentPods',
    ),
    (
        [not_running_curator_pod],
        'CuratorNotRunning',
    ),
    (
        [plain_curator_pod, plain_curator_pod],
        'TooManyCurators',
    ),
])
def test_get_curator_pods_fail(pods, expect_error):
    check = Curator()
    check.get_pods_for_component = lambda *_: pods
    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.run()
    assert excinfo.value.name == expect_error
