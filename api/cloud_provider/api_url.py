from django.urls import path
from rest_framework.routers import DefaultRouter
from rest_framework_nested import routers
from cloud_provider import api
from django.conf.urls import url

app_name = "cloudProvider_api"
router = DefaultRouter()

router.register('provider/template', api.CloudProviderTemplateViewSet, 'provider-template')
router.register('regions', api.RegionViewSet, 'regions')
router.register('zones', api.ZoneViewSet, 'zones')
router.register('plans', api.PlanViewSet, 'plans')

urlpatterns = [
                  url(r'cloud/region/', api.CloudRegionView.as_view(), name='cloud-region'),
                  path(r'cloud/<region>/zone/', api.CloudZoneView.as_view(), name='cloud-zone'),
              ] + router.urls
