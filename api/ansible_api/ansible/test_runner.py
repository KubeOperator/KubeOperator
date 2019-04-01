# -*- coding: utf-8 -*-
#

import unittest
import sys
import warnings
warnings.simplefilter("ignore", ResourceWarning)

sys.path.insert(0, '../..')
from ansible_ui.ansible.runner import AdHocRunner
from ansible_ui.ansible.inventory import BaseInventory


class TestAdHocRunner(unittest.TestCase):
    def setUp(self):
        data = {
            "hosts": [
                {
                    "hostname": "192.168.244.163",
                    "vars": {
                        "ansible_ssh_user": "root",
                        "ansible_ssh_pass": "redhat123"
                    }
                },
                {
                    "hostname": "centos",
                    "vars": {
                        "ansible_ssh_user": "web",
                        "ansible_ssh_pass": "gaga"
                    }
                }
            ]
        }

        inventory = BaseInventory(data)
        self.runner = AdHocRunner(inventory)

    def test_run(self):
        tasks = [
            {"action": {"module": "shell", "args": "ls"}, "name": "run_cmd"},
            {"action": {"module": "shell", "args": "whoami"}, "name": "run_whoami"},
        ]
        ret = self.runner.run(tasks, "all")
        print(ret.get("summary"))


if __name__ == '__main__':
    unittest.main(warnings='ignore')
