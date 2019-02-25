"""Check that recommended memory is available."""
from openshift_checks import OpenShiftCheck

MIB = 2**20
GIB = 2**30


class MemoryAvailability(OpenShiftCheck):
    """Check that recommended memory is available."""

    name = "memory_availability"
    tags = ["preflight"]

    # Values taken from the official installation documentation:
    # https://docs.okd.io/latest/install_config/install/prerequisites.html#system-requirements
    recommended_memory_bytes = {
        "oo_masters_to_config": 16 * GIB,
        "oo_nodes_to_config": 8 * GIB,
        "oo_etcd_to_config": 8 * GIB,
    }
    # https://access.redhat.com/solutions/3006511 physical RAM is partly reserved from memtotal
    memtotal_adjustment = 1 * GIB

    def is_active(self):
        """Skip hosts that do not have recommended memory requirements."""
        group_names = self.get_var("group_names", default=[])
        has_memory_recommendation = bool(set(group_names).intersection(self.recommended_memory_bytes))
        return super(MemoryAvailability, self).is_active() and has_memory_recommendation

    def run(self):
        group_names = self.get_var("group_names")
        total_memory_bytes = self.get_var("ansible_memtotal_mb") * MIB

        recommended_min = max(self.recommended_memory_bytes.get(name, 0) for name in group_names)
        configured_min = float(self.get_var("openshift_check_min_host_memory_gb", default=0)) * GIB
        min_memory_bytes = configured_min or recommended_min

        if total_memory_bytes + self.memtotal_adjustment < min_memory_bytes:
            return {
                'failed': True,
                'msg': (
                    'Available memory ({available:.1f} GiB) is too far '
                    'below recommended value ({recommended:.1f} GiB)'
                ).format(
                    available=float(total_memory_bytes) / GIB,
                    recommended=float(min_memory_bytes) / GIB,
                ),
            }

        return {}
