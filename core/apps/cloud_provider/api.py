import json
from django.http import HttpResponse
from django.shortcuts import get_object_or_404
from rest_framework import viewsets, status
from rest_framework.response import Response
from rest_framework.views import APIView
from ansible_api.permissions import IsSuperUser
from cloud_provider import serializers, get_cloud_client
from cloud_provider.compute_model import compute_models
from cloud_provider.models import CloudProviderTemplate, Region, Zone, Plan
from kubeops_api.models.cluster import Cluster


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
    lookup_value_regex = '[\u4e00-\u9fa50-9a-zA-Z._-]+'
    queryset = Region.objects.all()
    serializer_class = serializers.RegionSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        if instance.zone_size > 0:
            return Response(data={'msg': '区域: {} 下资源不为空'.format(instance.name)}, status=status.HTTP_400_BAD_REQUEST)
        return super().destroy(self, request, *args, **kwargs)


class ZoneViewSet(viewsets.ModelViewSet):
    lookup_value_regex = '[\u4e00-\u9fa50-9a-zA-Z._-]+'
    queryset = Zone.objects.all()
    serializer_class = serializers.ZoneSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        if instance.host_size > 0 or instance.has_plan():
            return Response(data={'msg': '可用区: {} 下资源不为空'.format(instance.name)}, status=status.HTTP_400_BAD_REQUEST)
        return super().destroy(self, request, *args, **kwargs)


class PlanViewSet(viewsets.ModelViewSet):
    lookup_value_regex = '[\u4e00-\u9fa50-9a-zA-Z._-]+'
    queryset = Plan.objects.all()
    serializer_class = serializers.PlanSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        query_set = Cluster.objects.filter(plan__name=instance.name)
        if len(query_set) > 0:
            return Response(data={'msg': '部署计划: {} 下资源不为空'.format(instance.name)}, status=status.HTTP_400_BAD_REQUEST)
        return super().destroy(self, request, *args, **kwargs)


class CloudRegionView(APIView):

    def post(self, request):
        vars = request.data.get('vars')
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
        region.vars['provider'] = region.template.name
        client = get_cloud_client(region.vars)
        data = client.list_zone(region.cloud_region)
        return HttpResponse(json.dumps(data))


class CloudFlavorView(APIView):

    def get(self, request, *args, **kwargs):
        region_name = kwargs.get('region')
        region = get_object_or_404(Region, name=region_name)
        client = get_cloud_client(region.vars)
        return HttpResponse(json.dumps(client.get_flavors(region.cloud_region)))
