#!/usr/bin/env python
"""
This is an Ansible dynamic inventory for OpenStack.

It requires your OpenStack credentials to be set in clouds.yaml or your shell
environment.

"""

import resources


def build_inventory():
    """Build the Ansible inventory for the current environment."""
    inventory = resources.build_inventory()
    inventory['nodes'] = inventory['openstack_nodes']
    inventory['masters'] = inventory['openstack_master_nodes']
    inventory['etcd'] = inventory['openstack_etcd_nodes']
    inventory['glusterfs'] = inventory['openstack_cns_nodes']
    return inventory


if __name__ == '__main__':
    resources.main(build_inventory)
