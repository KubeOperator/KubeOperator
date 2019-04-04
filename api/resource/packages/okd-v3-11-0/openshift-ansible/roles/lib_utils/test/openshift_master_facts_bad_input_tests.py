import copy
import os
import sys

from ansible.errors import AnsibleError
import pytest

sys.path.insert(1, os.path.join(os.path.dirname(__file__), os.pardir, "lookup_plugins"))

from openshift_master_facts_default_predicates import LookupModule  # noqa: E402


class TestOpenShiftMasterFactsBadInput(object):
    lookup = LookupModule()
    default_facts = {
        'openshift': {
            'common': {}
        }
    }

    def test_missing_openshift_facts(self):
        with pytest.raises(AnsibleError):
            facts = {}
            self.lookup.run(None, variables=facts)

    def test_missing_deployment_type(self):
        with pytest.raises(AnsibleError):
            facts = copy.deepcopy(self.default_facts)
            facts['openshift']['common']['short_version'] = '10.10'
            self.lookup.run(None, variables=facts)

    def test_missing_short_version_and_missing_openshift_release(self):
        with pytest.raises(AnsibleError):
            facts = copy.deepcopy(self.default_facts)
            facts['openshift']['common']['deployment_type'] = 'origin'
            self.lookup.run(None, variables=facts)

    def test_unknown_deployment_types(self):
        with pytest.raises(AnsibleError):
            facts = copy.deepcopy(self.default_facts)
            facts['openshift']['common']['short_version'] = '1.1'
            facts['openshift']['common']['deployment_type'] = 'bogus'
            self.lookup.run(None, variables=facts)

    def test_unknown_origin_version(self):
        with pytest.raises(AnsibleError):
            facts = copy.deepcopy(self.default_facts)
            facts['openshift']['common']['short_version'] = '0.1'
            facts['openshift']['common']['deployment_type'] = 'origin'
            self.lookup.run(None, variables=facts)

    def test_unknown_ocp_version(self):
        with pytest.raises(AnsibleError):
            facts = copy.deepcopy(self.default_facts)
            facts['openshift']['common']['short_version'] = '0.1'
            facts['openshift']['common']['deployment_type'] = 'openshift-enterprise'
            self.lookup.run(None, variables=facts)
