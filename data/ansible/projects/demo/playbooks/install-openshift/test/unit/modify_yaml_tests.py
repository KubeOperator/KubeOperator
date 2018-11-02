""" Tests for the modify_yaml Ansible module. """
# pylint: disable=missing-docstring,invalid-name

import os
import sys
import unittest

sys.path = [os.path.abspath(os.path.dirname(__file__) + "/../../library/")] + sys.path

# pylint: disable=import-error
from modify_yaml import set_key  # noqa: E402


class ModifyYamlTests(unittest.TestCase):

    def test_simple_nested_value(self):
        cfg = {"section": {"a": 1, "b": 2}}
        changes = set_key(cfg, 'section.c', 3)
        self.assertEquals(1, len(changes))
        self.assertEquals(3, cfg['section']['c'])

    # Tests a previous bug where property would land in section above where it should,
    # if the destination section did not yet exist:
    def test_nested_property_in_new_section(self):
        cfg = {
            "masterClients": {
                "externalKubernetesKubeConfig": "",
                "openshiftLoopbackKubeConfig": "openshift-master.kubeconfig",
            },
        }

        yaml_key = 'masterClients.externalKubernetesClientConnectionOverrides.acceptContentTypes'
        yaml_value = 'application/vnd.kubernetes.protobuf,application/json'
        set_key(cfg, yaml_key, yaml_value)
        self.assertEquals(yaml_value, cfg['masterClients']
                          ['externalKubernetesClientConnectionOverrides']
                          ['acceptContentTypes'])
