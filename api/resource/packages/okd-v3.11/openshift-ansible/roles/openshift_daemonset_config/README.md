OpenShift Daemonset Config
================================

This role creates a configmap, a secret, and then deploys a daemonset that uses these objects to apply configuration to hosts.

Requirements
------------

* Ansible 2.4
* One or more Master servers

Role Variables
--------------
There are many variables that are defined in this role that allow flexibility when deploying this daemonset.  The most important variables are the following:

This variable represents the main config container that will execute.  The default can be found in `defaults/main.yml`.
`openshift_daemonset_config_image: "centos:7"`

The container will start and perform a loop:
```
          while true; do

            # execute user defined script
            sh /opt/config/{{ openshift_daemonset_config_script }}

            # sleep for ${RESYNC_INTERVAL} minutes, then loop. if we fail Kubelet will restart us again
            echo "Success, sleeping for ${RESYNC_INTERVAL}s. Date: $(date)"
            sleep ${RESYNC_INTERVAL}

          # Return to perform the config
          done
```
As shown above, the config container will begin a configuration loop.  This loop will perform any actions defined in the 
`openshift_daemonset_config_script` variable.  This variable represents the script that will be called on the container's start up
`openshift_daemonset_config_script: config.sh`

Once this script has completed, the loop will enter a sleep state of `openshift_daemonset_config_interval`.  This defines the amount
of time between configuration that will occur.  The defaults can be found inside of the `defaults/main.yml`.

The next important set of variables is how the configuration files are supplied to the config container.  The config container will
receive a configmap and a secret defined by these variables:
- `openshift_daemonset_config_configmap_name`
- `openshift_daemonset_config_secret_name`

When the config container starts the configmap and secrets are mounted at `/opt/config` and `/opt/secrets` respectively.

The configuration files or secrets are then referenced at these mount points when the configuration scripts are running.  This allows the administrator to write configuration to the host and store the configuration management inside of Openshift.

The following variables are the interface to the role when creating the configuration.

This option allows an administrator to copy configuration files to disk and include them in the configmap. Recommended for large files or data.
```
openshift_daemonset_config_configmap_files: {}
```

This option will allow an administrator to pass in string content. The role will write this data to a file and pass in the filename so that it will be included inside of the configmap.  This is useful when contents are too large to pass to Openshift on the command line. Recommended for large files or data.
```
openshift_daemonset_config_configmap_contents_to_files: []
```

This option allows string contents for the configmap items and will be placed directly into the configmap. This does have a size limitation and is recommended for smaller string content.
```
openshift_daemonset_config_configmap_literals: {}
```

This option will place content passed into a secret.
```
openshift_daemonset_config_secrets: {}
```

Files created in this role are cleaned up after the configmap is created in attempt to ensure that artifacts are cleaned up.

See the example playbook below for examples of these variables and their usage.

Please see `defaults/main.yml` for an exhaustive list of variables.

Dependencies
------------

lib_openshift


Example Playbook
----------------

```
- import_role:
    name: openshift_daemonset_config
  vars:
    openshift_daemonset_config_daemonset_name: test-config

    openshift_daemonset_config_secrets:
    - path: api_credentials
      data: |
        user: abcdef
        password: 123456

    openshift_daemonset_config_configmap_contents_to_files:
    - path: /tmp/authorized_keys
      name: authorized_keys
      contents: |
        ssh-rsa AAAAB3... user1@domain
        ssh-rsa AAAAB3... user1@domain

    # config script to run
    openshift_daemonset_config_script: myconfig.sh

    # config files to include in config map
    openshift_daemonset_config_configmap_files:
      bashrc: /local/path/to/.bashrc

    # config values as strings
    openshift_daemonset_config_configmap_literals:
      # some config data
      config_data: 42

      # configuration script that will be called
      myconfig.sh: |
        #!/bin/bash

        # example from configmap content to files.
        # lay down the authorized_keys
        cp /opt/config/authorized_keys /host/root/.ssh/authorized_keys

        # example from configmap files
        # lay down the .bashrc
        cp /opt/config/bashrc /host/root/.bashrc

        # example from secrets
        # lay down the credentials for other files/scripts to use
        cp /opt/secrets/api_credentials /host/root/.creds

        # example from configmap literals
        # echo the answer to life
        echo "$(cat /opt/config/config_data)"

    # place this daemonset on the following nodes based on their node labels
    openshift_daemonset_config_node_selector:
      kubernetes.io/hostname=clusterhost123.xyz


Notes
-----

TODO

License
-------

Apache License, Version 2.0

Author Information
------------------

Openshift
