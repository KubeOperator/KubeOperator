#!/usr/bin/python
"""
Ansible module for rpm-based systems determining existing package version information in a host.
"""

from ansible.module_utils.basic import AnsibleModule
from ansible.module_utils.six import string_types

IMPORT_EXCEPTION = None
try:
    import rpm  # pylint: disable=import-error
except ImportError as err:
    IMPORT_EXCEPTION = err  # in tox test env, rpm import fails


class RpmVersionException(Exception):
    """Base exception class for package version problems"""
    def __init__(self, message, problem_pkgs=None):
        Exception.__init__(self, message)
        self.problem_pkgs = problem_pkgs


def main():
    """Entrypoint for this Ansible module"""
    module = AnsibleModule(
        argument_spec=dict(
            package_list=dict(type="list", required=True),
        ),
        supports_check_mode=True
    )

    if IMPORT_EXCEPTION:
        module.fail_json(msg="rpm_version module could not import rpm: %s" % IMPORT_EXCEPTION)

    # determine the packages we will look for
    pkg_list = module.params['package_list']
    if not pkg_list:
        module.fail_json(msg="package_list must not be empty")

    # get list of packages available and complain if any
    # of them are missing or if any errors occur
    try:
        pkg_versions = _retrieve_expected_pkg_versions(_to_dict(pkg_list))
        _check_pkg_versions(pkg_versions, _to_dict(pkg_list))
    except RpmVersionException as excinfo:
        module.fail_json(msg=str(excinfo))
    module.exit_json(changed=False)


def _to_dict(pkg_list):
    return {pkg["name"]: pkg for pkg in pkg_list}


def _retrieve_expected_pkg_versions(expected_pkgs_dict):
    """Search for installed packages matching given pkg names
    and versions. Returns a dictionary: {pkg_name: [versions]}"""

    transaction = rpm.TransactionSet()
    pkgs = {}

    for pkg_name in expected_pkgs_dict:
        matched_pkgs = transaction.dbMatch("name", pkg_name)
        if not matched_pkgs:
            continue

        for header in matched_pkgs:
            if header['name'] == pkg_name:
                if pkg_name not in pkgs:
                    pkgs[pkg_name] = []

                pkgs[pkg_name].append(header['version'])

    return pkgs


def _check_pkg_versions(found_pkgs_dict, expected_pkgs_dict):
    invalid_pkg_versions = {}
    not_found_pkgs = []

    for pkg_name, pkg in expected_pkgs_dict.items():
        if not found_pkgs_dict.get(pkg_name):
            not_found_pkgs.append(pkg_name)
            continue

        found_versions = [_parse_version(version) for version in found_pkgs_dict[pkg_name]]

        if isinstance(pkg["version"], string_types):
            expected_versions = [_parse_version(pkg["version"])]
        else:
            expected_versions = [_parse_version(version) for version in pkg["version"]]

        if not set(expected_versions) & set(found_versions):
            invalid_pkg_versions[pkg_name] = {
                "found_versions": found_versions,
                "required_versions": expected_versions,
            }

    if not_found_pkgs:
        raise RpmVersionException(
            '\n'.join([
                "The following packages were not found to be installed: {}".format('\n    '.join([
                    "{}".format(pkg)
                    for pkg in not_found_pkgs
                ]))
            ]),
            not_found_pkgs,
        )

    if invalid_pkg_versions:
        raise RpmVersionException(
            '\n    '.join([
                "The following packages were found to be installed with an incorrect version: {}".format('\n'.join([
                    "    \n{}\n    Required version: {}\n    Found versions: {}".format(
                        pkg_name,
                        ', '.join(pkg["required_versions"]),
                        ', '.join([version for version in pkg["found_versions"]]))
                    for pkg_name, pkg in invalid_pkg_versions.items()
                ]))
            ]),
            invalid_pkg_versions,
        )


def _parse_version(version_str):
    segs = version_str.split('.')
    if not segs or len(segs) <= 2:
        return version_str

    return '.'.join(segs[0:2])


if __name__ == '__main__':
    main()
