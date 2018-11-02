"""Check that available RPM packages match the required versions."""

from openshift_checks import OpenShiftCheck
from openshift_checks.mixins import NotContainerizedMixin


class PackageVersion(NotContainerizedMixin, OpenShiftCheck):
    """Check that available RPM packages match the required versions."""

    name = "package_version"
    tags = ["preflight"]

    def is_active(self):
        """Skip hosts that do not have package requirements."""
        group_names = self.get_var("group_names", default=[])
        master_or_node = 'oo_masters_to_config' in group_names or 'oo_nodes_to_config' in group_names
        return super(PackageVersion, self).is_active() and master_or_node

    def run(self):
        rpm_prefix = self.get_var("openshift_service_type")
        if self._templar is not None:
            rpm_prefix = self._templar.template(rpm_prefix)
        openshift_release = self.get_var("openshift_release", default='')
        deployment_type = self.get_var("openshift_deployment_type")
        check_multi_minor_release = deployment_type in ['openshift-enterprise']

        args = {
            "package_mgr": self.get_var("ansible_pkg_mgr"),
            "package_list": [
                {
                    "name": "{}".format(rpm_prefix),
                    "version": openshift_release,
                    "check_multi": check_multi_minor_release,
                },
                {
                    "name": "{}-master".format(rpm_prefix),
                    "version": openshift_release,
                    "check_multi": check_multi_minor_release,
                },
                {
                    "name": "{}-node".format(rpm_prefix),
                    "version": openshift_release,
                    "check_multi": check_multi_minor_release,
                },
            ],
        }

        return self.execute_module_with_retries("aos_version", args)
