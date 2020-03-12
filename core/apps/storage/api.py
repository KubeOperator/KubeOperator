from ansible_api.permissions import IsSuperUser
from storage import serializers
from storage.models import NfsStorage, CephStorage, ClusterCephStorage
from rest_framework.response import Response
from rest_framework import viewsets, status
from kubeops_api.models.item import Item
from kubeops_api.models.item_resource import ItemResource




class NfsStorageViewSet(viewsets.ModelViewSet):
    queryset = NfsStorage.objects.all()
    serializer_class = serializers.NfsStorageSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def list(self, request, *args, **kwargs):
        if request.query_params.get('itemName'):
            itemName = request.query_params.get('itemName')
            item = Item.objects.get(name=itemName)
            resource_ids = ItemResource.objects.filter(item_id=item.id).values_list("resource_id")
            self.queryset = NfsStorage.objects.filter(id__in=resource_ids)
            return super().list(self, request, *args, **kwargs)
        else:
            return super().list(self, request, *args, **kwargs)


class CephStorageViewSet(viewsets.ModelViewSet):
    queryset = CephStorage.objects.all()
    serializer_class = serializers.CephStorageSerializer

    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        ceph_storage = self.get_object()
        cluster_ceph_storage = ClusterCephStorage.objects.filter(ceph_storage_id=ceph_storage.id)
        if len(cluster_ceph_storage) > 0:
            return Response(data={'msg': '有集群使用此ceph 不可删除!'}, status=status.HTTP_400_BAD_REQUEST)
        else:
            return super().destroy(self, request, *args, **kwargs)

    def list(self, request, *args, **kwargs):
        if request.query_params.get('itemName'):
            itemName = request.query_params.get('itemName')
            item = Item.objects.get(name=itemName)
            resource_ids = ItemResource.objects.filter(item_id=item.id).values_list("resource_id")
            self.queryset = CephStorage.objects.filter(id__in=resource_ids)
            return super().list(self, request, *args, **kwargs)
        else:
            return super().list(self, request, *args, **kwargs)