import json
from django.http import HttpResponse
from django.shortcuts import get_object_or_404
from rest_framework import viewsets, status
from rest_framework.response import Response
from rest_framework.views import APIView
from cloud_provider import serializers, get_cloud_client
from cloud_provider.compute_model import compute_models
from cloud_provider.models import CloudProviderTemplate, Region, Zone, Plan
from kubeops_api.models.cluster import Cluster
from kubeops_api.models.item import Item
from kubeops_api.models.item_resource import ItemResource


class CloudProviderTemplateViewSet(viewsets.ModelViewSet):
    queryset = CloudProviderTemplate.objects.all()
    serializer_class = serializers.CloudProviderTemplateSerializer

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

    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        query_set = Cluster.objects.filter(plan__name=instance.name)
        if len(query_set) > 0:
            return Response(data={'msg': '部署计划: {} 下资源不为空'.format(instance.name)}, status=status.HTTP_400_BAD_REQUEST)
        return super().destroy(self, request, *args, **kwargs)

    def list(self, request, *args, **kwargs):
        if request.query_params.get('itemName'):
            itemName = request.query_params.get('itemName')
            item = Item.objects.get(name=itemName)
            resource_ids = ItemResource.objects.filter(item_id=item.id).values_list("resource_id")
            self.queryset = Plan.objects.filter(id__in=resource_ids)
            return super().list(self, request, *args, **kwargs)
        else:
            return super().list(self, request, *args, **kwargs)

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        self.perform_create(serializer)
        headers = self.get_success_headers(serializer.data)

        item_ids = request.data.get('item_id', None)
        if item_ids and len(item_ids) > 0:
            item_resources = []
            for item_id in item_ids:
                item_resources.append(ItemResource(item_id=item_id, resource_id=serializer.data['id'],
                                                   resource_type=ItemResource.RESOURCE_TYPE_PLAN))
            ItemResource.objects.bulk_create(item_resources)

        return Response(serializer.data, status=status.HTTP_201_CREATED, headers=headers)


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
        models = client.get_flavors(region.cloud_region)
        if not models:
            return Response(data={'msg': "没有合适的flavor规格！请添加大于4C8G60G的flavor."}, status=status.HTTP_400_BAD_REQUEST)
        return HttpResponse(json.dumps(models))

class CloudTemplateView(APIView):

    def get(self, request, *args, **kwargs):
        region_name = kwargs.get('region')
        region = get_object_or_404(Region, name=region_name)
        region.vars['provider'] = region.template.name
        client = get_cloud_client(region.vars)
        data = client.list_templates(region.cloud_region)
        return HttpResponse(json.dumps(data))
