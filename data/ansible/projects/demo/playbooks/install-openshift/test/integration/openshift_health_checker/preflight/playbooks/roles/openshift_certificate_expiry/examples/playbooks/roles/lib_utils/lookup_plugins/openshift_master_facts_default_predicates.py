# pylint: disable=missing-docstring

from ansible.errors import AnsibleError
from ansible.plugins.lookup import LookupBase


class LookupModule(LookupBase):
    # pylint: disable=too-many-branches,too-many-statements,too-many-arguments

    def run(self, terms, variables=None, regions_enabled=True, short_version=None,
            **kwargs):

        predicates = []

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

        # Predicates ordered according to OpenShift Origin source:
        # origin/vendor/k8s.io/kubernetes/plugin/pkg/scheduler/algorithmprovider/defaults/defaults.go

        if short_version in ['3.6']:
            predicates.extend([
                {'name': 'NoVolumeZoneConflict'},
                {'name': 'MaxEBSVolumeCount'},
                {'name': 'MaxGCEPDVolumeCount'},
                {'name': 'MatchInterPodAffinity'},
                {'name': 'NoDiskConflict'},
                {'name': 'GeneralPredicates'},
                {'name': 'PodToleratesNodeTaints'},
                {'name': 'CheckNodeMemoryPressure'},
                {'name': 'CheckNodeDiskPressure'},
            ])

        if short_version in ['3.7', '3.8']:
            predicates.extend([
                {'name': 'NoVolumeZoneConflict'},
                {'name': 'MaxEBSVolumeCount'},
                {'name': 'MaxGCEPDVolumeCount'},
                {'name': 'MaxAzureDiskVolumeCount'},
                {'name': 'MatchInterPodAffinity'},
                {'name': 'NoDiskConflict'},
                {'name': 'GeneralPredicates'},
                {'name': 'PodToleratesNodeTaints'},
                {'name': 'CheckNodeMemoryPressure'},
                {'name': 'CheckNodeDiskPressure'},
                {'name': 'NoVolumeNodeConflict'},
            ])

        if short_version in ['3.9', '3.10']:
            predicates.extend([
                {'name': 'NoVolumeZoneConflict'},
                {'name': 'MaxEBSVolumeCount'},
                {'name': 'MaxGCEPDVolumeCount'},
                {'name': 'MaxAzureDiskVolumeCount'},
                {'name': 'MatchInterPodAffinity'},
                {'name': 'NoDiskConflict'},
                {'name': 'GeneralPredicates'},
                {'name': 'PodToleratesNodeTaints'},
                {'name': 'CheckNodeMemoryPressure'},
                {'name': 'CheckNodeDiskPressure'},
                {'name': 'CheckVolumeBinding'},
            ])

        if regions_enabled:
            region_predicate = {
                'name': 'Region',
                'argument': {
                    'serviceAffinity': {
                        'labels': ['region']
                    }
                }
            }
            predicates.append(region_predicate)

        return predicates
