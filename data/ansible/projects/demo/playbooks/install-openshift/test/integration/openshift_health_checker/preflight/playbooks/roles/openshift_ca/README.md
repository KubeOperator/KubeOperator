OpenShift CA
============

This role delegates all tasks to the `openshift_ca_host` such that this role can be depended on by other OpenShift certificate roles.

Requirements
------------

Role Variables
--------------

From this role:

| Name                    | Default value                                 | Description                                                                 |
|-------------------------|-----------------------------------------------|-----------------------------------------------------------------------------|
| openshift_ca_host       | None (Required)                               | The hostname of the system where the OpenShift CA will be created.          |
| openshift_ca_config_dir | `{{ openshift.common.config_base }}/master`   | CA certificate directory.                                                   |
| openshift_ca_cert       | `{{ openshift_ca_config_dir }}/ca.crt`        | CA certificate path including CA certificate filename.                      |
| openshift_ca_key        | `{{ openshift_ca_config_dir }}/ca.key`        | CA key path including CA key filename.                                      |
| openshift_ca_serial     | `{{ openshift_ca_config_dir }}/ca.serial.txt` | CA serial path including CA serial filename.                                |
| openshift_master_cert_expire_days | `730` (2 years)                     | Validity of the certificates in days. Works only with OpenShift version 1.5 (3.5) and later. |
| openshift_ca_cert_expire_days     | `1825` (5 years)                    | Validity of the CA certificates in days. Works only with OpenShift version 1.5 (3.5) and later. |

Dependencies
------------

* openshift_repos
* openshift_cli

Example Playbook
----------------

```
- name: Create OpenShift CA
  hosts: localhost
  roles:
  - role: openshift_ca
```

License
-------

Apache License Version 2.0

Author Information
------------------

Jason DeTiberus (jdetiber@redhat.com)
