"""
Ansible action plugin to execute health checks in OpenShift clusters.
"""
import sys
import os
import base64
import traceback
import errno
import json
from collections import defaultdict

from ansible.plugins.action import ActionBase
from ansible.module_utils.six import string_types

try:
    from __main__ import display
except ImportError:
    # pylint: disable=ungrouped-imports; this is the standard way how to import
    # the default display object in Ansible action plugins.
    from ansible.utils.display import Display
    display = Display()

# Augment sys.path so that we can import checks from a directory relative to
# this callback plugin.
sys.path.insert(1, os.path.dirname(os.path.dirname(__file__)))

# pylint: disable=wrong-import-position; the import statement must come after
# the manipulation of sys.path.
from openshift_checks import OpenShiftCheck, OpenShiftCheckException, load_checks  # noqa: E402


class ActionModule(ActionBase):
    """Action plugin to execute health checks."""

    def run(self, tmp=None, task_vars=None):
        result = super(ActionModule, self).run(tmp, task_vars)
        task_vars = task_vars or {}

        # callback plugins cannot read Ansible vars, but we would like
        # zz_failure_summary to have access to certain values. We do so by
        # storing the information we need in the result.
        result['playbook_context'] = task_vars.get('r_openshift_health_checker_playbook_context')

        # if the user wants to write check results to files, they provide this directory:
        output_dir = task_vars.get("openshift_checks_output_dir")
        if output_dir:
            output_dir = os.path.join(output_dir, task_vars["ansible_host"])

        try:
            known_checks = self.load_known_checks(tmp, task_vars, output_dir)
            args = self._task.args
            requested_checks = normalize(args.get('checks', []))

            if not requested_checks:
                result['failed'] = True
                result['msg'] = list_known_checks(known_checks)
                return result

            resolved_checks = resolve_checks(requested_checks, known_checks.values())
        except OpenShiftCheckException as exc:
            result["failed"] = True
            result["msg"] = str(exc)
            return result

        if "openshift" not in task_vars:
            result["failed"] = True
            result["msg"] = "'openshift' is undefined, did 'openshift_facts' run?"
            return result

        result["checks"] = check_results = {}

        user_disabled_checks = normalize(task_vars.get('openshift_disable_check', []))
        # Automatically add docker_storage if only CRIO is used, as docker service would be stopped
        if task_vars.get('openshift_use_crio_only'):
            user_disabled_checks.append('docker_storage')

        for name in resolved_checks:
            display.banner("CHECK [{} : {}]".format(name, task_vars["ansible_host"]))
            check_results[name] = run_check(name, known_checks[name], user_disabled_checks, output_dir)

        result["changed"] = any(r.get("changed") for r in check_results.values())
        if any(r.get("failed") for r in check_results.values()):
            result["failed"] = True
            result["msg"] = "One or more checks failed"
        write_result_to_output_dir(output_dir, result)

        return result

    def load_known_checks(self, tmp, task_vars, output_dir=None):
        """Find all existing checks and return a mapping of names to instances."""
        load_checks()
        want_full_results = bool(output_dir)

        known_checks = {}
        for cls in OpenShiftCheck.subclasses():
            name = cls.name
            if name in known_checks:
                other_cls = known_checks[name].__class__
                raise OpenShiftCheckException(
                    "duplicate check name '{}' in: '{}' and '{}'"
                    "".format(name, full_class_name(cls), full_class_name(other_cls))
                )
            known_checks[name] = cls(
                execute_module=self._execute_module,
                tmp=tmp,
                task_vars=task_vars,
                want_full_results=want_full_results,
                templar=self._templar
            )
        return known_checks


def list_known_checks(known_checks):
    """Return text listing the existing checks and tags."""
    # TODO: we could include a description of each check by taking it from a
    # check class attribute (e.g., __doc__) when building the message below.
    msg = (
        'This playbook is meant to run health checks, but no checks were '
        'requested. Set the `openshift_checks` variable to a comma-separated '
        'list of check names or a YAML list. Available checks:\n  {}'
    ).format('\n  '.join(sorted(known_checks)))

    tags = describe_tags(known_checks.values())

    msg += (
        '\n\nTags can be used as a shortcut to select multiple '
        'checks. Available tags and the checks they select:\n  {}'
    ).format('\n  '.join(tags))

    return msg


def describe_tags(check_classes):
    """Return a sorted list of strings describing tags and the checks they include."""
    tag_checks = defaultdict(list)
    for cls in check_classes:
        for tag in cls.tags:
            tag_checks[tag].append(cls.name)
    tags = [
        '@{} = {}'.format(tag, ','.join(sorted(checks)))
        for tag, checks in tag_checks.items()
    ]
    return sorted(tags)


def resolve_checks(names, all_checks):
    """Returns a set of resolved check names.

    Resolving a check name expands tag references (e.g., "@tag") to all the
    checks that contain the given tag. OpenShiftCheckException is raised if
    names contains an unknown check or tag name.

    names should be a sequence of strings.

    all_checks should be a sequence of check classes/instances.
    """
    known_check_names = set(check.name for check in all_checks)
    known_tag_names = set(name for check in all_checks for name in check.tags)

    check_names = set(name for name in names if not name.startswith('@'))
    tag_names = set(name[1:] for name in names if name.startswith('@'))

    unknown_check_names = check_names - known_check_names
    unknown_tag_names = tag_names - known_tag_names

    if unknown_check_names or unknown_tag_names:
        msg = []
        if unknown_check_names:
            msg.append('Unknown check names: {}.'.format(', '.join(sorted(unknown_check_names))))
        if unknown_tag_names:
            msg.append('Unknown tag names: {}.'.format(', '.join(sorted(unknown_tag_names))))
        msg.append('Make sure there is no typo in the playbook and no files are missing.')
        # TODO: implement a "Did you mean ...?" when the input is similar to a
        # valid check or tag.
        msg.append('Known checks:')
        msg.append('  {}'.format('\n  '.join(sorted(known_check_names))))
        msg.append('Known tags:')
        msg.append('  {}'.format('\n  '.join(describe_tags(all_checks))))
        raise OpenShiftCheckException('\n'.join(msg))

    tag_to_checks = defaultdict(set)
    for check in all_checks:
        for tag in check.tags:
            tag_to_checks[tag].add(check.name)

    resolved = check_names.copy()
    for tag in tag_names:
        resolved.update(tag_to_checks[tag])

    return resolved


def normalize(checks):
    """Return a clean list of check names.

    The input may be a comma-separated string or a sequence. Leading and
    trailing whitespace characters are removed. Empty items are discarded.
    """
    if isinstance(checks, string_types):
        checks = checks.split(',')
    return [name.strip() for name in checks if name.strip()]


def run_check(name, check, user_disabled_checks, output_dir=None):
    """Run a single check if enabled and return a result dict."""

    # determine if we're going to run the check (not inactive or disabled)
    if name in user_disabled_checks or '*' in user_disabled_checks:
        return dict(skipped=True, skipped_reason="Disabled by user request")

    # pylint: disable=broad-except; capturing exceptions broadly is intentional,
    # to isolate arbitrary failures in one check from others.
    try:
        is_active = check.is_active()
    except Exception as exc:
        reason = "Could not determine if check should be run, exception: {}".format(exc)
        return dict(skipped=True, skipped_reason=reason, exception=traceback.format_exc())

    if not is_active:
        return dict(skipped=True, skipped_reason="Not active for this host")

    # run the check
    result = {}
    try:
        result = check.run()
    except OpenShiftCheckException as exc:
        check.register_failure(exc)
    except Exception as exc:
        check.register_failure("\n".join([str(exc), traceback.format_exc()]))

    # process the check state; compose the result hash, write files as needed
    if check.changed:
        result["changed"] = True
    if check.failures or result.get("failed"):
        if "msg" in result:  # failure result has msg; combine with any registered failures
            check.register_failure(result.get("msg"))
        result["failures"] = [(fail.name, str(fail)) for fail in check.failures]
        result["failed"] = True
        result["msg"] = "\n".join(str(fail) for fail in check.failures)
        write_to_output_file(output_dir, name + ".failures.json", result["failures"])
    if check.logs:
        write_to_output_file(output_dir, name + ".log.json", check.logs)
    if check.files_to_save:
        write_files_to_save(output_dir, check)

    return result


def prepare_output_dir(dirname):
    """Create the directory, including parents. Return bool for success/failure."""
    try:
        os.makedirs(dirname)
        return True
    except OSError as exc:
        # trying to create existing dir leads to error;
        # that error is fine, but for any other, assume the dir is not there
        return exc.errno == errno.EEXIST


def copy_remote_file_to_dir(check, file_to_save, output_dir, fname):
    """Copy file from remote host to local file in output_dir, if given."""
    if not output_dir or not prepare_output_dir(output_dir):
        return
    local_file = os.path.join(output_dir, fname)

    # pylint: disable=broad-except; do not need to do anything about failure to write dir/file
    # and do not want exceptions to break anything.
    try:
        # NOTE: it would have been nice to copy the file directly without loading it into
        # memory, but there does not seem to be a good way to do this via ansible.
        result = check.execute_module("slurp", dict(src=file_to_save), register=False)
        if result.get("failed"):
            display.warning("Could not retrieve file {}: {}".format(file_to_save, result.get("msg")))
            return

        content = result["content"]
        if result.get("encoding") == "base64":
            content = base64.b64decode(content)
        with open(local_file, "wb") as outfile:
            outfile.write(content)
    except Exception as exc:
        display.warning("Failed writing remote {} to local {}: {}".format(file_to_save, local_file, exc))
        return


def _no_fail(obj):
    # pylint: disable=broad-except; do not want serialization to fail for any reason
    try:
        return str(obj)
    except Exception:
        return "[not serializable]"


def write_to_output_file(output_dir, filename, data):
    """If output_dir provided, write data to file. Serialize as JSON if data is not a string."""

    if not output_dir or not prepare_output_dir(output_dir):
        return
    filename = os.path.join(output_dir, filename)
    try:
        with open(filename, 'w') as outfile:
            if isinstance(data, string_types):
                outfile.write(data)
            else:
                json.dump(data, outfile, sort_keys=True, indent=4, default=_no_fail)
    # pylint: disable=broad-except; do not want serialization/write to break for any reason
    except Exception as exc:
        display.warning("Could not write output file {}: {}".format(filename, exc))


def write_result_to_output_dir(output_dir, result):
    """If output_dir provided, write the result as json to result.json.

    Success/failure of the write is recorded as "output_files" in the result hash afterward.
    Otherwise this is much like write_to_output_file.
    """

    if not output_dir:
        return
    if not prepare_output_dir(output_dir):
        result["output_files"] = "Error creating output directory " + output_dir
        return

    filename = os.path.join(output_dir, "result.json")
    try:
        with open(filename, 'w') as outfile:
            json.dump(result, outfile, sort_keys=True, indent=4, default=_no_fail)
        result["output_files"] = "Check results for this host written to " + filename
    # pylint: disable=broad-except; do not want serialization/write to break for any reason
    except Exception as exc:
        result["output_files"] = "Error writing check results to {}:\n{}".format(filename, exc)


def write_files_to_save(output_dir, check):
    """Write files to check subdir in output dir."""
    if not output_dir:
        return
    output_dir = os.path.join(output_dir, check.name)
    seen_file = defaultdict(lambda: 0)
    for file_to_save in check.files_to_save:
        fname = file_to_save.filename
        while seen_file[fname]:  # just to be sure we never re-write a file, append numbers as needed
            seen_file[fname] += 1
            fname = "{}.{}".format(fname, seen_file[fname])
        seen_file[fname] += 1
        if file_to_save.remote_filename:
            copy_remote_file_to_dir(check, file_to_save.remote_filename, output_dir, fname)
        else:
            write_to_output_file(output_dir, fname, file_to_save.contents)


def full_class_name(cls):
    """Return the name of a class prefixed with its module name."""
    return '{}.{}'.format(cls.__module__, cls.__name__)
