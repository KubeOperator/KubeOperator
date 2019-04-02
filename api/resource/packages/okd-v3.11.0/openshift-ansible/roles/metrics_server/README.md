OpenShift metrics-server
========================

OpenShift metrics-server Installation

Requirements
------------
The following variables need to be set and will be validated:

- `openshift_metrics_server_project`: project (i.e. namespace) where the
  components will be deployed.

Role Variables
--------------

- `openshift_metrics_server_resolution`: How often metrics should be
  gathered.

Dependencies
------------
openshift_facts


Example Playbook
----------------

```
- name: Configure openshift-metrics-server
  hosts: oo_first_master
  roles:
  - role: openshift_metrics_server
```

License
-------

Apache License, Version 2.0

Author Information
------------------

Solly Ross <sross@redhat.com>
