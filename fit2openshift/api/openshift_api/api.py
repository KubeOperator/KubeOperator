from rest_framework import viewsets

<<<<<<< HEAD
from ansible_api.permissions import IsSuperUser
from . import serializers
from .models import Cluster, Node, Role, DeployExecution, Package
=======
from .models import Cluster, Node, Role, DeployExecution
from .serializers import (
    ClusterSerializer, NodeSerializer, RoleSerializer,
    DeployReadExecutionSerializer
)
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
from .mixin import ClusterResourceAPIMixin
from .tasks import start_deploy_execution


<<<<<<< HEAD
# 集群视图
class ClusterViewSet(viewsets.ModelViewSet):
    queryset = Cluster.objects.all()
    serializer_class = serializers.ClusterSerializer
    permission_classes = (IsSuperUser,)
=======
class ClusterViewSet(viewsets.ModelViewSet):
    queryset = Cluster.objects.all()
    serializer_class = ClusterSerializer
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


<<<<<<< HEAD
# 节点视图
class NodeViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Node.objects.all()
    serializer_class = serializers.NodeSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'
=======
class NodeViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Node.objects.all()
    serializer_class = NodeSerializer
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce


class RoleViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Role.objects.all()
<<<<<<< HEAD
    permission_classes = (IsSuperUser,)
    serializer_class = serializers.RoleSerializer
=======
    serializer_class = RoleSerializer
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce


class DeployExecutionViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = DeployExecution.objects.all()
<<<<<<< HEAD
    serializer_class = serializers.DeployExecutionSerializer
    permission_classes = (IsSuperUser,)
    read_serializer_class = serializers.DeployExecutionSerializer
=======
    serializer_class = DeployReadExecutionSerializer
    read_serializer_class = DeployReadExecutionSerializer
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce

    http_method_names = ['post', 'get', 'head', 'options']

    def perform_create(self, serializer):
        instance = serializer.save()
        start_deploy_execution.apply_async(
            args=(instance.id,), task_id=str(instance.id)
        )
        return instance
<<<<<<< HEAD


class PackageViewSet(viewsets.ModelViewSet):
    queryset = Package.objects.all()
    serializer_class = serializers.PackageSerializer
    permission_classes = (IsSuperUser,)
    http_method_names = ['get', 'head', 'options']
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def get_queryset(self):
        Package.lookup()
        return super().get_queryset()
=======
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
