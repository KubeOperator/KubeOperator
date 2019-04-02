"""
Ansible action plugin to ensure inventory variables are set
appropriately related to openshift_node_group_name
"""
from ansible.plugins.action import ActionBase
from ansible import errors

# Runs on first master
# Checks each openshift_node_group_name is found in openshift_node_groups
# Checks that master label is present in one of those groups
# Checks that node label is present in one of those groups


def get_or_fail(group, key):
    """Find a key in a group dictionary or fail"""
    res = group.get(key)
    if res is None:
        msg = "Each group in openshift_node_groups must have {} key".format(key)
        raise errors.AnsibleModuleError(msg)
    return res


def validate_labels(labels_found):
    """Ensure mandatory_labels are found in the labels we found, labels_found"""
    mandatory_labels = ('node-role.kubernetes.io/master=true',
                        'node-role.kubernetes.io/infra=true')
    for item in mandatory_labels:
        if item not in labels_found:
            msg = ("At least one group in openshift_node_groups requires the"
                   " {} label").format(item)
            raise errors.AnsibleModuleError(msg)


def process_group(group, groups_found, labels_found):
    """Validate format of each group in openshift_node_groups"""
    name = get_or_fail(group, 'name')
    if name in groups_found:
        msg = ("Duplicate definition of group {} in"
               " openshift_node_groups").format(name)
        raise errors.AnsibleModuleError(msg)
    groups_found.add(name)
    labels = get_or_fail(group, 'labels')
    if not issubclass(type(labels), list):
        msg = "labels value of each group in openshift_node_groups must be a list"
        raise errors.AnsibleModuleError(msg)
    labels_found.update(labels)


class ActionModule(ActionBase):
    """Action plugin to execute node_group_checks."""
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

    def get_node_group_name(self, hostvars, host):
        """Ensure openshift_node_group_name is defined for nodes"""
        group_name = self.template_var(hostvars, host, 'openshift_node_group_name')
        if not group_name:
            msg = "openshift_node_group_name must be defined for all nodes"
            raise errors.AnsibleModuleError(msg)
        return group_name

    def run_check(self, hostvars, host, groups_found):
        """Run the check for each host"""
        group_name = self.get_node_group_name(hostvars, host)
        if group_name not in groups_found:
            msg = "Group: {} not found in openshift_node_groups".format(group_name)
            raise errors.AnsibleModuleError(msg)

    def run(self, tmp=None, task_vars=None):
        """Run node_group_checks action plugin"""
        result = super(ActionModule, self).run(tmp, task_vars)
        result["changed"] = False
        result["failed"] = False
        result["msg"] = "Node group checks passed"
        # self.task_vars holds all in-scope variables.
        # Ignore settting self.task_vars outside of init.
        # pylint: disable=W0201
        self.task_vars = task_vars or {}

        # pylint: disable=W0201
        self.last_checked_host = "none"
        # pylint: disable=W0201
        self.last_checked_var = "none"

        # check_hosts is hard-set to oo_nodes_to_config
        check_hosts = self.task_vars['groups'].get('oo_nodes_to_config')
        if not check_hosts:
            result["msg"] = "skipping; oo_nodes_to_config is required for this check"
            return result

        # We need to access each host's variables
        hostvars = self.task_vars.get('hostvars')
        if not hostvars:
            msg = hostvars
            raise errors.AnsibleModuleError(msg)

        openshift_node_groups = self.task_vars.get('openshift_node_groups')
        if not openshift_node_groups:
            msg = "openshift_node_groups undefined"
            raise errors.AnsibleModuleError(msg)

        openshift_node_groups = self._templar.template(openshift_node_groups)
        groups_found = set()
        labels_found = set()
        # gather the groups and labels we believe should be present.
        for group in openshift_node_groups:
            process_group(group, groups_found, labels_found)

        if len(groups_found) == 0:
            msg = "No groups found in openshift_node_groups"
            raise errors.AnsibleModuleError(msg)

        validate_labels(labels_found)

        # We loop through each host in the provided list check_hosts
        for host in check_hosts:
            try:
                self.run_check(hostvars, host, groups_found)
            except Exception as uncaught_e:
                msg = "last_checked_host: {}, last_checked_var: {};"
                msg = msg.format(self.last_checked_host, self.last_checked_var)
                msg += str(uncaught_e)
                raise errors.AnsibleModuleError(msg)

        return result
