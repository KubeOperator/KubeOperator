OpenShift Persistent Local Volumes
==================================

OpenShift Persistent Local Volumes

Requirements
------------

Role Variables
--------------

| Name                                      | Default value                                            |                                                                           |
|-------------------------------------------|----------------------------------------------------------|---------------------------------------------------------------------------|
| persistentlocalstorage_project            | local-storage                                            | The namespace where the Persistent Local Volume Provider will be deployed |
| persistentlocalstorage_classes            | []                                                       | Storage classes that will be created                                      |
| persistentlocalstorage_path               | /mnt/local-storage                                       | Path on the hosts that will be used as base for the local storage classes |
| persistentlocalstorage_provisionner_image | quay.io/external_storage/local-volume-provisioner:v1.0.1 | Docker image for the persistent storage volume provisionner               |

Dependencies
------------


Example Playbook
----------------

```
- name: Create persistent Local Storage Provider
  hosts: oo_first_master
  vars:
    persistentlocalstorage_project: local-storage
    persistentlocalstorage_classes:
    - ssd
    - hdd
  roles:
  - role: openshift_persistentlocalstorage
```

License
-------

Apache License, Version 2.0

Author Information
------------------

Diego Abelenda (diego.abelenda@camptocamp.com)
