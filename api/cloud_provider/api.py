import json
from django.http import HttpResponse
from django.shortcuts import get_object_or_404
from rest_framework import viewsets
from rest_framework.views import APIView
from ansible_api.permissions import IsSuperUser
from cloud_provider import serializers, get_cloud_client
from cloud_provider.compute_model import compute_models
from cloud_provider.models import CloudProviderTemplate, Region, Zone, Plan


class CloudProviderTemplateViewSet(viewsets.ModelViewSet):
    queryset = CloudProviderTemplate.objects.all()
    serializer_class = serializers.CloudProviderTemplateSerializer
    permission_classes = (IsSuperUser,)
    http_method_names = ['get', 'head', 'options']
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def get_queryset(self):
        CloudProviderTemplate.lookup()
        return super().get_queryset()


class RegionViewSet(viewsets.ModelViewSet):
    queryset = Region.objects.all()
    serializer_class = serializers.RegionSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class ZoneViewSet(viewsets.ModelViewSet):
    queryset = Zone.objects.all()
    serializer_class = serializers.ZoneSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class PlanViewSet(viewsets.ModelViewSet):
    queryset = Plan.objects.all()
    serializer_class = serializers.PlanSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class CloudRegionView(APIView):

    def post(self, request):
        vars = request.data
        client = get_cloud_client(vars)
        data = client.list_region()
        return HttpResponse(json.dumps(data))


class ComputeModleView(APIView):

    def get(self, request, *args, **kwargs):
        return HttpResponse(json.dumps(compute_models))


class CloudZoneView(APIView):

    def get(self, request, *args, **kwargs):
        region_name = kwargs.get('region')
        region = get_object_or_404(Region, name=region_name)
        region.vars['provider'] = 'vsphere'
        client = get_cloud_client(region.vars)
        data = client.list_zone(region.cloud_region)
        return HttpResponse(json.dumps(data))
