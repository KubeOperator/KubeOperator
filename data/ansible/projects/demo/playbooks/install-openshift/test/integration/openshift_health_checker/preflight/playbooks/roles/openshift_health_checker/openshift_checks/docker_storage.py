"""Check Docker storage driver and usage."""
import json
import re
from openshift_checks import OpenShiftCheck, OpenShiftCheckException
from openshift_checks.mixins import DockerHostMixin


class DockerStorage(DockerHostMixin, OpenShiftCheck):
    """Check Docker storage driver compatibility.

    This check ensures that Docker is using a supported storage driver,
    and that loopback is not being used (if using devicemapper).
    Also that storage usage is not above threshold.
    """

    name = "docker_storage"
    tags = ["health", "preflight"]

    dependencies = ["python-docker-py"]
    storage_drivers = ["devicemapper", "overlay", "overlay2"]
    max_thinpool_data_usage_percent = 90.0
    max_thinpool_meta_usage_percent = 90.0
    max_overlay_usage_percent = 90.0

    # TODO(lmeyer): mention these in the output when check fails
    configuration_variables = [
        (
            "max_thinpool_data_usage_percent",
            "For 'devicemapper' storage driver, usage threshold percentage for data. "
            "Format: float. Default: {:.1f}".format(max_thinpool_data_usage_percent),
        ),
        (
            "max_thinpool_meta_usage_percent",
            "For 'devicemapper' storage driver, usage threshold percentage for metadata. "
            "Format: float. Default: {:.1f}".format(max_thinpool_meta_usage_percent),
        ),
        (
            "max_overlay_usage_percent",
            "For 'overlay' or 'overlay2' storage driver, usage threshold percentage. "
            "Format: float. Default: {:.1f}".format(max_overlay_usage_percent),
        ),
    ]

    def run(self):
        msg, failed = self.ensure_dependencies()
        if failed:
            return {
                "failed": True,
                "msg": "Some dependencies are required in order to query docker storage on host:\n" + msg
            }

        # attempt to get the docker info hash from the API
        docker_info = self.execute_module("docker_info", {})
        if docker_info.get("failed"):
            return {"failed": True,
                    "msg": "Failed to query Docker API. Is docker running on this host?"}
        if not docker_info.get("info"):  # this would be very strange
            return {"failed": True,
                    "msg": "Docker API query missing info:\n{}".format(json.dumps(docker_info))}
        docker_info = docker_info["info"]

        # check if the storage driver we saw is valid
        driver = docker_info.get("Driver", "[NONE]")
        if driver not in self.storage_drivers:
            msg = (
                "Detected unsupported Docker storage driver '{driver}'.\n"
                "Supported storage drivers are: {drivers}"
            ).format(driver=driver, drivers=', '.join(self.storage_drivers))
            return {"failed": True, "msg": msg}

        # driver status info is a list of tuples; convert to dict and validate based on driver
        driver_status = {item[0]: item[1] for item in docker_info.get("DriverStatus", [])}

        result = {}

        if driver == "devicemapper":
            result = self.check_devicemapper_support(driver_status)

        if driver in ['overlay', 'overlay2']:
            result = self.check_overlay_support(docker_info, driver_status)

        return result

    def check_devicemapper_support(self, driver_status):
        """Check if dm storage driver is supported as configured. Return: result dict."""
        if driver_status.get("Data loop file"):
            msg = (
                "Use of loopback devices with the Docker devicemapper storage driver\n"
                "(the default storage configuration) is unsupported in production.\n"
                "Please use docker-storage-setup to configure a backing storage volume.\n"
                "See http://red.ht/2rNperO for further information."
            )
            return {"failed": True, "msg": msg}
        result = self.check_dm_usage(driver_status)
        return result

    def check_dm_usage(self, driver_status):
        """Check usage thresholds for Docker dm storage driver. Return: result dict.
        Backing assumptions: We expect devicemapper to be backed by an auto-expanding thin pool
        implemented as an LV in an LVM2 VG. This is how docker-storage-setup currently configures
        devicemapper storage. The LV is "thin" because it does not use all available storage
        from its VG, instead expanding as needed; so to determine available space, we gather
        current usage as the Docker API reports for the driver as well as space available for
        expansion in the pool's VG.
        Usage within the LV is divided into pools allocated to data and metadata, either of which
        could run out of space first; so we check both.
        """
        vals = dict(
            vg_free=self.get_vg_free(driver_status.get("Pool Name")),
            data_used=driver_status.get("Data Space Used"),
            data_total=driver_status.get("Data Space Total"),
            metadata_used=driver_status.get("Metadata Space Used"),
            metadata_total=driver_status.get("Metadata Space Total"),
        )

        # convert all human-readable strings to bytes
        for key, value in vals.copy().items():
            try:
                vals[key + "_bytes"] = self.convert_to_bytes(value)
            except ValueError as err:  # unlikely to hit this from API info, but just to be safe
                return {
                    "failed": True,
                    "values": vals,
                    "msg": "Could not interpret {} value '{}' as bytes: {}".format(key, value, str(err))
                }

        # determine the threshold percentages which usage should not exceed
        for name, default in [("data", self.max_thinpool_data_usage_percent),
                              ("metadata", self.max_thinpool_meta_usage_percent)]:
            percent = self.get_var("max_thinpool_" + name + "_usage_percent", default=default)
            try:
                vals[name + "_threshold"] = float(percent)
            except ValueError:
                return {
                    "failed": True,
                    "msg": "Specified thinpool {} usage limit '{}' is not a percentage".format(name, percent)
                }

        # test whether the thresholds are exceeded
        messages = []
        for name in ["data", "metadata"]:
            vals[name + "_pct_used"] = 100 * vals[name + "_used_bytes"] / (
                vals[name + "_total_bytes"] + vals["vg_free_bytes"])
            if vals[name + "_pct_used"] > vals[name + "_threshold"]:
                messages.append(
                    "Docker thinpool {name} usage percentage {pct:.1f} "
                    "is higher than threshold {thresh:.1f}.".format(
                        name=name,
                        pct=vals[name + "_pct_used"],
                        thresh=vals[name + "_threshold"],
                    ))
                vals["failed"] = True

        vals["msg"] = "\n".join(messages or ["Thinpool usage is within thresholds."])
        return vals

    def get_vg_free(self, pool):
        """Determine which VG to examine according to the pool name. Return: size vgs reports.
        Pool name is the only indicator currently available from the Docker API driver info.
        We assume a name that looks like "vg--name-docker--pool";
        vg and lv names with inner hyphens doubled, joined by a hyphen.
        """
        match = re.match(r'((?:[^-]|--)+)-(?!-)', pool)  # matches up to the first single hyphen
        if not match:  # unlikely, but... be clear if we assumed wrong
            raise OpenShiftCheckException(
                "This host's Docker reports it is using a storage pool named '{}'.\n"
                "However this name does not have the expected format of 'vgname-lvname'\n"
                "so the available storage in the VG cannot be determined.".format(pool)
            )
        vg_name = match.groups()[0].replace("--", "-")
        vgs_cmd = "/sbin/vgs --noheadings -o vg_free --units g --select vg_name=" + vg_name
        # should return free space like "  12.00g" if the VG exists; empty if it does not

        ret = self.execute_module("command", {"_raw_params": vgs_cmd})
        if ret.get("failed") or ret.get("rc", 0) != 0:
            raise OpenShiftCheckException(
                "Is LVM installed? Failed to run /sbin/vgs "
                "to determine docker storage usage:\n" + ret.get("msg", "")
            )
        size = ret.get("stdout", "").strip()
        if not size:
            raise OpenShiftCheckException(
                "This host's Docker reports it is using a storage pool named '{pool}'.\n"
                "which we expect to come from local VG '{vg}'.\n"
                "However, /sbin/vgs did not find this VG. Is Docker for this host"
                "running and using the storage on the host?".format(pool=pool, vg=vg_name)
            )
        return size

    @staticmethod
    def convert_to_bytes(string):
        """Convert string like "10.3 G" to bytes (binary units assumed). Return: float bytes."""
        units = dict(
            b=1,
            k=1024,
            m=1024**2,
            g=1024**3,
            t=1024**4,
            p=1024**5,
        )
        string = string or ""
        match = re.match(r'(\d+(?:\.\d+)?)\s*(\w)?', string)  # float followed by optional unit
        if not match:
            raise ValueError("Cannot convert to a byte size: " + string)

        number, unit = match.groups()
        multiplier = 1 if not unit else units.get(unit.lower())
        if not multiplier:
            raise ValueError("Cannot convert to a byte size: " + string)

        return float(number) * multiplier

    def check_overlay_support(self, docker_info, driver_status):
        """Check if overlay storage driver is supported for this host. Return: result dict."""
        # check for xfs as backing store
        backing_fs = driver_status.get("Backing Filesystem", "[NONE]")
        if backing_fs != "xfs":
            msg = (
                "Docker storage drivers 'overlay' and 'overlay2' are only supported with\n"
                "'xfs' as the backing storage, but this host's storage is type '{fs}'."
            ).format(fs=backing_fs)
            return {"failed": True, "msg": msg}

        # check support for OS and kernel version
        o_s = docker_info.get("OperatingSystem", "[NONE]")
        if "Red Hat Enterprise Linux" in o_s or "CentOS" in o_s:
            # keep it simple, only check enterprise kernel versions; assume everyone else is good
            kernel = docker_info.get("KernelVersion", "[NONE]")
            kernel_arr = [int(num) for num in re.findall(r'\d+', kernel)]
            if kernel_arr < [3, 10, 0, 514]:  # rhel < 7.3
                msg = (
                    "Docker storage drivers 'overlay' and 'overlay2' are only supported beginning with\n"
                    "kernel version 3.10.0-514; but Docker reports kernel version {version}."
                ).format(version=kernel)
                return {"failed": True, "msg": msg}
            # NOTE: we could check for --selinux-enabled here but docker won't even start with
            # that option until it's supported in the kernel so we don't need to.

        return self.check_overlay_usage(docker_info)

    def check_overlay_usage(self, docker_info):
        """Check disk usage on OverlayFS backing store volume. Return: result dict."""
        path = docker_info.get("DockerRootDir", "/var/lib/docker") + "/" + docker_info["Driver"]

        threshold = self.get_var("max_overlay_usage_percent", default=self.max_overlay_usage_percent)
        try:
            threshold = float(threshold)
        except ValueError:
            return {
                "failed": True,
                "msg": "Specified 'max_overlay_usage_percent' is not a percentage: {}".format(threshold),
            }

        mount = self.find_ansible_mount(path)
        try:
            free_bytes = mount['size_available']
            total_bytes = mount['size_total']
            usage = 100.0 * (total_bytes - free_bytes) / total_bytes
        except (KeyError, ZeroDivisionError):
            return {
                "failed": True,
                "msg": "The ansible_mount found for path {} is invalid.\n"
                       "This is likely to be an Ansible bug. The record was:\n"
                       "{}".format(path, json.dumps(mount, indent=2)),
            }

        if usage > threshold:
            return {
                "failed": True,
                "msg": (
                    "For Docker OverlayFS mount point {path},\n"
                    "usage percentage {pct:.1f} is higher than threshold {thresh:.1f}."
                ).format(path=mount["mount"], pct=usage, thresh=threshold)
            }

        return {}
