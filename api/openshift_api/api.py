from rest_framework import viewsets
from rest_framework.response import Response
from django.db import transaction

from openshift_api.models.auth import AuthTemplate
from openshift_api.models.host import Host
from ansible_api.permissions import IsSuperUser
from openshift_api.models.cluster import Cluster
from openshift_api.models.deploy import DeployExecution
from openshift_api.models.host import Volume, HostInfo
from openshift_api.models.node import Node
from openshift_api.models.package import Package
from openshift_api.models.role import Role
from openshift_api.models.setting import Setting
from openshift_api.models.storage import StorageTemplate, Storage, StorageNode
from . import serializers
from .mixin import ClusterResourceAPIMixin, StorageResourceAPIMixin
from .tasks import start_deploy_execution
from django.db.models import Q


# 集群视图
class ClusterViewSet(viewsets.ModelViewSet):
    queryset = Cluster.objects.all()
    serializer_class = serializers.ClusterSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class StorageViewSet(viewsets.ModelViewSet):
    queryset = Storage.objects.all()
    serializer_class = serializers.StorageSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


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


class StorageTemplateViewSet(viewsets.ModelViewSet):
    queryset = StorageTemplate.objects.all()
    serializer_class = serializers.StorageTemplateSerializer
    permission_classes = (IsSuperUser,)
    http_method_names = ['get', 'head', 'options']
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def get_queryset(self):
        StorageTemplate.lookup()
        return super().get_queryset()


class AuthViewSet(viewsets.ModelViewSet):
    queryset = AuthTemplate.objects.all()
    serializer_class = serializers.AuthTemplateSerializer
    permission_classes = (IsSuperUser,)
    http_method_names = ['get', 'head', 'options']
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def get_queryset(self):
        AuthTemplate.lookup()
        return super().get_queryset()


class RoleViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Role.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = serializers.RoleSerializer
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class NodeViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Node.objects.all()
    serializer_class = serializers.NodeSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class StorageNodeViewSet(StorageResourceAPIMixin, viewsets.ModelViewSet):
    queryset = StorageNode.objects.all()
    serializer_class = serializers.StorageNodeSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class VolumeViewSet(viewsets.ModelViewSet):
    queryset = Volume.objects.all()
    serializers_class = serializers.VolumeSerializer
    permission_classes = (IsSuperUser,)
    http_method_names = ['get']
    lookup_field = 'host'
    lookup_url_kwarg = 'host_id'


class HostViewSet(viewsets.ModelViewSet):
    queryset = Host.objects.all()
    serializer_class = serializers.HostSerializer
    permission_classes = (IsSuperUser,)

    def perform_create(self, serializer):
        instance = serializer.save()
        transaction.on_commit(lambda: instance.gather_info())


class ClusterConfigViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    serializer_class = serializers.ClusterConfigSerializer
    permission_classes = (IsSuperUser,)
    cluster = None
    lookup_url_kwarg = 'key'

    def dispatch(self, request, *args, **kwargs):
        cluster_name = kwargs.get('cluster_name')
        self.cluster = Cluster.objects.get(name=cluster_name)
        resp = super().dispatch(request, *args, **kwargs)
        return resp

    def retrieve(self, request, *args, **kwargs):
        key = self.kwargs.get('key')
        config = self.cluster.get_config(key) or {}
        serializer = self.serializer_class(config)
        return Response(serializer.data, status=200)

    def update(self, request, *args, **kwargs):
        key = kwargs.get('key')
        data = {k: v for k, v in request.data.items()}
        data['key'] = key
        serializer = self.serializer_class(data=data)
        serializer.is_valid(raise_exception=True)
        data = serializer.validated_data
        self.cluster.set_config(key, data['value'])
        return Response(data=data, status=200)

    def create(self, request, *args, **kwargs):
        serializer = self.serializer_class(data=self.request.data)
        serializer.is_valid(raise_exception=True)
        data = serializer.validated_data
        self.cluster.set_config(data['key'], data['value'])
        return Response(data=serializer.data, status=201)

    def list(self, request, *args, **kwargs):
        configs = self.cluster.configs()
        serializer = self.serializer_class(configs, many=True)
        return Response(serializer.data)

    def destroy(self, request, *args, **kwargs):
        key = self.kwargs.get('key')
        self.cluster.del_config(key)
        return Response(status=204)


# 节点视图


class DeployExecutionViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = DeployExecution.objects.all()
    serializer_class = serializers.DeployExecutionSerializer
    permission_classes = (IsSuperUser,)
    read_serializer_class = serializers.DeployExecutionSerializer

    http_method_names = ['post', 'get', 'head', 'options']

    def perform_create(self, serializer):
        instance = serializer.save()
        transaction.on_commit(lambda: start_deploy_execution.apply_async(
            args=(instance.id,), task_id=str(instance.id)
        ))
        return instance


class HostInfoViewSet(viewsets.ModelViewSet):
    queryset = HostInfo.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = serializers.HostInfoSerializer
    http_method_names = ['head', 'options', 'post']

    def perform_create(self, serializer):
        instance = serializer.save()
        instance.gather_info()


class SettingViewSet(viewsets.ModelViewSet):
    queryset = Setting.objects.all()
    permission_classes = (IsSuperUser,)
    serializer_class = serializers.SettingSerializer
    http_method_names = ['get', 'head', 'options', 'put', 'patch']
    lookup_field = 'key'
    lookup_url_kwarg = 'key'
