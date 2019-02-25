# OpenShift Metering

This playbook runs the [Openshift Metering role](../../roles/openshift_metering).
See the role for more information.

## Prequisites:

This playbook requires Openshift Cluster Monitoring, which is installed by default.
If Openshift Cluster Monitoring is not installed, check that the variable
openshift\_cluster\_monitoring\_operator\_install is set to true.

```yaml
openshift_cluster_monitoring_operator_install: true
```

## Installation

To install Openshift Metering, run the install playbook:

```bash
ansible-playbook playbooks/openshift-metering/config.yml
```

To uninstall, run the uninstall playbook:

```bash
ansible-playbook playbooks/openshift-metering/uninstall.yml
```

