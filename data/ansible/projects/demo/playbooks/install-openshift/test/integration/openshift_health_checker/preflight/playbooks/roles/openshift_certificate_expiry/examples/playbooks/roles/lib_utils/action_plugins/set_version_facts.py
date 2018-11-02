"""
Ansible action plugin to set version facts
"""

# pylint: disable=no-name-in-module, import-error, wrong-import-order
from distutils.version import LooseVersion
from ansible.plugins.action import ActionBase


# pylint: disable=too-many-statements
def set_version_facts_if_unset(version):
    """ Set version facts. This currently includes common.version and
        common.version_gte_3_x

        Args:
            version (string): version of openshift installed/to install
        Returns:
            dict: the facts dict updated with version facts.
    """
    facts = {}
    if version and version != "latest":
        version = LooseVersion(version)
        version_gte_3_10 = version >= LooseVersion('3.10')
        version_gte_3_11 = version >= LooseVersion('3.11')
    else:
        # 'Latest' version is set to True, 'Next' versions set to False
        version_gte_3_10 = True
        version_gte_3_11 = False
    facts['openshift_version_gte_3_10'] = version_gte_3_10
    facts['openshift_version_gte_3_11'] = version_gte_3_11

    if version_gte_3_11:
        examples_content_version = 'v3.11'
    else:
        examples_content_version = 'v3.10'

    facts['openshift_examples_content_version'] = examples_content_version

    return facts


# pylint: disable=too-few-public-methods
class ActionModule(ActionBase):
    """Action plugin to set version facts"""

    def run(self, tmp=None, task_vars=None):
        """Run set_version_facts"""
        result = super(ActionModule, self).run(tmp, task_vars)
        # Ignore settting self.task_vars outside of init.
        # pylint: disable=W0201
        self.task_vars = task_vars or {}

        result["changed"] = False
        result["failed"] = False
        result["msg"] = "Version facts set"

        version = self._task.args.get('version')
        result["ansible_facts"] = set_version_facts_if_unset(version)
        return result
