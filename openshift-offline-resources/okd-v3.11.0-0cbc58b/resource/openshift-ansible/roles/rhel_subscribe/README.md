RHEL Subscribe
==============

Subscribes the RHEL servers and add the OpenShift enterprise repos.

Role variables
--------------

**NOTE**: `rhsub_user`/`rhsub_pass` and `rhsub_ak`/`rhsub_orgid` are mutually exclusive:
* If you want to use user/password to register the instance, it is required to
configure the `rhsub_user` and `rhsub_pass`.
* If you want to use an activation key to register the instance, it is required to
configure the `rhsub_ak` and `rhsub_orgid`.

### `rhsub_user`

Username for the subscription-manager.

### `rhsub_pass`

Password for the subscription-manager.

### `rhsub_ak`

Activation key for the subscription-manager.

### `rhsub_orgid`

Organization ID for the subscription-manager.

### `rhsub_pool`

Name of the pool to attach (optional).

### `rhsub_server`

Custom hostname for the Satellite server (optional).

### `openshift_release`

Version for the OpenShift Container Platform repositories.

Example: `3.6`
