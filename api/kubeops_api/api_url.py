from rest_framework.routers import DefaultRouter
from rest_framework_nested import routers
from kubeops_api import api
from django.urls import path

app_name = "kubeops_api"
router = DefaultRouter()

router.register('clusters', api.ClusterViewSet, 'cluster')

# 注册离线包路由
router.register('packages', api.PackageViewSet, 'package')
router.register('credential', api.CredentialViewSet, 'credential')
router.register('host', api.HostViewSet, 'host')
router.register('setting', api.SettingViewSet, 'setting')
router.register('auth', api.AuthViewSet, 'auth')
router.register('backupStorage', api.BackupStorageViewSet, 'backupStorage')

cluster_router = routers.NestedDefaultRouter(router, r'clusters', lookup='cluster')
cluster_router.register(r'configs', api.ClusterConfigViewSet, 'cluster-config')
cluster_router.register(r'nodes', api.NodeViewSet, 'cluster-node')
cluster_router.register(r'roles', api.RoleViewSet, 'cluster-role')
cluster_router.register(r'executions', api.DeployExecutionViewSet, 'cluster-deploy-execution')

urlpatterns = [
                  path('cluster/<uuid:pk>/download/', api.DownloadView.as_view()),
                  path('cluster/<uuid:pk>/token/', api.GetClusterTokenView.as_view()),
                  path('cluster/config', api.GetClusterConfigView.as_view()),
                  path('version/', api.VersionView.as_view())
              ] + router.urls + cluster_router.urls
