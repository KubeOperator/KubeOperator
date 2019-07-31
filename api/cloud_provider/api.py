import json

from django.http import HttpResponse, JsonResponse
from rest_framework import viewsets
from rest_framework.views import APIView

from ansible_api.permissions import IsSuperUser
from cloud_provider import serializers, get_cloud_client
from cloud_provider.models import CloudProviderTemplate, Region


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


class CloudRegionView(APIView):

    def post(self, request):
        vars = request.data
        vars['provider'] = 'vsphere'
        client = get_cloud_client(vars)
        data = client.list_region()
        return HttpResponse(json.dumps(data))
