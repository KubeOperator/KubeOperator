import kubernetes

from kubeops_api.kubernetes import KubernetesApi
from kubeops_api.models import HealthCheck, HealthChecker, Condition
from kubeops_api.models.node import Node

type = "Node health check"


class APIChecker(HealthChecker):
    def __init__(self, cluster):
        self.cluster = cluster

    def check(self):
        api = KubernetesApi(self.cluster)
        client = api.get_api_client()
        core = kubernetes.client.CoreV1Api(client)
        items = core.list_node().items
        self.cluster.change_to()
        nodes = Node.objects.all()

        for item in items:
            for node in nodes:
                conditions = []
                if node.name == item.metadata.labels['kubernetes.io/hostname']:
                    for condition in item.status.conditions:
                        cond = Condition(
                            message=condition.message,
                            reason=condition.reason,
                            status=condition.status,
                            type=type
                        )
                        cond.save()
                        conditions.append(cond)
                    info = item.status.node_info
                    node.info = {
                        "container_runtime_version": info.container_runtime_version,
                        "kernel_version": info.kernel_version,
                        "kube_proxy_version": info.kube_proxy_version,
                        "kubelet_version": info.kubelet_version,
                        "os_image": info.os_image
                    }
                    node.save()
                    node.conditions.set(conditions)


class NodeHealthCheck(HealthCheck):
    def __init__(self, cluster):
        self.cluster = cluster

    def run(self):
        checkers = [
            APIChecker(self.cluster)
        ]
        for checker in checkers:
            checker.check()
