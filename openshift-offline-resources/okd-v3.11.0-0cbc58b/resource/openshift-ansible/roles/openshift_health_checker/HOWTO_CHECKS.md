# OpenShift health checks

This Ansible role contains health checks to diagnose problems in OpenShift
environments.

Checks are typically implemented as two parts:

1. a Python module in [openshift_checks/](openshift_checks), with a class that
   inherits from `OpenShiftCheck`.
2. a custom Ansible module in [library/](library), for cases when the modules
   shipped with Ansible do not provide the required functionality.

The checks are called from Ansible playbooks via the `openshift_health_check`
action plugin. See
[playbooks/openshift-checks/pre-install.yml](../../playbooks/openshift-checks/pre-install.yml)
for an example.

The action plugin dynamically discovers all checks and executes only those
selected in the play.

Checks can determine when they are active by implementing the method
`is_active`. Inactive checks are skipped. This is similar to the `when`
instruction in Ansible plays.

Checks may have tags, which are a way to group related checks together. For
instance, to run all preflight checks, pass in the group `'@preflight'` to
`openshift_health_check`.

Groups are automatically computed from tags.

Groups and individual check names can be used together in the argument list to
`openshift_health_check`.

Look at existing checks for the implementation details.
