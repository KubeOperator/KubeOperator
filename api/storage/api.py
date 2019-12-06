from ansible_api.permissions import IsSuperUser
from storage import serializers
from storage.models import NfsStorage, CephStorage, ClusterCephStorage
from rest_framework.response import Response
from rest_framework import viewsets, status



class NfsStorageViewSet(viewsets.ModelViewSet):
    queryset = NfsStorage.objects.all()
    serializer_class = serializers.NfsStorageSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'


class CephStorageViewSet(viewsets.ModelViewSet):
    queryset = CephStorage.objects.all()
    serializer_class = serializers.CephStorageSerializer
    permission_classes = (IsSuperUser,)
    lookup_field = 'name'
    lookup_url_kwarg = 'name'

    def destroy(self, request, *args, **kwargs):
        ceph_storage = self.get_object()
        cluster_ceph_storage = ClusterCephStorage.objects.filter(ceph_storage_id=ceph_storage.id)
        if len(cluster_ceph_storage) > 0:
            return Response(data={'msg': '有集群使用此ceph 不可删除!'}, status=status.HTTP_400_BAD_REQUEST)
        else:
            return super().destroy(self, request, *args, **kwargs)