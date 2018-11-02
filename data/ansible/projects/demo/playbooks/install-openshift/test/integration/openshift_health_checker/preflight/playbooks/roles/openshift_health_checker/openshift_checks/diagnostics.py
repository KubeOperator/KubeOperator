"""
A check to run relevant diagnostics via `oc adm diagnostics`.
"""

import os

from openshift_checks import OpenShiftCheck, OpenShiftCheckException


DIAGNOSTIC_LIST = (
    "AggregatedLogging ClusterRegistry ClusterRoleBindings ClusterRoles "
    "ClusterRouter DiagnosticPod NetworkCheck"
).split()


class DiagnosticCheck(OpenShiftCheck):
    """A check to run relevant diagnostics via `oc adm diagnostics`."""

    name = "diagnostics"
    tags = ["health"]

    def is_active(self):
        return super(DiagnosticCheck, self).is_active() and self.is_first_master()

    def run(self):
        if self.exec_diagnostic("ConfigContexts"):
            # only run the other diagnostics if that one succeeds (otherwise, all will fail)
            diagnostics = self.get_var("openshift_check_diagnostics", default=DIAGNOSTIC_LIST)
            for diagnostic in self.normalize(diagnostics):
                self.exec_diagnostic(diagnostic)
        return {}

    def exec_diagnostic(self, diagnostic):
        """
        Execute an 'oc adm diagnostics' command on the remote host.
        Raises OcNotFound or registers OcDiagFailed.
        Returns True on success or False on failure (non-zero rc).
        """
        config_base = self.get_var("openshift.common.config_base")
        args = {
            "config_file": os.path.join(config_base, "master", "admin.kubeconfig"),
            "cmd": "adm diagnostics",
            "extra_args": [diagnostic],
        }

        result = self.execute_module("ocutil", args, save_as_name=diagnostic + ".failure.json")
        self.register_file(diagnostic + ".txt", result['result'])
        if result.get("failed"):
            if result['result'] == '[Errno 2] No such file or directory':
                raise OpenShiftCheckException(
                    "OcNotFound",
                    "This host is supposed to be a master but does not have the `oc` command where expected.\n"
                    "Has an installation been run on this host yet?"
                )

            self.register_failure(OpenShiftCheckException(
                'OcDiagFailed',
                'The {diag} diagnostic reported an error:\n'
                '{error}'.format(diag=diagnostic, error=result['result'])
            ))
            return False
        return True
