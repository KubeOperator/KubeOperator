from rest_framework.routers import DefaultRouter
from rest_framework_nested import routers
from kubeops_api import api
from django.urls import path
from django.conf.urls import url

from kubeops_api.apis import host
from kubeops_api.apis import item
from kubeops_api.apis import grade
from kubeops_api.apis import file

app_name = "kubeops_api"
router = DefaultRouter()

router.register('clusters', api.ClusterViewSet, 'cluster')
router.register('packages', api.PackageViewSet, 'package')
router.register('credential', api.CredentialViewSet, 'credential')
router.register('host', host.HostViewSet, 'host')
router.register('backupStorage', api.BackupStorageViewSet, 'backupStorage')
router.register('backupStrategy', api.BackupStrategyViewSet, 'backupStrategy')
router.register('clusterBackup', api.ClusterBackupViewSet, 'clusterBackup')
router.register('items', item.ItemViewSet, 'item')
router.register('item/profiles', item.ItemUserViewSet, 'item-profiles')

cluster_router = routers.NestedDefaultRouter(router, r'clusters', lookup='cluster')
cluster_router.register(r'configs', api.ClusterConfigViewSet, 'cluster-config')
cluster_router.register(r'nodes', api.NodeViewSet, 'cluster-node')
cluster_router.register(r'roles', api.RoleViewSet, 'cluster-role')
cluster_router.register(r'executions', api.DeployExecutionViewSet, 'cluster-deploy-execution')

urlpatterns = [
                  path('host/import/', host.HostImportAPIView.as_view()),
                  path('file/upload/', file.FileUploadAPIView.as_view()),
                  path('cluster/<uuid:pk>/download/', api.DownloadView.as_view()),
                  path('cluster/<uuid:pk>/token/', api.GetClusterTokenView.as_view()),
                  path('cluster/<uuid:pk>/webkubectl/token/', api.WebKubeCtrlToken.as_view()),
                  path('cluster/<cluster_name>/grade/', grade.GradeRetrieveAPIView.as_view()),
                  path('cluster/config', api.GetClusterConfigView.as_view()),
                  path('version/', api.VersionView.as_view()),
                  path('version/', api.VersionView.as_view()),
                  path('backupStorage/check', api.CheckStorageView.as_view()),
                  path('backupStorage/getBuckets', api.GetBucketsView.as_view()),
                  path('clusterBackup/<uuid:project_id>/', api.ClusterBackupList.as_view()),
                  path('clusterBackup/<uuid:id>/delete/', api.ClusterBackupDelete.as_view()),
                  path('clusterBackup/restore/', api.ClusterBackupRestore.as_view()),
                  path('cluster/<project_name>/health/<namespace>/', api.ClusterHealthView.as_view()),
                  path('cluster/<project_name>/component/', api.ClusterComponentView.as_view()),
                  path('cluster/<project_name>/namespace/', api.ClusterNamespaceView.as_view()),
                  path('cluster/<project_name>/storage/', api.ClusterStorageView.as_view()),
                  path('cluster/<project_name>/event/', api.ClusterEventView.as_view()),
                  path('cluster/<project_name>/checkNodes/', api.CheckNodeView.as_view()),
                  path('cluster/<project_name>/syncNodeTime/', api.SyncHostTimeView.as_view()),
                  path('clusterHealthHistory/<project_id>/', api.ClusterHealthHistoryView.as_view()),
                  path('dashboard/<project_name>/<item_name>/', api.DashBoardView.as_view()),
                  path('resource/<item_name>/', item.ItemResourceView.as_view()),
                  path('resource/item/clusters/', item.ItemResourceClusterView.as_view()),
                  path('resource/<item_name>/<resource_type>/', item.ResourceView.as_view()),
                  path('resource/<item_name>/<resource_type>/<resource_id>/', item.ItemResourceDeleteView.as_view()),
                  url('settings', api.SettingView.as_view(), name='settings'),
              ] + router.urls + cluster_router.urls
