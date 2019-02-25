"""Check that required RPM packages are available."""

from openshift_checks import OpenShiftCheck
from openshift_checks.mixins import NotContainerizedMixin


class PackageAvailability(NotContainerizedMixin, OpenShiftCheck):
    """Check that required RPM packages are available."""

    name = "package_availability"
    tags = ["preflight"]

    def is_active(self):
        """Run only when yum is the package manager as the code is specific to it."""
        return super(PackageAvailability, self).is_active() and self.get_var("ansible_pkg_mgr") == "yum"

    def run(self):
        rpm_prefix = self.get_var("openshift_service_type")
        if self._templar is not None:
            rpm_prefix = self._templar.template(rpm_prefix)
        group_names = self.get_var("group_names", default=[])

        packages = set()

        if "oo_masters_to_config" in group_names:
            packages.update(self.master_packages(rpm_prefix))
        if "oo_nodes_to_config" in group_names:
            packages.update(self.node_packages(rpm_prefix))

        args = {"packages": sorted(set(packages))}
        return self.execute_module_with_retries("check_yum_update", args)

    @staticmethod
    def master_packages(rpm_prefix):
        """Return a list of RPMs that we expect a master install to have available."""
        return [
            "{rpm_prefix}".format(rpm_prefix=rpm_prefix),
            "{rpm_prefix}-clients".format(rpm_prefix=rpm_prefix),
            "{rpm_prefix}-hyperkube".format(rpm_prefix=rpm_prefix),
            "bash-completion",
            "httpd-tools",
        ]

    @staticmethod
    def node_packages(rpm_prefix):
        """Return a list of RPMs that we expect a node install to have available."""
        return [
            "{rpm_prefix}".format(rpm_prefix=rpm_prefix),
            "{rpm_prefix}-node".format(rpm_prefix=rpm_prefix),
            "bind",
            "ceph-common",
            "dnsmasq",
            "docker",
            "firewalld",
            "flannel",
            "glusterfs-fuse",
            "iptables-services",
            "iptables",
            "iscsi-initiator-utils",
            "libselinux-python",
            "nfs-utils",
            "ntp",
            "openssl",
            "pyparted",
            "python-httplib2",
            "PyYAML",
            "yum-utils",
        ]
