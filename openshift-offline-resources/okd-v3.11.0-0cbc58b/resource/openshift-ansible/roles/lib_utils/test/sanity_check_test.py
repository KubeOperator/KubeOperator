import os
import pytest
import sys

from ansible.playbook.play_context import PlayContext
from ansible.template import Templar
from ansible import errors

sys.path.insert(1, os.path.join(os.path.dirname(__file__), os.pardir, "action_plugins"))
from sanity_checks import ActionModule  # noqa: E402


@pytest.mark.parametrize('hostvars, host, varname, result', [
    ({"example.com": {"param": 3.11}}, "example.com", "param", 3.11),
    ({"example.com": {"param": 3.11}}, "example.com", "another_param", None)
])
def test_template_var(hostvars, host, varname, result):
    task = FakeTask('sanity_checks', {'checks': []})
    plugin = ActionModule(task, None, PlayContext(), None, Templar(None, None, None), None)
    check = plugin.template_var(hostvars, host, varname)
    assert check == result


@pytest.mark.parametrize('hostvars, host, result', [
    ({"example.com": {"openshift_pkg_version": "-3.6.0"}}, "example.com", None),
    ({"example.com": {"openshift_pkg_version": "-3.7.0-0.126.0.git.0.9351aae.el7"}}, "example.com", None),
    ({"example.com": {"openshift_pkg_version": "-3.9.0-2.fc28"}}, "example.com", None),
    ({"example.com": {"openshift_pkg_version": "-3.11*"}}, "example.com", None),
    ({"example.com": {"openshift_pkg_version": "-3"}}, "example.com", None),
])
def test_valid_check_pkg_version_format(hostvars, host, result):
    task = FakeTask('sanity_checks', {'checks': []})
    plugin = ActionModule(task, None, PlayContext(), None, Templar(None, None, None), None)
    check = plugin.check_pkg_version_format(hostvars, host)
    assert check == result


@pytest.mark.parametrize('hostvars, host, result', [
    ({"example.com": {"openshift_pkg_version": "3.11.0"}}, "example.com", None),
    ({"example.com": {"openshift_pkg_version": "v3.11.0"}}, "example.com", None),
])
def test_invalid_check_pkg_version_format(hostvars, host, result):
    with pytest.raises(errors.AnsibleModuleError):
        task = FakeTask('sanity_checks', {'checks': []})
        plugin = ActionModule(task, None, PlayContext(), None, Templar(None, None, None), None)
        plugin.check_pkg_version_format(hostvars, host)


@pytest.mark.parametrize('hostvars, host, result', [
    ({"example.com": {"openshift_release": "v3"}}, "example.com", None),
    ({"example.com": {"openshift_release": "v3.11"}}, "example.com", None),
    ({"example.com": {"openshift_release": "v3.11.0"}}, "example.com", None),
    ({"example.com": {"openshift_release": "3.11"}}, "example.com", None),
])
def test_valid_check_release_format(hostvars, host, result):
    task = FakeTask('sanity_checks', {'checks': []})
    plugin = ActionModule(task, None, PlayContext(), None, Templar(None, None, None), None)
    check = plugin.check_release_format(hostvars, host)
    assert check == result


@pytest.mark.parametrize('hostvars, host, result', [
    ({"example.com": {"openshift_release": "-3.11.0"}}, "example.com", None),
    ({"example.com": {"openshift_release": "-3.7.0-0.126.0.git.0.9351aae.el7"}}, "example.com", None),
    ({"example.com": {"openshift_release": "3.1.2.3"}}, "example.com", None),
])
def test_invalid_check_release_format(hostvars, host, result):
    with pytest.raises(errors.AnsibleModuleError):
        task = FakeTask('sanity_checks', {'checks': []})
        plugin = ActionModule(task, None, PlayContext(), None, Templar(None, None, None), None)
        plugin.check_release_format(hostvars, host)


@pytest.mark.parametrize('hostvars, host, result', [
    ({"example.com": {"openshift_builddefaults_json": "{}"}}, "example.com", None),
    ({"example.com": {"openshift_builddefaults_json": '[]'}}, "example.com", None),
    ({"example.com": {"openshift_builddefaults_json": '{"a": []}'}}, "example.com", None),
    ({"example.com": {"openshift_builddefaults_json": '{"a": [], "b": "c"}'}}, "example.com", None),
    ({"example.com": {"openshift_builddefaults_json": '{"a": [], "b": {"c": "d"}}'}}, "example.com", None),
    ({"example.com": {"openshift_builddefaults_json": '["a", "b", "c"]'}}, "example.com", None),
    ({"example.com": {"NOT_IN_JSON_FORMAT_VARIABLES": '{"invalid"}'}}, "example.com", None),
])
def test_valid_valid_json_format_vars(hostvars, host, result):
    task = FakeTask('sanity_checks', {'checks': []})
    plugin = ActionModule(task, None, PlayContext(), None, Templar(None, None, None), None)
    check = plugin.validate_json_format_vars(hostvars, host)
    assert check == result


@pytest.mark.parametrize('hostvars, host, result', [
    ({"example.com": {"openshift_builddefaults_json": '{"a"}'}}, "example.com", None),
    ({"example.com": {"openshift_builddefaults_json": '{"a": { '}}, "example.com", None),
    ({"example.com": {"openshift_builddefaults_json": '{"a": [ }'}}, "example.com", None),
])
def test_invalid_valid_json_format_vars(hostvars, host, result):
    with pytest.raises(errors.AnsibleModuleError):
        task = FakeTask('sanity_checks', {'checks': []})
        plugin = ActionModule(task, None, PlayContext(), None, Templar(None, None, None), None)
        plugin.validate_json_format_vars(hostvars, host)


def fake_execute_module(*args):
    raise AssertionError('this function should not be called')


class FakeTask(object):
    def __init__(self, action, args):
        self.action = action
        self.args = args
        self.async = 0
