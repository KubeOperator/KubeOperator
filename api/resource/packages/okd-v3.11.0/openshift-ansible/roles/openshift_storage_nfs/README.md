OpenShift NFS Server
====================

OpenShift NFS Server Installation

Requirements
------------

* Ansible 2.2
* This role is intended to be applied to the [nfs] host group which is
separate from OpenShift infrastructure components.
* Requires access to the 'nfs-utils' package.

Role Variables
--------------

From this role:

| Name                                            | Default value         |                                                             |
|-------------------------------------------------|-----------------------|-------------------------------------------------------------|
| openshift_hosted_registry_storage_nfs_directory | /exports              | Root export directory.                                      |
| openshift_hosted_registry_storage_volume_name   | registry              | Registry volume within openshift_hosted_registry_volume_dir |
| openshift_hosted_registry_storage_nfs_options   | *(rw,root_squash)     | NFS options for configured exports.                         |

Dependencies
------------

* os_firewall
* openshift_facts
* openshift_repos

Example Playbook
----------------

```
- name: Configure nfs hosts
  hosts: oo_nfs_to_config
  roles:
  - role: openshift_storage_nfs
```

License
-------

Apache License, Version 2.0

Author Information
------------------

Andrew Butcher (abutcher@redhat.com)
