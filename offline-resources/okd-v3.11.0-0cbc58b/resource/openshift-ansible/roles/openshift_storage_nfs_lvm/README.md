# openshift_storage_nfs_lvm

This role is useful to create and export nfs disks for openshift persistent volumes.
It does so by creating lvm partitions on an already setup pv/vg, creating xfs
filesystem on each partition, mounting the partitions, exporting the mounts via NFS
and creating a json file for each mount that an openshift master can use to
create persistent volumes.

## Requirements

* Ansible 2.2
* NFS server with NFS, iptables, and everything setup
* A lvm volume group created on the nfs server (default: openshiftvg)
* The lvm volume needs to have as much free space as you are allocating

## Role Variables

```
# Options of NFS exports.
osnl_nfs_export_options: "*(rw,sync,all_squash)"

# Directory, where the created partitions should be mounted. They will be
# mounted as <osnl_mount_dir>/<lvm volume name>
osnl_mount_dir: /exports/openshift

# Volume Group to use.
# This role always assumes that there is enough free space on the volume
#   group for all the partitions you will be making
osnl_volume_group: openshiftvg

# volume names
# volume names are {{osnl_volume_prefix}}{{osnl_volume_size}}g{{volume number}}
#   example: stg5g0004

# osnl_volume_prefix
# Useful if you are using the nfs server for more than one cluster
osnl_volume_prefix: "stg"

# osnl_volume_size
# Size of the volumes/partitions in Gigabytes.
osnl_volume_size: 5

# osnl_volume_num_start
# Where to start the volume number numbering.
osnl_volume_num_start: 3

# osnl_number_of_volumes
# How many volumes/partitions to build, with the size we stated.
osnl_number_of_volumes: 2

# osnl_volume_reclaim_policy
# Volume reclaim policy of a PersistentVolume tells the cluster
# what to do with the volume after it is released.
#
# Valid values are "Retain" or "Recycle" (default).
osnl_volume_reclaim_policy: "Recycle"

```

## Dependencies

None

## Example Playbook

With this playbook, 2 5Gig lvm partitions are created, named stg5g0003 and stg5g0004
Both of them are mounted into `/exports/openshift` directory.  Both directories are
exported via NFS.  json files are created in /root.

    - hosts: nfsservers
      remote_user: root
      gather_facts: no
      roles:
        - role: openshift_storage_nfs_lvm
          osnl_mount_dir: /exports/openshift
          osnl_volume_prefix: "stg"
          osnl_volume_size: 5
          osnl_volume_num_start: 3
          osnl_number_of_volumes: 2
          osnl_volume_reclaim_policy: "Recycle"


## Full example


* Create an `inventory` file:
    ```
    [nfsservers]
    10.0.0.1
    10.0.0.2
    ```

* Create an ansible playbook, say `setupnfs.yaml`:
    ```
    - hosts: nfsservers
      remote_user: root
      gather_facts: no
      roles:
        - role: openshift_storage_nfs_lvm
          osnl_mount_dir: /exports/stg
          osnl_volume_prefix: "stg"
          osnl_volume_size: 5
          osnl_volume_num_start: 3
          osnl_number_of_volumes: 2
          osnl_volume_reclaim_policy: "Recycle"

* Run the playbook:
    ```
    ansible-playbook -i inventory setupnfs.yml
    ```

## License

Apache 2.0
