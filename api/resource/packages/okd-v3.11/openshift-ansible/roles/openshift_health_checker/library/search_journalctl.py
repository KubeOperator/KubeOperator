#!/usr/bin/python
"""Interface to journalctl."""

from time import time
import json
import re
import subprocess

from ansible.module_utils.basic import AnsibleModule


class InvalidMatcherRegexp(Exception):
    """Exception class for invalid matcher regexp."""
    pass


class InvalidLogEntry(Exception):
    """Exception class for invalid / non-json log entries."""
    pass


class LogInputSubprocessError(Exception):
    """Exception class for errors that occur while executing a subprocess."""
    pass


def main():
    """Scan a given list of "log_matchers" for journalctl messages containing given patterns.
    "log_matchers" is a list of dicts consisting of three keys that help fine-tune log searching:
    'start_regexp', 'regexp', and 'unit'.

    Sample "log_matchers" list:

    [
      {
        'start_regexp': r'Beginning of systemd unit',
        'regexp': r'the specific log message to find',
        'unit': 'etcd',
      }
    ]
    """
    module = AnsibleModule(
        argument_spec=dict(
            log_count_limit=dict(type="int", default=500),
            log_matchers=dict(type="list", required=True),
        ),
    )

    timestamp_limit_seconds = time() - 60 * 60  # 1 hour

    log_count_limit = module.params["log_count_limit"]
    log_matchers = module.params["log_matchers"]

    matched_regexp, errors = get_log_matches(log_matchers, log_count_limit, timestamp_limit_seconds)

    module.exit_json(
        changed=False,
        failed=bool(errors),
        errors=errors,
        matched=matched_regexp,
    )


def get_log_matches(matchers, log_count_limit, timestamp_limit_seconds):
    """Return a list of up to log_count_limit matches for each matcher.

    Log entries are only considered if newer than timestamp_limit_seconds.
    """
    matched_regexp = []
    errors = []

    for matcher in matchers:
        try:
            log_output = get_log_output(matcher)
        except LogInputSubprocessError as err:
            errors.append(str(err))
            continue

        try:
            matched = find_matches(log_output, matcher, log_count_limit, timestamp_limit_seconds)
            if matched:
                matched_regexp.append(matcher.get("regexp", ""))
        except InvalidMatcherRegexp as err:
            errors.append(str(err))
        except InvalidLogEntry as err:
            errors.append(str(err))

    return matched_regexp, errors


def get_log_output(matcher):
    """Return an iterator on the logs of a given matcher."""
    try:
        cmd_output = subprocess.Popen(list([
            '/bin/journalctl',
            '-ru', matcher.get("unit", ""),
            '--output', 'json',
        ]), stdout=subprocess.PIPE)

        return iter(cmd_output.stdout.readline, '')

    except subprocess.CalledProcessError as exc:
        msg = "Could not obtain journalctl logs for the specified systemd unit: {}: {}"
        raise LogInputSubprocessError(msg.format(matcher.get("unit", "<missing>"), str(exc)))
    except OSError as exc:
        raise LogInputSubprocessError(str(exc))


def find_matches(log_output, matcher, log_count_limit, timestamp_limit_seconds):
    """Return log messages matched in iterable log_output by a given matcher.

    Ignore any log_output items older than timestamp_limit_seconds.
    """
    try:
        regexp = re.compile(matcher.get("regexp", ""))
        start_regexp = re.compile(matcher.get("start_regexp", ""))
    except re.error as err:
        msg = "A log matcher object was provided with an invalid regular expression: {}"
        raise InvalidMatcherRegexp(msg.format(str(err)))

    matched = None

    for log_count, line in enumerate(log_output):
        if log_count >= log_count_limit:
            break

        try:
            obj = json.loads(line)

            # don't need to look past the most recent service restart
            if start_regexp.match(obj["MESSAGE"]):
                break

            log_timestamp_seconds = float(obj["__REALTIME_TIMESTAMP"]) / 1000000
            if log_timestamp_seconds < timestamp_limit_seconds:
                break

            if regexp.match(obj["MESSAGE"]):
                matched = line
                break

        except ValueError:
            msg = "Log entry for systemd unit {} contained invalid json syntax: {}"
            raise InvalidLogEntry(msg.format(matcher.get("unit"), line))

    return matched


if __name__ == '__main__':
    main()
