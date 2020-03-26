import json
import logging
import os
import yaml
import kubeops_api.cluster_backup_utils
import kubeops_api.cluster_monitor
import kubeops_api.cluster_backup_utils
import log.es
from django.db import transaction
from django.db.models import Q
from django.http import HttpResponse, JsonResponse
from django.shortcuts import get_object_or_404
from rest_framework import viewsets, status
from rest_framework.response import Response
from rest_framework.views import APIView
from kubeoperator.settings import VERSION_DIR, CLUSTER_CONFIG_DIR
from kubeops_api.cluster_monitor import ClusterMonitor
from kubeops_api.models.backup_storage import BackupStorage
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.credential import Credential
from kubeops_api.models.deploy import DeployExecution
from kubeops_api.models.host import Host
from kubeops_api.models.node import Node
from kubeops_api.models.package import Package
from kubeops_api.models.role import Role
from kubeops_api.models.setting import Setting
from . import s_serializers as serializers
from .mixin import ClusterResourceAPIMixin
from .tasks import start_deploy_execution
from kubeops_api.storage_client import StorageClient
from kubeops_api.models.backup_strategy import BackupStrategy
from kubeops_api.models.cluster_backup import ClusterBackup
from rest_framework import generics
from kubeops_api.models.cluster_health_history import ClusterHealthHistory
from storage.models import ClusterCephStorage
from kubeops_api.models.item import Item, ItemRoleMapping
from kubeops_api.models.item_resource import ItemResource

logger = logging.getLogger('kubeops')


class ClusterViewSet(viewsets.ModelViewSet):
    queryset = Cluster.objects.all()
    serializer_class = serializers.ClusterSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        if request.user.is_superuser is False:
            can_delete = False
            item_resource = ItemResource.objects.filter(resource_id=instance.id,
                                                        resource_type=ItemResource.RESOURCE_TYPE_CLUSTER)
            if item_resource:
                role_mappings = ItemRoleMapping.objects.filter(profile__id=request.user.profile.id,
                                                               item_id=item_resource[0].item_id)
                if len(role_mappings) > 0 and role_mappings[0].role == ItemRoleMapping.ITEM_ROLE_MANAGER:
                    can_delete = True
            if can_delete is False:
                return Response(data={'msg': '当前用户没有删除权限'}, status=status.HTTP_400_BAD_REQUEST)
        if not instance.status == Cluster.CLUSTER_STATUS_READY and not instance.status == Cluster.CLUSTER_STATUS_ERROR:
            return Response(data={'msg': '集群处于: {} 状态,不可删除'.format(instance.status)},
                            status=status.HTTP_400_BAD_REQUEST)
        response = super().destroy(self, request, *args, **kwargs)
        if response.status_code == 204:
            BackupStrategy.objects.filter(project_id=instance.id).delete()
            ClusterBackup.objects.filter(project_id=instance.id).delete()
            kubeops_api.cluster_monitor.delete_cluster_redis_data(instance.name)
            ClusterHealthHistory.objects.filter(project_id=instance.id).delete()
            ClusterCephStorage.objects.filter(cluster_id=instance.id).delete()
            ItemResource.objects.filter(resource_id=instance.id).delete()
        return response

    def list(self, request, *args, **kwargs):
        user = request.user
        if request.query_params.get('itemName'):
            itemName = request.query_params.get('itemName')
            if itemName == 'all' and user.profile.items:
                item_ids = []
                for item in user.profile.items:
                    item_ids.append(item.id)
                resource_ids = ItemResource.objects.filter(item_id__in=item_ids).values_list("resource_id")
            elif itemName == 'all' and user.is_superuser:
                item_ids = Item.objects.all().values_list("id")
                resource_ids = ItemResource.objects.filter(item_id__in=item_ids).values_list("resource_id")
            else:
                try:
                    item = Item.objects.get(name=itemName)
                    resource_ids = ItemResource.objects.filter(item_id=item.id).values_list("resource_id")
                except Item.DoesNotExist:
                    resource_ids = []
            self.queryset = Cluster.objects.filter(id__in=resource_ids)
            return super().list(self, request, *args, **kwargs)
        elif user.profile.items:
            item_ids = []
            for item in user.profile.items:
                item_ids.append(item.id)
            resource_ids = ItemResource.objects.filter(item_id__in=item_ids).values_list("resource_id")
            self.queryset = Cluster.objects.filter(id__in=resource_ids)
            return super().list(self, request, *args, **kwargs)
        elif user.is_superuser:
            return super().list(self, request, *args, **kwargs)
        else:
            self.queryset = []
            return super().list(self, request, *args, **kwargs)

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        self.perform_create(serializer)
        headers = self.get_success_headers(serializer.data)
        if request.data.get('item_name'):
            item = Item.objects.get(name=request.data.get('item_name'))
            cluster = Cluster.objects.get(name=request.data.get('name'))
            itemResource = ItemResource(item_id=item.id, resource_id=cluster.id,
                                        resource_type=ItemResource.RESOURCE_TYPE_CLUSTER)
            itemResource.save()
        return Response(serializer.data, status=status.HTTP_201_CREATED, headers=headers)


class PackageViewSet(viewsets.ModelViewSet):
    queryset = Package.objects.all()
    serializer_class = serializers.PackageSerializer

    http_method_names = ['get', 'head', 'options']
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def get_queryset(self):
        try:
            Package.lookup()
        except Exception as e:
            logger.error(e, exc_info=True)
        return super().get_queryset()


class RoleViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Role.objects.all()

    serializer_class = serializers.RoleSerializer
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class NodeViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    queryset = Node.objects.all()
    serializer_class = serializers.NodeSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class CredentialViewSet(viewsets.ModelViewSet):
    queryset = Credential.objects.all()
    serializer_class = serializers.CredentialSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        query_set = Host.objects.filter(credential__name=instance.name)
        if len(query_set) > 0:
            return Response(data={'msg': '凭据: {} 下资源不为空'.format(instance.name)}, status=status.HTTP_400_BAD_REQUEST)
        return super().destroy(self, request, *args, **kwargs)


class ClusterConfigViewSet(ClusterResourceAPIMixin, viewsets.ModelViewSet):
    serializer_class = serializers.ClusterConfigSerializer

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
        configs = self.cluster.get_configs()
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

    read_serializer_class = serializers.DeployExecutionSerializer

    http_method_names = ['post', 'get', 'head', 'options']

    def create(self, request, *args, **kwargs):
        cluster_name = kwargs.get('cluster_name')
        self.mark(cluster_name)
        cluster = Cluster.objects.get(name=cluster_name)
        if cluster.deploy_type == Cluster.CLUSTER_DEPLOY_TYPE_AUTOMATIC:
            operation = request.data['operation']
            cluster.change_to()
            nodes = Node.objects.all()
            if operation == 'install':
                if cluster.worker_size > cluster.plan.count_ip_available() + len(nodes):
                    return Response(data={'msg': ': Ip 资源不足！'}, status=status.HTTP_400_BAD_REQUEST)
            if operation == 'scale':
                num = request.data['params']['num']
                if num > cluster.worker_size:
                    if num - cluster.worker_size > cluster.plan.count_ip_available():
                        return Response(data={'msg': ': Ip 资源不足！'}, status=status.HTTP_400_BAD_REQUEST)
        return super().create(request, *args, **kwargs)

    def mark(self, cluster_name):
        cluster = Cluster.objects.get(name=cluster_name)
        last = cluster.current_execution
        if last and last.state == last.STATE_STARTED:
            last.mark_state(last.STATE_FAILURE)

    def perform_create(self, serializer):
        instance = serializer.save()
        transaction.on_commit(lambda: start_deploy_execution.apply_async(
            args=(instance.id,), task_id=str(instance.id)
        ))
        return instance


class VersionView(APIView):

    def get(self, request, **kwargs):
        with open(VERSION_DIR) as f:
            response = HttpResponse()
            result = yaml.load(f)
            response.write(json.dumps(result))
            return response


class DownloadView(APIView):

    def get(self, request, **kwargs):
        pk = kwargs.get("pk")
        cluster = get_object_or_404(Cluster, pk=pk)
        file_name = cluster.fetch_config()
        with open(file_name) as f:
            response = HttpResponse(f)
            response["content_type"] = 'application/octet-stream'
            response['Content-Disposition'] = "attachment; filename= {}".format(cluster.name + '-kube-config')
            return response


class GetClusterTokenView(APIView):

    def get(self, request, **kwargs):
        pk = kwargs.get("pk")
        cluster = get_object_or_404(Cluster, pk=pk)
        token = cluster.get_cluster_token()
        result = {
            "token": token
        }
        response = HttpResponse()
        response.write(json.dumps(result))
        return response


class GetClusterConfigView(APIView):
    def get(self, request, **kwargs):
        config_file = os.path.join(CLUSTER_CONFIG_DIR, "config.yml")
        with open(config_file) as f:
            data = yaml.load(f)
            return JsonResponse(data)


class BackupStorageViewSet(viewsets.ModelViewSet):
    queryset = BackupStorage.objects.all()
    serializer_class = serializers.BackupStorageSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        backup_storage_id = BackupStorage.objects.get(name=self.kwargs['name']).id
        result = BackupStrategy.objects.filter(backup_storage_id=backup_storage_id)
        if len(result) > 0:
            return Response(data={'msg': ': 有集群使用此备份账号!删除失败!'}, status=status.HTTP_400_BAD_REQUEST)
        else:
            return super().destroy(self, request, *args, **kwargs)

    def list(self, request, *args, **kwargs):
        if request.query_params.get('itemName'):
            itemName = request.query_params.get('itemName')
            item = Item.objects.get(name=itemName)
            resource_ids = ItemResource.objects.filter(item_id=item.id,
                                                       resource_type=ItemResource.RESOURCE_TYPE_BACKUP_STORAGE).values_list(
                "resource_id")
            self.queryset = BackupStorage.objects.filter(id__in=resource_ids)
            return super().list(self, request, *args, **kwargs)
        else:
            return super().list(self, request, *args, **kwargs)


class CheckStorageView(APIView):

    def post(self, request, **kwargs):
        client = StorageClient(request.data)
        valid = client.check_valid()
        response = HttpResponse()
        result = {
            "message": '验证成功!',
            "success": True
        }
        if valid:
            response.write(json.dumps(result))
        else:
            result['message'] = '验证失败！'
            result['success'] = False
            response.write(json.dumps(result))
        return response


class GetBucketsView(APIView):

    def post(self, request):
        client = StorageClient(request.data)
        buckets = client.list_buckets()
        response = HttpResponse()
        result = {
            "message": '验证成功!',
            "success": True,
            "data": buckets
        }
        response.write(json.dumps(result))
        return response


class BackupStrategyViewSet(viewsets.ModelViewSet):
    queryset = BackupStrategy.objects.all()
    serializer_class = serializers.BackupStrategySerializer

    lookup_field = 'project_id'
    lookup_url_kwarg = 'project_id'


class ClusterBackupViewSet(viewsets.ModelViewSet):
    queryset = ClusterBackup.objects.all()
    serializer_class = serializers.ClusterBackupSerializer

    lookup_field = 'project_id'
    lookup_url_kwarg = 'project_id'


class ClusterBackupList(generics.ListAPIView):
    serializer_class = serializers.ClusterBackupSerializer

    def get_queryset(self):
        project_id = self.kwargs['project_id']
        return ClusterBackup.objects.filter(project_id=project_id)


class ClusterBackupDelete(generics.DestroyAPIView):
    serializer_class = serializers.ClusterBackupSerializer

    def destroy(self, request, *args, **kwargs):
        id = self.kwargs['id']
        ok = kubeops_api.cluster_backup_utils.delete_backup(id)
        result = {
            "message": '删除成功!',
            "success": True
        }
        response = HttpResponse()
        if ok:
            response.write(json.dumps(result))
        else:
            result['message'] = '删除失败！'
            result['success'] = False
            response.write(json.dumps(result))
        return response


class ClusterBackupRestore(generics.UpdateAPIView):
    serializer_class = serializers.ClusterBackupSerializer

    def put(self, request, *args, **kwargs):
        ok = kubeops_api.cluster_backup_utils.run_restore(request.data['id'])
        result = {
            "message": '恢复成功!',
            "success": True
        }
        response = HttpResponse()
        if ok:
            response.write(json.dumps(result))
        else:
            result['message'] = '恢复失败！'
            result['success'] = False
            response.write(json.dumps(result))
        return response


class ClusterHealthHistoryView(generics.ListAPIView):
    serializer_class = serializers.ClusterHeathHistorySerializer

    def get_queryset(self):
        project_id = str(self.kwargs['project_id'])
        return ClusterHealthHistory.objects.filter(project_id=str(project_id),
                                                   date_type=ClusterHealthHistory.CLUSTER_HEALTH_HISTORY_DATE_TYPE_DAY).order_by(
            '-date_created')


class ClusterHealthView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = self.kwargs['project_name']
        namespace = self.kwargs['namespace']
        cluster = Cluster.objects.get(name=project_name)
        response = HttpResponse(content_type='application/json')
        if cluster.status == Cluster.CLUSTER_STATUS_READY or cluster.status == Cluster.CLUSTER_STATUS_INSTALLING:
            return Response(data={'msg': ': 集群未创建'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        cluster_monitor = ClusterMonitor(cluster)
        try:
            result = cluster_monitor.list_pod_status(namespace)
        except Exception as e:
            logger.error(e, exc_info=True)
            return Response(data={'msg': ': 数据读取失败！'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        response.write(json.dumps(result))
        return response


class WebKubeCtrlToken(APIView):

    def get(self, request, *args, **kwargs):
        pk = kwargs.get('pk')
        cluster = get_object_or_404(Cluster, pk=pk)
        return JsonResponse({'token': cluster.get_webkubectl_token()})


class DashBoardView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        item_name = kwargs['item_name']
        cluster_data = []
        restart_pods = []
        warn_containers = []
        error_loki_containers = []
        error_pods = []

        if len(request.user.profile.items) == 0 and request.user.is_superuser is False:
            return Response(data={'data': cluster_data, 'warnContainers': warn_containers, 'restartPods': restart_pods,
                                  'errorLokiContainers': error_loki_containers, 'errorPods': error_pods})

        if project_name == 'all':
            item_ids = []
            if item_name == 'all':
                item_ids = Item.objects.all().values_list('id')
            else:
                item = Item.objects.get(name=item_name)
                item_ids.append(item.id)
            resourceIds = ItemResource.objects.filter(item_id__in=item_ids).values_list('resource_id')
            clusters = Cluster.objects.filter(~Q(status=Cluster.CLUSTER_STATUS_READY),
                                              ~Q(status=Cluster.CLUSTER_STATUS_INSTALLING),
                                              ~Q(status=Cluster.CLUSTER_STATUS_DELETING)).filter(id__in=resourceIds)
            for c in clusters:
                cluster_monitor = ClusterMonitor(c)
                res = cluster_monitor.list_cluster_data()
                if len(res) != 0:
                    restart_pods = restart_pods + res.get('restart_pods', [])
                    warn_containers = warn_containers + res.get('warn_containers', [])
                    error_loki_containers = error_loki_containers + res.get('error_loki_containers', [])
                    error_pods = error_pods + res.get('error_pods', [])
                    cluster_data.append(json.dumps(res))
            restart_pods = kubeops_api.cluster_monitor.quick_sort_pods(restart_pods)
            error_loki_containers = kubeops_api.cluster_monitor.quick_sort_error_loki_container(error_loki_containers)
        else:
            cluster = Cluster.objects.get(name=project_name)
            if cluster.status != Cluster.CLUSTER_STATUS_READY and cluster.status != Cluster.CLUSTER_STATUS_INSTALLING and cluster.status != Cluster.CLUSTER_STATUS_DELETING:
                cluster_monitor = ClusterMonitor(cluster)
                res = cluster_monitor.list_cluster_data()
                if len(res) != 0:
                    restart_pods = res.get('restart_pods', [])
                    warn_containers = res.get('warn_containers', [])
                    error_loki_containers = res.get('error_loki_containers', [])
                    error_pods = res.get('error_pods', [])
                    cluster_data.append(json.dumps(res))
        return Response(data={'data': cluster_data, 'warnContainers': warn_containers, 'restartPods': restart_pods,
                              'errorLokiContainers': error_loki_containers, 'errorPods': error_pods})


class SettingView(APIView):

    def get(self, request, *args, **kwargs):
        params = request.query_params
        tab = params.get("tab", None)
        return JsonResponse(Setting.get_settings(tab))

    def post(self, request, *args, **kwargs):
        params = request.query_params
        tab = params.get("tab", None)
        settings = request.data
        Setting.set_settings(settings, tab)
        return Response(settings, status=status.HTTP_201_CREATED)


class ClusterStorageView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        cluster = Cluster.objects.get(name=project_name)
        if cluster.status == Cluster.CLUSTER_STATUS_READY or cluster.status == Cluster.CLUSTER_STATUS_INSTALLING:
            return Response(data={'msg': ': 集群未创建'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        cluster_monitor = ClusterMonitor(cluster)
        result = cluster_monitor.list_storage_class()
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(result))
        return response


class ClusterEventView(APIView):

    def post(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        params = request.data
        result = log.es.search_event(params, project_name)
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(result))
        return response


class ClusterNamespaceView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        cluster = Cluster.objects.get(name=project_name)
        if cluster.status == Cluster.CLUSTER_STATUS_READY or cluster.status == Cluster.CLUSTER_STATUS_INSTALLING:
            return Response(data={'msg': ': 集群未创建'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        cluster_monitor = ClusterMonitor(cluster)
        result = cluster_monitor.list_namespace()
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(result))
        return response


class ClusterComponentView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        cluster = Cluster.objects.get(name=project_name)
        if cluster.status == Cluster.CLUSTER_STATUS_READY or cluster.status == Cluster.CLUSTER_STATUS_INSTALLING:
            return Response(data={'msg': ': 集群未创建'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        cluster_monitor = ClusterMonitor(cluster)
        result = cluster_monitor.get_component_status()
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(result))
        return response


class CheckNodeView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        cluster = Cluster.objects.get(name=project_name)
        if cluster.status == Cluster.CLUSTER_STATUS_READY or cluster.status == Cluster.CLUSTER_STATUS_INSTALLING:
            return Response(data={'msg': ': 集群未创建'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        kubeops_api.cluster_monitor.delete_unused_node(cluster)
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps({'msg': '检查完毕'}))
        return response


class SyncHostTimeView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        cluster = Cluster.objects.get(name=project_name)
        if cluster.status == Cluster.CLUSTER_STATUS_READY or cluster.status == Cluster.CLUSTER_STATUS_INSTALLING:
            return Response(data={'msg': ': 集群未创建'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        result = kubeops_api.cluster_monitor.sync_node_time(cluster)
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(result))
        return response


class ClusterNamespaceView(APIView):

    def get(self, request, *args, **kwargs):
        project_name = kwargs['project_name']
        cluster = Cluster.objects.get(name=project_name)
        if cluster.status == Cluster.CLUSTER_STATUS_READY or cluster.status == Cluster.CLUSTER_STATUS_INSTALLING:
            return Response(data={'msg': ': 集群未创建'}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
        cluster_monitor = ClusterMonitor(cluster)
        result = cluster_monitor.list_namespace()
        response = HttpResponse(content_type='application/json')
        response.write(json.dumps(result))
        return response
