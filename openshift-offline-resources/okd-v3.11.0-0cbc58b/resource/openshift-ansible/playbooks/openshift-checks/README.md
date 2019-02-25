# OpenShift health checks

This directory contains Ansible playbooks for detecting potential problems prior
to an install, as well as health checks to run on existing OpenShift clusters.

Ansible's default operation mode is to fail fast, on the first error. However,
when performing checks, it is useful to gather as much information about
problems as possible in a single run.

Thus, the playbooks run a battery of checks against the inventory hosts and
gather intermediate errors, giving a more complete diagnostic of the state of
each host. If any check failed, the playbook run will be marked as failed.

To facilitate understanding the problems that were encountered, a custom
callback plugin summarizes execution errors at the end of a playbook run.

## Available playbooks

1. Pre-install playbook ([pre-install.yml](pre-install.yml)) - verifies system
   requirements and look for common problems that can prevent a successful
   installation of a production cluster.

2. Diagnostic playbook ([health.yml](health.yml)) - check an existing cluster
   for known signs of problems.

3. Certificate expiry playbooks ([certificate_expiry](certificate_expiry)) -
   check that certificates in use are valid and not expiring soon.

4. Adhoc playbook ([adhoc.yml](adhoc.yml)) - use it to run adhoc checks or to
   list existing checks.
   See the [next section](#the-adhoc-playbook) for a usage example.

## Running

With a [recent installation of Ansible](../../../README.md#setup), run the playbook
against your inventory file. Here is the step-by-step:

1. If you haven't done it yet, clone this repository:

    ```console
    $ git clone https://github.com/openshift/openshift-ansible
    $ cd openshift-ansible
    ```

2. Install the [dependencies](../../../README.md#setup)

3. Run the appropriate playbook:

    ```console
    $ ansible-playbook -i <inventory file> playbooks/openshift-checks/pre-install.yml
    ```

    or

    ```console
    $ ansible-playbook -i <inventory file> playbooks/openshift-checks/health.yml
    ```

    or

    ```console
    $ ansible-playbook -i <inventory file> playbooks/openshift-checks/certificate_expiry/default.yaml -v
    ```

### The adhoc playbook

The adhoc playbook gives flexibility to run any check or a custom group of
checks. What will be run is determined by the `openshift_checks` variable,
which, among other ways supported by Ansible, can be set on the command line
using the `-e` flag.

For example, to run the `docker_storage` check:

```console
$ ansible-playbook -i <inventory file> playbooks/openshift-checks/adhoc.yml -e openshift_checks=docker_storage
```

To run more checks, use a comma-separated list of check names:

```console
$ ansible-playbook -i <inventory file> playbooks/openshift-checks/adhoc.yml -e openshift_checks=docker_storage,disk_availability
```

To run an entire class of checks, use the name of a check group tag, prefixed by `@`. This will run all checks tagged `preflight`:

```console
$ ansible-playbook -i <inventory file> playbooks/openshift-checks/adhoc.yml -e openshift_checks=@preflight
```

It is valid to specify multiple check tags and individual check names together
in a comma-separated list.

To list all of the available checks and tags, run the adhoc playbook without
setting the `openshift_checks` variable:

```console
$ ansible-playbook -i <inventory file> playbooks/openshift-checks/adhoc.yml
```

## Running in a container

This repository is built into a Docker image including Ansible so that it can
be run anywhere Docker is available, without the need to manually install dependencies.
Instructions for doing so may be found [in the README](../../../README_CONTAINER_IMAGE.md).
