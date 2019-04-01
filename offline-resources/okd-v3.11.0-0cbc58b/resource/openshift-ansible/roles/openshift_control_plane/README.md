OpenShift Control Plane
==================================

Installs the services that comprise the OpenShift control plane onto nodes that are preconfigured for
bootstrapping.

Requirements
------------

* Ansible 2.2
* A RHEL 7.1 host pre-configured with access to the rhel-7-server-rpms,
rhel-7-server-extras-rpms, and rhel-7-server-ose-3.0-rpms repos.

Role Variables
--------------

From this role:

| Name                                             | Default value         |                                                                               |
|---------------------------------------------------|-----------------------|-------------------------------------------------------------------------------|
| openshift_node_ips                                | []                    | List of the openshift node ip addresses to pre-register when master starts up |
| oreg_url                                          | UNDEF                 | Default docker registry to use                                                |                                                                               |
| openshift_master_console_port                     | UNDEF                 |                                                                               |
| openshift_master_api_url                          | UNDEF                 |                                                                               |
| openshift_master_console_url                      | UNDEF                 |                                                                               |
| openshift_persistentlocalstorage_enabled          | false                 | Enable the persistent local storage                                           |
| openshift_master_public_api_url                   | UNDEF                 |                                                                               |
| openshift_master_public_console_url               | UNDEF                 |                                                                               |
| openshift_master_saconfig_limit_secret_references | false                 |                                                                               |


Dependencies
------------


Example Playbook
----------------

TODO

License
-------

Apache License, Version 2.0

Author Information
------------------

TODO
