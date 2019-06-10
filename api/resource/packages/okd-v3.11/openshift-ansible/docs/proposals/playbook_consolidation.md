# OpenShift-Ansible Playbook Consolidation

## Description
The designation of `byo` is no longer applicable due to being able to deploy on
physical hardware or cloud resources using the playbooks in the `byo` directory.
Consolidation of these directories will make maintaining the code base easier
and provide a more straightforward project for users and developers.

The main points of this proposal are:
* Consolidate initialization playbooks into one set of playbooks in
  `playbooks/init`. 
* Collapse the `playbooks/byo` and `playbooks/common` into one set of
  directories at `playbooks/openshift-*`.

This consolidation effort may be more appropriate when the project moves to
using a container as the default installation method.

## Design

### Initialization Playbook Consolidation
Currently there are two separate sets of initialization playbooks:
* `playbooks/byo/openshift-cluster/initialize_groups.yml`
* `playbooks/common/openshift-cluster/std_include.yml`

Although these playbooks are located in the `openshift-cluster` directory they
are shared by all of the `openshift-*` areas.  These playbooks would be better
organized in a `playbooks/init` directory collocated with all their related
playbooks.

In the example below, the following changes have been made:
* `playbooks/byo/openshift-cluster/initialize_groups.yml` renamed to
  `playbooks/init/initialize_host_groups.yml`
* `playbooks/common/openshift-cluster/std_include.yml` renamed to
  `playbooks/init/main.yml`
* `- include: playbooks/init/initialize_host_groups.yml` has been added to the
  top of `playbooks/init/main.yml`
* All other related files for initialization have been moved to `playbooks/init`

The `initialize_host_groups.yml` playbook is only one play with one task for
importing variables for inventory group conversions.  This task could be further
consolidated with the play in `evaluate_groups.yml`.

The new standard initialization playbook would be
`playbooks/init/main.yml`.


```
 
> $ tree openshift-ansible/playbooks/init
.
├── evaluate_groups.yml
├── initialize_facts.yml
├── initialize_host_groups.yml
├── initialize_openshift_repos.yml
├── initialize_openshift_version.yml
├── main.yml
├── roles -> ../../roles
├── validate_hostnames.yml
└── vars
    └── cluster_hosts.yml
```

```yaml
# openshift-ansible/playbooks/init/main.yml
---
- include: initialize_host_groups.yml

- include: evaluate_groups.yml

- include: initialize_facts.yml

- include: validate_hostnames.yml

- include: initialize_openshift_repos.yml

- include: initialize_openshift_version.yml
```

### `byo` and `common` Playbook Consolidation
Historically, the `byo` directory coexisted with other platform directories
which contained playbooks that then called into `common` playbooks to perform
common installation steps for all platforms.  Since the other platform
directories have been removed this separation is no longer necessary.

In the example below, the following changes have been made:
* `playbooks/byo/openshift-master` renamed to
  `playbooks/openshift-master`
* `playbooks/common/openshift-master` renamed to
  `playbooks/openshift-master/private`
* Original `byo` entry point playbooks have been updated to include their
  respective playbooks from `private/`.
* Symbolic links have been updated as necessary

All user consumable playbooks are in the root of `openshift-master` and no entry
point playbooks exist in the `private` directory.  Maintaining the separation
between entry point playbooks and the private playbooks allows individual pieces
of the deployments to be used as needed by other components.

```
openshift-ansible/playbooks/openshift-master 
> $ tree
.
├── config.yml
├── private
│   ├── additional_config.yml
│   ├── config.yml
│   ├── filter_plugins -> ../../../filter_plugins
│   ├── library -> ../../../library
│   ├── lookup_plugins -> ../../../lookup_plugins
│   ├── restart_hosts.yml
│   ├── restart_services.yml
│   ├── restart.yml
│   ├── roles -> ../../../roles
│   ├── scaleup.yml
│   └── validate_restart.yml
├── restart.yml
└── scaleup.yml
```

```yaml
# openshift-ansible/playbooks/openshift-master/config.yml
---
- include: ../init/main.yml

- include: private/config.yml
```

With the consolidation of the directory structure and component installs being
removed from `openshift-cluster`, that directory is no longer necessary.  To
deploy an entire OpenShift cluster, a playbook would be created to tie together
all of the different components.  The following example shows how multiple
components would be combined to perform a complete install.

```yaml
# openshift-ansible/playbooks/deploy_cluster.yml
---
- include: init/main.yml

- include: openshift-etcd/private/config.yml

- include: openshift-nfs/private/config.yml

- include: openshift-loadbalancer/private/config.yml

- include: openshift-master/private/config.yml

- include: openshift-node/private/config.yml

- include: openshift-glusterfs/private/config.yml

- include: openshift-hosted/private/config.yml

- include: openshift-service-catalog/private/config.yml
```

## User Story
As a developer of OpenShift-Ansible,
I want simplify the playbook directory structure
so that users can easily find deployment playbooks and developers know where new
features should be developed.

## Implementation
Given the size of this refactoring effort, it should be broken into smaller
steps which can be completed independently while still maintaining a functional
project.

Steps:
1. Update and merge consolidation of the initialization playbooks.
2. Update each merge consolidation of each `openshift-*` component area
3. Update and merge consolidation of `openshift-cluster` 

## Acceptance Criteria
* Verify that all entry points playbooks install or configure as expected.
* Verify that CI is updated for testing new playbook locations.
* Verify that repo documentation is updated
* Verify that user documentation is updated

## References
