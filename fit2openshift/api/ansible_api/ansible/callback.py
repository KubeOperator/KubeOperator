# ~*~ coding: utf-8 ~*~
import datetime
from collections import defaultdict
from ansible import constants as C

from ansible.plugins.callback.default import CallbackModule
from ansible.plugins.callback.minimal import CallbackModule as CMDCallBackModule


class CallbackMixin:
    def __init__(self):
        # result_raw example: {
        #   "ok": {"hostname": {"task_name": {}，...},..},
        #   "failed": {"hostname": {"task_name": {}..}, ..},
        #   "unreachable: {"hostname": {"task_name": {}, ..}},
        #   "skipped": {"hostname": {"task_name": {}, ..}, ..},
        # }
        # results_summary example: {
        #   "contacted": {"hostname": {"task_name": {}}, "hostname": {}},
        #   "dark": {"hostname": {"task_name": {}, "task_name": {}},...,},
        #   "success": True
        # }
        self.results_raw = dict(
            ok=defaultdict(dict),
            failed=defaultdict(dict),
            unreachable=defaultdict(dict),
            skippe=defaultdict(dict),
        )
        self.results_summary = dict(
            contacted=defaultdict(dict),
            dark=defaultdict(dict),
            success=True
        )
        self.results = {
            'raw': self.results_raw,
            'summary': self.results_summary,
        }
        super().__init__()
        self._display.columns = 79

    def display(self, msg):
        self._display.display(msg)

    def gather_result(self, t, result):
        self._clean_results(result._result, result._task.action)
        host = result._host.get_name()
        task_name = result.task_name
        task_result = result._result

        self.results_raw[t][host][task_name] = task_result
        self.clean_result(t, host, task_name, task_result)


class AdHocResultCallback(CallbackMixin, CallbackModule, CMDCallBackModule):
    """
    Task result Callback
    """
    def clean_result(self, t, host, task_name, task_result):
        contacted = self.results_summary["contacted"]
        dark = self.results_summary["dark"]

        if task_result.get('rc') is not None:
            cmd = task_result.get('cmd')
            if isinstance(cmd, list):
                cmd = " ".join(cmd)
            else:
                cmd = str(cmd)
            detail = {
                'cmd': cmd,
                'stderr': task_result.get('stderr'),
                'stdout': task_result.get('stdout'),
                'rc': task_result.get('rc'),
                'delta': task_result.get('delta'),
                'msg': task_result.get('msg', '')
            }
        else:
            detail = {
                "changed": task_result.get('changed', False),
                "msg": task_result.get('msg', '')
            }

        if t in ("ok", "skipped"):
            contacted[host][task_name] = detail
        else:
            dark[host][task_name] = detail

    def v2_runner_on_failed(self, result, ignore_errors=False):
        self.results_summary['success'] = False
        self.gather_result("failed", result)

        if result._task.action in C.MODULE_NO_JSON:
            super(CMDCallBackModule, self).v2_runner_on_failed(
                result, ignore_errors=ignore_errors
            )
        else:
            super(CallbackModule, self).v2_runner_on_failed(
                result, ignore_errors=ignore_errors
            )

    def v2_runner_on_ok(self, result):
        self.gather_result("ok", result)
        if result._task.action in C.MODULE_NO_JSON:
            CMDCallBackModule.v2_runner_on_ok(self, result)
        else:
            super().v2_runner_on_ok(result)

    def v2_runner_on_skipped(self, result):
        self.gather_result("skipped", result)
        super().v2_runner_on_skipped(result)

    def v2_runner_on_unreachable(self, result):
        self.results_summary['success'] = False
        self.gather_result("unreachable", result)
        super().v2_runner_on_unreachable(result)

    def on_playbook_start(self, name):
        date_start = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        self.display(
            "{} Start task: {}\r\n".format(date_start, name)
        )

    def on_playbook_end(self, name):
        date_finished = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        self.display(
            "{} Task finish\r\n".format(date_finished)
        )

    def display_skipped_hosts(self):
        pass

    def display_ok_hosts(self):
        pass


class PlaybookResultCallBack(AdHocResultCallback):
    """
    Custom callback model for handlering the output data of
    execute playbook file,
    Base on the build-in callback plugins of ansible which named `json`.
    """
    pass



