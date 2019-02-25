# OpenShift Availability Monitoring

### **NOTE: This component is unsupported in OCP at this time.**

These are black box tests of several components which elable SLA verification by
capturing metrics for SLIs. All of the tests in this role should be used to
measure availability. Each test should expose metrics endpoints for scraping by
the monitoring platform (Prometheus).

All of the test applications are installed into the `openshift-monitor-availability` namespace and are enabled/disabled using the `openshift_monitor_availability_install` variable.

## Adding a new application

To add a new application to the installer:

1. Add an OpenShift Template to the `files` directory which can be used with `oc process | oc apply` to install the application.
2. Create an Ansible task in the `tasks` directory, e.g. `install_{APP}.yaml`. The task should install the application into the `openshift-monitor-availability` namespace.
3. Include the new task in `install.yaml`:

        - import_tasks: install_{APP}.yaml


## Guidelines

Here are some guidelines for applications:

* App metrics endpoints **must be secured**. Use the [oauth-proxy](https://github.com/openshift/oauth-proxy) or [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy).
* Templates should be usable outside the Ansible role (e.g. directly via `oc process`); avoid Jinjia templates if possible.
* As with all other Ansible roles in the installer, app tasks must be idempotent.
* Minimize configuration, be opinionated.

# License

Apache License, Version 2.0
