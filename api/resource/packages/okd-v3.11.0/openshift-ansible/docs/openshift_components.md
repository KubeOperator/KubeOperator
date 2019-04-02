# OpenShift-Ansible Components

>**TL;DR: Look at playbooks/openshift-web-console as an example**

## General Guidelines

Components in OpenShift-Ansible consist of two main parts:
* Entry point playbook(s)
* Ansible role
* OWNERS files in both the playbooks and roles associated with the component

When writing playbooks and roles, follow these basic guidelines to ensure
success and maintainability. 

### Idempotency

Definition:

>_an idempotent operation is one that has no additional effect if it is called
more than once with the same input parameters_

Ansible playbooks and roles should be written such that when the playbook is run
again with the same configuration, no tasks should report `changed` as well as
no material changes should be made to hosts in the inventory.  Playbooks should
be re-runnable, but also be idempotent.

### Other advice for success

* Try not to leave artifacts like files or directories
* Avoid using `failed_when:` where ever possible
* Always `name:` your tasks
* Document complex logic or code in tasks
* Set role defaults in `defaults/main.yml`
* Avoid the use of `set_fact:`

## Building Component Playbooks

Component playbooks are divided between the root of the component directory and
the `private` directory.  This allows other parts of openshift-ansible to import
component playbooks without also running the common initialization playbooks
unnecessarily.

Entry point playbooks are located in the `playbooks` directory and follow the
following structure:

```
playbooks/openshift-component_name
├── config.yml                          Entry point playbook
├── private
│   ├── config.yml                      Included by the Cluster Installer
│   └── roles -> ../../roles            Don't forget to create this symlink
├── OWNERS                              Assign 2-3 approvers and reviewers
└── README.md                           Tell us what this component does
```

### Entry point config playbook

The primary component entry point playbook will at a minimum run the common
initialization playbooks and then import the private playbook.

```yaml
# playbooks/openshift-component_name/config.yml
---
- import_playbook: ../init/main.yml

- import_playbook: private/config.yml

```

### Private config playbook

The private component playbook will run the component role against the intended
host groups and provide any required variables.  This playbook is also called
during cluster installs and upgrades.  Think of this as the shareable portion of
the component playbooks.

```yaml
# playbooks/openshift-component_name/private/config.yml
---

- name: OpenShift Component_Name Installation
  hosts: oo_first_master
  tasks:
  - import_role:
      name: openshift_component_name
```

NOTE: The private playbook may also include wrapper plays for the Installer
Checkpoint plugin which will be discussed later.

## Building Component Roles

Component roles contain all of the necessary files and logic to install and
configure the component.  The install portion of the role should also support
performing upgrades on the component.

Ansible roles are located in the `roles` directory and follow the following
structure:

```
roles/openshift_component_name
├── defaults
│   └── main.yml                        Defaults for variables used in the role
│                                           which can be overridden by the user
├── files
│   ├── component-config.yml
│   ├── component-rbac-template.yml
│   └── component-template.yml
├── handlers
│   └── main.yml
├── meta
│   └── main.yml
├── OWNERS                              Assign 2-3 approvers and reviewers
├── README.md
├── tasks
│   └── main.yml                        Default playbook used when calling the role
├── templates
└── vars
    └── main.yml                        Internal roles variables
```
### Component Installation

Where possible, Ansible modules should be used to perform idempotent operations
with the OpenShift API.  Avoid using the `command` or `shell` modules with the
`oc` cli unless the required operation is not available through either the
`lib_openshift` modules or Ansible core modules.

The following is a basic flow of Ansible tasks for installation. 

- Create the project (oc_project)
- Create a temp directory for processing files
- Copy the client config to temp
- Copy templates to temp
- Read existing config map
- Copy existing config map to temp
- Generate/update config map
- Reconcile component RBAC (oc_process)
- Apply component template (oc_process)
- Poll healthz and wait for it to come up
- Log status of deployment
- Clean up temp

### Component Removal

- Remove the project (oc_project)

## Enabling the Installer Checkpoint callback

- Add the wrapper plays to the entry point playbook
- Update the installer_checkpoint callback plugin

Details can be found in the installer_checkpoint role.
