from rest_framework.routers import DefaultRouter
from storage import api

app_name = "storage"
router = DefaultRouter()

router.register('storage/nfs', api.NfsStorageViewSet, 'nfs')

urlpatterns = [
              ] + router.urls
