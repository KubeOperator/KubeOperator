from rest_framework import viewsets

from .models import Cluster, Node, Role, DeployExecution
from .serializers import (
    ClusterSerializer, NodeSerializer, RoleSerializer,
    DeployReadExecutionSerializer
)
from .mixin import ClusterResourceAPIMixin
from .tasks import start_deploy_execution


class ClusterViewSet(viewsets.ModelViewSet):
    queryset = Cluster.objects.all()
    serializer_class = ClusterSerializer
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class NodeViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Node.objects.all()
    serializer_class = NodeSerializer


class RoleViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Role.objects.all()
    serializer_class = RoleSerializer


class DeployExecutionViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = DeployExecution.objects.all()
    serializer_class = DeployReadExecutionSerializer
    read_serializer_class = DeployReadExecutionSerializer

    http_method_names = ['post', 'get', 'head', 'options']

    def perform_create(self, serializer):
        instance = serializer.save()
        start_deploy_execution.apply_async(
            args=(instance.id,), task_id=str(instance.id)
        )
        return instance
