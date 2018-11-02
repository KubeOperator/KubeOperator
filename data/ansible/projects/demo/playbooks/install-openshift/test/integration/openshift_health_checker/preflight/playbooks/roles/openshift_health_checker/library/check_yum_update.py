#!/usr/bin/python
'''
Ansible module to test whether a yum update or install will succeed,
without actually performing it or running yum.
parameters:
  packages: (optional) A list of package names to install or update.
            If omitted, all installed RPMs are considered for updates.
'''

import sys

import yum  # pylint: disable=import-error

from ansible.module_utils.basic import AnsibleModule


def main():  # pylint: disable=missing-docstring,too-many-branches
    module = AnsibleModule(
        argument_spec=dict(
            packages=dict(type='list', default=[])
        ),
        supports_check_mode=True
    )

    def bail(error):  # pylint: disable=missing-docstring
        module.fail_json(msg=error)

    yb = yum.YumBase()  # pylint: disable=invalid-name
    yb.conf.disable_excludes = ["all"]  # assume the openshift excluder will be managed, ignore current state
    # determine if the existing yum configuration is valid
    try:
        yb.repos.populateSack(mdtype='metadata', cacheonly=1)
    # for error of type:
    #   1. can't reach the repo URL(s)
    except yum.Errors.NoMoreMirrorsRepoError as e:  # pylint: disable=invalid-name
        bail('Error getting data from at least one yum repository: %s' % e)
    #   2. invalid repo definition
    except yum.Errors.RepoError as e:  # pylint: disable=invalid-name
        bail('Error with yum repository configuration: %s' % e)
    #   3. other/unknown
    #    * just report the problem verbatim
    except:  # pylint: disable=bare-except; # noqa
        bail('Unexpected error with yum repository: %s' % sys.exc_info()[1])

    packages = module.params['packages']
    no_such_pkg = []
    for pkg in packages:
        try:
            yb.install(name=pkg)
        except yum.Errors.InstallError as e:  # pylint: disable=invalid-name
            no_such_pkg.append(pkg)
        except:  # pylint: disable=bare-except; # noqa
            bail('Unexpected error with yum install/update: %s' %
                 sys.exc_info()[1])
    if not packages:
        # no packages requested means test a yum update of everything
        yb.update()
    elif no_such_pkg:
        # wanted specific packages to install but some aren't available
        user_msg = 'Cannot install all of the necessary packages. Unavailable:\n'
        for pkg in no_such_pkg:
            user_msg += '  %s\n' % pkg
        user_msg += 'You may need to enable one or more yum repositories to make this content available.'
        bail(user_msg)

    try:
        txn_result, txn_msgs = yb.buildTransaction()
    except:  # pylint: disable=bare-except; # noqa
        bail('Unexpected error during dependency resolution for yum update: \n %s' %
             sys.exc_info()[1])

    # find out if there are any errors with the update/install
    if txn_result == 0:  # 'normal exit' meaning there's nothing to install/update
        pass
    elif txn_result == 1:  # error with transaction
        user_msg = 'Could not perform a yum update.\n'
        if len(txn_msgs) > 0:
            user_msg += 'Errors from dependency resolution:\n'
            for msg in txn_msgs:
                user_msg += '  %s\n' % msg
            user_msg += 'You should resolve these issues before proceeding with an install.\n'
            user_msg += 'You may need to remove or downgrade packages or enable/disable yum repositories.'
        bail(user_msg)
    # TODO: it would be nice depending on the problem:
    #   1. dependency for update not found
    #    * construct the dependency tree
    #    * find the installed package(s) that required the missing dep
    #    * determine if any of these packages matter to openshift
    #    * build helpful error output
    #   2. conflicts among packages in available content
    #    * analyze dependency tree and build helpful error output
    #   3. other/unknown
    #    * report the problem verbatim
    #    * add to this list as we come across problems we can clearly diagnose
    elif txn_result == 2:  # everything resolved fine
        pass
    else:
        bail('Unknown error(s) from dependency resolution. Exit Code: %d:\n%s' %
             (txn_result, txn_msgs))

    module.exit_json(changed=False)


if __name__ == '__main__':
    main()
