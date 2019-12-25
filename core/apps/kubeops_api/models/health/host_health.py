from kubeops_api.models.health_check import HealthChecker, Condition, HealthCheck
from kubeops_api.models.host import Host
from kubeops_api.utils.health import test_ssh

type = "Host health check"
reason = "health check {} result {}"


class SSHHealthChecker(HealthChecker):
    def __init__(self, host):
        self.host = host

    def check(self):
        ssh_config = self.host.to_ssh_config()
        result, error = test_ssh(ssh_config)
        return Condition(
            type=type,
            status=result,
            message=error,
            reason=reason.format("ssh", "success" if not error else "fail")
        )


class HostHealthCheck(HealthCheck):
    type = "Host health check"

    def __init__(self, host, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.host = host

    def run(self):
        conditions = []
        checkers = [
            SSHHealthChecker(self.host)
        ]
        for checker in checkers:
            cond = checker.check()
            conditions.append(cond)
            cond.save()
            if not cond.status:
                self.host.status = Host.HOST_STATUS_UNKNOWN
        self.host.conditions.set(conditions)
