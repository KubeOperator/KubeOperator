# OpenStack Provisioning

This directory contains [Ansible][ansible] playbooks and roles to create
OpenStack resources (servers, networking, volumes, security groups,
etc.). The result is an environment ready for OpenShift installation
via [openshift-ansible].

We provide everything necessary to be able to install OpenShift on
OpenStack. In addition we work on providing integration with the
OpenStack-native services (storage, lbaas, baremetal as a service, dns, etc.).


## Requirements

In order to run these Ansible playbooks, you'll need an Ansible host and an
OpenStack environment.

### Ansible Host

Start by choosing a host from which you'll run [Ansible][ansible]. This can
be the computer you read this guide on or an OpenStack VM you'll create
specifically for this purpose.

The required dependencies for the Ansible host are:

* [Ansible](https://pypi.python.org/pypi/ansible) version >=2.4.1
* [jinja2](http://jinja.pocoo.org/docs/2.9/) version >= 2.10
* [shade](https://pypi.python.org/pypi/shade) version >= 1.26
* python-jmespath / [jmespath](https://pypi.python.org/pypi/jmespath)
* python-dns / [dnspython](https://pypi.python.org/pypi/dnspython)
* Become (sudo) is *not* required.

Optional dependencies include:

* `python-openstackclient`
* `python-heatclient`

There are a few OS-specific instructions:

* RHEL: The `rhel-7-server-openstack-10-rpms` repository is required in order to install these openstack clients.
* CentOS: Run `yum install -y centos-release-openstack-pike`

Once the dependencies are installed, clone the [openshift-ansible][openshift-ansible]
repository:

```
$ git clone https://github.com/openshift/openshift-ansible
```

### OpenStack Environment

Before you start the installation, you'll need an OpenStack environment.
Options include:

* [Devstack][devstack]
* [Packstack][packstack]
* [TripleO][tripleo] (Overcloud)

You can also use a public cloud or an OpenStack within your organization.

The OpenStack environment must satisfy these requirements:

* It must be Newton (equivalent to Red hat OpenStack 10) or newer
* Heat (Orchestration) must be available
* The deployment image (CentOS 7.4 or RHEL 7) must be loaded
* The deployment flavor must be available to your user
  - `m1.medium` / 4GB RAM + 40GB disk should be enough for testing
  - look at
    the [Minimum Hardware Requirements page][hardware-requirements]
    for production
* The keypair for SSH must be available in OpenStack
* You must have a`keystonerc` file that lets you talk to the OpenStack services
* In order to install an OpenShift cluster and deploy services, we recommend
that you have a minimum of 30 security groups, 200 security group rules, and 200
ports available in your quota.

It is also strongly recommended that you configure an external Neutron network
with a floating IP address pool.


## Configuration

Configuration is done through an Ansible inventory directory. You can switch
between multiple inventories to test multiple configurations.

Start by copying the sample inventory to your inventory directory.

```
$ cp -r openshift-ansible/playbooks/openstack/sample-inventory/ inventory
```

The sample inventory contains defaults that will do the following:

* create VMs for an OpenShift cluster with 1 Master node, 1 Infra node, and 2 App nodes
* create a new Neutron network and assign floating IP addresses to the VMs

You may have to perform further configuration in order to match the inventory
to your environment.

### OpenStack Configuration

The OpenStack configuration file is `inventory/group_vars/all.yml`.

Open the file and plug in the image, flavor and network configuration
corresponding to your OpenStack installation.

```bash
$ vi inventory/group_vars/all.yml
```

* `openshift_openstack_keypair_name` Set your OpenStack keypair name.
   - See `openstack keypair list` to find the keypairs registered with
   OpenShift.
   - This must correspond to your private SSH key in `~/.ssh/id_rsa`
* `openshift_openstack_external_network_name` Set the floating IP
   network of your OpenStack.
   - See `openstack network list` for the list of networks.
   - Often called `public`, `external` or `ext-net`.
* `openshift_openstack_default_image_name` Set the image you want your
   OpenShift VMs to run.
   - See `openstack image list` for the list of available images.
* `openshift_openstack_default_flavor` Set the flavor you want your
   OpenShift VMs to use.
   - See `openstack flavor list` for the list of available flavors.


### OpenShift Configuration

The OpenShift configuration file is `inventory/group_vars/OSEv3.yml`.

The default options will mostly work, but openshift-ansible's hardware check
may fail unless you specified a large flavor suitable for a production-ready
environment.

You can disable those checks by adding this line to `inventory/group_vars/OSEv3.yml`:

```yaml
openshift_disable_check: disk_availability,memory_availability,docker_storage
```

**Important**: The default authentication method will allow **any username
and password** in! If you're running this in a public place, you need
to set up access control by [configuring authentication][configure-authentication].


### Advanced Configuration

The [Configuration page][configuration] details several
additional options. These include:

* Set Up Authentication (TODO)
* [Multiple Masters with a load balancer][loadbalancer]
* [External DNS][external-dns]
* [Kuryr SDN][kuryr-sdn]
* Multiple Clusters (TODO)
* [Cinder Registry][cinder-registry]

Read the [Configuration page][configuration] for a full listing of
configuration options.


## Installation

Before running the installation playbook, you may want to create an `ansible.cfg`
file with useful defaults:

```bash
$ cp openshift-ansible/ansible.cfg ansible.cfg
```

We recommend adding an additional option:

```cfg
any_errors_fatal = true
```

This will abort the Ansible playbook execution as soon as any error is
encountered.

If you want, you can [Build the OpenShift node images at this
point][build-images].

Now, run the provision + install playbook. This will create OpenStack resources
and deploy an OpenShift cluster on top of them:

```bash
$ ansible-playbook --user openshift \
  -i openshift-ansible/playbooks/openstack/inventory.py \
  -i inventory \
  openshift-ansible/playbooks/openstack/openshift-cluster/provision_install.yml
```

* If you're using multiple inventories, make sure you pass the path to
the right one to `-i`.
* If your SSH private key is not in `~/.ssh/id_rsa`, use the `--private-key`
option to specify the correct path.
* Note that we must pass in the [dynamic inventory][dynamic] --
`openshift-ansible/playbooks/openstack/inventory.py`. This is a script that
looks for OpenStack resources and enables Ansible to reference them.


## Post-Install

Once installation completes, a few additional steps may be required or useful.

* [Configure DNS][configure-dns]
* [Get the `oc` Client][get-the-oc-client]
* [Log in Using the Command Line][log-in-using-the-command-line]
* [Access the UI][access-the-ui]

Read the [Post-Install page][post-install] for a full list of options.


## Uninstall

The installation process not only creates a Heat stack, but can also
perform actions such as writing DNS records or subscribing a host to RHN.
In order to do a clean uninstall, run this command:

```bash
$ ansible-playbook --user openshift \
  -i openshift-ansible/playbooks/openstack/inventory.py \
  -i inventory \
  openshift-ansible/playbooks/openstack/openshift-cluster/uninstall.yml
```

[ansible]: https://www.ansible.com/
[openshift-ansible]: https://github.com/openshift/openshift-ansible
[openshift-ansible-setup]: https://github.com/openshift/openshift-ansible#setup
[devstack]: https://docs.openstack.org/devstack/
[tripleo]: http://tripleo.org/
[packstack]: https://www.rdoproject.org/install/packstack/
[configure-authentication]: https://docs.okd.io/latest/install_config/configuring_authentication.html
[hardware-requirements]: https://docs.okd.io/latest/install_config/install/prerequisites.html#hardware
[origin]: https://www.openshift.org/
[centos7]: https://www.centos.org/
[sample-openshift-inventory]: https://github.com/openshift/openshift-ansible/blob/master/inventory/hosts.example
[configuration]: ./configuration.md
[loadbalancer]: ./configuration.md#multi-master-configuration
[external-dns]: ./configuration.md#dns-configuration
[cinder-registry]: ./configuration.md#cinder-backed-registry-configuration
[post-install]: ./post-install.md
[configure-dns]: ./post-install.md#configure-dns
[get-the-oc-client]: ./post-install.md#get-the-oc-client
[log-in-using-the-command-line]: ./post-install.md#log-in-using-the-command-line
[access-the-ui]: ./post-install.md#access-the-ui
[dynamic]: http://docs.ansible.com/ansible/latest/intro_dynamic_inventory.html
[kuryr-sdn]: ./configuration.md#kuryr-networking-configuration
[build-images]: ./configuration.md#building-node-images
