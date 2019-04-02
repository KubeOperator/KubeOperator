OpenShift Excluder
==================

Manages the excluder packages which add yum and dnf exclusions ensuring that
the packages we care about are not inadvertently updated. See
https://github.com/openshift/origin/tree/master/contrib/excluder

Requirements
------------

None

Inventory Variables
-------------------

| Name                                 | Default Value              | Description                            |
---------------------------------------|----------------------------|----------------------------------------|
| openshift_enable_excluders           | True                       | Enable all excluders                   |
| openshift_enable_docker_excluder     | openshift_enable_excluders | Enable docker excluder. If not set, the docker excluder is ignored. |
| openshift_enable_openshift_excluder  | openshift_enable_excluders | Enable openshift excluder. If not set, the openshift excluder is ignored. |

Role Variables
--------------

| Name                                      | Default | Choices         | Description                                                               |
|-------------------------------------------|---------|-----------------|---------------------------------------------------------------------------|
| r_openshift_excluder_action               | enable  | enable, disable | Action to perform when calling this role                                  |
| r_openshift_excluder_verify_upgrade       | false   | true, false     | When upgrading, this variable should be set to true when calling the role |
| r_openshift_excluder_package_state        | present | present, latest | Use 'latest' to upgrade openshift_excluder package                        |
| r_openshift_excluder_docker_package_state | present | present, latest | Use 'latest' to upgrade docker_excluder package                           |
| r_openshift_excluder_service_type         | None    |                 | (Required) Defined as openshift_service_type e.g. atomic-openshift        |
| r_openshift_excluder_upgrade_target       | None    |                 | Required when r_openshift_excluder_verify_upgrade is true, defined as openshift_upgrade_target by Upgrade playbooks e.g. '3.6'|

Dependencies
------------

- lib_utils

Example Playbook
----------------

```yaml
- name: Demonstrate OpenShift Excluder usage
  hosts: oo_masters_to_config:oo_nodes_to_config
  roles:
  # Disable all excluders
  - role: openshift_excluder
    r_openshift_excluder_action: disable
  # Enable all excluders
  - role: openshift_excluder
    r_openshift_excluder_action: enable
  # Disable all excluders and verify appropriate excluder packages are available for upgrade
  - role: openshift_excluder
    r_openshift_excluder_action: disable
    r_openshift_excluder_verify_upgrade: true
    r_openshift_excluder_upgrade_target: "{{ openshift_upgrade_target }}"
    r_openshift_excluder_package_state: latest
    r_openshift_excluder_docker_package_state: latest
```

TODO
----

It should be possible to manage the two excluders independently though that's not a hard requirement. However it should be done to manage docker on RHEL Containerized hosts.

License
-------

Apache License, Version 2.0

Author Information
------------------

Scott Dodson (sdodson@redhat.com)
