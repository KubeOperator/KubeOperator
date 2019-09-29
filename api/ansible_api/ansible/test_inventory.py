# -*- coding: utf-8 -*-
#
#
import sys
import unittest
import warnings

sys.path.insert(0, '../..')
from ansible_api.ansible.inventory import BaseInventory

warnings.simplefilter("ignore", ResourceWarning)


class TestJMSInventory(unittest.TestCase):
    def setUp(self):
        host_list = [{
            "hostname": "testserver1",
            "vars": {
                "ansible_ssh_host": "102.1.1.1",
                "ansible_ssh_port": 22,
                "ansible_ssh_user": "root",
                "ansible_ssh_pass": "password",
                "service": "mysql"
            },
        }, {
            "hostname": "testserver2",
            "vars": {
                "ansible_ssh_host": "102.1.1.2",
                "ansible_ssh_port": 22,
                "ansible_ssh_user": "root",
                "ansible_ssh_pass": "password",
                "service": "web"
            }
        },
        ]

        group_list = [{
            "name": "group1",
            "hosts": ["testserver1",],
            "vars": {"service": "app"}
        }, {
            "name": "group2",
            "hosts": ["testserver2",],
            "vars": {"service": "gaga"}
        }, {
            "name": "group3",
            "children": ["group1", "group2"]
        }
        ]
        data = {"hosts": host_list, "groups": group_list}
        self.inventory = BaseInventory(data=data)

    def test_hosts(self):
        print("#"*10 + "Hosts" + "#"*10)
        for host in self.inventory.hosts:
            print(host)

    def test_groups(self):
        print("#" * 10 + "Groups" + "#" * 10)
        for group in self.inventory.groups:
            print(group)

    def test_group_all(self):
        print("#" * 10 + "all group hosts" + "#" * 10)
        group = self.inventory.get_group('all')
        print(group.hosts)


if __name__ == '__main__':
    unittest.main()
