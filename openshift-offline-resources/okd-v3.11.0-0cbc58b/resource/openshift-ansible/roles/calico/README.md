# Calico

Configure Calico components for the Master host.

## Requirements

* Ansible 2.2

## Installation

To install, set the following inventory configuration parameters:

* `openshift_use_calico=True`
* `openshift_use_openshift_sdn=False`
* `os_sdn_network_plugin_name='cni'`

By default, Calico will share the etcd used by OpenShift.
To configure Calico to use a separate instance of etcd, place etcd SSL client certs on your master,
then set the following variables in your inventory.ini:

* `calico_etcd_ca_cert_file=/path/to/etcd-ca.crt`
* `calico_etcd_cert_file=/path/to/etcd-client.crt`
* `calico_etcd_key_file=/path/to/etcd-client.key`
* `calico_etcd_endpoints=https://etcd:2379`

## Upgrading

OpenShift-Ansible installs Calico as a self-hosted install. Previously, Calico ran as a systemd service. Running Calico
in this manner is now deprecated, and must be upgraded to a hosted cluster. Please run the Legacy Upgrade playbook to
upgrade your existing Calico deployment to a hosted deployment:

        ansible-playbook -i inventory.ini playbooks/byo/calico/legacy_upgrade.yml

## Additional Calico/Node and Felix Configuration Options

Additional parameters that can be defined in the inventory are:


| Environment | Description | Schema | Default |   
|---------|----------------------|---------|---------|
| CALICO_IPV4POOL_IPIP | IPIP Mode to use for the IPv4 POOL created at start up.	| off, always, cross-subnet	| always |
| CALICO_LOG_DIR | Directory on the host machine where Calico Logs are written.| String	| /var/log/calico |

### Contact Information

Author: Dan Osborne <dan@projectcalico.org>

For support, join the `#openshift` channel on the [calico users slack](calicousers.slack.com).
