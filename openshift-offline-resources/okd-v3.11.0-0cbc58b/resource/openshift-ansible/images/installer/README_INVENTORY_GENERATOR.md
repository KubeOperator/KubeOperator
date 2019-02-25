Dynamic Inventory Generation
============================

Script within the openshift-ansible image that can dynamically
generate an Ansible inventory file from an existing cluster.

## Configure

User configuration helps to provide additional details when creating an inventory file.
The default location of this file is in `/etc/inventory-generator-config.yaml`. The
following configuration values are either expected or default to the given values when omitted:

- `master_config_path`:
  - specifies where to look for the bind-mounted `master-config.yaml` file in the container
  - if omitted or a `null` value is given, its value is defaulted to `/opt/app-root/src/master-config.yaml`

- `admin_kubeconfig_path`:
  - specifies where to look for the bind-mounted `admin.kubeconfig` file in the container
  - if omitted or a `null` value is given, its value is defaulted to `/opt/app-root/src/.kube/config`

- `ansible_ssh_user`:
  - specifies the ssh user to be used by Ansible when running the specified `PLAYBOOK_FILE` (see `README_CONTAINER_IMAGE.md` for additional information on this environment variable).
  - if omitted, its value is defaulted to `root`

- `ansible_become_user`:
  - specifies a user to "become" on the remote host. Used for privilege escalation.
  - If a non-null value is specified, `ansible_become` is implicitly set to `yes` in the resulting inventory file.

See the supplied sample user configuration file in [`root/etc/inventory-generator-config.yaml`](./root/etc/inventory-generator-config.yaml) for additional optional inventory variables that may be specified.

## Build

See `README_CONTAINER_IMAGE.md` for information on building this image.

## Run

Given a master node's `master-config.yaml` file, a user configuration file (see "Configure" section), and an `admin.kubeconfig` file, the command below will:

1. Use `oc` to query the host about additional node information (using the supplied `kubeconfig` file)
2. Generate an inventory file based on information retrieved from `oc get nodes` and the given `master-config.yaml` file.
3. run the specified [openshift-ansible](https://github.com/openshift/openshift-ansible) `health.yml` playbook using the generated inventory file from the previous step

```
docker run -u `id -u` \
       -v $HOME/.ssh/id_rsa:/opt/app-root/src/.ssh/id_rsa:Z,ro \
       -v /tmp/origin/master/admin.kubeconfig:/opt/app-root/src/.kube/config:Z \
       -v /tmp/origin/master/master-config.yaml:/opt/app-root/src/master-config.yaml:Z \
       -e OPTS="-v --become-user root" \
       -e PLAYBOOK_FILE=playbooks/openshift-checks/health.yml \
       -e GENERATE_INVENTORY=true \
       -e USER=`whoami` \
       docker.io/openshift/origin-ansible

```

**Note** In the command above, specifying the `GENERATE_INVENTORY` environment variable will automatically generate the inventory file in an expected location.
An `INVENTORY_FILE` variable (or any other inventory location) does not need to be supplied when generating an inventory.

## Debug

To debug the `generate` script, run the above script interactively
and manually execute `/usr/local/bin/generate`:

```
...
docker run -u `id -u` \
       -v ...
       ...
       -it docker.io/openshift/origin-ansible /bin/bash

---

bash-4.2$ cd $HOME
bash-4.2$ ls
master-config.yaml
bash-4.2$ /usr/local/bin/generate $HOME/generated_hosts
bash-4.2$ ls
generated_hosts  master-config.yaml
bash-4.2$ less generated_hosts
...
```

## Notes

See `README_CONTAINER_IMAGE.md` for additional information about this image.
