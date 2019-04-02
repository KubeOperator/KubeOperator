#!/usr/bin/env python
"""
This is an Ansible dynamic inventory for OpenStack, specifically for use with
the scaling playbooks.

It requires your OpenStack credentials to be set in clouds.yaml or your shell
environment.

"""

import resources


if __name__ == '__main__':
    resources.main(resources.build_inventory)
