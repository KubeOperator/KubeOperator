import pytest
import rpm_version

expected_pkgs = {
    "spam": {
        "name": "spam",
        "version": "3.2.1",
    },
    "eggs": {
        "name": "eggs",
        "version": "3.2.1",
    },
}


@pytest.mark.parametrize('pkgs, expect_not_found', [
    (
        {},
        ["spam", "eggs"],  # none found
    ),
    (
        {"spam": ["3.2.1", "4.5.1"]},
        ["eggs"],  # completely missing
    ),
    (
        {
            "spam": ["3.2.1", "4.5.1"],
            "eggs": ["3.2.1"],
        },
        [],  # all found
    ),
])
def test_check_pkg_found(pkgs, expect_not_found):
    if expect_not_found:
        with pytest.raises(rpm_version.RpmVersionException) as e:
            rpm_version._check_pkg_versions(pkgs, expected_pkgs)

        assert "not found to be installed" in str(e.value)
        assert set(expect_not_found) == set(e.value.problem_pkgs)
    else:
        rpm_version._check_pkg_versions(pkgs, expected_pkgs)


@pytest.mark.parametrize('pkgs, expect_not_found', [
    (
        {
            'spam': ['3.2.1'],
            'eggs': ['3.3.2'],
        },
        {
            "eggs": {
                "required_versions": ["3.2"],
                "found_versions": ["3.3"],
            }
        },  # not the right version
    ),
    (
        {
            'spam': ['3.1.2', "3.3.2"],
            'eggs': ['3.3.2', "1.2.3"],
        },
        {
            "eggs": {
                "required_versions": ["3.2"],
                "found_versions": ["3.3", "1.2"],
            },
            "spam": {
                "required_versions": ["3.2"],
                "found_versions": ["3.1", "3.3"],
            }
        },  # not the right version
    ),
])
def test_check_pkg_version_found(pkgs, expect_not_found):
    if expect_not_found:
        with pytest.raises(rpm_version.RpmVersionException) as e:
            rpm_version._check_pkg_versions(pkgs, expected_pkgs)

        assert "found to be installed with an incorrect version" in str(e.value)
        assert expect_not_found == e.value.problem_pkgs
    else:
        rpm_version._check_pkg_versions(pkgs, expected_pkgs)
