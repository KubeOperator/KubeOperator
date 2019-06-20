OpenShift Master Certificates
========================

This role determines if OpenShift master certificates must be created, delegates certificate creation to the `openshift_ca_host` and then deploys those certificates to master hosts which this role is being applied to. If this role is applied to the `openshift_ca_host`, certificate deployment will be skipped.

Requirements
------------

Role Variables
--------------

From `openshift_ca`:

| Name                                  | Default value                                                             | Description                                                                                                                   |
|---------------------------------------|---------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| openshift_ca_host                     | None (Required)                                                           | The hostname of the system where the OpenShift CA will be (or has been) created.                                              |

From this role:

| Name                                  | Default value                                                             | Description                                                                                                                   |
|---------------------------------------|---------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| openshift_generated_configs_dir       | `{{ openshift.common.config_base }}/generated-configs`                    | Directory in which per-master generated config directories will be created on the `openshift_ca_host`.                        |
| openshift_master_cert_subdir          | `master-{{ openshift.common.hostname }}`                                  | Directory within `openshift_generated_configs_dir` where per-master configurations will be placed on the `openshift_ca_host`. |
| openshift_master_cert_expire_days     | `730` (2 years)                                                           | Validity of the certificates in days. Works only with OpenShift version 1.5 (3.5) and later.                                  |
| openshift_master_generated_config_dir | `{{ openshift_generated_configs_dir }}/{{ openshift_master_cert_subdir }` | Full path to the per-master generated config directory.                                                                       |

Dependencies
------------

* openshift_ca

Example Playbook
----------------

```
- name: Create OpenShift Master Certificates
  hosts: masters
  roles:
  - role: openshift_master_certificates
    openshift_ca_host: master1.example.com
```

License
-------

Apache License Version 2.0

Author Information
------------------

Jason DeTiberus (jdetiber@redhat.com)
