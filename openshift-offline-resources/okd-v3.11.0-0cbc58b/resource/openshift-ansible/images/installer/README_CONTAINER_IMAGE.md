ORIGIN-ANSIBLE IMAGE INSTALLER
===============================

Contains Dockerfile information for building an openshift/origin-ansible image
based on `centos:7` or `rhel7.3:7.3-released`.

Read additional setup information for this image at: https://hub.docker.com/r/openshift/origin-ansible/

Read additional information about the `openshift/origin-ansible` at: https://github.com/openshift/openshift-ansible/blob/master/README_CONTAINER_IMAGE.md

Also contains necessary components for running the installer using an Atomic System Container.


System container installer
==========================

These files are needed to run the installer using an [Atomic System container](http://www.projectatomic.io/blog/2016/09/intro-to-system-containers/).
These files can be found under `root/exports`:

* config.json.template - Template of the configuration file used for running containers.

* manifest.json - Used to define various settings for the system container, such as the default values to use for the installation.

* service.template - Template file for the systemd service.

* tmpfiles.template - Template file for systemd-tmpfiles.

These files can be found under `root/usr/local/bin`:

* run-system-container.sh - Entrypoint to the container.

## Options

These options may be set via the ``atomic`` ``--set`` flag. For defaults see ``root/exports/manifest.json``

* OPTS - Additional options to pass to ansible when running the installer

* VAR_LIB_OPENSHIFT_INSTALLER - Full path of the installer code to mount into the container

* VAR_LOG_OPENSHIFT_LOG - Full path of the log file to mount into the container

* PLAYBOOK_FILE - Full path of the playbook inside the container

* HOME_ROOT - Full path on host to mount as the root home directory inside the container (for .ssh/, etc..)

* ANSIBLE_CONFIG - Full path for the ansible configuration file to use inside the container

* INVENTORY_FILE - Full path for the inventory to use from the host

* INVENTORY_DIR - Full path for the inventory directory to use (e.g. for use with a hybrid dynamic/static inventory)
