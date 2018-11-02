import pytest


# Predicates ordered according to OpenShift Origin source:
# origin/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithmprovider/defaults/defaults.go

DEFAULT_PREDICATES_3_6 = [
    {'name': 'NoVolumeZoneConflict'},
    {'name': 'MaxEBSVolumeCount'},
    {'name': 'MaxGCEPDVolumeCount'},
    {'name': 'MatchInterPodAffinity'},
    {'name': 'NoDiskConflict'},
    {'name': 'GeneralPredicates'},
    {'name': 'PodToleratesNodeTaints'},
    {'name': 'CheckNodeMemoryPressure'},
    {'name': 'CheckNodeDiskPressure'},
]

DEFAULT_PREDICATES_3_7 = [
    {'name': 'NoVolumeZoneConflict'},
    {'name': 'MaxEBSVolumeCount'},
    {'name': 'MaxGCEPDVolumeCount'},
    {'name': 'MaxAzureDiskVolumeCount'},
    {'name': 'MatchInterPodAffinity'},
    {'name': 'NoDiskConflict'},
    {'name': 'GeneralPredicates'},
    {'name': 'PodToleratesNodeTaints'},
    {'name': 'CheckNodeMemoryPressure'},
    {'name': 'CheckNodeDiskPressure'},
    {'name': 'NoVolumeNodeConflict'},
]

DEFAULT_PREDICATES_3_8 = DEFAULT_PREDICATES_3_7

DEFAULT_PREDICATES_3_9 = [
    {'name': 'NoVolumeZoneConflict'},
    {'name': 'MaxEBSVolumeCount'},
    {'name': 'MaxGCEPDVolumeCount'},
    {'name': 'MaxAzureDiskVolumeCount'},
    {'name': 'MatchInterPodAffinity'},
    {'name': 'NoDiskConflict'},
    {'name': 'GeneralPredicates'},
    {'name': 'PodToleratesNodeTaints'},
    {'name': 'CheckNodeMemoryPressure'},
    {'name': 'CheckNodeDiskPressure'},
    {'name': 'CheckVolumeBinding'},
]

DEFAULT_PREDICATES_3_10 = DEFAULT_PREDICATES_3_9

REGION_PREDICATE = {
    'name': 'Region',
    'argument': {
        'serviceAffinity': {
            'labels': ['region']
        }
    }
}

TEST_VARS = [
    ('3.6', DEFAULT_PREDICATES_3_6),
    ('3.7', DEFAULT_PREDICATES_3_7),
    ('3.8', DEFAULT_PREDICATES_3_8),
    ('3.9', DEFAULT_PREDICATES_3_9),
    ('3.10', DEFAULT_PREDICATES_3_10),
]


def assert_ok(predicates_lookup, default_predicates, regions_enabled, **kwargs):
    results = predicates_lookup.run(None, regions_enabled=regions_enabled, **kwargs)
    if regions_enabled:
        assert results == default_predicates + [REGION_PREDICATE]
    else:
        assert results == default_predicates


def test_openshift_version(predicates_lookup, openshift_version_fixture, regions_enabled):
    facts, default_predicates = openshift_version_fixture
    assert_ok(predicates_lookup, default_predicates, variables=facts, regions_enabled=regions_enabled)


@pytest.fixture(params=TEST_VARS)
def openshift_version_fixture(request, facts):
    version, default_predicates = request.param
    version += '.1'
    facts['openshift_version'] = version
    return facts, default_predicates


def test_openshift_release(predicates_lookup, openshift_release_fixture, regions_enabled):
    facts, default_predicates = openshift_release_fixture
    assert_ok(predicates_lookup, default_predicates, variables=facts, regions_enabled=regions_enabled)


@pytest.fixture(params=TEST_VARS)
def openshift_release_fixture(request, facts, release_mod):
    release, default_predicates = request.param
    facts['openshift_release'] = release_mod(release)
    return facts, default_predicates


def test_short_version_kwarg(predicates_lookup, short_version_kwarg_fixture, regions_enabled):
    facts, short_version, default_predicates = short_version_kwarg_fixture
    assert_ok(
        predicates_lookup, default_predicates, variables=facts,
        regions_enabled=regions_enabled, short_version=short_version)


@pytest.fixture(params=TEST_VARS)
def short_version_kwarg_fixture(request, facts):
    short_version, default_predicates = request.param
    return facts, short_version, default_predicates
