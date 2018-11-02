"""Check that there is enough disk space in predefined paths."""

import tempfile
import os.path

from openshift_checks import OpenShiftCheck, OpenShiftCheckException


class DiskAvailability(OpenShiftCheck):
    """Check that recommended disk space is available before a first-time install."""

    name = "disk_availability"
    tags = ["preflight"]

    # Values taken from the official installation documentation:
    # https://docs.openshift.org/latest/install_config/install/prerequisites.html#system-requirements
    recommended_disk_space_bytes = {
        '/var': {
            'oo_masters_to_config': 40 * 10**9,
            'oo_nodes_to_config': 15 * 10**9,
            'oo_etcd_to_config': 20 * 10**9,
        },
        # Used to copy client binaries into,
        # see roles/lib_utils/library/openshift_container_binary_sync.py.
        '/usr/local/bin': {
            'oo_masters_to_config': 1 * 10**9,
            'oo_nodes_to_config': 1 * 10**9,
            'oo_etcd_to_config': 1 * 10**9,
        },
        # Used as temporary storage in several cases.
        tempfile.gettempdir(): {
            'oo_masters_to_config': 1 * 10**9,
            'oo_nodes_to_config': 1 * 10**9,
            'oo_etcd_to_config': 1 * 10**9,
        },
    }

    # recommended disk space for each location under an upgrade context
    recommended_disk_upgrade_bytes = {
        '/var': {
            'oo_masters_to_config': 10 * 10**9,
            'oo_nodes_to_config': 5 * 10 ** 9,
            'oo_etcd_to_config': 5 * 10 ** 9,
        },
    }

    def is_active(self):
        """Skip hosts that do not have recommended disk space requirements."""
        group_names = self.get_var("group_names", default=[])
        active_groups = set()
        for recommendation in self.recommended_disk_space_bytes.values():
            active_groups.update(recommendation.keys())
        has_disk_space_recommendation = bool(active_groups.intersection(group_names))
        return super(DiskAvailability, self).is_active() and has_disk_space_recommendation

    def run(self):
        group_names = self.get_var("group_names")
        user_config = self.get_var("openshift_check_min_host_disk_gb", default={})
        try:
            # For backwards-compatibility, if openshift_check_min_host_disk_gb
            # is a number, then it overrides the required config for '/var'.
            number = float(user_config)
            user_config = {
                '/var': {
                    'oo_masters_to_config': number,
                    'oo_nodes_to_config': number,
                    'oo_etcd_to_config': number,
                },
            }
        except TypeError:
            # If it is not a number, then it should be a nested dict.
            pass

        self.register_log("recommended thresholds", self.recommended_disk_space_bytes)
        if user_config:
            self.register_log("user-configured thresholds", user_config)

        # TODO: as suggested in
        # https://github.com/openshift/openshift-ansible/pull/4436#discussion_r122180021,
        # maybe we could support checking disk availability in paths that are
        # not part of the official recommendation but present in the user
        # configuration.
        for path, recommendation in self.recommended_disk_space_bytes.items():
            free_bytes = self.free_bytes(path)
            recommended_bytes = max(recommendation.get(name, 0) for name in group_names)

            config = user_config.get(path, {})
            # NOTE: the user config is in GB, but we compare bytes, thus the
            # conversion.
            config_bytes = max(config.get(name, 0) for name in group_names) * 10**9
            recommended_bytes = config_bytes or recommended_bytes

            # if an "upgrade" context is set, update the minimum disk requirement
            # as this signifies an in-place upgrade - the node might have the
            # required total disk space, but some of that space may already be
            # in use by the existing OpenShift deployment.
            context = self.get_var("r_openshift_health_checker_playbook_context", default="")
            if context == "upgrade":
                recommended_upgrade_paths = self.recommended_disk_upgrade_bytes.get(path, {})
                if recommended_upgrade_paths:
                    recommended_bytes = config_bytes or max(recommended_upgrade_paths.get(name, 0)
                                                            for name in group_names)

            if free_bytes < recommended_bytes:
                free_gb = float(free_bytes) / 10**9
                recommended_gb = float(recommended_bytes) / 10**9
                msg = (
                    'Available disk space in "{}" ({:.1f} GB) '
                    'is below minimum recommended ({:.1f} GB)'
                ).format(path, free_gb, recommended_gb)

                # warn if check failed under an "upgrade" context
                # due to limits imposed by the user config
                if config_bytes and context == "upgrade":
                    msg += ('\n\nMake sure to account for decreased disk space during an upgrade\n'
                            'due to an existing OpenShift deployment. Please check the value of\n'
                            '  openshift_check_min_host_disk_gb={}\n'
                            'in your Ansible inventory, and lower the recommended disk space availability\n'
                            'if necessary for this upgrade.').format(config_bytes)

                self.register_failure(msg)

        return {}

    def find_ansible_submounts(self, path):
        """Return a list of ansible_mounts that are below the given path."""
        base = os.path.join(path, "")
        return [
            mount
            for mount in self.get_var("ansible_mounts")
            if mount["mount"].startswith(base)
        ]

    def free_bytes(self, path):
        """Return the size available in path based on ansible_mounts."""
        submounts = sum(mnt.get('size_available', 0) for mnt in self.find_ansible_submounts(path))
        mount = self.find_ansible_mount(path)
        try:
            return mount['size_available'] + submounts
        except KeyError:
            raise OpenShiftCheckException(
                'Unable to retrieve disk availability for "{path}".\n'
                'Ansible facts included a matching mount point for this path:\n'
                '  {mount}\n'
                'however it is missing the size_available field.\n'
                'To investigate, you can inspect the output of `ansible -m setup <host>`'
                ''.format(path=path, mount=mount)
            )
