# openshift_expand_partition

This role is useful to expand a partition, and it's file system to
fully utilize the disk it is on.  It does this by first expanding the
partition, and then expanding the file system on the partition.

## Requirements

* A machine with a disk that is not fully utilized

* cloud-utils-growpart rpm (either installed or avialable via yum or dnf)

* The partition you are expanding needs to be at the end of the partition list

## Role Variables

```
# The following variables are if you want to expand
#   /dev/xvda3 that has a filesystem xfs

# oep_drive
# Drive that has the partition we wish to expand.
oep_drive: "/dev/xvda"

# oep_partition
# Partition that we wish to expand.
oep_partition: 3

# oep_file_system
# What file system is on the partition
#   Currently only xfs, and ext(2,3,4) are supported
#   For ext2, ext3, or ext4 just use ext
oep_file_system: "xfs"

```

## Dependencies

growpart

## Example Playbook

With this playbook, the partition /dev/xvda3 will expand to fill the free
space on /dev/xvda, and the file system will be expanded to fill the new
partition space.

    - hosts: mynodes
      remote_user: root
      gather_facts: no
      roles:
        - role: openshift_expand_partition
          oep_drive: "/dev/xvda"
          oep_partition: 3
          oep_file_system: "xfs"


## Full example


* Create an `inventory` file:
    ```
    [mynodes]
    10.0.0.1
    10.0.0.2
    ```

* Create an ansible playbook, say `expandvar.yaml`:
    ```
    - hosts: mynodes
      remote_user: root
      gather_facts: no
      roles:
        - role: openshift_expand_partition
          oep_drive: "/dev/xvda"
          oep_partition: 3
          oep_file_system: "xfs"

* Run the playbook:
    ```
    ansible-playbook -i inventory expandvar.yml
    ```

## License

Apache 2.0
