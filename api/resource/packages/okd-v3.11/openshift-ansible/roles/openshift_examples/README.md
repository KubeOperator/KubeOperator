OpenShift Examples
================

Installs example image streams, db-templates, and quickstart-templates by copying
examples from this module to your first master and importing them with oc create -n into the openshift namespace

The examples-sync.sh script can be used to pull the latest content from github
and stage it for updating the ansible repo. This script is not used directly by
ansible.

Requirements
------------

Facts
-----

| Name                       | Default Value | Description                            |
-----------------------------|---------------|----------------------------------------|
| openshift_install_examples | true          | Runs the role with the below variables |

Role Variables
--------------

| Name                                | Default value                                                  |                                          |
|-------------------------------------|----------------------------------------------------------------|------------------------------------------|
| openshift_examples_load_centos      | true when openshift_deployment_type not 'openshift-enterprise' | Load centos image streams                |
| openshift_examples_load_rhel        | true if openshift_deployment_type is 'openshift-enterprise'    | Load rhel image streams                  |
| openshift_examples_load_db_templates| true                                                           | Loads database templates                 |
| openshift_examples_load_quickstarts | true                                                           | Loads quickstarts ie: nodejs, rails, etc |
| openshift_examples_load_xpaas       | false                                                          | Loads xpass streams and templates        |


Dependencies
------------

Example Playbook
----------------

TODO
----
Currently we use `oc create -f` against various files and we accept non zero return code as a success
if (and only if) stderr also contains the string 'already exists'. This means that if one object in the file exists already
but others fail to create you won't be aware of the failure. This also means that we do not currently support
updating existing objects.

We should add the ability to compare existing image streams against those we're being asked to load and update if necessary.

License
-------

Apache License, Version 2.0

Author Information
------------------

Scott Dodson (sdodson@redhat.com)
