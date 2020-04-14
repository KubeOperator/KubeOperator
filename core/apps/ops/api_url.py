from rest_framework.routers import DefaultRouter
from ops.apis.script import ScriptViewSet, ScriptExecutionViewSet

app_name = "ops"
router = DefaultRouter()

router.register('scripts', ScriptViewSet, 'scripts')
router.register('scripts/execution', ScriptExecutionViewSet, 'script-execution')

urlpatterns = [] + router.urls
