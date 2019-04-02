"""
Mixin classes meant to be used with subclasses of OpenShiftCheck.
"""


class NotContainerizedMixin(object):
    """Mixin for checks that are only active when not in containerized mode."""
    # permanent # pylint: disable=too-few-public-methods
    # Reason: The mixin is not intended to stand on its own as a class.

    def is_active(self):
        """Only run on non-containerized hosts."""
        openshift_is_atomic = self.get_var("openshift_is_atomic")
        return super(NotContainerizedMixin, self).is_active() and not openshift_is_atomic


class DockerHostMixin(object):
    """Mixin for checks that are only active on hosts that require Docker."""

    dependencies = []

    def is_active(self):
        """Only run on hosts that depend on Docker."""
        group_names = set(self.get_var("group_names", default=[]))
        needs_docker = set(["oo_nodes_to_config"])
        return super(DockerHostMixin, self).is_active() and bool(group_names.intersection(needs_docker))

    def ensure_dependencies(self):
        """
        Ensure that docker-related packages exist, but not on atomic hosts
        (which would not be able to install but should already have them).
        Returns: msg, failed
        """
        if self.get_var("openshift_is_atomic"):
            return "", False

        # NOTE: we would use the "package" module but it's actually an action plugin
        # and it's not clear how to invoke one of those. This is about the same anyway:
        result = self.execute_module_with_retries(
            self.get_var("ansible_pkg_mgr", default="yum"),
            {"name": self.dependencies, "state": "present"},
        )
        msg = result.get("msg", "")
        if result.get("failed"):
            if "No package matching" in msg:
                msg = "Ensure that all required dependencies can be installed via `yum`.\n"
            msg = (
                "Unable to install required packages on this host:\n"
                "    {deps}\n{msg}"
            ).format(deps=',\n    '.join(self.dependencies), msg=msg)
        failed = result.get("failed", False) or result.get("rc", 0) != 0
        return msg, failed
