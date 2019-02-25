# OpenShift-Ansible Git Hooks

## Introduction

This `hack` sub-directory holds
[git commit hooks](https://www.atlassian.com/git/tutorials/git-hooks#conceptual-overview)
you may use when working on openshift-ansible contributions. See the
README in each sub-directory for an overview of what each hook does
and if the hook has any specific usage or setup instructions.

## Usage

Basic git hook usage is simple:

1) Copy (or symbolic link) the hook to the `$REPO_ROOT/.git/hooks/` directory
2) Make the hook executable (`chmod +x $PATH_TO_HOOK`)

## Multiple Hooks of the Same Type

If you want to install multiple hooks of the same type, for example:
multiple `pre-commit` hooks, you will need some kind of *hook
dispatcher*. For an example of an easy to use hook dispatcher check
out this gist by carlos-jenkins:

* [multihooks.py](https://gist.github.com/carlos-jenkins/89da9dcf9e0d528ac978311938aade43)

## Contributing Hooks

If you want to contribute a new hook there are only a few criteria
that must be met:

* The hook **MUST** include a README describing the purpose of the hook
* The README **MUST** describe special setup instructions if they are required
* The hook **MUST** be in a sub-directory of this directory
* The hook file **MUST** be named following the standard git hook
  naming pattern (i.e., pre-commit hooks **MUST** be called
  `pre-commit`)
