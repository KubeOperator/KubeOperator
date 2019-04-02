import pytest
import search_journalctl


def canned_search_journalctl(get_log_output=None):
    """Create a search_journalctl object with canned get_log_output method"""
    module = search_journalctl
    if get_log_output:
        module.get_log_output = get_log_output
    return module


DEFAULT_TIMESTAMP = 1496341364


def get_timestamp(modifier=0):
    return DEFAULT_TIMESTAMP + modifier


def get_timestamp_microseconds(modifier=0):
    return get_timestamp(modifier) * 1000000


def create_test_log_object(stamp, msg):
    return '{{"__REALTIME_TIMESTAMP": "{}", "MESSAGE": "{}"}}'.format(stamp, msg)


@pytest.mark.parametrize('name,matchers,log_input,expected_matches,expected_errors', [
    (
        'test with valid params',
        [
            {
                "start_regexp": r"Sample Logs Beginning",
                "regexp": r"test log message",
                "unit": "test",
            },
        ],
        [
            create_test_log_object(get_timestamp_microseconds(), "test log message"),
            create_test_log_object(get_timestamp_microseconds(), "Sample Logs Beginning"),
        ],
        ["test log message"],
        [],
    ),
    (
        'test with invalid json in log input',
        [
            {
                "start_regexp": r"Sample Logs Beginning",
                "regexp": r"test log message",
                "unit": "test-unit",
            },
        ],
        [
            '{__REALTIME_TIMESTAMP: ' + str(get_timestamp_microseconds()) + ', "MESSAGE": "test log message"}',
        ],
        [],
        [
            ["invalid json", "test-unit", "test log message"],
        ],
    ),
    (
        'test with invalid regexp',
        [
            {
                "start_regexp": r"Sample Logs Beginning",
                "regexp": r"test [ log message",
                "unit": "test",
            },
        ],
        [
            create_test_log_object(get_timestamp_microseconds(), "test log message"),
            create_test_log_object(get_timestamp_microseconds(), "sample log message"),
            create_test_log_object(get_timestamp_microseconds(), "fake log message"),
            create_test_log_object(get_timestamp_microseconds(), "dummy log message"),
            create_test_log_object(get_timestamp_microseconds(), "Sample Logs Beginning"),
        ],
        [],
        [
            ["invalid regular expression"],
        ],
    ),
], ids=lambda argval: argval[0])
def test_get_log_matches(name, matchers, log_input, expected_matches, expected_errors):
    def get_log_output(matcher):
        return log_input

    module = canned_search_journalctl(get_log_output)
    matched_regexp, errors = module.get_log_matches(matchers, 500, 60 * 60)

    assert set(matched_regexp) == set(expected_matches)
    assert len(expected_errors) == len(errors)

    for idx, partial_err_set in enumerate(expected_errors):
        for partial_err_msg in partial_err_set:
            assert partial_err_msg in errors[idx]


@pytest.mark.parametrize('name,matcher,log_count_lim,stamp_lim_seconds,log_input,expected_match', [
    (
        'test with matching log message, but out of bounds of log_count_lim',
        {
            "start_regexp": r"Sample Logs Beginning",
            "regexp": r"dummy log message",
            "unit": "test",
        },
        3,
        get_timestamp(-100 * 60 * 60),
        [
            create_test_log_object(get_timestamp_microseconds(), "test log message"),
            create_test_log_object(get_timestamp_microseconds(), "sample log message"),
            create_test_log_object(get_timestamp_microseconds(), "fake log message"),
            create_test_log_object(get_timestamp_microseconds(), "dummy log message"),
            create_test_log_object(get_timestamp_microseconds(), "Sample Logs Beginning"),
        ],
        None,
    ),
    (
        'test with matching log message, but with timestamp too old',
        {
            "start_regexp": r"Sample Logs Beginning",
            "regexp": r"dummy log message",
            "unit": "test",
        },
        100,
        get_timestamp(-10),
        [
            create_test_log_object(get_timestamp_microseconds(), "test log message"),
            create_test_log_object(get_timestamp_microseconds(), "sample log message"),
            create_test_log_object(get_timestamp_microseconds(), "fake log message"),
            create_test_log_object(get_timestamp_microseconds(-1000), "dummy log message"),
            create_test_log_object(get_timestamp_microseconds(-1000), "Sample Logs Beginning"),
        ],
        None,
    ),
    (
        'test with matching log message, and timestamp within time limit',
        {
            "start_regexp": r"Sample Logs Beginning",
            "regexp": r"dummy log message",
            "unit": "test",
        },
        100,
        get_timestamp(-1010),
        [
            create_test_log_object(get_timestamp_microseconds(), "test log message"),
            create_test_log_object(get_timestamp_microseconds(), "sample log message"),
            create_test_log_object(get_timestamp_microseconds(), "fake log message"),
            create_test_log_object(get_timestamp_microseconds(-1000), "dummy log message"),
            create_test_log_object(get_timestamp_microseconds(-1000), "Sample Logs Beginning"),
        ],
        create_test_log_object(get_timestamp_microseconds(-1000), "dummy log message"),
    ),
], ids=lambda argval: argval[0])
def test_find_matches_skips_logs(name, matcher, log_count_lim, stamp_lim_seconds, log_input, expected_match):
    match = search_journalctl.find_matches(log_input, matcher, log_count_lim, stamp_lim_seconds)
    assert match == expected_match
