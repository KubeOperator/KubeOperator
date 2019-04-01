# oVirt Playbooks
## Provisioning
This subdirectory contains the Ansible playbooks used to deploy
an OpenShift Container Platform environment on oVirt.

### Where do I start?
Choose a host from which Ansible plays will be executed. This host must have
the ability to access the web interface of the oVirt cluster engine and the
network on which the OpenShift nodes will be installed. We will refer to
this host as the *bastion*.

#### oVirt Ansible Roles
The oVirt project maintains Ansible roles for managing an oVirt cluster.
These should be installed on the *bastion* host according to the instructions
at the [oVirt Ansible Roles page](https://github.com/ovirt/ovirt-ansible/).

#### DNS Server
An external DNS server is required to provide name resolution to nodes and
applications. See the
[OpenShift Installation Documentation](https://docs.openshift.com/container-platform/latest/install_config/install/prerequisites.html#prereq-dns)
for details.

### Let's Provision!
#### High-level overview
After populating inventory and variables files with the proper values,
(see [The OpenShift Advanced Installation Documentation](https://docs.openshift.com/container-platform/latest/install_config/install/advanced_install.html)
) a series of Ansible playbooks from this subdirectory will provision a set of
nodes on the oVirt cluster, prepare them for OpenShift installation,
and deploy an OpenShift cluster on them.

#### Step 1 Inventory
The [`inventory.example`](inventory.example) file here is provided as an example of a three master, three inventory
environment. It is up to the user to add additional OpenShift specific variables to this file to configure
required elements such as the registry, storage, authentication, and networking.

One required variable added for this environment is the `openshift_ovirt_dns_zone`. As this is used to construct
hostnames during VM creation, it is essential that this be set to the default dns zone for those nodes' hostnames.

#### Step 2 oVirt Provisioning Variables

Fill out a provisioning variables file (example [`provisioning-vars.yaml.example`](provisioning-vars.yaml.example)
with values from your oVirt environment, making sure to fill in all commented values.

*oVirt Engine internal Certificate*

A copy of the `/etc/pki/ovirt-engine/ca.pem` from the oVirt engine will need to
be downloaded to the *bastion* and its location set in the `engine_cafile` variable. Replace the
example server in the following command to download the certificate:

```
$ curl --output ca.pem 'http://engine.example.com/ovirt-engine/services/pki-resource?resource=ca-certificate&format=X509-PEM-CA'

```

#### Step 3 Provision Virtual Machines in oVirt
Once all the variables in the `provisioning_vars.yaml` file are set, use the
[`ovirt-vm-infra.yml`](openshift-cluster/ovirt-vm-infra.yml) playbook to begin
provisioning.

```
ansible-playbook -i inventory -e@provisioning_vars.yml ${PATH_TO_OPENSHIFT_ANSIBLE}/playbooks/ovirt/openshift-cluster/ovirt-vm-infra.yml
```

#### Step 4 Update DNS

At this stage, ensure DNS is set up properly for the following access:

* Nodes are available to each other by their hostnames.
* The nodes running router services (typically the infrastructure nodes) are reachable by the wildcard entry.
* The load balancer node is reachable as the openshift-master host entry for console access.

#### Step 5 Install Prerequisite Services
```
ansible-playbook -i inventory ${PATH_TO_OPENSHIFT_ANSIBLE}/playbooks/prerequisites.yml
```

#### Step 6 Deploy OpenShift
```
ansible-playbook -i inventory ${PATH_TO_OPENSHIFT_ANSIBLE}/playbooks/deploy_cluster.yml
```

### Ready To Work!

## Uninstall / Deprovisioning
In case of a failed installation due to a missing variable, it is occasionally necessary to start from a fresh set of virtual machines. Uninstalling the virtual machines and reprovisioning them may be perfomed by running the [`openshift-cluster/unregister-vms.yaml`](openshift-cluster/unregister-vms.yaml) playbook (to recover RHSM entitlements) followed by the [`openshift-cluster/ovirt-vm-uninstall.yaml`](openshift-cluster/ovirt-vm-uninstall.yaml) playbook.
