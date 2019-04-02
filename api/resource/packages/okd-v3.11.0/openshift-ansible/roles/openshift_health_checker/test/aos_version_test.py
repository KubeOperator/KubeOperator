import pytest
import aos_version

from collections import namedtuple
Package = namedtuple('Package', ['name', 'version'])

expected_pkgs = {
    "spam": {
        "name": "spam",
        "version": "3.2.1",
        "check_multi": False,
    },
    "eggs": {
        "name": "eggs",
        "version": "3.2.1",
        "check_multi": False,
    },
}


@pytest.mark.parametrize('pkgs,expected_pkgs_dict', [
    (
        # all found
        [Package('spam', '3.2.1'), Package('eggs', '3.2.1')],
        expected_pkgs,
    ),
    (
        # found with more specific version
        [Package('spam', '3.2.1'), Package('eggs', '3.2.1.5')],
        expected_pkgs,
    ),
    (
        [Package('ovs', '2.6'), Package('ovs', '2.4')],
        {
            "ovs": {
                "name": "ovs",
                "version": ["2.6", "2.7"],
                "check_multi": False,
            }
        },
    ),
    (
        [Package('ovs', '2.7')],
        {
            "ovs": {
                "name": "ovs",
                "version": ["2.6", "2.7"],
                "check_multi": False,
            }
        },
    ),
])
def test_check_precise_version_found(pkgs, expected_pkgs_dict):
    aos_version._check_precise_version_found(pkgs, expected_pkgs_dict)


@pytest.mark.parametrize('pkgs,expect_not_found', [
    (
        [],
        {
            "spam": {
                "name": "spam",
                "version": "3.2.1",
                "check_multi": False,
            },
            "eggs": {
                "name": "eggs",
                "version": "3.2.1",
                "check_multi": False,
            }
        },  # none found
    ),
    (
        [Package('spam', '3.2.1')],
        {
            "eggs": {
                "name": "eggs",
                "version": "3.2.1",
                "check_multi": False,
            }
        },  # completely missing
    ),
    (
        [Package('spam', '3.2.1'), Package('eggs', '3.3.2')],
        {
            "eggs": {
                "name": "eggs",
                "version": "3.2.1",
                "check_multi": False,
            }
        },  # not the right version
    ),
    (
        [Package('eggs', '1.2.3'), Package('eggs', '3.2.1.5')],
        {
            "spam": {
                "name": "spam",
                "version": "3.2.1",
                "check_multi": False,
            }
        },  # eggs found with multiple versions
    ),
])
def test_check_precise_version_found_fail(pkgs, expect_not_found):
    with pytest.raises(aos_version.PreciseVersionNotFound) as e:
        aos_version._check_precise_version_found(pkgs, expected_pkgs)
    assert list(expect_not_found.values()) == e.value.problem_pkgs


@pytest.mark.parametrize('pkgs,expected_pkgs_dict', [
    (
        [],
        expected_pkgs,
    ),
    (
        # more precise but not strictly higher
        [Package('spam', '3.2.1.9')],
        expected_pkgs,
    ),
    (
        [Package('ovs', '2.7')],
        {
            "ovs": {
                "name": "ovs",
                "version": ["2.6", "2.7"],
                "check_multi": False,
            }
        },
    ),
])
def test_check_higher_version_found(pkgs, expected_pkgs_dict):
    aos_version._check_higher_version_found(pkgs, expected_pkgs_dict)


@pytest.mark.parametrize('pkgs,expected_pkgs_dict,expect_higher', [
    (
        [Package('spam', '3.3')],
        expected_pkgs,
        ['spam-3.3'],  # lower precision, but higher
    ),
    (
        [Package('spam', '3.2.1'), Package('eggs', '3.3.2')],
        expected_pkgs,
        ['eggs-3.3.2'],  # one too high
    ),
    (
        [Package('eggs', '1.2.3'), Package('eggs', '3.2.1.5'), Package('eggs', '3.4')],
        expected_pkgs,
        ['eggs-3.4'],  # multiple versions, one is higher
    ),
    (
        [Package('eggs', '3.2.1'), Package('eggs', '3.4'), Package('eggs', '3.3')],
        expected_pkgs,
        ['eggs-3.4'],  # multiple versions, two are higher
    ),
    (
        [Package('ovs', '2.8')],
        {
            "ovs": {
                "name": "ovs",
                "version": ["2.6", "2.7"],
                "check_multi": False,
            }
        },
        ['ovs-2.8'],
    ),
])
def test_check_higher_version_found_fail(pkgs, expected_pkgs_dict, expect_higher):
    with pytest.raises(aos_version.FoundHigherVersion) as e:
        aos_version._check_higher_version_found(pkgs, expected_pkgs_dict)
    assert set(expect_higher) == set(e.value.problem_pkgs)


@pytest.mark.parametrize('pkgs', [
    [],
    [Package('spam', '3.2.1')],
    [Package('spam', '3.2.1'), Package('eggs', '3.2.2')],
])
def test_check_multi_minor_release(pkgs):
    aos_version._check_multi_minor_release(pkgs, expected_pkgs)


@pytest.mark.parametrize('pkgs,expect_to_flag_pkgs', [
    (
        [Package('spam', '3.2.1'), Package('spam', '3.3.2')],
        ['spam'],
    ),
    (
        [Package('eggs', '1.2.3'), Package('eggs', '3.2.1.5'), Package('eggs', '3.4')],
        ['eggs'],
    ),
])
def test_check_multi_minor_release_fail(pkgs, expect_to_flag_pkgs):
    with pytest.raises(aos_version.FoundMultiRelease) as e:
        aos_version._check_multi_minor_release(pkgs, expected_pkgs)
    assert set(expect_to_flag_pkgs) == set(e.value.problem_pkgs)
