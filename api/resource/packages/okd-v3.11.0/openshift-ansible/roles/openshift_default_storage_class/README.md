openshift_master_storage_class
=========

A role that deploys configurations for Openshift StorageClass

Requirements
------------

None

Role Variables
--------------

openshift_storageclass_name: Name of the storage class to create
openshift_storageclass_provisioner: The kubernetes provisioner to use
openshift_storageclass_type: type of storage to use. This is different among clouds/providers

Dependencies
------------


Example Playbook
----------------

- role: openshift_default_storage_class
  openshift_storageclass_name: awsEBS
  openshift_storageclass_provisioner: kubernetes.io/aws-ebs
  openshift_storageclass_type: gp2


License
-------

Apache

Author Information
------------------

Openshift Operations
