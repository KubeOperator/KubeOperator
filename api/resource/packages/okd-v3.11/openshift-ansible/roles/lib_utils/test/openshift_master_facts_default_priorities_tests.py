import pytest


DEFAULT_PRIORITIES_3_6 = [
    {'name': 'SelectorSpreadPriority', 'weight': 1},
    {'name': 'InterPodAffinityPriority', 'weight': 1},
    {'name': 'LeastRequestedPriority', 'weight': 1},
    {'name': 'BalancedResourceAllocation', 'weight': 1},
    {'name': 'NodePreferAvoidPodsPriority', 'weight': 10000},
    {'name': 'NodeAffinityPriority', 'weight': 1},
    {'name': 'TaintTolerationPriority', 'weight': 1}
]
DEFAULT_PRIORITIES_3_8 = DEFAULT_PRIORITIES_3_7 = DEFAULT_PRIORITIES_3_6
DEFAULT_PRIORITIES_3_11 = DEFAULT_PRIORITIES_3_10 = DEFAULT_PRIORITIES_3_9 = DEFAULT_PRIORITIES_3_8

ZONE_PRIORITY = {
    'name': 'Zone',
    'argument': {
        'serviceAntiAffinity': {
            'label': 'zone'
        }
    },
    'weight': 2
}

TEST_VARS = [
    ('3.6', DEFAULT_PRIORITIES_3_6),
    ('3.7', DEFAULT_PRIORITIES_3_7),
    ('3.8', DEFAULT_PRIORITIES_3_8),
    ('3.9', DEFAULT_PRIORITIES_3_9),
    ('3.10', DEFAULT_PRIORITIES_3_10),
    ('3.11', DEFAULT_PRIORITIES_3_11),
]


def assert_ok(priorities_lookup, default_priorities, zones_enabled, **kwargs):
    results = priorities_lookup.run(None, zones_enabled=zones_enabled, **kwargs)
    if zones_enabled:
        assert results == default_priorities + [ZONE_PRIORITY]
    else:
        assert results == default_priorities


def test_openshift_version(priorities_lookup, openshift_version_fixture, zones_enabled):
    facts, default_priorities = openshift_version_fixture
    assert_ok(priorities_lookup, default_priorities, variables=facts, zones_enabled=zones_enabled)


@pytest.fixture(params=TEST_VARS)
def openshift_version_fixture(request, facts):
    version, default_priorities = request.param
    version += '.1'
    facts['openshift_version'] = version
    return facts, default_priorities


def test_openshift_release(priorities_lookup, openshift_release_fixture, zones_enabled):
    facts, default_priorities = openshift_release_fixture
    assert_ok(priorities_lookup, default_priorities, variables=facts, zones_enabled=zones_enabled)


@pytest.fixture(params=TEST_VARS)
def openshift_release_fixture(request, facts, release_mod):
    release, default_priorities = request.param
    facts['openshift_release'] = release_mod(release)
    return facts, default_priorities


def test_short_version_kwarg(priorities_lookup, short_version_kwarg_fixture, zones_enabled):
    facts, short_version, default_priorities = short_version_kwarg_fixture
    assert_ok(
        priorities_lookup, default_priorities, variables=facts,
        zones_enabled=zones_enabled, short_version=short_version)


@pytest.fixture(params=TEST_VARS)
def short_version_kwarg_fixture(request, facts):
    short_version, default_priorities = request.param
    return facts, short_version, default_priorities
