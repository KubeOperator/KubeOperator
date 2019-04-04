Openshift cluster autoscaler
================================

Install the cluster autoscaler

Requirements
------------

* One or more Master servers
* A cloud provider that supports the cluster-autoscaler

Role Variables
--------------
Check defaults/main.yml

Dependencies
------------


Example Playbook
----------------

#!/usr/bin/ansible-playbook
```
---
- hosts: masters
  gather_facts: no
  remote_user: root
  tasks:
  - name: include role autoscaler
    import_role:
      name: openshift_cluster_autoscaler
    vars:
      openshift_clusterid: opstest
      openshift_cluster_autoscaler_aws_key: <aws_key>
      openshift_cluster_autoscaler_aws_secret_key: <aws_secret_key>
```


Notes
-----

This is currently experimental software.  This role allows users to install the cluster-autoscaler and the necessary authorization pieces that allow the autoscaler to function.


This feature requires cloud provider credentials or a serviceaccount that has access to scale up/down nodes within the scaling groups.

https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler

License
-------

Apache License, Version 2.0

Author Information
------------------

Openshift
