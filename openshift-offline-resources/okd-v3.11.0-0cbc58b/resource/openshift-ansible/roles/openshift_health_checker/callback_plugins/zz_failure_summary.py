"""Ansible callback plugin to print a nicely formatted summary of failures.

The file / module name is prefixed with `zz_` to make this plugin be loaded last
by Ansible, thus making its output the last thing that users see.
"""

from collections import defaultdict
import traceback

from ansible.plugins.callback import CallbackBase
from ansible import constants as C
from ansible.utils.color import stringc
from ansible.module_utils.six import string_types, PY2


FAILED_NO_MSG = u'Failed without returning a message.'


class CallbackModule(CallbackBase):
    """This callback plugin stores task results and summarizes failures."""

    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'aggregate'
    CALLBACK_NAME = 'failure_summary'
    CALLBACK_NEEDS_WHITELIST = False

    def __init__(self):
        super(CallbackModule, self).__init__()
        self.__failures = []
        self.__playbook_file = ''

    def v2_playbook_on_start(self, playbook):
        super(CallbackModule, self).v2_playbook_on_start(playbook)
        # pylint: disable=protected-access; Ansible gives us no public API to
        # get the file name of the current playbook from a callback plugin.
        self.__playbook_file = playbook._file_name

    def v2_runner_on_failed(self, result, ignore_errors=False):
        super(CallbackModule, self).v2_runner_on_failed(result, ignore_errors)
        if not ignore_errors:
            self.__failures.append(result)

    def v2_playbook_on_stats(self, stats):
        super(CallbackModule, self).v2_playbook_on_stats(stats)
        # pylint: disable=broad-except; capturing exceptions broadly is
        # intentional, to isolate arbitrary failures in this callback plugin.
        try:
            if self.__failures:
                self._display.display(failure_summary(self.__failures, self.__playbook_file))
        except Exception:
            msg = stringc(
                u'An error happened while generating a summary of failures:\n'
                u'{}'.format(traceback.format_exc()), C.COLOR_WARN)
            self._display.v(msg)


def failure_summary(failures, playbook):
    """Return a summary of failed tasks, including details on health checks."""
    if not failures:
        return u''

    # NOTE: because we don't have access to task_vars from callback plugins, we
    # store the playbook context in the task result when the
    # openshift_health_check action plugin is used, and we use this context to
    # customize the error message.
    # pylint: disable=protected-access; Ansible gives us no sufficient public
    # API on TaskResult objects.
    context = next((
        context for context in
        (failure._result.get('playbook_context') for failure in failures)
        if context
    ), None)

    failures = [failure_to_dict(failure) for failure in failures]
    failures = deduplicate_failures(failures)

    summary = [u'', u'', u'Failure summary:', u'']

    width = len(str(len(failures)))
    initial_indent_format = u'  {{:>{width}}}. '.format(width=width)
    initial_indent_len = len(initial_indent_format.format(0))
    subsequent_indent = u' ' * initial_indent_len
    subsequent_extra_indent = u' ' * (initial_indent_len + 10)

    for i, failure in enumerate(failures, 1):
        entries = format_failure(failure)
        summary.append(u'\n{}{}'.format(initial_indent_format.format(i), entries[0]))
        for entry in entries[1:]:
            if PY2:
                entry = entry.decode('utf8')
            entry = entry.replace(u'\n', u'\n' + subsequent_extra_indent)
            indented = u'{}{}'.format(subsequent_indent, entry)
            summary.append(indented)

    failed_checks = set()
    for failure in failures:
        failed_checks.update(name for name, message in failure['checks'])
    if failed_checks:
        summary.append(check_failure_footer(failed_checks, context, playbook))

    return u'\n'.join(summary)


def failure_to_dict(failed_task_result):
    """Extract information out of a failed TaskResult into a dict.

    The intent is to transform a TaskResult object into something easier to
    manipulate. TaskResult is ansible.executor.task_result.TaskResult.
    """
    # pylint: disable=protected-access; Ansible gives us no sufficient public
    # API on TaskResult objects.
    _result = failed_task_result._result
    return {
        'host': failed_task_result._host.get_name(),
        'play': play_name(failed_task_result._task),
        'task': failed_task_result.task_name,
        'msg': _result.get('msg', FAILED_NO_MSG),
        'checks': tuple(
            (name, result.get('msg', FAILED_NO_MSG))
            for name, result in sorted(_result.get('checks', {}).items())
            if result.get('failed')
        ),
    }


def play_name(obj):
    """Given a task or block, return the name of its parent play.

    This is loosely inspired by ansible.playbook.base.Base.dump_me.
    """
    # pylint: disable=protected-access; Ansible gives us no sufficient public
    # API to implement this.
    if not obj:
        return ''
    if hasattr(obj, '_play'):
        return obj._play.get_name()
    return play_name(getattr(obj, '_parent'))


def deduplicate_failures(failures):
    """Group together similar failures from different hosts.

    Returns a new list of failures such that identical failures from different
    hosts are grouped together in a single entry. The relative order of failures
    is preserved.

    If failures is unhashable, the original list of failures is returned.
    """
    groups = defaultdict(list)
    for failure in failures:
        group_key = tuple(sorted((key, value) for key, value in failure.items() if key != 'host'))
        try:
            groups[group_key].append(failure)
        except TypeError:
            # abort and return original list of failures when failures has an
            # unhashable type.
            return failures

    result = []
    for failure in failures:
        group_key = tuple(sorted((key, value) for key, value in failure.items() if key != 'host'))
        if group_key not in groups:
            continue
        failure['host'] = tuple(sorted(g_failure['host'] for g_failure in groups.pop(group_key)))
        result.append(failure)
    return result


def format_failure(failure):
    """Return a list of pretty-formatted text entries describing a failure, including
    relevant information about it. Expect that the list of text entries will be joined
    by a newline separator when output to the user."""
    if isinstance(failure['host'], string_types):
        host = failure['host']
    else:
        host = u', '.join(failure['host'])
    play = failure['play']
    task = failure['task']
    msg = failure['msg']
    if not isinstance(msg, string_types):
        msg = str(msg)
    checks = failure['checks']
    fields = (
        (u'Hosts', host),
        (u'Play', play),
        (u'Task', task),
        (u'Message', stringc(msg, C.COLOR_ERROR)),
    )
    if checks:
        fields += ((u'Details', format_failed_checks(checks)),)
    row_format = '{:10}{}'
    return [row_format.format(header + u':', body.encode('utf8')) if PY2 else body for header, body in fields]


def format_failed_checks(checks):
    """Return pretty-formatted text describing checks that failed."""
    messages = []
    for name, message in checks:
        messages.append(u'check "{}":\n{}'.format(name, message))
    return stringc(u'\n\n'.join(messages), C.COLOR_ERROR)


def check_failure_footer(failed_checks, context, playbook):
    """Return a textual explanation about checks depending on context.

    The purpose of specifying context is to vary the output depending on what
    the user was expecting to happen (based on which playbook they ran). The
    only use currently is to vary the message depending on whether the user was
    deliberately running checks or was trying to install/upgrade and checks are
    just included. Other use cases may arise.
    """
    checks = ','.join(sorted(failed_checks))
    summary = [u'']
    if context in ['pre-install', 'health', 'adhoc']:
        # User was expecting to run checks, less explanation needed.
        summary.extend([
            u'You may configure or disable checks by setting Ansible '
            u'variables. To disable those above, set:',
            u'    openshift_disable_check={checks}'.format(checks=checks),
            u'Consult check documentation for configurable variables.',
        ])
    else:
        # User may not be familiar with the checks, explain what checks are in
        # the first place.
        summary.extend([
            u'The execution of "{playbook}" includes checks designed to fail '
            u'early if the requirements of the playbook are not met. One or '
            u'more of these checks failed. To disregard these results,'
            u'explicitly disable checks by setting an Ansible variable:'.format(playbook=playbook),
            u'   openshift_disable_check={checks}'.format(checks=checks),
            u'Failing check names are shown in the failure details above. '
            u'Some checks may be configurable by variables if your requirements '
            u'are different from the defaults; consult check documentation.',
        ])
    summary.append(
        u'Variables can be set in the inventory or passed on the command line '
        u'using the -e flag to ansible-playbook.'
    )
    return u'\n'.join(summary)
