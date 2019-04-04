# OpenShift GlusterFS Playbooks

These playbooks are intended to enable the use of GlusterFS volumes by pods in
OpenShift. While they try to provide a sane set of defaults they do cover a
variety of scenarios and configurations, so read carefully. :)

## Playbook: config.yml

This is the main playbook that integrates GlusterFS into a new or existing
OpenShift cluster. It will also, if specified, configure a hosted Docker
registry with GlusterFS backend storage.

This playbook requires the `glusterfs` group to exist in the Ansible inventory
file. The hosts in this group are the nodes of the GlusterFS cluster.

 * If this is a newly configured cluster each host must have a
   `glusterfs_devices` variable defined, each of which must be a list of block
   storage devices intended for use only by the GlusterFS cluster. If this is
   also an external GlusterFS cluster, you must specify
   `openshift_storage_glusterfs_is_native=False`. If the cluster is to be
   managed by an external heketi service you must also specify
   `openshift_storage_glusterfs_heketi_is_native=False` and
   `openshift_storage_glusterfs_heketi_url=<URL>` with the URL to the heketi
   service. All these variables are specified in `[OSEv3:vars]`,
 * If this is an existing cluster you do not need to specify a list of block
   devices but you must specify the following variables in `[OSEv3:vars]`:
   * `openshift_storage_glusterfs_is_missing=False`
   * `openshift_storage_glusterfs_heketi_is_missing=False`
 * If GlusterFS will be running natively, the target hosts must also be listed
   in the `nodes` group. They must also already be configured as OpenShift
   nodes before this playbook runs.

By default, pods for a native GlusterFS cluster will be created in the
`default` namespace. To change this, specify
`openshift_storage_glusterfs_namespace=<other namespace>` in `[OSEv3:vars]`.

To configure the deployment of a Docker registry with GlusterFS backend
storage, specify `openshift_hosted_registry_storage_kind=glusterfs` in
`[OSEv3:vars]`. To create a separate GlusterFS cluster for use only by the
registry, specify a `glusterfs_registry` group that is populated as the
`glusterfs` is with the nodes for the separate cluster. If no
`glusterfs_registry` group is specified, the cluster defined by the `glusterfs`
group will be used.

To swap an existing hosted registry's backend storage for a GlusterFS volume,
specify `openshift_hosted_registry_storage_glusterfs_swap=True`. To
additoinally copy any existing contents from an existing hosted registry,
specify `openshift_hosted_registry_storage_glusterfs_swapcopy=True`.

**NOTE:** For each namespace that is to have access to GlusterFS volumes an
Enpoints resource pointing to the GlusterFS cluster nodes and a corresponding
Service resource must be created. If dynamic provisioning using StorageClasses
is configure, these resources are created automatically in the namespaces that
require them. This playbook also takes care of creating these resources in the
namespaces used for deployment.

An example of a minimal inventory file:
```
[OSEv3:children]
masters
nodes
glusterfs

[OSEv3:vars]
ansible_ssh_user=root
openshift_deployment_type=origin

[masters]
master

[nodes]
node0
node1
node2

[glusterfs]
node0 glusterfs_devices='[ "/dev/sdb" ]'
node1 glusterfs_devices='[ "/dev/sdb", "/dev/sdc" ]'
node2 glusterfs_devices='[ "/dev/sdd" ]'
```

## Playbook: registry.yml

This playbook is intended for admins who want to deploy a hosted Docker
registry with GlusterFS backend storage on an existing OpenShift cluster. It
has all the same requirements and behaviors as `config.yml`.

## Playbook: uninstall.yml

This playbook is intended to uninstall all GlusterFS related resources
on an existing OpenShift cluster.
It has all the same requirements and behaviors as `config.yml`.

If the variable `openshift_storage_glusterfs_wipe` is set as True,
it clears the backend data as well.

## Role: openshift_storage_glusterfs

The bulk of the work is done by the `openshift_storage_glusterfs` role. This
role can handle the deployment of GlusterFS (if it is to be hosted on the
OpenShift cluster), the registration of GlusterFS nodes (hosted or standalone),
and (if specified) integration as backend storage for a hosted Docker registry.

See the documentation in the role's directory for further details.

## Role: openshift_hosted

The `openshift_hosted` role recognizes `glusterfs` as a possible storage
backend for a hosted docker registry. It will also, if configured, handle the
swap of an existing registry's backend storage to a GlusterFS volume.
