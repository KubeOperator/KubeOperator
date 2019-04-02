# Repository structure

### Ansible

```
.
├── inventory           Contains dynamic inventory scripts, and examples of
│                       Ansible inventories.
├── playbooks           Contains Ansible playbooks targeting multiple use cases.
└── roles               Contains Ansible roles, units of shared behavior among
                        playbooks.
```

#### Ansible shared libraries and plugins

Shared libraries and plugins are located in the `lib_utils` role.

#### Ansible playbooks

The `playbooks` directory is organized such that entry point playbooks are
located in either component sub directories or cloud provisioning subdirectories.

_Cloud Provisioning_
- aws
- gcp
- openstack

_OpenShift Components_
- openshift-etcd
- openshift-master
- openshift-node
- openshift-<component_name>

### Scripts

```
.
└── utils               Contains the `atomic-openshift-installer` command, an
                        interactive CLI utility to install OpenShift across a
                        set of hosts.
```

### Documentation

```
.
└── docs                Contains documentation for this repository.
```

### Tests

```
.
└── test                Contains tests.
```

### CI

These files are used by [PAPR](https://github.com/projectatomic/papr),
It is very similar in workflow to Travis, with the test
environment and test scripts defined in a YAML file.

```
.
├── .papr.yml
├── .papr.sh
└── .papr.inventory
├── .papr.all-in-one.inventory
└── .papr-master-ha.inventory
```
