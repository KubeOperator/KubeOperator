import pytest

from ansible.playbook.play_context import PlayContext

from openshift_health_check import ActionModule, resolve_checks
from openshift_health_check import copy_remote_file_to_dir, write_result_to_output_dir, write_to_output_file
from openshift_checks import OpenShiftCheckException, FileToSave


def fake_check(name='fake_check', tags=None, is_active=True, run_return=None, run_exception=None,
               run_logs=None, run_files=None, changed=False, get_var_return=None):
    """Returns a new class that is compatible with OpenShiftCheck for testing."""

    _name, _tags = name, tags

    class FakeCheck(object):
        name = _name
        tags = _tags or []

        def __init__(self, **_):
            self.changed = False
            self.failures = []
            self.logs = run_logs or []
            self.files_to_save = run_files or []

        def is_active(self):
            if isinstance(is_active, Exception):
                raise is_active
            return is_active

        def run(self):
            self.changed = changed
            if run_exception is not None:
                raise run_exception
            return run_return

        def get_var(*args, **_):
            return get_var_return

        def register_failure(self, exc):
            self.failures.append(OpenShiftCheckException(str(exc)))
            return

    return FakeCheck


# Fixtures


@pytest.fixture
def plugin():
    task = FakeTask('openshift_health_check', {'checks': ['fake_check']})
    plugin = ActionModule(task, None, PlayContext(), None, None, None)
    return plugin


class FakeTask(object):
    def __init__(self, action, args):
        self.action = action
        self.args = args
        self.async = 0


@pytest.fixture
def task_vars():
    return dict(openshift=dict(), ansible_host='unit-test-host')


# Assertion helpers


def failed(result, msg_has=None):
    if msg_has is not None:
        assert 'msg' in result
        for term in msg_has:
            assert term.lower() in result['msg'].lower()
    return result.get('failed', False)


def changed(result):
    return result.get('changed', False)


# tests whether task is skipped, not individual checks
def skipped(result):
    return result.get('skipped', False)


# Tests


@pytest.mark.parametrize('task_vars', [
    None,
    {},
])
def test_action_plugin_missing_openshift_facts(plugin, task_vars, monkeypatch):
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {})
    monkeypatch.setattr('openshift_health_check.resolve_checks', lambda *args: ['fake_check'])
    result = plugin.run(tmp=None, task_vars=task_vars)

    assert failed(result, msg_has=['openshift_facts'])


def test_action_plugin_cannot_load_checks_with_the_same_name(plugin, task_vars, monkeypatch):
    FakeCheck1 = fake_check('duplicate_name')
    FakeCheck2 = fake_check('duplicate_name')
    checks = [FakeCheck1, FakeCheck2]
    monkeypatch.setattr('openshift_checks.OpenShiftCheck.subclasses', classmethod(lambda cls: checks))

    result = plugin.run(tmp=None, task_vars=task_vars)

    assert failed(result, msg_has=['duplicate', 'duplicate_name', 'FakeCheck'])


@pytest.mark.parametrize('is_active, skipped_reason', [
    (False, "Not active for this host"),
    (Exception("borked"), "exception"),
])
def test_action_plugin_skip_non_active_checks(is_active, skipped_reason, plugin, task_vars, monkeypatch):
    checks = [fake_check(is_active=is_active)]
    monkeypatch.setattr('openshift_checks.OpenShiftCheck.subclasses', classmethod(lambda cls: checks))

    result = plugin.run(tmp=None, task_vars=task_vars)

    assert result['checks']['fake_check'].get('skipped')
    assert skipped_reason in result['checks']['fake_check'].get('skipped_reason')
    assert not failed(result)
    assert not changed(result)
    assert not skipped(result)


@pytest.mark.parametrize('to_disable', [
    'fake_check',
    ['fake_check', 'spam'],
    '*,spam,eggs',
])
def test_action_plugin_skip_disabled_checks(to_disable, plugin, task_vars, monkeypatch):
    checks = [fake_check('fake_check', is_active=True)]
    monkeypatch.setattr('openshift_checks.OpenShiftCheck.subclasses', classmethod(lambda cls: checks))

    task_vars['openshift_disable_check'] = to_disable
    result = plugin.run(tmp=None, task_vars=task_vars)

    assert result['checks']['fake_check'] == dict(skipped=True, skipped_reason="Disabled by user request")
    assert not failed(result)
    assert not changed(result)
    assert not skipped(result)


def test_action_plugin_run_list_checks(monkeypatch):
    task = FakeTask('openshift_health_check', {'checks': []})
    plugin = ActionModule(task, None, PlayContext(), None, None, None)
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {})
    result = plugin.run()

    assert failed(result, msg_has="Available checks")
    assert not changed(result)
    assert not skipped(result)


def test_action_plugin_run_check_ok(plugin, task_vars, monkeypatch):
    check_return_value = {'ok': 'test'}
    check_class = fake_check(run_return=check_return_value, run_files=[None])
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {'fake_check': check_class()})
    monkeypatch.setattr('openshift_health_check.resolve_checks', lambda *args: ['fake_check'])

    result = plugin.run(tmp=None, task_vars=task_vars)

    assert result['checks']['fake_check'] == check_return_value
    assert not failed(result)
    assert not changed(result)
    assert not skipped(result)


def test_action_plugin_run_check_changed(plugin, task_vars, monkeypatch):
    check_return_value = {'ok': 'test'}
    check_class = fake_check(run_return=check_return_value, changed=True)
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {'fake_check': check_class()})
    monkeypatch.setattr('openshift_health_check.resolve_checks', lambda *args: ['fake_check'])

    result = plugin.run(tmp=None, task_vars=task_vars)

    assert result['checks']['fake_check'] == check_return_value
    assert changed(result['checks']['fake_check'])
    assert not failed(result)
    assert changed(result)
    assert not skipped(result)


def test_action_plugin_run_check_fail(plugin, task_vars, monkeypatch):
    check_return_value = {'failed': True, 'msg': 'this is a failure'}
    check_class = fake_check(run_return=check_return_value)
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {'fake_check': check_class()})
    monkeypatch.setattr('openshift_health_check.resolve_checks', lambda *args: ['fake_check'])

    result = plugin.run(tmp=None, task_vars=task_vars)

    assert result['checks']['fake_check'] == check_return_value
    assert failed(result, msg_has=['failed'])
    assert not changed(result)
    assert not skipped(result)


@pytest.mark.parametrize('exc_class, expect_traceback', [
    (OpenShiftCheckException, False),
    (Exception, True),
])
def test_action_plugin_run_check_exception(plugin, task_vars, exc_class, expect_traceback, monkeypatch):
    exception_msg = 'fake check has an exception'
    run_exception = exc_class(exception_msg)
    check_class = fake_check(run_exception=run_exception, changed=True)
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {'fake_check': check_class()})
    monkeypatch.setattr('openshift_health_check.resolve_checks', lambda *args: ['fake_check'])

    result = plugin.run(tmp=None, task_vars=task_vars)

    assert failed(result['checks']['fake_check'], msg_has=exception_msg)
    assert expect_traceback == ("Traceback" in result['checks']['fake_check']['msg'])
    assert failed(result, msg_has=['failed'])
    assert changed(result['checks']['fake_check'])
    assert changed(result)
    assert not skipped(result)


def test_action_plugin_run_check_output_dir(plugin, task_vars, tmpdir, monkeypatch):
    check_class = fake_check(
        run_return={},
        run_logs=[('thing', 'note')],
        run_files=[
            FileToSave('save.file', 'contents', None),
            FileToSave('save.file', 'duplicate', None),
            FileToSave('copy.file', None, 'foo'),  # note: copy runs execute_module => exception
        ],
    )
    task_vars['openshift_checks_output_dir'] = str(tmpdir)
    check_class.get_var = lambda self, name, **_: task_vars.get(name)
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {'fake_check': check_class()})
    monkeypatch.setattr('openshift_health_check.resolve_checks', lambda *args: ['fake_check'])

    plugin.run(tmp=None, task_vars=task_vars)
    assert any(path.basename == task_vars['ansible_host'] for path in tmpdir.listdir())
    assert any(path.basename == 'fake_check.log.json' for path in tmpdir.visit())
    assert any(path.basename == 'save.file' for path in tmpdir.visit())
    assert any(path.basename == 'save.file.2' for path in tmpdir.visit())


def test_action_plugin_resolve_checks_exception(plugin, task_vars, monkeypatch):
    monkeypatch.setattr(plugin, 'load_known_checks', lambda *_: {})

    result = plugin.run(tmp=None, task_vars=task_vars)

    assert failed(result, msg_has=['unknown', 'name'])
    assert not changed(result)
    assert not skipped(result)


@pytest.mark.parametrize('names,all_checks,expected', [
    ([], [], set()),
    (
        ['a', 'b'],
        [
            fake_check('a'),
            fake_check('b'),
        ],
        set(['a', 'b']),
    ),
    (
        ['a', 'b', '@group'],
        [
            fake_check('from_group_1', ['group', 'another_group']),
            fake_check('not_in_group', ['another_group']),
            fake_check('from_group_2', ['preflight', 'group']),
            fake_check('a'),
            fake_check('b'),
        ],
        set(['a', 'b', 'from_group_1', 'from_group_2']),
    ),
])
def test_resolve_checks_ok(names, all_checks, expected):
    assert resolve_checks(names, all_checks) == expected


@pytest.mark.parametrize('names,all_checks,words_in_exception', [
    (
        ['testA', 'testB'],
        [],
        ['check', 'name', 'testA', 'testB'],
    ),
    (
        ['@group'],
        [],
        ['tag', 'name', 'group'],
    ),
    (
        ['testA', 'testB', '@group'],
        [],
        ['check', 'name', 'testA', 'testB', 'tag', 'group'],
    ),
    (
        ['testA', 'testB', '@group'],
        [
            fake_check('from_group_1', ['group', 'another_group']),
            fake_check('not_in_group', ['another_group']),
            fake_check('from_group_2', ['preflight', 'group']),
        ],
        ['check', 'name', 'testA', 'testB'],
    ),
])
def test_resolve_checks_failure(names, all_checks, words_in_exception):
    with pytest.raises(Exception) as excinfo:
        resolve_checks(names, all_checks)
    for word in words_in_exception:
        assert word in str(excinfo.value)


@pytest.mark.parametrize('give_output_dir, result, expect_file', [
    (False, None, False),
    (True, dict(content="c3BhbQo=", encoding="base64"), True),
    (True, dict(content="encoding error", encoding="base64"), False),
    (True, dict(content="spam", no_encoding=None), True),
    (True, dict(failed=True, msg="could not slurp"), False),
])
def test_copy_remote_file_to_dir(give_output_dir, result, expect_file, tmpdir):
    check = fake_check()()
    check.execute_module = lambda *args, **_: result
    copy_remote_file_to_dir(check, "remote_file", str(tmpdir) if give_output_dir else "", "local_file")
    assert expect_file == any(path.basename == "local_file" for path in tmpdir.listdir())


def test_write_to_output_exceptions(tmpdir, monkeypatch, capsys):

    class Spam(object):
        def __str__(self):
            raise Exception("break str")

    test = {1: object(), 2: Spam()}
    test[3] = test
    write_result_to_output_dir(str(tmpdir), test)
    assert "Error writing" in test["output_files"]

    output_dir = tmpdir.join("eggs")
    output_dir.write("spam")  # so now it's not a dir
    write_to_output_file(str(output_dir), "somefile", "somedata")
    assert "Could not write" in capsys.readouterr()[1]

    monkeypatch.setattr("openshift_health_check.prepare_output_dir", lambda *_: False)
    write_result_to_output_dir(str(tmpdir), test)
    assert "Error creating" in test["output_files"]
