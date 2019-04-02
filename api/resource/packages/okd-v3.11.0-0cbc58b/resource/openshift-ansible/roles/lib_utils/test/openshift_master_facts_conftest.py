import os
import sys

import pytest

sys.path.insert(1, os.path.join(os.path.dirname(__file__), os.pardir, "lookup_plugins"))

from openshift_master_facts_default_predicates import LookupModule as PredicatesLookupModule  # noqa: E402
from openshift_master_facts_default_priorities import LookupModule as PrioritiesLookupModule  # noqa: E402


@pytest.fixture()
def predicates_lookup():
    return PredicatesLookupModule()


@pytest.fixture()
def priorities_lookup():
    return PrioritiesLookupModule()


@pytest.fixture()
def facts():
    return {
        'openshift': {
            'common': {}
        }
    }


@pytest.fixture(params=[True, False])
def regions_enabled(request):
    return request.param


@pytest.fixture(params=[True, False])
def zones_enabled(request):
    return request.param


def v_prefix(release):
    """Prefix a release number with 'v'."""
    return "v" + release


def minor(release):
    """Add a suffix to release, making 'X.Y' become 'X.Y.Z'."""
    return release + ".1"


@pytest.fixture(params=[str, v_prefix, minor])
def release_mod(request):
    """Modifies a release string to alternative valid values."""
    return request.param
