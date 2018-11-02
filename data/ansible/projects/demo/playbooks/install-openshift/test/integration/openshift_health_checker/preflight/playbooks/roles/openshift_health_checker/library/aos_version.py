#!/usr/bin/python
"""
Ansible module for yum-based systems determining if multiple releases
of an OpenShift package are available, and if the release requested
(if any) is available down to the given precision.

For Enterprise, multiple releases available suggest that multiple repos
are enabled for the different releases, which may cause installation
problems. With Origin, however, this is a normal state of affairs as
all the releases are provided in a single repo with the expectation that
only the latest can be installed.

Code in the openshift_version role contains a lot of logic to pin down
the exact package and image version to use and so does some validation
of release availability already. Without duplicating all that, we would
like the user to have a helpful error message if we detect things will
not work out right. Note that if openshift_release is not specified in
the inventory, the version comparison checks just pass.
"""

from ansible.module_utils.basic import AnsibleModule
# NOTE: because of the dependency on yum (Python 2-only), this module does not
# work under Python 3. But since we run unit tests against both Python 2 and
# Python 3, we use six for cross compatibility in this module alone:
from ansible.module_utils.six import string_types

YUM_IMPORT_EXCEPTION = None
DNF_IMPORT_EXCEPTION = None
try:
    import yum  # pylint: disable=import-error
except ImportError as err:
    YUM_IMPORT_EXCEPTION = err

try:
    import dnf  # pylint: disable=import-error
except ImportError as err:
    DNF_IMPORT_EXCEPTION = err


class AosVersionException(Exception):
    """Base exception class for package version problems"""
    def __init__(self, message, problem_pkgs=None):
        Exception.__init__(self, message)
        self.problem_pkgs = problem_pkgs


def main():
    """Entrypoint for this Ansible module"""
    module = AnsibleModule(
        argument_spec=dict(
            package_list=dict(type="list", required=True),
            package_mgr=dict(type="str", required=True),
        ),
        supports_check_mode=True
    )

    # determine the package manager to use
    package_mgr = module.params['package_mgr']
    if package_mgr not in ('yum', 'dnf'):
        module.fail_json(msg="package_mgr must be one of: yum, dnf")
    pkg_mgr_exception = dict(yum=YUM_IMPORT_EXCEPTION, dnf=DNF_IMPORT_EXCEPTION)[package_mgr]
    if pkg_mgr_exception:
        module.fail_json(
            msg="aos_version module could not import {}: {}".format(package_mgr, pkg_mgr_exception)
        )

    # determine the packages we will look for
    package_list = module.params['package_list']
    if not package_list:
        module.fail_json(msg="package_list must not be empty")

    # generate set with only the names of expected packages
    expected_pkg_names = [p["name"] for p in package_list]

    # gather packages that require a multi_minor_release check
    multi_minor_pkgs = [p for p in package_list if p["check_multi"]]

    # generate list of packages with a specified (non-empty) version
    # should look like a version string with possibly many segments e.g. "3.4.1"
    versioned_pkgs = [p for p in package_list if p["version"]]

    # get the list of packages available and complain if anything is wrong
    try:
        pkgs = _retrieve_available_packages(package_mgr, expected_pkg_names)
        if versioned_pkgs:
            _check_precise_version_found(pkgs, _to_dict(versioned_pkgs))
            _check_higher_version_found(pkgs, _to_dict(versioned_pkgs))
        if multi_minor_pkgs:
            _check_multi_minor_release(pkgs, _to_dict(multi_minor_pkgs))
    except AosVersionException as excinfo:
        module.fail_json(msg=str(excinfo))
    module.exit_json(changed=False)


def _to_dict(pkg_list):
    return {pkg["name"]: pkg for pkg in pkg_list}


def _retrieve_available_packages(pkg_mgr, expected_pkgs):
    # The openshift excluder prevents unintended updates to openshift
    # packages by setting yum excludes on those packages. See:
    # https://wiki.centos.org/SpecialInterestGroup/PaaS/OpenShift-Origin-Control-Updates
    # Excludes are then disabled during an install or upgrade, but
    # this check will most likely be running outside either. When we
    # attempt to determine what packages are available via yum they may
    # be excluded. So, for our purposes here, disable excludes to see
    # what will really be available during an install or upgrade.

    if pkg_mgr == "yum":
        # search for package versions available for openshift pkgs
        yb = yum.YumBase()  # pylint: disable=invalid-name

        yb.conf.disable_excludes = ['all']

        try:
            pkgs = yb.rpmdb.returnPackages(patterns=expected_pkgs)
            pkgs += yb.pkgSack.returnPackages(patterns=expected_pkgs)
        except yum.Errors.PackageSackError as excinfo:
            # you only hit this if *none* of the packages are available
            raise AosVersionException('\n'.join([
                'Unable to find any OpenShift packages.',
                'Check your subscription and repo settings.',
                str(excinfo),
            ]))
    elif pkg_mgr == "dnf":
        dbase = dnf.Base()  # pyling: disable=invalid-name

        dbase.conf.disable_excludes = ['all']
        dbase.read_all_repos()
        dbase.fill_sack(load_system_repo=False, load_available_repos=True)

        dquery = dbase.sack.query()
        aquery = dquery.available()
        iquery = dquery.installed()

        available_pkgs = list(aquery.filter(name=expected_pkgs))
        installed_pkgs = list(iquery.filter(name=expected_pkgs))
        pkgs = available_pkgs + installed_pkgs

        if not pkgs:
            # pkgs list is empty, raise because no expected packages found
            raise AosVersionException('\n'.join([
                'Unable to find any OpenShift packages.',
                'Check your subscription and repo settings.',
            ]))

    return pkgs


class PreciseVersionNotFound(AosVersionException):
    """Exception for reporting packages not available at given version"""
    def __init__(self, not_found):
        msg = ['Not all of the required packages are available at their requested version']
        msg += ['{}:{} '.format(pkg["name"], pkg["version"]) for pkg in not_found]
        msg += ['Please check your subscriptions and enabled repositories.']
        AosVersionException.__init__(self, '\n'.join(msg), not_found)


def _check_precise_version_found(pkgs, expected_pkgs_dict):
    # see if any packages couldn't be found at requested release version
    # we would like to verify that the latest available pkgs have however specific a version is given.
    # so e.g. if there is a package version 3.4.1.5 the check passes; if only 3.4.0, it fails.

    pkgs_precise_version_found = set()
    for pkg in pkgs:
        if pkg.name not in expected_pkgs_dict:
            continue
        expected_pkg_versions = expected_pkgs_dict[pkg.name]["version"]
        if isinstance(expected_pkg_versions, string_types):
            expected_pkg_versions = [expected_pkg_versions]
        for expected_pkg_version in expected_pkg_versions:
            # does the version match, to the precision requested?
            # and, is it strictly greater, at the precision requested?
            match_version = '.'.join(pkg.version.split('.')[:expected_pkg_version.count('.') + 1])
            if match_version == expected_pkg_version:
                pkgs_precise_version_found.add(pkg.name)

    not_found = []
    for name, pkg in expected_pkgs_dict.items():
        if name not in pkgs_precise_version_found:
            not_found.append(pkg)

    if not_found:
        raise PreciseVersionNotFound(not_found)


class FoundHigherVersion(AosVersionException):
    """Exception for reporting that a higher version than requested is available"""
    def __init__(self, higher_found):
        msg = ['Some required package(s) are available at a version',
               'that is higher than requested']
        msg += ['  ' + name for name in higher_found]
        msg += ['This will prevent installing the version you requested.']
        msg += ['Please check your enabled repositories or adjust openshift_release.']
        AosVersionException.__init__(self, '\n'.join(msg), higher_found)


def _check_higher_version_found(pkgs, expected_pkgs_dict):
    expected_pkg_names = list(expected_pkgs_dict)

    # see if any packages are available in a version higher than requested
    higher_version_for_pkg = {}
    for pkg in pkgs:
        if pkg.name not in expected_pkg_names:
            continue
        expected_pkg_versions = expected_pkgs_dict[pkg.name]["version"]
        if isinstance(expected_pkg_versions, string_types):
            expected_pkg_versions = [expected_pkg_versions]
        # NOTE: the list of versions is assumed to be sorted so that the highest
        # desirable version is the last.
        highest_desirable_version = expected_pkg_versions[-1]
        req_release_arr = [int(segment) for segment in highest_desirable_version.split(".")]
        version = [int(segment) for segment in pkg.version.split(".")]
        too_high = version[:len(req_release_arr)] > req_release_arr
        higher_than_seen = version > higher_version_for_pkg.get(pkg.name, [])
        if too_high and higher_than_seen:
            higher_version_for_pkg[pkg.name] = version

    if higher_version_for_pkg:
        higher_found = []
        for name, version in higher_version_for_pkg.items():
            higher_found.append(name + '-' + '.'.join(str(segment) for segment in version))
        raise FoundHigherVersion(higher_found)


class FoundMultiRelease(AosVersionException):
    """Exception for reporting multiple minor releases found for same package"""
    def __init__(self, multi_found):
        msg = ['Multiple minor versions of these packages are available']
        msg += ['  ' + name for name in multi_found]
        msg += ["There should only be one OpenShift release repository enabled at a time."]
        AosVersionException.__init__(self, '\n'.join(msg), multi_found)


def _check_multi_minor_release(pkgs, expected_pkgs_dict):
    # see if any packages are available in more than one minor version
    pkgs_by_name_version = {}
    for pkg in pkgs:
        # keep track of x.y (minor release) versions seen
        minor_release = '.'.join(pkg.version.split('.')[:2])
        if pkg.name not in pkgs_by_name_version:
            pkgs_by_name_version[pkg.name] = set()
        pkgs_by_name_version[pkg.name].add(minor_release)

    multi_found = []
    for name in expected_pkgs_dict:
        if name in pkgs_by_name_version and len(pkgs_by_name_version[name]) > 1:
            multi_found.append(name)

    if multi_found:
        raise FoundMultiRelease(multi_found)


if __name__ == '__main__':
    main()
