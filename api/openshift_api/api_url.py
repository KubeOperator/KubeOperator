from rest_framework.routers import DefaultRouter
from rest_framework_nested import routers
from openshift_api import api

app_name = "openshift_api"
router = DefaultRouter()

router.register('clusters', api.ClusterViewSet, 'cluster')
router.register('storage', api.StorageViewSet, 'storage')

# 注册离线包路由
router.register('packages', api.PackageViewSet, 'package')
router.register('template', api.StorageTemplateViewSet, 'template')
router.register('host', api.HostViewSet, 'host')
router.register('volume', api.VolumeViewSet, 'volume')
router.register('setting', api.SettingViewSet, 'setting')
router.register('hostInfo', api.HostInfoViewSet, 'hostInfo')

cluster_router = routers.NestedDefaultRouter(router, r'clusters', lookup='cluster')
cluster_router.register(r'configs', api.ClusterConfigViewSet, 'cluster-config')
cluster_router.register(r'nodes', api.NodeViewSet, 'cluster-node')
cluster_router.register(r'roles', api.RoleViewSet, 'cluster-role')
cluster_router.register(r'executions', api.DeployExecutionViewSet, 'cluster-deploy-execution')

storage_router = routers.NestedDefaultRouter(router, r'storage', lookup='storage')
storage_router.register(r'nodes', api.StorageNodeViewSet, 'storage-node')

urlpatterns = [
              ] + router.urls + cluster_router.urls + storage_router.urls
