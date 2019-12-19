# -*- coding: utf-8 -*-
#

import unittest
import sys
import warnings
warnings.simplefilter("ignore", ResourceWarning)

sys.path.insert(0, '../..')
from ansible_api.ansible.runner import AdHocRunner
from ansible_api.ansible.inventory import BaseInventory


class TestAdHocRunner(unittest.TestCase):
    def setUp(self):
        data = {
            "hosts": [
                {
                    "hostname": "192.168.0.77",
                    "vars": {
                        "ansible_ssh_user": "root",
                        "ansible_ssh_pass": "KubeOperator@2019"
                    }
                },
                {
                    "hostname": "192.168.0.75",
                    "vars": {
                        "ansible_ssh_user": "root",
                        "ansible_ssh_pass": "KubeOperator@2019"
                    }
                },
                {
                    "hostname": "172.190.92.62",
                    "vars": {
                        "ansible_ssh_user": "root",
                        "ansible_ssh_pass": "123"
                    }
                }
            ]
        }

        inventory = BaseInventory(data)
        print("inventory")
        print(inventory.__dict__)
        self.runner = AdHocRunner(inventory)

    def test_run(self):
        tasks = [
            {'action': {'module': 'setup', 'args': ''}},
            # {"action": {"module": "shell", "args": "ls"}, "name": "run_cmd"},
            # {"action": {"module": "shell", "args": "whoami"}, "name": "run_whoami"},
        ]
        pattern = "172.190.92.62"
        ret = self.runner.run(tasks, pattern)
        print("result: ")
        print(ret["raw"]["ok"][pattern]["setup"]["ansible_facts"])
        print(ret.get("summary"))


if __name__ == '__main__':
    unittest.main(warnings='ignore')
