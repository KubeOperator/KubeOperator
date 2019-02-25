"""Check that required Docker images are available."""

from pipes import quote
from ansible.module_utils import six
from openshift_checks import OpenShiftCheck
from openshift_checks.mixins import DockerHostMixin


class DockerImageAvailability(DockerHostMixin, OpenShiftCheck):
    """Check that required Docker images are available.

    Determine docker images that an install would require and check that they
    are either present in the host's docker index, or available for the host to pull
    with known registries as defined in our inventory file (or defaults).
    """

    name = "docker_image_availability"
    tags = ["preflight"]
    # we use python-docker-py to check local docker for images, and skopeo
    # to look for images available remotely without actually pulling them.

    # command for checking if remote registries have an image, without docker pull
    skopeo_command = "{proxyvars} timeout 10 skopeo inspect --tls-verify={tls} {creds} docker://{image}"
    skopeo_example_command = "skopeo inspect [--tls-verify=false] [--creds=<user>:<pass>] docker://<registry>/<image>"

    def ensure_list(self, registry_param):
        """Return the task var as a list."""
        # https://bugzilla.redhat.com/show_bug.cgi?id=1497274
        # If the result was a string type, place it into a list. We must do this
        # as using list() on a string will split the string into its characters.
        # Otherwise cast to a list as was done previously.
        registry = self.get_var(registry_param, default=[])
        if not isinstance(registry, six.string_types):
            return list(registry)
        return self.normalize(registry)

    def __init__(self, *args, **kwargs):
        super(DockerImageAvailability, self).__init__(*args, **kwargs)

        self.registries_insecure = set(self.ensure_list(
            "openshift_docker_insecure_registries"))

        # Retrieve and template registry credentials, if provided
        self.skopeo_command_creds = None
        oreg_auth_user = self.get_var('oreg_auth_user', default='')
        oreg_auth_password = self.get_var('oreg_auth_password', default='')
        if oreg_auth_user != '' and oreg_auth_password != '':
            oreg_auth_user = self.template_var(oreg_auth_user)
            oreg_auth_password = self.template_var(oreg_auth_password)
            self.skopeo_command_creds = quote("--creds={}:{}".format(oreg_auth_user, oreg_auth_password))

        # take note of any proxy settings needed
        proxies = []
        for var in ['http_proxy', 'https_proxy', 'no_proxy']:
            # ansible vars are openshift_http_proxy, openshift_https_proxy, openshift_no_proxy
            value = self.get_var("openshift_" + var, default=None)
            if value:
                proxies.append(var.upper() + "=" + quote(self.template_var(value)))
        self.skopeo_proxy_vars = " ".join(proxies)

    def is_image_local(self, image):
        """Check if image is already in local docker index."""
        result = self.execute_module("docker_image_facts", {"name": image})
        return bool(result.get("images")) and not result.get("failed")

    def local_images(self, images):
        """Filter a list of images and return those available locally."""
        found_images = []
        for image in images:
            if self.is_image_local(image):
                found_images.append(image)
        return found_images

    def is_available_skopeo_image(self, image):
        """Use Skopeo to determine if required image exists"""
        if six.PY2:
            image = image.encode('utf8')
        use_insecure = False
        for insec_reg in self.registries_insecure:
            if insec_reg in image:
                use_insecure = True
        args = dict(
            proxyvars=self.skopeo_proxy_vars,
            tls="false" if use_insecure else "true",
            creds=self.skopeo_command_creds if self.skopeo_command_creds else "",
            image=quote(image),
        )

        result = self.execute_module_with_retries("command", {
            "_uses_shell": True,
            "_raw_params": self.skopeo_command.format(**args),
        })
        if result.get("rc", 0) == 0 and not result.get("failed"):
            return True

        return False

    def available_images(self, images):
        """Search remotely for images. Returns: list of images found."""
        return [
            image for image in images
            if self.is_available_skopeo_image(image)
        ]

    def run(self):
        '''Run this check'''
        required_images = self.template_var(
            self.get_var('openshift_health_check_required_images'))

        missing_images = set(required_images) - set(self.local_images(required_images))

        # exit early if all images were found locally
        if not missing_images:
            return {}

        available_images = self.available_images(missing_images)
        unavailable_images = set(missing_images) - set(available_images)

        if unavailable_images:
            missing = u",\n    ".join(sorted(unavailable_images))

            msg = (
                u"One or more required container images are not available:\n    {missing}\n"
                "Checked with: {cmd}\n"
            ).format(
                missing=missing,
                cmd=self.skopeo_example_command,
            )

            return dict(failed=True, msg=msg.encode('utf8') if six.PY2 else msg)

        return {}
