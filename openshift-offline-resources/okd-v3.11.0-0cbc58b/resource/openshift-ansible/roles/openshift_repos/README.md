OpenShift Repos
================

Configures repositories for an OpenShift installation

Requirements
------------

A RHEL 7.1 host pre-configured with access to the rhel-7-server-rpms,
rhel-7-server-extra-rpms, and rhel-7-server-ose-3.0-rpms repos.

Role Variables
--------------

| Name                          | Default value |                                              |
|-------------------------------|---------------|----------------------------------------------|
| openshift_deployment_type     | None          | Possible values openshift-enterprise, origin |
| openshift_additional_repos    | {}            | TODO                                         |

Dependencies
------------

None.

Example Playbook
----------------

TODO

License
-------

Apache License, Version 2.0

Author Information
------------------

TODO
