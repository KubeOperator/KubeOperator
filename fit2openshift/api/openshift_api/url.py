from rest_framework.routers import DefaultRouter
from rest_framework_nested import routers

<<<<<<< HEAD

=======
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
from . import api

app_name = "openshift_api"
router = DefaultRouter()

<<<<<<< HEAD

router.register('clusters', api.ClusterViewSet, 'cluster')
#注册离线包路由
router.register('packages', api.PackageViewSet, 'package')
=======
router.register('clusters', api.ClusterViewSet, 'cluster')
>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce

cluster_router = routers.NestedDefaultRouter(router, r'clusters', lookup='cluster')
cluster_router.register(r'nodes', api.NodeViewSet, 'cluster-node')
cluster_router.register(r'roles', api.RoleViewSet, 'cluster-role')
cluster_router.register(r'executions', api.DeployExecutionViewSet, 'cluster-deploy-execution')

urlpatterns = [

] + router.urls + cluster_router.urls

<<<<<<< HEAD
=======

>>>>>>> 9c76263301cfc6cf73a3338535563cc4b44211ce
