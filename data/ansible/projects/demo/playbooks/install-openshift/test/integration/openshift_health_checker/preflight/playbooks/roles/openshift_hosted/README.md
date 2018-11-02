OpenShift Hosted
================

OpenShift Hosted Resources

* OpenShift Router
* OpenShift Registry

Requirements
------------

This role requires a running OpenShift cluster.

Role Variables
--------------

From this role:

| Name                                  | Default value                            | Description                                                                                                              |
|---------------------------------------|------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|
| openshift_hosted_router_certificate   | None                                     | Dictionary containing "certfile", "keyfile" and "cafile" keys with values containing paths to local certificate files.   |
| openshift_hosted_router_registryurl   | 'registry.access.redhat.com/openshift3/ose-${component}:${version}' | The image to base the OpenShift router on.                                                                               |
| openshift_hosted_router_replicas      | Number of nodes matching selector        | The number of replicas to configure.                                                                                     |
| openshift_hosted_router_selector      | node-role.kubernetes.io/infra=true       | Node selector used when creating router. The OpenShift router will only be deployed to nodes matching this selector.     |
| openshift_hosted_router_name          | router                                   | The name of the router to be created.                                                                                    |
| openshift_hosted_registry_registryurl | 'registry.access.redhat.com/openshift3/ose-${component}:${version}' | The image to base the OpenShift registry on.                                                                             |
| openshift_hosted_registry_replicas    | Number of nodes matching selector        | The number of replicas to configure.                                                                                     |
| openshift_hosted_registry_selector    | node-role.kubernetes.io/infra=true                   | Node selector used when creating registry. The OpenShift registry will only be deployed to nodes matching this selector. |
| openshift_hosted_registry_cert_expire_days | `730` (2 years)                     | Validity of the certificates in days. Works only with OpenShift version 1.5 (3.5) and later.                             |
| openshift_hosted_registry_clusterip   | None                                     | Cluster IP for registry service                                                                                          |

If you specify `openshift_hosted_registry_kind=glusterfs`, the following
variables also control configuration behavior:

| Name                                         | Default value | Description                                                                  |
|----------------------------------------------|---------------|------------------------------------------------------------------------------|
| openshift_hosted_registry_storage_glusterfs_endpoints | glusterfs-registry-endpoints | The name for the Endpoints resource that will point the registry to the GlusterFS nodes
| openshift_hosted_registry_storage_glusterfs_path      | glusterfs-registry-volume    | The name for the GlusterFS volume that will provide registry storage
| openshift_hosted_registry_storage_glusterfs_readonly  | False                        | Whether the GlusterFS volume should be read-only
| openshift_hosted_registry_storage_glusterfs_swap      | False                        | Whether to swap an existing registry's storage volume for a GlusterFS volume
| openshift_hosted_registry_storage_glusterfs_swapcopy  | True                         | If swapping, copy the contents of the pre-existing registry storage to the new GlusterFS volume
| openshift_hosted_registry_storage_glusterfs_ips       | `[]`                         | A list of IP addresses of the nodes of the GlusterFS cluster to use for hosted registry storage

**NOTE:** Configuring a value for
`openshift_hosted_registry_storage_glusterfs_ips` with a `glusterfs_registry`
host group is not allowed. Specifying a `glusterfs_registry` host group
indicates that a new GlusterFS cluster should be configured, whereas
specifying `openshift_hosted_registry_storage_glusterfs_ips` indicates wanting
to use a pre-configured GlusterFS cluster for the registry storage.

_

Dependencies
------------

* openshift_persistent_volumes

Example Playbook
----------------

```
- name: Create hosted resources
  hosts: oo_first_master
  roles:
  - role: openshift_hosted
    openshift_hosted_router_certificate:
      certfile: /path/to/my-router.crt
      keyfile: /path/to/my-router.key
      cafile: /path/to/my-router-ca.crt
    openshift_hosted_router_registryurl: 'registry.access.redhat.com/openshift3/ose-haproxy-router:v3.0.2.0'
    openshift_hosted_router_selector: 'type=infra'
    openshift_hosted_registry_storage_kind=glusterfs
    openshift_hosted_registry_storage_glusterfs_path=external_glusterfs_volume_name
    openshift_hosted_registry_storage_glusterfs_ips=['192.168.20.239','192.168.20.96','192.168.20.114']

```

License
-------

Apache License, Version 2.0

Author Information
------------------

Red Hat openshift@redhat.com
