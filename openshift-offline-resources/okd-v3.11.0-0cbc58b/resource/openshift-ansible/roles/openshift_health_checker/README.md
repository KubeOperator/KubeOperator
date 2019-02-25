OpenShift Health Checker
========================

This role detects common problems with OpenShift installations or with
environments prior to install.

For more information about creating new checks, see [HOWTO_CHECKS.md](HOWTO_CHECKS.md).

Requirements
------------

* Ansible 2.2+

Role Variables
--------------

None

Dependencies
------------

- openshift_facts

Example Playbook
----------------

```yaml
---
- hosts: OSEv3
  name: run OpenShift health checks
  roles:
    - openshift_health_checker
  post_tasks:
    - action: openshift_health_check
```

License
-------

Apache License Version 2.0

Author Information
------------------

Customer Success team (dev@lists.openshift.redhat.com)
