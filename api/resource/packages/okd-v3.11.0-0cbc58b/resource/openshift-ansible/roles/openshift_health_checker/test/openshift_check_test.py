import pytest

from openshift_checks import OpenShiftCheck, OpenShiftCheckException
from openshift_checks import load_checks


# Fixtures


@pytest.fixture()
def task_vars():
    return dict(foo=42, bar=dict(baz="openshift"))


@pytest.fixture(params=[
    ("notfound",),
    ("multiple", "keys", "not", "in", "task_vars"),
])
def missing_keys(request):
    return request.param


# Tests


def test_OpenShiftCheck_init():
    class TestCheck(OpenShiftCheck):
        name = "test_check"
        run = NotImplemented

    # execute_module required at init if it will be used
    with pytest.raises(RuntimeError) as excinfo:
        TestCheck().execute_module("foo")
    assert 'execute_module' in str(excinfo.value)

    execute_module = object()

    # initialize with positional argument
    check = TestCheck(execute_module)
    assert check._execute_module == execute_module

    # initialize with keyword argument
    check = TestCheck(execute_module=execute_module)
    assert check._execute_module == execute_module

    assert check.task_vars == {}
    assert check.tmp is None


def test_subclasses():
    """OpenShiftCheck.subclasses should find all subclasses recursively."""
    class TestCheck1(OpenShiftCheck):
        pass

    class TestCheck2(OpenShiftCheck):
        pass

    class TestCheck1A(TestCheck1):
        pass

    local_subclasses = set([TestCheck1, TestCheck1A, TestCheck2])
    known_subclasses = set(OpenShiftCheck.subclasses())

    assert local_subclasses - known_subclasses == set(), "local_subclasses should be a subset of known_subclasses"


def test_load_checks():
    """Loading checks should load and return Python modules."""
    modules = load_checks()
    assert modules


def dummy_check(task_vars):
    class TestCheck(OpenShiftCheck):
        name = "dummy"
        run = NotImplemented

    return TestCheck(task_vars=task_vars)


@pytest.mark.parametrize("keys,expected", [
    (("foo",), 42),
    (("bar", "baz"), "openshift"),
    (("bar.baz",), "openshift"),
])
def test_get_var_ok(task_vars, keys, expected):
    assert dummy_check(task_vars).get_var(*keys) == expected


def test_get_var_error(task_vars, missing_keys):
    with pytest.raises(OpenShiftCheckException):
        dummy_check(task_vars).get_var(*missing_keys)


def test_get_var_default(task_vars, missing_keys):
    default = object()
    assert dummy_check(task_vars).get_var(*missing_keys, default=default) == default


@pytest.mark.parametrize("keys, convert, expected", [
    (("foo",), str, "42"),
    (("foo",), float, 42.0),
    (("bar", "baz"), bool, False),
])
def test_get_var_convert(task_vars, keys, convert, expected):
    assert dummy_check(task_vars).get_var(*keys, convert=convert) == expected


def convert_oscexc(_):
    raise OpenShiftCheckException("known failure")


def convert_exc(_):
    raise Exception("failure unknown")


@pytest.mark.parametrize("keys, convert, expect_text", [
    (("bar", "baz"), int, "Cannot convert"),
    (("bar.baz",), float, "Cannot convert"),
    (("foo",), "bogus", "TypeError"),
    (("foo",), lambda a, b: 1, "TypeError"),
    (("foo",), lambda a: 1 / 0, "ZeroDivisionError"),
    (("foo",), convert_oscexc, "known failure"),
    (("foo",), convert_exc, "failure unknown"),
])
def test_get_var_convert_error(task_vars, keys, convert, expect_text):
    with pytest.raises(OpenShiftCheckException) as excinfo:
        dummy_check(task_vars).get_var(*keys, convert=convert)
    assert expect_text in str(excinfo.value)


def test_register(task_vars):
    check = dummy_check(task_vars)

    check.register_failure(OpenShiftCheckException("spam"))
    assert "spam" in str(check.failures[0])

    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.register_file("spam")  # no file contents specified
    assert "not specified" in str(excinfo.value)

    # normally execute_module registers the result file; test disabling that
    check._execute_module = lambda *args, **_: dict()
    check.execute_module("eggs", module_args={}, register=False)
    assert not check.files_to_save
