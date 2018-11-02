# pylint: disable=missing-docstring

from ansible.errors import AnsibleError
from ansible.plugins.lookup import LookupBase


class LookupModule(LookupBase):
    # pylint: disable=too-many-branches,too-many-statements,too-many-arguments

    def run(self, terms, variables=None, zones_enabled=True, short_version=None,
            **kwargs):

        priorities = []

        if short_version is None:
            if 'openshift_release' in variables:
                release = variables['openshift_release']
                if release.startswith('v'):
                    short_version = release[1:]
                else:
                    short_version = release
                short_version = '.'.join(short_version.split('.')[0:2])
            elif 'openshift_version' in variables:
                version = variables['openshift_version']
                short_version = '.'.join(version.split('.')[0:2])
            else:
                # pylint: disable=line-too-long
                raise AnsibleError("Either OpenShift needs to be installed or openshift_release needs to be specified")

        if short_version not in ['3.6', '3.7', '3.8', '3.9', '3.10', 'latest']:
            raise AnsibleError("Unknown short_version %s" % short_version)

        if short_version == 'latest':
            short_version = '3.10'

        if short_version in ['3.6', '3.7', '3.8', '3.9', '3.10']:
            priorities.extend([
                {'name': 'SelectorSpreadPriority', 'weight': 1},
                {'name': 'InterPodAffinityPriority', 'weight': 1},
                {'name': 'LeastRequestedPriority', 'weight': 1},
                {'name': 'BalancedResourceAllocation', 'weight': 1},
                {'name': 'NodePreferAvoidPodsPriority', 'weight': 10000},
                {'name': 'NodeAffinityPriority', 'weight': 1},
                {'name': 'TaintTolerationPriority', 'weight': 1}
            ])

        if zones_enabled:
            zone_priority = {
                'name': 'Zone',
                'argument': {
                    'serviceAntiAffinity': {
                        'label': 'zone'
                    }
                },
                'weight': 2
            }
            priorities.append(zone_priority)

        return priorities
