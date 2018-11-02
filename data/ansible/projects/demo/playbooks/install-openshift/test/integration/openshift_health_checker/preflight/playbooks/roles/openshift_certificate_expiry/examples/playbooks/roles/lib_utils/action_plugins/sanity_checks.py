"""
Ansible action plugin to ensure inventory variables are set
appropriately and no conflicting options have been provided.
"""
import re

from ansible.plugins.action import ActionBase
from ansible import errors

# Valid values for openshift_deployment_type
VALID_DEPLOYMENT_TYPES = ('origin', 'openshift-enterprise')

# Tuple of variable names and default values if undefined.
NET_PLUGIN_LIST = (('openshift_use_openshift_sdn', True),
                   ('openshift_use_flannel', False),
                   ('openshift_use_nuage', False),
                   ('openshift_use_contiv', False),
                   ('openshift_use_calico', False),
                   ('openshift_use_kuryr', False),
                   ('openshift_use_nsx', False))

ENTERPRISE_TAG_REGEX_ERROR = """openshift_image_tag must be in the format
v#.#[.#[.#]]. Examples: v1.2, v3.4.1, v3.5.1.3,
v3.5.1.3.4, v1.2-1, v1.2.3-4, v1.2.3-4.5, v1.2.3-4.5.6
You specified openshift_image_tag={}"""

ORIGIN_TAG_REGEX_ERROR = """openshift_image_tag must be in the format
v#.#[.#-optional.#]. Examples: v1.2.3, v3.5.1-alpha.1
You specified openshift_image_tag={}"""

ORIGIN_TAG_REGEX = {'re': '(^v?\\d+\\.\\d+.*)',
                    'error_msg': ORIGIN_TAG_REGEX_ERROR}
ENTERPRISE_TAG_REGEX = {'re': '(^v\\d+\\.\\d+(\\.\\d+)*(-\\d+(\\.\\d+)*)?$)',
                        'error_msg': ENTERPRISE_TAG_REGEX_ERROR}
IMAGE_TAG_REGEX = {'origin': ORIGIN_TAG_REGEX,
                   'openshift-enterprise': ENTERPRISE_TAG_REGEX}

PKG_VERSION_REGEX_ERROR = """openshift_pkg_version must be in the format
-[optional.release]. Examples: -3.6.0, -3.7.0-0.126.0.git.0.9351aae.el7 -3.11*
You specified openshift_pkg_version={}"""
PKG_VERSION_REGEX = {'re': '(^-.*)',
                     'error_msg': PKG_VERSION_REGEX_ERROR}

RELEASE_REGEX_ERROR = """openshift_release must be in the format
v#[.#[.#]]. Examples: v3.9, v3.10.0
You specified openshift_release={}"""
RELEASE_REGEX = {'re': '(^v?\\d+(\\.\\d+(\\.\\d+)?)?$)',
                 'error_msg': RELEASE_REGEX_ERROR}

STORAGE_KIND_TUPLE = (
    'openshift_hosted_registry_storage_kind',
    'openshift_loggingops_storage_kind',
    'openshift_logging_storage_kind',
    'openshift_metrics_storage_kind',
    'openshift_prometheus_alertbuffer_storage_kind',
    'openshift_prometheus_alertmanager_storage_kind',
    'openshift_prometheus_storage_kind')

REMOVED_VARIABLES = (
    # TODO(michaelgugino): Remove these in 3.11
    ('openshift_metrics_image_prefix', 'openshift_metrics_<component>_image'),
    ('openshift_metrics_image_version', 'openshift_metrics_<component>_image'),
    ('openshift_grafana_proxy_image_prefix', 'openshift_grafana_proxy_image'),
    ('openshift_grafana_proxy_image_version', 'openshift_grafana_proxy_image'),
    ('openshift_logging_image_prefix', 'openshift_logging_image'),
    ('openshift_logging_image_verion', 'openshift_logging_image'),
    ('openshift_logging_curator_image_prefix', 'openshift_logging_curator_image'),
    ('openshift_logging_curator_image_version', 'openshift_logging_curator_image'),
    ('openshift_logging_elasticsearch_image_prefix', 'openshift_logging_elasticsearch_image'),
    ('openshift_logging_elasticsearch_image_version', 'openshift_logging_elasticsearch_image'),
    ('openshift_logging_elasticsearch_proxy_image_prefix', 'openshift_logging_elasticsearch_proxy_image'),
    ('openshift_logging_elasticsearch_proxy_image_version', 'openshift_logging_elasticsearch_proxy_image'),
    ('openshift_logging_fluentd_image_prefix', 'openshift_logging_fluentd_image'),
    ('openshift_logging_fluentd_image_version', 'openshift_logging_fluentd_image'),
    ('openshift_logging_kibana_image_prefix', 'openshift_logging_kibana_image'),
    ('openshift_logging_kibana_image_version', 'openshift_logging_kibana_image'),
    ('openshift_logging_kibana_proxy_image_prefix', 'openshift_logging_kibana_proxy_image'),
    ('openshift_logging_kibana_proxy_image_version', 'openshift_logging_kibana_proxy_image'),
    ('openshift_logging_mux_image_prefix', 'openshift_logging_mux_image'),
    ('openshift_logging_mux_image_version', 'openshift_logging_mux_image'),
    ('openshift_prometheus_image_prefix', 'openshift_prometheus_image'),
    ('openshift_prometheus_image_version', 'openshift_prometheus_image'),
    ('openshift_prometheus_proxy_image_prefix', 'openshift_prometheus_proxy_image'),
    ('openshift_prometheus_proxy_image_version', 'openshift_prometheus_proxy_image'),
    ('openshift_prometheus_altermanager_image_prefix', 'openshift_prometheus_alertmanager_image'),
    # A typo was introduced at some point, need to warn for this older version.
    ('openshift_prometheus_altermanager_image_prefix', 'openshift_prometheus_alertmanager_image'),
    ('openshift_prometheus_alertmanager_image_version', 'openshift_prometheus_alertmanager_image'),
    ('openshift_prometheus_alertbuffer_image_prefix', 'openshift_prometheus_alertbuffer_image'),
    ('openshift_prometheus_alertbuffer_image_version', 'openshift_prometheus_alertbuffer_image'),
    ('openshift_prometheus_node_exporter_image_prefix', 'openshift_prometheus_node_exporter_image'),
    ('openshift_prometheus_node_exporter_image_version', 'openshift_prometheus_node_exporter_image'),
    ('openshift_descheduler_image_prefix', 'openshift_descheduler_image'),
    ('openshift_descheduler_image_version', 'openshift_descheduler_image'),
    ('openshift_docker_gc_version', 'openshift_docker_gc_image'),
    ('openshift_web_console_prefix', 'openshift_web_console_image'),
    ('openshift_web_console_version', 'openshift_web_console_image'),
    ('openshift_web_console_image_name', 'openshift_web_console_image'),
    ('openshift_storage_glusterfs_version', 'openshift_storage_glusterfs_image'),
    ('openshift_storage_glusterfs_block_version', 'openshift_storage_glusterfs_block_image'),
    ('openshift_storage_glusterfs_s3_version', 'openshift_storage_glusterfs_s3_image'),
    ('openshift_storage_glusterfs_heketi_version', 'openshift_storage_glusterfs_heketi_image'),
    ('openshift_storage_glusterfs_registry_version', 'openshift_storage_glusterfs_registry_image'),
    ('openshift_storage_glusterfs_registry_block_version', 'openshift_storage_glusterfs_registry_block_image'),
    ('openshift_storage_glusterfs_registry_s3_version', 'openshift_storage_glusterfs_registry_s3_image'),
    ('openshift_storage_glusterfs_registry_heketi_version', 'openshift_storage_glusterfs_registry_heketi_image'),
    # TODO(michaelgugino): Remove in 3.12
    ('openshift_cockpit_deployer_prefix', 'openshift_cockpit_deployer_image'),
    ('openshift_cockpit_deployer_basename', 'openshift_cockpit_deployer_image'),
    ('openshift_cockpit_deployer_version', 'openshift_cockpit_deployer_image'),
    ('openshift_hostname', 'Removed: See documentation'),
)

# TODO(michaelgugino): Remove in 3.11
CHANGED_IMAGE_VARS = (
    'openshift_storage_glusterfs_image',
    'openshift_storage_glusterfs_block_image',
    'openshift_storage_glusterfs_s3_image',
    'openshift_storage_glusterfs_heketi_image',
    'openshift_storage_glusterfs_registry_image',
    'openshift_storage_glusterfs_registry_block_image',
    'openshift_storage_glusterfs_registry_s3_image',
    'openshift_storage_glusterfs_registry_heketi_image',
)


def to_bool(var_to_check):
    """Determine a boolean value given the multiple
       ways bools can be specified in ansible."""
    # http://yaml.org/type/bool.html
    yes_list = (True, 1, "True", "1", "true", "TRUE",
                "Yes", "yes", "Y", "y", "YES",
                "on", "ON", "On")
    return var_to_check in yes_list


def check_for_removed_vars(hostvars, host):
    """Fails if removed variables are found"""
    found_removed = []
    for item in REMOVED_VARIABLES:
        if item in hostvars[host]:
            found_removed.append(item)

    if found_removed:
        msg = "Found removed variables: "
        for item in found_removed:
            msg += "{} is replaced by {}; ".format(item[0], item[1])
        raise errors.AnsibleModuleError(msg)
    return None


class ActionModule(ActionBase):
    """Action plugin to execute sanity checks."""
    def template_var(self, hostvars, host, varname):
        """Retrieve a variable from hostvars and template it.
           If undefined, return None type."""
        # We will set the current host and variable checked for easy debugging
        # if there are any unhandled exceptions.
        # pylint: disable=W0201
        self.last_checked_var = varname
        # pylint: disable=W0201
        self.last_checked_host = host
        res = hostvars[host].get(varname)
        if res is None:
            return None
        return self._templar.template(res)

    def check_openshift_deployment_type(self, hostvars, host):
        """Ensure a valid openshift_deployment_type is set"""
        openshift_deployment_type = self.template_var(hostvars, host,
                                                      'openshift_deployment_type')
        if openshift_deployment_type not in VALID_DEPLOYMENT_TYPES:
            type_strings = ", ".join(VALID_DEPLOYMENT_TYPES)
            msg = "openshift_deployment_type must be defined and one of {}".format(type_strings)
            raise errors.AnsibleModuleError(msg)
        return openshift_deployment_type

    def check_python_version(self, hostvars, host, distro):
        """Ensure python version is 3 for Fedora and python 2 for others"""
        ansible_python = self.template_var(hostvars, host, 'ansible_python')
        if distro == "Fedora":
            if ansible_python['version']['major'] != 3:
                msg = "openshift-ansible requires Python 3 for {};".format(distro)
                msg += " For information on enabling Python 3 with Ansible,"
                msg += " see https://docs.ansible.com/ansible/python_3_support.html"
                raise errors.AnsibleModuleError(msg)
        else:
            if ansible_python['version']['major'] != 2:
                msg = "openshift-ansible requires Python 2 for {};".format(distro)

    def check_image_tag_format(self, hostvars, host, openshift_deployment_type):
        """Ensure openshift_image_tag is formatted correctly"""
        openshift_image_tag = self.template_var(hostvars, host, 'openshift_image_tag')
        if not openshift_image_tag or openshift_image_tag == 'latest':
            return None
        regex_to_match = IMAGE_TAG_REGEX[openshift_deployment_type]['re']
        res = re.match(regex_to_match, str(openshift_image_tag))
        if res is None:
            msg = IMAGE_TAG_REGEX[openshift_deployment_type]['error_msg']
            msg = msg.format(str(openshift_image_tag))
            raise errors.AnsibleModuleError(msg)

    def check_pkg_version_format(self, hostvars, host):
        """Ensure openshift_pkg_version is formatted correctly"""
        openshift_pkg_version = self.template_var(hostvars, host, 'openshift_pkg_version')
        if not openshift_pkg_version:
            return None
        regex_to_match = PKG_VERSION_REGEX['re']
        res = re.match(regex_to_match, str(openshift_pkg_version))
        if res is None:
            msg = PKG_VERSION_REGEX['error_msg']
            msg = msg.format(str(openshift_pkg_version))
            raise errors.AnsibleModuleError(msg)

    def check_release_format(self, hostvars, host):
        """Ensure openshift_release is formatted correctly"""
        openshift_release = self.template_var(hostvars, host, 'openshift_release')
        if not openshift_release:
            return None
        regex_to_match = RELEASE_REGEX['re']
        res = re.match(regex_to_match, str(openshift_release))
        if res is None:
            msg = RELEASE_REGEX['error_msg']
            msg = msg.format(str(openshift_release))
            raise errors.AnsibleModuleError(msg)

    def network_plugin_check(self, hostvars, host):
        """Ensure only one type of network plugin is enabled"""
        res = []
        # Loop through each possible network plugin boolean, determine the
        # actual boolean value, and append results into a list.
        for plugin, default_val in NET_PLUGIN_LIST:
            res_temp = self.template_var(hostvars, host, plugin)
            if res_temp is None:
                res_temp = default_val
            res.append(to_bool(res_temp))

        if sum(res) not in (0, 1):
            plugin_str = list(zip([x[0] for x in NET_PLUGIN_LIST], res))

            msg = "ClusterHost Checked: {} Only one of must be true. Found: {}".format(host, plugin_str)
            raise errors.AnsibleModuleError(msg)

    def check_hostname_vars(self, hostvars, host):
        """Checks to ensure openshift_kubelet_name_override
           and openshift_public_hostname
           conform to the proper length of 63 characters or less"""
        for varname in ('openshift_public_hostname', 'openshift_kubelet_name_override'):
            var_value = self.template_var(hostvars, host, varname)
            if var_value and len(var_value) > 63:
                msg = '{} must be 63 characters or less'.format(varname)
                raise errors.AnsibleModuleError(msg)

    def check_session_auth_secrets(self, hostvars, host):
        """Checks session_auth_secrets is correctly formatted"""
        sas = self.template_var(hostvars, host,
                                'openshift_master_session_auth_secrets')
        ses = self.template_var(hostvars, host,
                                'openshift_master_session_encryption_secrets')
        # This variable isn't mandatory, only check if set.
        if sas is None and ses is None:
            return None

        if not (
                issubclass(type(sas), list) and issubclass(type(ses), list)
        ) or len(sas) != len(ses):
            raise errors.AnsibleModuleError(
                'Expects openshift_master_session_auth_secrets and '
                'openshift_master_session_encryption_secrets are equal length lists')

        for secret in sas:
            if len(secret) < 32:
                raise errors.AnsibleModuleError(
                    'Invalid secret in openshift_master_session_auth_secrets. '
                    'Secrets must be at least 32 characters in length.')

        for secret in ses:
            if len(secret) not in [16, 24, 32]:
                raise errors.AnsibleModuleError(
                    'Invalid secret in openshift_master_session_encryption_secrets. '
                    'Secrets must be 16, 24, or 32 characters in length.')
        return None

    def check_unsupported_nfs_configs(self, hostvars, host):
        """Fails if nfs storage is in use for any components. This check is
           ignored if openshift_enable_unsupported_configurations=True"""

        enable_unsupported = self.template_var(
            hostvars, host, 'openshift_enable_unsupported_configurations')

        if to_bool(enable_unsupported):
            return None

        for storage in STORAGE_KIND_TUPLE:
            kind = self.template_var(hostvars, host, storage)
            if kind == 'nfs':
                raise errors.AnsibleModuleError(
                    'nfs is an unsupported type for {}. '
                    'openshift_enable_unsupported_configurations=True must'
                    'be specified to continue with this configuration.'
                    ''.format(storage))
        return None

    def check_htpasswd_provider(self, hostvars, host):
        """Fails if openshift_master_identity_providers contains an entry of
        kind HTPasswdPasswordIdentityProvider and
        openshift_master_manage_htpasswd is False"""

        manage_pass = self.template_var(
            hostvars, host, 'openshift_master_manage_htpasswd')
        # default for openshift_master_manage_htpasswd is True.
        if to_bool(manage_pass) or manage_pass is None:
            # If we manage the file, we can just generate in the new path.
            return None
        idps = self.template_var(
            hostvars, host, 'openshift_master_identity_providers')
        if not idps:
            # If we don't find any identity_providers, nothing for us to do.
            return None
        old_keys = ('file', 'fileName', 'file_name', 'filename')
        if not isinstance(idps, list):
            raise errors.AnsibleModuleError("| not a list")
        for idp in idps:
            if idp['kind'] == 'HTPasswdPasswordIdentityProvider':
                for old_key in old_keys:
                    if old_key in idp is not None:
                        raise errors.AnsibleModuleError(
                            'openshift_master_identity_providers contains a '
                            'provider of kind==HTPasswdPasswordIdentityProvider '
                            'and {} is set.  Please migrate your htpasswd '
                            'files to /etc/origin/master/htpasswd and update your '
                            'existing master configs, and remove the {} key'
                            'before proceeding.'.format(old_key, old_key))

    def check_contains_version(self, hostvars, host):
        """Fails if variable has old format not containing image version"""
        found_incorrect = []
        for img_var in CHANGED_IMAGE_VARS:
            img_string = self.template_var(hostvars, host, img_var)
            if not img_string:
                return None
            # split the image string by '/' to account for something like docker://
            img_string_parts = img_string.split('/')
            if ':' not in img_string_parts[-1]:
                found_incorrect.append((img_var, img_string))

        if found_incorrect:
            msg = ("Found image variables without version.  Please ensure any image"
                   "defined contains ':someversion'; ")
            for item in found_incorrect:
                msg += "{} was: {}; ".format(item[0], item[1])
            raise errors.AnsibleModuleError(msg)
        return None

    def run_checks(self, hostvars, host):
        """Execute the hostvars validations against host"""
        distro = self.template_var(hostvars, host, 'ansible_distribution')
        odt = self.check_openshift_deployment_type(hostvars, host)
        self.check_python_version(hostvars, host, distro)
        self.check_image_tag_format(hostvars, host, odt)
        self.check_pkg_version_format(hostvars, host)
        self.check_release_format(hostvars, host)
        self.network_plugin_check(hostvars, host)
        self.check_hostname_vars(hostvars, host)
        self.check_session_auth_secrets(hostvars, host)
        self.check_unsupported_nfs_configs(hostvars, host)
        self.check_htpasswd_provider(hostvars, host)
        check_for_removed_vars(hostvars, host)
        self.check_contains_version(hostvars, host)

    def run(self, tmp=None, task_vars=None):
        result = super(ActionModule, self).run(tmp, task_vars)

        # self.task_vars holds all in-scope variables.
        # Ignore settting self.task_vars outside of init.
        # pylint: disable=W0201
        self.task_vars = task_vars or {}

        # pylint: disable=W0201
        self.last_checked_host = "none"
        # pylint: disable=W0201
        self.last_checked_var = "none"

        # self._task.args holds task parameters.
        # check_hosts is a parameter to this plugin, and should provide
        # a list of hosts.
        check_hosts = self._task.args.get('check_hosts')
        if not check_hosts:
            msg = "check_hosts is required"
            raise errors.AnsibleModuleError(msg)

        # We need to access each host's variables
        hostvars = self.task_vars.get('hostvars')
        if not hostvars:
            msg = hostvars
            raise errors.AnsibleModuleError(msg)

        # We loop through each host in the provided list check_hosts
        for host in check_hosts:
            try:
                self.run_checks(hostvars, host)
            except Exception as uncaught_e:
                msg = "last_checked_host: {}, last_checked_var: {};"
                msg = msg.format(self.last_checked_host, self.last_checked_var)
                msg += str(uncaught_e)
                raise errors.AnsibleModuleError(msg)

        result["changed"] = False
        result["failed"] = False
        result["msg"] = "Sanity Checks passed"

        return result
