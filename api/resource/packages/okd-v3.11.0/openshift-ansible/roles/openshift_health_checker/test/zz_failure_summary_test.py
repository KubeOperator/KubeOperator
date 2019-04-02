from zz_failure_summary import deduplicate_failures

import pytest


@pytest.mark.parametrize('failures,deduplicated', [
    (
        [
            {
                'host': 'master1',
                'msg': 'One or more checks failed',
            },
        ],
        [
            {
                'host': ('master1',),
                'msg': 'One or more checks failed',
            },
        ],
    ),
    (
        [
            {
                'host': 'master1',
                'msg': 'One or more checks failed',
            },
            {
                'host': 'node1',
                'msg': 'One or more checks failed',
            },
        ],
        [
            {
                'host': ('master1', 'node1'),
                'msg': 'One or more checks failed',
            },
        ],
    ),
    (
        [
            {
                'host': 'node1',
                'msg': 'One or more checks failed',
                'checks': (('test_check', 'error message'),),
            },
            {
                'host': 'master2',
                'msg': 'Some error happened',
            },
            {
                'host': 'master1',
                'msg': 'One or more checks failed',
                'checks': (('test_check', 'error message'),),
            },
        ],
        [
            {
                'host': ('master1', 'node1'),
                'msg': 'One or more checks failed',
                'checks': (('test_check', 'error message'),),
            },
            {
                'host': ('master2',),
                'msg': 'Some error happened',
            },
        ],
    ),
    # if a failure contain an unhashable value, it will not be deduplicated
    (
        [
            {
                'host': 'master1',
                'msg': {'unhashable': 'value'},
            },
        ],
        [
            {
                'host': 'master1',
                'msg': {'unhashable': 'value'},
            },
        ],
    ),
])
def test_deduplicate_failures(failures, deduplicated):
    assert deduplicate_failures(failures) == deduplicated
