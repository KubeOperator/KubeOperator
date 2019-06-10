# OpenShift Auto-heal Service

The OpenShift Auto-heal Service recevies alert notifications from the
[Prometheus alert manager](https://prometheus.io/docs/alerting/alertmanager) and
tries to solve the root cause executing Ansible
[Tower](https://www.ansible.com/products/tower) or
[AWX](https://github.com/ansible/awx) jobs.

# Installation

See the [installation playbook](../../playbooks/openshift-autoheal) uses the
following variables:

- `openshift_autoheal_deploy`: `true` - install/update. `false` - uninstall.
  Defaults to `false`.

- `openshift_autoheal_config`: The content of the configuration of the
  auto-heal service, as described in the [documentation](https://github.com/openshift/autoheal)
  of the service and in this [example](https://github.com/openshift/autoheal/blob/master/autoheal.yml).

# Requirements

Ansible 2.4.

## License

Apache license, version 2.0.
