import pytest

from openshift_checks.diagnostics import DiagnosticCheck, OpenShiftCheckException


@pytest.fixture()
def task_vars():
    return dict(
        openshift=dict(
            common=dict(config_base="/etc/origin/")
        )
    )


def test_module_succeeds(task_vars):
    check = DiagnosticCheck(lambda *_: {"result": "success"}, task_vars)
    check.is_first_master = lambda: True
    assert check.is_active()
    check.exec_diagnostic("spam")
    assert not check.failures


def test_oc_not_there(task_vars):
    def exec_module(*_):
        return {"failed": True, "result": "[Errno 2] No such file or directory"}

    check = DiagnosticCheck(exec_module, task_vars)
    with pytest.raises(OpenShiftCheckException) as excinfo:
        check.exec_diagnostic("spam")
    assert excinfo.value.name == "OcNotFound"


def test_module_fails(task_vars):
    def exec_module(*_):
        return {"failed": True, "result": "something broke"}

    check = DiagnosticCheck(exec_module, task_vars)
    check.exec_diagnostic("spam")
    assert check.failures and check.failures[0].name == "OcDiagFailed"


def test_names_executed(task_vars):
    task_vars["openshift_check_diagnostics"] = diagnostics = "ConfigContexts,spam,,eggs"

    def exec_module(module, args, *_):
        assert "extra_args" in args
        assert args["extra_args"][0] in diagnostics
        return {"result": "success"}

    DiagnosticCheck(exec_module, task_vars).run()
