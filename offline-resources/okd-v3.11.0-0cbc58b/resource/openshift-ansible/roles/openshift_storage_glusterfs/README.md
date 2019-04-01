OpenShift GlusterFS Cluster
===========================

OpenShift GlusterFS Cluster Configuration

This role handles the configuration of GlusterFS clusters. It can handle
two primary configuration scenarios:

* Configuring a new, natively-hosted GlusterFS cluster. In this scenario,
  GlusterFS pods are deployed on nodes in the OpenShift cluster which are
  configured to provide storage.
* Configuring a new, external GlusterFS cluster. In this scenario, the
  cluster nodes have the GlusterFS software pre-installed but have not
  been configured yet. The installer will take care of configuring the
  cluster(s) for use by OpenShift applications.
* Using existing GlusterFS clusters. In this scenario, one or more
  GlusterFS clusters are assumed to be already setup. These clusters can
  be either natively-hosted or external, but must be managed by a
  [heketi service](https://github.com/heketi/heketi).

As part of the configuration, a particular GlusterFS cluster may be
specified to provide backend storage for a natively-hosted Docker
registry.

Unless configured otherwise, a StorageClass will be automatically
created for each non-registry GlusterFS cluster. This will allow
applications which can mount PersistentVolumes to request
dynamically-provisioned GlusterFS volumes.

Requirements
------------

* Ansible 2.2

Host Groups
-----------

The following group is expected to be populated for this role to run:

* `[glusterfs]`

Additionally, the following group may be specified either in addition to or
instead of the above group to deploy a GlusterFS cluster for use by a natively
hosted Docker registry:

* `[glusterfs_registry]`

Host Variables
--------------

For configuring new clusters, the following role variables are available.

Each host in either of the above groups must have the following variable
defined:

| Name              | Default value | Description                             |
|-------------------|---------------|-----------------------------------------|
| glusterfs_devices | None          | A list of block devices that will be completely managed as part of a GlusterFS cluster. There must be at least one device listed. Each device must be bare, e.g. no partitions or LVM PVs. **Example:** '[ "/dev/sdb" ]' **NOTE:** You MUST set this as a host variable on each node host. For some reason, if you set this as a group variable it gets interpreted as a string rather than an array. See https://github.com/openshift/openshift-ansible/issues/5071

In addition, each host may specify the following variables to further control
their configuration as GlusterFS nodes:

| Name               | Default value             | Description                             |
|--------------------|---------------------------|-----------------------------------------|
| glusterfs_cluster  | 1                         | The ID of the cluster this node should belong to. This is useful when a single heketi service is expected to manage multiple distinct clusters. **NOTE:** For natively-hosted clusters, all pods will be in the same OpenShift namespace
| glusterfs_hostname | l_kubelet_node_name  | A hostname (or IP address) that will be used for internal GlusterFS communication
| glusterfs_ip       | openshift.common.ip       | An IP address that will be used by pods to communicate with the GlusterFS node. **NOTE:** Required for external GlusterFS nodes
| glusterfs_zone     | 1                         | A zone number for the node. Zones are used within the cluster for determining how to distribute the bricks of GlusterFS volumes. heketi will try to spread each volumes' bricks as evenly as possible across all zones

Role Variables
--------------

This role has the following variables that control the integration of a
GlusterFS cluster into a new or existing OpenShift cluster:

| Name                                                   | Default value           | Description                             |
|--------------------------------------------------------|-------------------------|-----------------------------------------|
| openshift_storage_glusterfs_timeout                    | 300                     | Seconds to wait for pods to become ready
| openshift_storage_glusterfs_namespace                  | 'glusterfs'             | Namespace/project in which to create GlusterFS resources
| openshift_storage_glusterfs_is_native                  | True                    | GlusterFS should be containerized
| openshift_storage_glusterfs_name                       | 'storage'               | A name to identify the GlusterFS cluster, which will be used in resource names
| openshift_storage_glusterfs_nodeselector               | 'glusterfs=storage-host'| Selector to determine which nodes will host GlusterFS pods in native mode. **NOTE:** The label value is taken from the cluster name
| openshift_storage_glusterfs_use_default_selector       | False                   | Whether to use a default node selector for the GlusterFS namespace/project. If False, the namespace/project will have no restricting node selector. If True, uses pre-existing or default (e.g. osm_default_node_selector) node selectors. **NOTE:** If True, nodes which will host GlusterFS pods must already have the additional labels.
| openshift_storage_glusterfs_storageclass               | True                    | Automatically create a GlusterFS StorageClass for this group
| openshift_storage_glusterfs_storageclass_default       | False                   | Sets the GlusterFS StorageClass for this group as cluster-wide default
| openshift_storage_glusterfs_image                      | 'gluster/gluster-centos'| Container image to use for GlusterFS pods, enterprise default is 'rhgs3/rhgs-server-rhel7'
| openshift_storage_glusterfs_block_deploy               | True                    | Deploy glusterblock provisioner service
| openshift_storage_glusterfs_block_image                | 'gluster/glusterblock-provisioner'| Container image to use for glusterblock-provisioner pod, enterprise default is 'rhgs3/rhgs-gluster-block-prov-rhel7'
| openshift_storage_glusterfs_block_host_vol_create      | True                    | Automatically create GlusterFS volumes to host glusterblock volumes. **NOTE:** If this is False, block-hosting volumes will need to be manually created before glusterblock volumes can be provisioned
| openshift_storage_glusterfs_block_host_vol_size        | 100                     | Size, in GB, of GlusterFS volumes that will be automatically created to host glusterblock volumes if not enough space is available for a glusterblock volume create request. **NOTE:** This value is effectively an upper limit on the size of glusterblock volumes unless you manually create larger GlusterFS block-hosting volumes
| openshift_storage_glusterfs_block_host_vol_max         | 15                      | Max number of GlusterFS volumes to host glusterblock volumes
| openshift_storage_glusterfs_block_storageclass         | False                   | Automatically create a StorageClass for each glusterblock cluster
| openshift_storage_glusterfs_block_storageclass_default | False                   | Sets the glusterblock StorageClass for this group as cluster-wide default
| openshift_storage_glusterfs_s3_deploy                  | True                    | Deploy gluster-s3 service
| openshift_storage_glusterfs_s3_image                   | 'gluster/gluster-object'| Container image to use for gluster-s3 pod, enterprise default is 'rhgs3/rhgs-s3-server-rhel7'
| openshift_storage_glusterfs_s3_account                 | Undefined               | S3 account name for the S3 service, required for S3 service deployment
| openshift_storage_glusterfs_s3_user                    | Undefined               | S3 user name for the S3 service, required for S3 service deployment
| openshift_storage_glusterfs_s3_password                | Undefined               | S3 user password for the S3 service, required for S3 service deployment
| openshift_storage_glusterfs_s3_pvc                     | Dynamic                 | Name of the GlusterFS-backed PVC which will be used for S3 object data storage, generated from the cluster name and S3 account by default
| openshift_storage_glusterfs_s3_pvc_size                | "2Gi"                   | Size, in Gi, of the GlusterFS-backed PVC which will be used for S3 object data storage
| openshift_storage_glusterfs_s3_meta_pvc                | Dynamic                 | Name of the GlusterFS-backed PVC which will be used for S3 object metadata storage, generated from the cluster name and S3 account by default
| openshift_storage_glusterfs_s3_meta_pvc_size           | "1Gi"                   | Size, in Gi, of the GlusterFS-backed PVC which will be used for S3 object metadata storage
| openshift_storage_glusterfs_wipe                       | False                   | Destroy any existing GlusterFS resources and wipe storage devices. **WARNING: THIS WILL DESTROY ANY DATA ON THOSE DEVICES.**
| openshift_storage_glusterfs_heketi_is_native           | True                    | heketi should be containerized
| openshift_storage_glusterfs_heketi_cli                 | 'heketi-cli'            | Command/Path to invoke the heketi-cli tool **NOTE:** Change this only for **non-native heketi** if heketi-cli is not in the global `$PATH` of the machine running openshift-ansible
| openshift_storage_glusterfs_heketi_image               | 'heketi/heketi'         | Container image to use for heketi pods, enterprise default is 'rhgs3/rhgs-volmanager-rhel7'
| openshift_storage_glusterfs_heketi_admin_key           | auto-generated          | String to use as secret key for performing heketi commands as admin
| openshift_storage_glusterfs_heketi_user_key            | auto-generated          | String to use as secret key for performing heketi commands as user that can only view or modify volumes
| openshift_storage_glusterfs_heketi_topology_load       | True                    | Load the GlusterFS topology information into heketi
| openshift_storage_glusterfs_heketi_url                 | Undefined               | When heketi is native, this sets the hostname portion of the final heketi route URL. When heketi is external, this is the FQDN or IP address to the heketi service.
| openshift_storage_glusterfs_heketi_port                | 8080                    | TCP port for external heketi service **NOTE:** This has no effect in native mode
| openshift_storage_glusterfs_heketi_executor            | 'kubernetes'            | Selects how a native heketi service will manage GlusterFS nodes: 'kubernetes' for native nodes, 'ssh' for external nodes
| openshift_storage_glusterfs_heketi_ssh_port            | 22                      | SSH port for external GlusterFS nodes via native heketi
| openshift_storage_glusterfs_heketi_ssh_user            | 'root'                  | SSH user for external GlusterFS nodes via native heketi
| openshift_storage_glusterfs_heketi_ssh_sudo            | False                   | Whether to sudo (if non-root user) for SSH to external GlusterFS nodes via native heketi
| openshift_storage_glusterfs_heketi_ssh_keyfile         | Undefined               | Path to a private key file for use with SSH connections to external GlusterFS nodes via native heketi **NOTE:** This must be an absolute path
| openshift_storage_glusterfs_heketi_fstab               | '/var/lib/heketi/fstab' | When heketi is native, sets the path to the fstab file on the GlusterFS nodes to update on LVM volume mounts, changes to '/etc/fstab/' when the heketi executor is 'ssh' **NOTE:** This should not need to be changed
| openshift_storage_glusterfs_heketi_wipe                | False                   | Destroy any existing heketi resources, defaults to the value of `openshift_storage_glusterfs_wipe`

Each role variable also has a corresponding variable to optionally configure a
separate GlusterFS cluster for use as storage for an integrated Docker
registry. These variables start with the prefix
`openshift_storage_glusterfs_registry_` and, for the most part, default to the
values in their corresponding non-registry variables. The following variables
are an exception:

| Name                                                            | Default value         | Description                             |
|-----------------------------------------------------------------|-----------------------|-----------------------------------------|
| openshift_storage_glusterfs_registry_namespace                  | registry namespace    | Default is to use the hosted registry's namespace, otherwise 'glusterfs'
| openshift_storage_glusterfs_registry_name                       | 'registry'            | This allows for the logical separation of the registry group from other Gluster groups
| openshift_storage_glusterfs_registry_storageclass               | False                 | It is recommended to not create a StorageClass for this group, so as to avoid noisy neighbor complications
| openshift_storage_glusterfs_registry_storageclass_default       | False                 | Separate from the above
| openshift_storage_glusterfs_registry_block_storageclass         | False                 | Only enable this for use by Logging and Metrics
| openshift_storage_glusterfs_registry_block_storageclass_default | False                 | Separate from the above
| openshift_storage_glusterfs_registry_heketi_admin_key           | auto-generated        | Separate from the above
| openshift_storage_glusterfs_registry_heketi_user_key            | auto-generated        | Separate from the above

Additionally, this role's behavior responds to several registry-specific variables in the [openshift_hosted role](../openshift_hosted/README.md):

| Name                                                  | Default value                | Description                             |
|-------------------------------------------------------|------------------------------|-----------------------------------------|
| openshift_hosted_registry_storage_glusterfs_endpoints | glusterfs-registry-endpoints | The name for the Endpoints resource that will point the registry to the GlusterFS nodes
| openshift_hosted_registry_storage_glusterfs_path      | glusterfs-registry-volume    | The name for the GlusterFS volume that will provide registry storage
| openshift_hosted_registry_storage_glusterfs_readonly  | False                        | Whether the GlusterFS volume should be read-only
| openshift_hosted_registry_storage_glusterfs_swap      | False                        | Whether to swap an existing registry's storage volume for a GlusterFS volume
| openshift_hosted_registry_storage_glusterfs_swapcopy  | True                         | If swapping, copy the contents of the pre-existing registry storage to the new GlusterFS volume

Dependencies
------------

* os_firewall
* openshift_repos
* lib_openshift

Example Playbook
----------------

```
- name: Configure GlusterFS hosts
  hosts: oo_first_master
  roles:
  - role: openshift_storage_glusterfs
    when: groups.oo_glusterfs_to_config | default([]) | count > 0
```

License
-------

Apache License, Version 2.0

Author Information
------------------

Jose A. Rivera (jarrpa@redhat.com)
