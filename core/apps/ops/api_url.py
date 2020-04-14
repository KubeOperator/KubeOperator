from rest_framework.routers import DefaultRouter
from ops.apis.script import ScriptViewSet

app_name = "ops"
router = DefaultRouter()

router.register('scripts', ScriptViewSet, 'scripts')

urlpatterns = [] + router.urls
